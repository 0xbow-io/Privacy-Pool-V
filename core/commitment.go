package core

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// TODO: Add Documentation
var (
	RecoverCommitmentKeys = Program{
		Identity: "RecoverCommitmentKeys",
		Src: `
		template RecoverCommitmentKeys(){
            input signal privateKey;
            input signal saltPublicKey[2];

            output signal publicKey[2];
            output signal secretKey[2];
            output signal encryptionKey[2];

            var computedPublicKey[2] = BabyPrivToPubKey()(privateKey);

            publicKey <== computedPublicKey;
            secretKey <== Ecdh()(privateKey, computedPublicKey);
            encryptionKey <== Ecdh()(privateKey, saltPublicKey);
        }
	`}

	DecryptCommitment = Program{
		Identity: "DecryptCommitment",
		Src: `
		template DecryptCommitment(cipherLen, tupleLen){
            input signal encryptionKey[2];               // ecdh shared secret key
            input signal nonce;                          // nonce value for Poseidon decryption
            input signal ciphertext[cipherLen];          // encrypted commitment tuple

            output signal tuple[tupleLen];
            output signal hash;

            var decryptor[cipherLen-1] = PoseidonDecryptWithoutCheck(tupleLen)(
                [
                    ciphertext[0], ciphertext[1], ciphertext[2], ciphertext[3],
                    ciphertext[4], ciphertext[5], ciphertext[6]
                ],
                nonce,
                encryptionKey
            );

            var recovered[tupleLen];
            for (var i = 0; i < tupleLen; i++) {
                recovered[i] = decryptor[i];
            }
            tuple <== recovered;
            hash <== POSEIDON_STD(tupleLen)(recovered);
        }
		`,
	}
	// TODO: Fragment template into smaller components
	CommitmentOwnershipProof = Program{
		Identity: "CommitmentOwnershipProof",
		Src: `
		template CommitmentOwnershipProof(cipherLen, tupleLen){
            input signal scope;
            input signal privateKey;                // EdDSA private key
            input signal saltPublicKey[2];          // used to derive the encryptionKey
            input signal nonce;                     // nonce value used for Poseidon decryption
            input signal ciphertext[cipherLen];     // encrypted commitment tuple

            output signal value;
            output signal nullRoot;
            output signal commitmentRoot;
            output signal commitmentHash;

            //  [publicKey, secretKey, encryptionKey]
            var (
                publicKey[2],
                secretKey[2],
                encryptionKey[2]
            ) = RecoverCommitmentKeys()
                (
                    privateKey,
                    saltPublicKey
                );

            // null root is the computed root of all secrets/keys
            // that were involved with the commitment.
            // As it only contains private elements (aside
            // from the saltPublicKey which is public).
            // It's utilised as a nullifier to the commitmentRoot.
            nullRoot <== ComputeMerkleTreeRoot(3)(
                [
                    publicKey[0], publicKey[1],
                    secretKey[0], secretKey[1],
                    saltPublicKey[0], saltPublicKey[1],
                    encryptionKey[0], encryptionKey[1]
                ]);

            //  [value, scope, secret.x, secret.y]
            var (tuple[tupleLen],hash) = DecryptCommitment(cipherLen, tupleLen)(
                    encryptionKey, nonce, ciphertext
                );

            value <== tuple[0];
            commitmentHash <== hash;

            // Verify contents and
            // Compute commitment root
            // CommitmentRoot can be verified outside of the circuit
            // as ciphertext & commitmenthash are public values.
            var (
                    scopeEqCheck,
                    secret_xEqCheck,
                    secret_yEqCheck,
                    computedCommitmentRoot
                ) = (
                    IsEqual()([scope, tuple[1]]),        // match scope
                    IsEqual()([secretKey[0],tuple[2]]), // match secret.x component
                    IsEqual()([secretKey[1],tuple[3]]), // match secret.y component
                    // compute commitment root
                    ComputeMerkleTreeRoot(3)(
                    [
                        ciphertext[0], ciphertext[1],
                        ciphertext[2], ciphertext[3],
                        ciphertext[4], ciphertext[5],
                        ciphertext[6], hash
                    ])
                );

            // invalidate the root if ownership is invalid
            // necessary to invalidate membership proofs
            var ownershipValidityCheck = IsEqual()(
                    [scopeEqCheck + secret_xEqCheck + secret_yEqCheck, 3]
                );
            commitmentRoot <== computedCommitmentRoot * ownershipValidityCheck;
        }
		`,
	}
	CommitmentMembershipProof = Program{
		Identity: "CommitmentMembershipProof",
		Src: `
		template CommitmentMembershipProof(maxTreeDepth){
            input signal actualTreeDepth;
            input signal commitmentRoot;
            input signal index;
            input signal siblings[maxTreeDepth];

            output signal root;
            var computedRoot = LeanIMTInclusionProof( maxTreeDepth )(
                    commitmentRoot,index,siblings,actualTreeDepth
                );
            root <== computedRoot;
        }
		`,
	}
)

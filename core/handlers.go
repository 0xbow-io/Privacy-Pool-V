package core

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// TODO: Add Documentation
var (
	HandleExistingCommitment = Program{
		Identity: "HandleExistingCommitment",
		Src: `
		template HandleExistingCommitment(maxTreeDepth, cipherLen, tupleLen){
            input signal scope, stateRoot, actualTreeDepth;
            input signal privateKey, nonce;
            input signal saltPublicKey[2], ciphertext[cipherLen];
            input signal index, siblings[maxTreeDepth];

            // aggregate (nullRoot, commitmentRoot, hash, value)
            // into 1 array output signal
            output signal out[4];

            var (value, nullRoot, commitmentRoot, hash) = CommitmentOwnershipProof(cipherLen, tupleLen)(
                scope, privateKey, saltPublicKey, nonce, ciphertext
            );

            var computedStateRoot = CommitmentMembershipProof(maxTreeDepth)(
                actualTreeDepth, commitmentRoot, index, siblings
            );

            var isVoidCheck = IsZero()(value);
            var stateRootEqCheck = IsZero()(computedStateRoot - stateRoot);

            signal isInvalid <== NOR()(isVoidCheck, stateRootEqCheck);

            out[0] <== nullRoot;
            out[1] <== commitmentRoot * isInvalid;
            out[2] <== hash * isInvalid;
            out[3] <== value * ( 1- isInvalid);
        }
	`}

	HandleNewCommitment = Program{
		Identity: "HandleNewCommitment",
		Src: `
		template HandleNewCommitment(cipherLen, tupleLen){
            input signal scope;
            input signal privateKey, nonce;
            input signal saltPublicKey[2], ciphertext[cipherLen];

            // aggregate (nullRoot, commitmentRoot, hash, value)
            // into 1 array output signal
            output signal out[4];

            var (value, nullRoot, commitmentRoot, hash) = CommitmentOwnershipProof(cipherLen, tupleLen)(
                scope, privateKey, saltPublicKey, nonce, ciphertext
            );

            signal invalidCommitmentRootCheck <== IsZero()(commitmentRoot);

            // If invalid ownership then commitmentRoot is 0
            // don't stop the circuit.
            // Output nullRoot if so, as this nullifies the new commitment.
            out[0] <==  nullRoot * invalidCommitmentRootCheck;
            out[1] <== commitmentRoot;
            out[2] <== hash;
            // null value if commitmentRoot is invalid
            out[3] <== value * (1 - invalidCommitmentRootCheck);
        }
		`,
	}
)

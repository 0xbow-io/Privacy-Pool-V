package privacypool

import (
	_ "embed"

	"github.com/0xBow-io/privacy-pool-veritas/core"

	. "github.com/0xBow-io/veritas"
)

// PrivacyPoolCircuitPkg is the package containg all core circuit
// blocks required to build a complete PrivacyPool circuit
//
// You will still need to import the following circuit pkgs
// to meet dependencies:
// - BinSumCircuitPkg && BitifyCircuitPkg
// - CircuitUtilsPkg
// - ComparatorsCircuitPkg & SafeComparatorsCircuitPkg
// - PoseidonCircuitPkg with PoseidonDecrypt
// - MerkleTreeCircuitPkg
// - MultiplexerCircuitPkg
// - GatesCircuitPkg
// - BabyJubCircuitPkg
// - MontGomeryCircuitPkg
// - EcdhCircuitPkg
// - EscalarMulCircuitPkg

var PrivacyPoolCircuitPkg = CircuitPkg{
	TargetVersion: "2.2.0",
	Field:         "bn128",
	Programs: []Program{
		PrivacyPool,
		// Core Circuit Blocks
		core.RecoverCommitmentKeys,
		core.DecryptCommitment,
		core.CommitmentOwnershipProof,
		core.CommitmentMembershipProof,
		core.HandleExistingCommitment,
		core.HandleNewCommitment,
		core.PoseidonDecryptWithoutCheck,
		core.PoseidonDecryptIterations,
	},
}

// TODO:
// - Add Documentation
// - Fragment template into smaller components
// - Utilise Buses to tidy up signals
var (
	PrivacyPool = Program{
		Identity: "PrivacyPool",
		Src: `
		template PrivacyPool(maxTreeDepth, cipherLen, tupleLen, nExisting, nNew) {
            /// **** Public Signals ****

            // Scope is the domain identifier
            // i.e. Keccak256(chainID, contractAddress)
            input signal scope;
            // The depth of the State Tree
            // at which the merkleproofs
            // were generated
            input signal actualTreeDepth;

            input signal context;
            // external input values to existing commitments
            // external output values from new commitments
            input signal externIO[2];

            input signal existingStateRoot;
            input signal newSaltPublicKey[nNew][2];
            input signal newCiphertext[nNew][cipherLen];

            /// **** End Of Public Signals ****

            /// **** Private Signals ****

            input signal privateKey[nExisting+nNew];
            input signal nonce[nExisting+nNew];

            input signal exSaltPublicKey[nExisting][2];
            input signal exCiphertext[nExisting][cipherLen];
            input signal exIndex[nExisting];
            input signal exSiblings[nExisting][maxTreeDepth];

            /// **** End Of Private Signals ****

            output signal newNullRoot[nExisting+nNew];
            output signal newCommitmentRoot[nExisting+nNew];
            output signal newCommitmentHash[nExisting+nNew];

            // ensure that External Input & Output
            // fits within the 252 bits
            var n2bIO[2][252];
            n2bIO[0] = Num2Bits(252)(externIO[0]);
            n2bIO[1] = Num2Bits(252)(externIO[1]);

            signal _newNullRootOut[nNew+nExisting];
            signal _newCommitmentRootOut[nNew+nExisting];
            signal _newCommitmentHashOut[nNew+nExisting];

            // get ownership & membership proofs for existing commitments
            // and compute total sum
            signal totalEx[nExisting+1];
            totalEx[0] <== externIO[0];
            for (var i = 0; i < nExisting; i++) {
                var out[4] = HandleExistingCommitment(
                                maxTreeDepth,
                                cipherLen,
                                tupleLen
                            )(
                                scope,
                                existingStateRoot,
                                actualTreeDepth,
                                privateKey[i],
                                nonce[i],
                                exSaltPublicKey[i],
                                exCiphertext[i],
                                exIndex[i],
                                exSiblings[i]
                            );
                _newNullRootOut[i] <== out[0];
                _newCommitmentRootOut[i] <== out[1];
                _newCommitmentHashOut[i] <== out[2];
                totalEx[i+1] <== totalEx[i] + out[3];
            }

            // get ownership for new commitments
            // and compute total sum
            signal totalNew[nNew+1];
            totalNew[0] <== externIO[1];
            var k = nExisting; // offset for new commitments
            for (var i = 0; i < nNew; i++) {

                var out[4] = HandleNewCommitment(
                                cipherLen,
                                tupleLen
                            )(
                                scope,
                                privateKey[k],
                                nonce[k],
                                newSaltPublicKey[i],
                                newCiphertext[i]
                            );
                _newNullRootOut[k] <== out[0];
                _newCommitmentRootOut[k] <== out[1];
                _newCommitmentHashOut[k] <== out[2];
                totalNew[i+1] <== totalNew[i] + out[3];
                k++;
            }

            // lastly ensure that all total sums are equal
            signal sumEqCheck <== IsEqual()(
                                [
                                    totalEx[nExisting],
                                    totalNew[nNew]
                                ]
                            );
            sumEqCheck === 1;

            newNullRoot <== _newNullRootOut;
            newCommitmentRoot <== _newCommitmentRootOut;
            newCommitmentHash <== _newCommitmentHashOut;

            // constraint on context
            signal contextSqrd <== context * context;
        }
	`}
)

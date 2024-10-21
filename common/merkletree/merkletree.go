package merkletree

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// MerkleTreeCircuitPkg contains circuit blocks
// to support Merkle tree operations
//
// These pkgs are required but not included here:
// - BitifyCircuitPkg
// - PoseidonCircuitPkg
// - MultiplexerCircuitPkg
// - SafeComparatorsCircuitPkg
// - CircuitUtilsPkg
var MerkleTreeCircuitPkg = CircuitPkg{
	TargetVersion: "2.2.0",
	Field:         "bn128",
	Programs: []Program{
		MerkleGeneratePathIndices,
		LeanIMTInclusionProof,
		ComputeMerkleTreeRoot,
	},
}

// Program below are dependent on Poseidon hash function
// Will need to import PoseidonCircuitPkg.
var (
	// MerkleGeneratePathIndices generates the indices of the path to the root
	// of a Merkle tree
	MerkleGeneratePathIndices = Program{
		Identity: "MerkleGeneratePathIndices",
		Src: `
    		template MerkleGeneratePathIndices(levels) {
                var BASE = 2;

                input signal in;
                output signal out[levels];

                var m = in;
                var computedResults[levels];

                for (var i = 0; i < levels; i++) {
                    // circom's best practices suggests to avoid using <-- unless you
                    // are aware of what's going on. This is the only way to do modulo operation.
                    out[i] <-- m % BASE;
                    m = m \ BASE;

                    // Check that each output element is less than the base.
                    var computedIsOutputElementLessThanBase = SafeLessThan(3)([out[i], BASE]);
                    computedIsOutputElementLessThanBase === 1;

                    // Re-compute the total sum.
                    computedResults[i] = out[i] * (BASE ** i);
                }

                // Check that the total sum matches the index.
                var computedCalculateTotal = CalculateTotal(levels)(computedResults);

                computedCalculateTotal === in;
            }
        `,
	}
	// LeanIMTInclusionProof verifies the inclusion proof of a leaf in an
	// incremental Merkle tree
	LeanIMTInclusionProof = Program{
		Identity: "LeanIMTInclusionProof",
		Src: `
		template LeanIMTInclusionProof(maxDepth) {
            input signal leaf, leafIndex, siblings[maxDepth], actualDepth;
            output signal out;

            var indices[maxDepth] = MerkleGeneratePathIndices(maxDepth)(leafIndex);

            signal nodes[maxDepth + 1];
            signal roots[maxDepth];

            // let node = leaf
            nodes[0] <== leaf;

            var root = 0;
            for (var i = 0; i < maxDepth; i++) {
                var isDepth = IsEqual()([actualDepth, i]);

                roots[i] <== isDepth * nodes[i];
                root += roots[i];

                var c[2][2] = [ [nodes[i], siblings[i]], [siblings[i], nodes[i]] ];
                var childNodes[2] = MultiMux1(2)(c, indices[i]);
                nodes[i + 1] <== POSEIDON_STD(2)(childNodes);
            }

            var isDepth = IsEqual()([actualDepth, maxDepth]);
            out <== root + isDepth * nodes[maxDepth];
        }
		`,
	}
	// ComputeMerkleTreeRoot  computes the root of a Merkle tree
	// given the leaves of the tree
	ComputeMerkleTreeRoot = Program{
		Identity: "ComputeMerkleTreeRoot",
		Src: `
		template ComputeMerkleTreeRoot(levels) {
            var totalLeaves = 2 ** levels;
            var numLeafHashers = totalLeaves / 2;
            var numIntermediateHashers = numLeafHashers - 1;

            // Array of leaf values input to the circuit.
            input signal leaves[totalLeaves];

            // Output signal for the Merkle root that results from hashing all the input leaves.
            output signal root;

            // Total number of hashers used in constructing the tree, one less than the total number of leaves,
            // since each level of the tree combines two elements into one.
            var numHashers = totalLeaves - 1;
            var computedLevelHashers[numHashers];

            // Initialize hashers for the leaves, each taking two adjacent leaves as inputs.
            for (var i = 0; i < numLeafHashers; i++){
                computedLevelHashers[i] = POSEIDON_STD(2)([leaves[i*2], leaves[i*2+1]]);
            }

            // Initialize hashers for intermediate levels, each taking the outputs of two hashers from the previous level.
            var k = 0;
            for (var i = numLeafHashers; i < numLeafHashers + numIntermediateHashers; i++) {
                computedLevelHashers[i] = POSEIDON_STD(2)([computedLevelHashers[k*2], computedLevelHashers[k*2+1]]);
                k++;
            }

            // Connect the output of the final hasher in the array to the root output signal.
            root <== computedLevelHashers[numHashers-1];
        }
	`,
	}
)

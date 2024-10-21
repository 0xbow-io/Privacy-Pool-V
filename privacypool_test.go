package privacypool

import (
	"testing"

	"github.com/0xBow-io/privacy-pool-veritas/common/babyjub"
	"github.com/0xBow-io/privacy-pool-veritas/common/bit"
	"github.com/0xBow-io/privacy-pool-veritas/common/comparators"
	"github.com/0xBow-io/privacy-pool-veritas/common/ecdh"
	"github.com/0xBow-io/privacy-pool-veritas/common/logic"
	"github.com/0xBow-io/privacy-pool-veritas/common/merkletree"
	"github.com/0xBow-io/privacy-pool-veritas/common/multiplexer"
	"github.com/0xBow-io/privacy-pool-veritas/common/poseidon"
	"github.com/0xBow-io/privacy-pool-veritas/common/scalar"
	"github.com/0xBow-io/privacy-pool-veritas/common/utils"
	. "github.com/0xBow-io/veritas"
	"github.com/test-go/testify/require"
)

// Test to see if all pkgs have been considered
// Will stil fail if any warning messages are present
// Even though all pkgs exists.
func Test_Compile_PrivacyPool(t *testing.T) {
	var (
		lib  = NewEmptyLibrary()
		main = CircuitPkg{
			TargetVersion: "2.2.0",
			Field:         "bn128",
			Programs: []Program{
				{
					Identity: "main",
					Src:      "component main {public[scope, actualTreeDepth, context, externIO, existingStateRoot, newSaltPublicKey, newCiphertext]} = PrivacyPool(32, 7, 4, 2, 2);",
				},
			},
		}
	)
	defer lib.Burn()

	reports, err := lib.Compile(main,
		PrivacyPoolCircuitPkg,
		babyjub.BabyJubCircuitPkg,
		babyjub.MontGomeryCircuitPkg,
		poseidon.PoseidonCircuitPkg,
		bit.BinSumCircuitPkg,
		bit.BitifyCircuitPkg,
		comparators.ComparatorsCircuitPkg,
		comparators.SafeComparatorsCircuitPkg,
		ecdh.EcdhCircuitPkg,
		logic.GatesCircuitPkg,
		merkletree.MerkleTreeCircuitPkg,
		multiplexer.MultiplexerCircuitPkg,
		scalar.EscalarMulCircuitPkg,
		utils.CircuitUtilsPkg,
	)
	require.Nil(t, err)
	// TO-DO: Clean up Warning messages
	if reports != nil && len(reports) > 0 {
		println(reports.String())
		t.FailNow()
	}

}

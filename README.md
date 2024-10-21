> [!Caution]
> This project is in active development.
>
> Circuits are yet to be audited and are not production-ready.

# Rationale:

This repository compiles all of the Circom Circuit blocks
(i.e. templates, functions) required for Privacy Pool Zk circuit into [Veritas](https://github.com/0xbow-io/Veritas)
Circuit packages in order to improve auditability and maintainability of the Privacy Pool Zk Circuit,

# Contents:

## Common Circuit Packages:

Packages found in the common/ directory include commonly used circuits sourced from Circomlib, Zk-Kit or Maci.
**Only the Poseidon Package has been built-up from scratch.**

## Core Circuit Packages:

The core/ directory contains the core circuit blocks required for the Privacy Pool Zk Circuit.
These are dependent on the common/ circuit packages.

# How to Use:

You are able to import these packages into your own circuit library by simply importing the associated go-module:

```Bash
go get github.com/0xBow-io/privacy-pool-veritas
```

And then importing the individuals into your own circuit library:

```Go
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
)

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
	if reports != nil && len(reports) > 0 {
		println(reports.String())
		t.FailNow()
	}
}

```

## TODO:

-   [ ] Refactor & Clean up warning reports.
-   [ ] Fragment Bigger templates into smaller manageable templates.
-   [ ] Add logic & constraint tests for core Privacy Pool Circuit components.
-   [ ] Technical Documentation

## How to Contribute:

Submit a PR with your changes.
We will prioritise those that align with the TODO list or is a critical bug fix.
See CONTRIBUTING.md for more details.

All contributions are welcome and contributers will be credited & recognised.

package comparators

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// SafeComparators is the package containg all
// circuit blocks to support safe comparison operations
// @NOTICE: Bitfiy pkg is required but not included here
var SafeComparatorsCircuitPkg = CircuitPkg{
	TargetVersion: "2.2.0",
	Field:         "bn128",
	Programs: []Program{
		SafeLessThan,
		SafeLessEqThan,
		SafeGreaterThan,
		SafeGreaterEqThan,
	},
}

// Below Implementations are from zk-kit.circom
// See: https://github.com/privacy-scaling-explorations/zk-kit.circom
// packages/utils/src/safe-comparators.circom
var (
	// template for safely comparing if one input is less than another,
	// ensuring inputs are within a specified bit-length.
	SafeLessThan = Program{
		Identity: "SafeLessThan",
		Src: `
    		template SafeLessThan(n) {
                // Ensure the bit-length does not exceed 252 bits.
                assert(n <= 252);

                input signal in[2];
                output signal out;

                // Additional conversion to handle arithmetic operation and capture the comparison result.
                var n2b[254];
                n2b = Num2Bits_strict()(in[0] + (1<<n) - in[1]);

                // Determine if in[0] is less than in[1] based on the most significant bit.
                out <== 1 - n2b[n];
            }
	`}
	// template to check if one input is less than or equal to another.
	SafeLessEqThan = Program{
		Identity: "SAFE_LESS_EQ_THAN",
		Src: `
		    template SafeLessEqThan(n) {
                input signal in[2];
                output signal out;

                // Use SafeLessThan to determine if in[0] is less than in[1] + 1.
                out <== SafeLessThan(n)([in[0], in[1] + 1]);
            }
	`}
	// template for safely comparing if one input is greater than another.
	SafeGreaterThan = Program{
		Identity: "SafeGreaterThan",
		Src: `
            template SafeGreaterThan(n) {
                // Two inputs to compare.
                input signal in[2];
                // Output signal indicating comparison result.
                output signal out;

                // Invert the inputs for SafeLessThan to check if in[1] is less than in[0].
                out <== SafeLessThan(n)([in[1], in[0]]);
            }
	`}
	// template for safely comparing if one input is greater than another.
	SafeGreaterEqThan = Program{
		Identity: " SafeGreaterEqThan",
		Src: `
			template SafeGreaterEqThan(n) {
                // Two inputs to compare.
                input signal in[2];
                // Output signal indicating comparison result.
                output signal out;

                // Invert the inputs and adjust for equality in SafeLessThan to
                // check if in[1] is less than or equal to in[0].
                out <== SafeLessThan(n)([in[1], in[0] + 1]);
            }
	`}
)

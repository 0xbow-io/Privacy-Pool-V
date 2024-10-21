package comparators

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// ComparatorsCircuitPkg is the package containg all
// circuit blocks to support comparison operations
// @NOTICE: Bitify & Binsum pkgs are required but not included here
var ComparatorsCircuitPkg = CircuitPkg{
	TargetVersion: "2.2.0",
	Field:         "bn128",
	Programs: []Program{
		IsZero,
		IsEqual,
		ForceEqualIfEnabled,
		LessThan,
		LessEqThan,
		GreaterThan,
		GreaterEqThan,
	},
}

// Below Implementations are from circomlib
// See: https://github.com/iden3/circomlib/
// circuits/comparators.circom
var (
	IsZero = Program{
		Identity: "IsZero",
		Src: `
		template IsZero() {
            input signal in;
            output signal out;

            signal inv;

            inv <-- in!=0 ? 1/in : 0;

            out <== -in*inv +1;
            in*out === 0;
        }
	`}
	IsEqual = Program{
		Identity: "IsEqual",
		Src: `
        template IsEqual() {
            input signal in[2];
            output signal out;

            component isz = IsZero();

            in[1] - in[0] ==> isz.in;

            isz.out ==> out;
        }
	`}
	ForceEqualIfEnabled = Program{
		Identity: "ForceEqualIfEnabled",
		Src: `
		template ForceEqualIfEnabled() {
            input signal enabled;
            input signal in[2];

            component isz = IsZero();

            in[1] - in[0] ==> isz.in;

            (1 - isz.out)*enabled === 0;
        }
	`}
	LessThan = Program{
		Identity: " LessThan",
		Src: `
		template LessThan(n) {
             assert(n <= 252);
            input signal in[2];
            output signal out;

            component n2b = Num2Bits(n+1);

            n2b.in <== in[0]+ (1<<n) - in[1];

            out <== 1-n2b.out[n];
        }
	`}
	LessEqThan = Program{
		Identity: " LessEqThan",
		Src: `
		template LessEqThan(n) {
      		input signal in[2];
            output signal out;

            component lt = LessThan(n);

            lt.in[0] <== in[0];
            lt.in[1] <== in[1]+1;
            lt.out ==> out;
        }
	`}
	GreaterThan = Program{
		Identity: " GreaterThan",
		Src: `
		template GreaterThan(n) {
            input signal in[2];
            output signal out;

            component lt = LessThan(n);

            lt.in[0] <== in[1];
            lt.in[1] <== in[0];
            lt.out ==> out;
        }
	`}
	GreaterEqThan = Program{
		Identity: " GreaterEqThan",
		Src: `
		template GreaterEqThan(n) {
      		input signal in[2];
            output signal out;

            component lt = LessThan(n);

            lt.in[0] <== in[1];
            lt.in[1] <== in[0]+1;
            lt.out ==> out;
        }
	`}
)

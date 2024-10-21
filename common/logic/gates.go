package logic

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// GatesCircuitPkg contains common logic gate circuits
var GatesCircuitPkg = CircuitPkg{
	TargetVersion: "2.2.0",
	Field:         "bn128",
	Programs: []Program{
		XOR,
		AND,
		OR,
		NOT,
		NAND,
		NOR,
		MultiAND,
	},
}

// Below Implementations are from circomlib
// See: https://github.com/iden3/circomlib/
var (
	XOR = Program{
		Identity: "XOR",
		Src: `
		template XOR() {
            input signal a;
            input signal b;
            output signal out;

            out <== a + b - 2*a*b;
        }
	`}
	AND = Program{
		Identity: "AND",
		Src: `
		template AND() {
            input signal a;
            input signal b;
            output signal out;

            out <== a*b;
        }
	`}
	OR = Program{
		Identity: "OR",
		Src: `
		template OR() {
            input signal a;
            input signal b;
            output signal out;

            out <== a + b - a*b;
        }
	`}
	NOT = Program{
		Identity: "NOT",
		Src: `
		template NOT() {
            input signal in;
            output signal out;

            out <== 1 + in - 2*in;
        }
	`}
	NAND = Program{
		Identity: "NAND",
		Src: `
		template NAND() {
            input signal a;
            input signal b;
            output signal out;

            out <== 1 - a*b;
        }
	`}
	NOR = Program{
		Identity: "NOR",
		Src: `
        template NOR() {
            input signal a;
            input signal b;
            output signal out;

            out <== a*b + 1 - a - b;
        }
	`}
	MultiAND = Program{
		Identity: "MultiAND",
		Src: `
		template MultiAND(n) {
            input signal in[n];
            output signal out;
            component and1;
            component and2;
            component ands[2];
            if (n==1) {
                out <== in[0];
            } else if (n==2) {
                and1 = AND();
                and1.a <== in[0];
                and1.b <== in[1];
                out <== and1.out;
            } else {
                and2 = AND();
                var n1 = n\2;
                var n2 = n-n\2;
                ands[0] = MultiAND(n1);
                ands[1] = MultiAND(n2);
                var i;
                for (i=0; i<n1; i++) ands[0].in[i] <== in[i];
                for (i=0; i<n2; i++) ands[1].in[i] <== in[n1+i];
                and2.a <== ands[0].out;
                and2.b <== ands[1].out;
                out <== and2.out;
            }
        }
	`}
)

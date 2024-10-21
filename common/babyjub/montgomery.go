package babyjub

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// MontGomeryCircuitPkg contains circuit blocks
// to suppor Edwards & Montgomery curve operationss
var MontGomeryCircuitPkg = CircuitPkg{
	TargetVersion: "2.2.0",
	Field:         "bn128",
	Programs: []Program{
		Edwards2Montgomery,
		Montgomery2Edwards,
		MontgomeryAdd,
		MontgomeryDouble,
	},
}

// Below Implementations are from circomlib
// See: https://github.com/iden3/circomlib/
var (
	Edwards2Montgomery = Program{
		Identity: "Edwards2Montgomery",
		Src: `
		template Edwards2Montgomery() {
            input signal in[2];
            output signal out[2];

            out[0] <-- (1 + in[1]) / (1 - in[1]);
            out[1] <-- out[0] / in[0];


            out[0] * (1-in[1]) === (1 + in[1]);
            out[1] * in[0] === out[0];
        }
    `}

	Montgomery2Edwards = Program{
		Identity: "Montgomery2Edwardss",
		Src: `
		template Montgomery2Edwards() {
            input signal in[2];
            output signal out[2];

            out[0] <-- in[0] / in[1];
            out[1] <-- (in[0] - 1) / (in[0] + 1);

            out[0] * in[1] === in[0];
            out[1] * (in[0] + 1) === in[0] - 1;
        }
    `}

	MontgomeryAdd = Program{
		Identity: "MontgomeryAdd",
		Src: `
		template MontgomeryAdd() {
            input signal in1[2];
            input signal in2[2];
            output signal out[2];

            var a = 168700;
            var d = 168696;

            var A = (2 * (a + d)) / (a - d);
            var B = 4 / (a - d);

            signal lamda;

            lamda <-- (in2[1] - in1[1]) / (in2[0] - in1[0]);
            lamda * (in2[0] - in1[0]) === (in2[1] - in1[1]);

            out[0] <== B*lamda*lamda - A - in1[0] -in2[0];
            out[1] <== lamda * (in1[0] - out[0]) - in1[1];
        }
    `}

	MontgomeryDouble = Program{
		Identity: "MontgomeryDouble",
		Src: `
		template MontgomeryDouble() {
            input signal in[2];
            output signal out[2];

            var a = 168700;
            var d = 168696;

            var A = (2 * (a + d)) / (a - d);
            var B = 4 / (a - d);

            signal lamda;
            signal x1_2;

            x1_2 <== in[0] * in[0];

            lamda <-- (3*x1_2 + 2*A*in[0] + 1 ) / (2*B*in[1]);
            lamda * (2*B*in[1]) === (3*x1_2 + 2*A*in[0] + 1 );

            out[0] <== B*lamda*lamda - A - 2*in[0];
            out[1] <== lamda * (in[0] - out[0]) - in[1];
        }
    `}
)

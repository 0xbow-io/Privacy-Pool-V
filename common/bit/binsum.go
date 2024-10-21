package bit

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// BinSum is the package containing the binsum template
// @NOTICE: No Dependencies required
var BinSumCircuitPkg = CircuitPkg{
	TargetVersion: "2.2.0",
	Field:         "bn128",
	Programs: []Program{
		Nbits,
		BinSum,
	},
}

// Below Implementations are from circomlib
// See: https://github.com/iden3/circomlib/
// circuits/binsum.circom
var (
	Nbits = Program{
		Identity: "nbits",
		Src: `
		function nbits(a) {
            var n = 1;
            var r = 0;
            while (n-1<a) {
                r++;
                n *= 2;
            }
            return r;
        }
	`}
	BinSum = Program{
		Identity: "BinSum",
		Src: `
		template BinSum(n, ops) {
            var nout = nbits((2**n -1)*ops);
            input signal in[ops][n];
            output signal out[nout];

            var lin = 0;
            var lout = 0;

            var k;
            var j;

            var e2;

            e2 = 1;
            for (k=0; k<n; k++) {
                for (j=0; j<ops; j++) {
                    lin += in[j][k] * e2;
                }
                e2 = e2 + e2;
            }

            e2 = 1;
            for (k=0; k<nout; k++) {
                out[k] <-- (lin >> k) & 1;

                // Ensure out is binary
                out[k] * (out[k] - 1) === 0;

                lout += out[k] * e2;

                e2 = e2+e2;
            }

            // Ensure the sum;

            lin === lout;
        }
	`}
)

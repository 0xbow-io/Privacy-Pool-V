package poseidon

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// Defenitions for Poseidon.circom
var (
	// SBox implements the SubWords function
	// S-box(x) = x^a where a>= 3, default is 5
	//
	// From https://eprint.iacr.org/2019/458.pdf:
	// S-box x^5 is suitable for two of the most popular prime fields in ZK applications,
	// concretely the prime subfields of the scalar field of the BLS12-381 and BN254 curves
	SBOX = Program{
		Identity: "SBOX",
		Src: `
		template SBOX() {
            input signal x;
            output signal z;

            signal a <== x*x;
            signal b <== a*a;

            z <== b*x;
        }
	`}
	// MULTISBOX applies SBOX
	// for each element in x.
	MULTISBOX = Program{
		Identity: "MULTI_SBOX",
		Src: `
		template MULTI_SBOX(t) {
		    input signal x[t];
			output signal z[t];

			for (var i=0; i<t; i++) {
			     z[i] <== SBOX()(x[i]);
			}
        }
	`}
	// Arc implements AddRoundConstants
	// For every element in x add the corresponding element in c
	// t - state size
	// r - round number
	// x - state
	// c - round constants
	ARC = Program{
		Identity: "ARC",
		Src: `
		template ARC(C, r, t) {
		    input signal in[t];
            output signal out[t];

            for (var i=0; i<t; i++) {
                out[i] <== in[i] + C[i + r];
            }
        }
	`}
	// MIXM implements the MDS (maximum distance separable)
	// mixing with the matrix m / p
	MIXM = Program{
		Identity: "MIXM",
		Src: `
		template MIXM(m, t) {
            input signal x[t];
            output signal z[t];

            var mul = 0;
            for (var i=0; i<t; i++) {
                mul = 0;
                for (var j=0; j<t; j++) {
                    mul += m[j][i] * x[j];
                }
                z[i] <== mul;
            }
        }
	`}
	// MIXS implements the MDS (maximum distance separable)
	// mixing with the sparse-matrix (s)
	MIXS = Program{
		Identity: "MIXS",
		Src: `
		template MIXS(s, t, r) {
		    input signal x[t];
            output signal z[t];

            var sum = 0;
            for (var i=0; i<t; i++) {
                sum += s[(t*2-1)*r+i]*x[i];
            }

            for (var i=1; i<t; i++) {
                z[i] <== s[(t*2-1)*r+t+i-1] * x[0] + x[i];
            }
            z[0] <== sum;
        }
	`}
)

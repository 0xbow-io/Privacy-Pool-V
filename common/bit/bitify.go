package bit

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// BitifyCircuitPkg is the package containg all
// circuit blocks to support bit conversions
// @NOTICE: Comparators & Aliascheck pkgs are required but not included here
var BitifyCircuitPkg = CircuitPkg{
	TargetVersion: "2.2.0",
	Field:         "bn128",
	Programs: []Program{
		Num2Bits,
		Num2Bits_strict,
		Bits2Num,
		Bits2Num_strict,
		Num2BitsNeg,
	},
}

// Below Implementations are from circomlib
// See: https://github.com/iden3/circomlib/
// circuits/bitify.circom
var (
	Num2Bits = Program{
		Identity: "Num2Bits",
		Src: `
		template Num2Bits(n) {
            input signal in;
            output signal out[n];
            var lc1=0;

            var e2=1;
            for (var i = 0; i<n; i++) {
                out[i] <-- (in >> i) & 1;
                out[i] * (out[i] -1 ) === 0;
                lc1 += out[i] * e2;
                e2 = e2+e2;
            }

            lc1 === in;
        }
	`}
	Num2Bits_strict = Program{
		Identity: "Num2Bits_strict",
		Src: `
            template Num2Bits_strict() {
                input signal in;
                output signal out[254];

                component aliasCheck = AliasCheck();
                component n2b = Num2Bits(254);
                in ==> n2b.in;

                for (var i=0; i<254; i++) {
                    n2b.out[i] ==> out[i];
                    n2b.out[i] ==> aliasCheck.in[i];
                }
            }
	`}
	Bits2Num = Program{
		Identity: "Bits2Num",
		Src: `
			template Bits2Num(n) {
                input signal in[n];
                output signal out;
                var lc1=0;

                var e2 = 1;
                for (var i = 0; i<n; i++) {
                    lc1 += in[i] * e2;
                    e2 = e2 + e2;
                }

                lc1 ==> out;
            }
	`}
	Bits2Num_strict = Program{
		Identity: " Bits2Num_strict",
		Src: `
    		template Bits2Num_strict() {
                input signal in[254];
                output signal out;

                component aliasCheck = AliasCheck();
                component b2n = Bits2Num(254);

                for (var i=0; i<254; i++) {
                    in[i] ==> b2n.in[i];
                    in[i] ==> aliasCheck.in[i];
                }

                b2n.out ==> out;
            }
	`}
	Num2BitsNeg = Program{
		Identity: " Num2BitsNeg",
		Src: `
		     template Num2BitsNeg(n) {
                input signal in;
                output signal out[n];
                var lc1=0;

                component isZero;

                isZero = IsZero();

                var neg = n == 0 ? 0 : 2**n - in;

                for (var i = 0; i<n; i++) {
                    out[i] <-- (neg >> i) & 1;
                    out[i] * (out[i] -1 ) === 0;
                    lc1 += out[i] * 2**i;
                }

                in ==> isZero.in;
                lc1 + isZero.out * 2**n === 2**n - in;
            }
	`}
)

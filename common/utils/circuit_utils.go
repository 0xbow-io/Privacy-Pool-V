package utils

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// CircuitUtilsPkg contains common circuit
// utility templates & functions
// @NOTICE: Bitify pkg is required but not included here
var CircuitUtilsPkg = CircuitPkg{
	TargetVersion: "2.2.0",
	Field:         "bn128",
	Programs: []Program{
		AliasCheck,
		CompConstant,
		CalculateTotal,
	},
}

var (
	AliasCheck = Program{
		Identity: "AliasCheck",
		Src: `
		template AliasCheck() {
            input signal in[254];
            component  compConstant = CompConstant(-1);
            for (var i=0; i<254; i++) in[i] ==> compConstant.in[i];
            compConstant.out === 0;
        }
	`}
	CompConstant = Program{
		Identity: "CompConstant",
		Src: `
		template CompConstant(ct) {
            input signal in[254];
            output signal out;

            signal parts[127];
            signal sout;

            var clsb;
            var cmsb;
            var slsb;
            var smsb;

            var sum=0;

            var b = (1 << 128) -1;
            var a = 1;
            var e = 1;
            var i;

            for (i=0;i<127; i++) {
                clsb = (ct >> (i*2)) & 1;
                cmsb = (ct >> (i*2+1)) & 1;
                slsb = in[i*2];
                smsb = in[i*2+1];

                if ((cmsb==0)&&(clsb==0)) {
                    parts[i] <== -b*smsb*slsb + b*smsb + b*slsb;
                } else if ((cmsb==0)&&(clsb==1)) {
                    parts[i] <== a*smsb*slsb - a*slsb + b*smsb - a*smsb + a;
                } else if ((cmsb==1)&&(clsb==0)) {
                    parts[i] <== b*smsb*slsb - a*smsb + a;
                } else {
                    parts[i] <== -a*smsb*slsb + a;
                }

                sum = sum + parts[i];

                b = b -e;
                a = a +e;
                e = e*2;
            }

            sout <== sum;

            component num2bits = Num2Bits(135);

            num2bits.in <== sout;

            out <== num2bits.out[127];
        }
	`}

	/**
	 * From MACI (https://github.com/privacy-scaling-explorations/maci/)
	 *     packages/circuits/circom/utils/calculateTotal.circom
	 *
	 * Computes the cumulative sum of an array of n input signals.
	 * It iterates through each input, aggregating the sum up to that point,
	 * and outputs the total sum of all inputs. This template is useful for
	 * operations requiring the total sum of multiple signals, ensuring the
	 * final output reflects the cumulative total of the inputs provided.
	 */
	CalculateTotal = Program{
		Identity: "CalculateTotal",
		Src: `
        template CalculateTotal(n) {
            input signal nums[n];
            output signal sum;

            signal sums[n];
            sums[0] <== nums[0];

            for (var i = 1; i < n; i++) {
                sums[i] <== sums[i - 1] + nums[i];
            }

            sum <== sums[n - 1];
        }
    `}
)

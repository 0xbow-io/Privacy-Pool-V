package babyjub

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// BabyJubCircuitPkg contains circuit blocks
// to support babyjub curve operations
// @NOTICE: Bitify and Escalarmulfix pkgs
// are required but not included here
var BabyJubCircuitPkg = CircuitPkg{
	TargetVersion: "2.2.0",
	Field:         "bn128",
	Programs: []Program{
		BabyAdd,
		BabyDbl,
		BabyCheck,
		BabyPrivToPubKey,
	},
}

// Below Implementations are from circomlib
// See: https://github.com/iden3/circomlib/
var (
	BabyAdd = Program{
		Identity: "BabyAdd",
		Src: `
		template BabyAdd() {
            input signal x1;
            input signal y1;
            input signal x2;
            input signal y2;
            output signal xout;
            output signal yout;

            signal beta;
            signal gamma;
            signal delta;
            signal tau;

            var a = 168700;
            var d = 168696;

            beta <== x1*y2;
            gamma <== y1*x2;
            delta <== (-a*x1+y1)*(x2 + y2);
            tau <== beta * gamma;

            xout <-- (beta + gamma) / (1+ d*tau);
            (1+ d*tau) * xout === (beta + gamma);

            yout <-- (delta + a*beta - gamma) / (1-d*tau);
            (1-d*tau)*yout === (delta + a*beta - gamma);
    }`}

	BabyDbl = Program{
		Identity: "BabyDbl",
		Src: `
		template BabyDbl() {
		    input signal x;
            input signal y;
            output signal xout;
            output signal yout;

            component adder = BabyAdd();
            adder.x1 <== x;
            adder.y1 <== y;
            adder.x2 <== x;
            adder.y2 <== y;

            adder.xout ==> xout;
            adder.yout ==> yout;
        }
    `}

	BabyCheck = Program{
		Identity: "BabyCheck",
		Src: `
		template BabyCheck() {
            input signal x;
            input signal y;

            signal x2;
            signal y2;

            var a = 168700;
            var d = 168696;

            x2 <== x*x;
            y2 <== y*y;

            a*x2 + y2 === 1 + d*x2*y2;
        }
    `}

	// Taken from MACI project
	// replaces BabyPbk template
	BabyPrivToPubKey = Program{
		Identity: "BabyPrivToPubKey",
		Src: `
		template BabyPrivToPubKey() {
            // The base point of the BabyJubJub curve.
            var BASE8[2] = [
                5299619240641551281634865583518297030282874472190772894086521144482721001553,
                16950150798460657717958625567821834550301663161624707787222815936182638968203
            ];

            // Prime subgroup order 'l'.
            var l = 2736030358979909402780800718157159386076813972158567259200215660948447373041;

            input signal privKey;
            output signal pubKey[2];

            // Check if private key is in the prime subgroup order 'l'
            var isLessThan = LessThan(251)([privKey, l]);
            isLessThan === 1;

            // Convert the private key to bits.
            var computedPrivBits[253] = Num2Bits(253)(privKey);

            // Perform scalar multiplication with the basepoint.
            var computedEscalarMulFix[2] = EscalarMulFix(253, BASE8)(computedPrivBits);

            pubKey <== computedEscalarMulFix;
        }
    `}
)

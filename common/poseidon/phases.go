package poseidon

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// Definitions for the phases involved in
// Optimized Poseidon using static
// preprocessed constants (c, s, m, p)
//
// Based on neptunes/filecoin implementation:
// https://github.com/argumentcomputer/neptune/blob/main/spec/poseidonspec.pdf

var (
	// PRE_ROUND (Phase 0)
	PRE_ROUND = Program{
		Identity: "PRE_ROUND_PHASE",
		Src: `
		template PRE_ROUND_PHASE(C, S, M , P, t, phase, nRoundsF, nRoundsP) {
		    assert(phase == 0);

		    input signal x[t];
            output signal z[t];

            z <== ARC(C, 0, t)(x);
        }
	`}

	// PARTIAL_ROUND (Phase 2)
	// for r in nRoundsP:
	// apply sbox to first element in x
	// and add round constant to result of sbox
	// mix x & sparse matrix s
	PARTIAL_ROUND = Program{
		Identity: "PARTIAL_ROUND_PHASE",
		Src: `

		template _PRE_PARTIAL_ROUND(C, P, rt, t){
            input signal x[t];
      		output signal z[t];

            var sbox[t] = MULTI_SBOX(t)(x);
            var arc[t] = ARC(C,rt,t)(sbox);

            z <== MIXM(P, t)(arc);
		}

		template _PARTIAL_ROUND(C, S, t, r, nRoundsF){
            var _t = (nRoundsF/2+1)*t;

            input signal x[t];
    		output signal z[t];

            var mixs[t];
            mixs[0] = SBOX()(x[0]);
            mixs[0] = mixs[0] + C[_t+r];
            for (var i = 1; i < t; i++) {
                mixs[i] = x[i];
            }

            z <== MIXS(S, t, r)(mixs);
		}

		template PARTIAL_ROUND_PHASE(C, S, M , P, t, phase, nRoundsF, nRoundsP) {
            assert(phase == 2);

    		input signal x[t];
    		output signal z[t];

            var rounds[t] = _PRE_PARTIAL_ROUND(C, P, nRoundsF/2*t, t)(x);
    		for (var r = 0; r < nRoundsP; r++) {
                var _rounds[t] = _PARTIAL_ROUND(C, S, t, r, nRoundsF)(rounds);
                for (var i = 0; i < t; i++) {
                    rounds[i] = _rounds[i];
                }
            }

            z <== rounds;
        }
	`}

	// FULL_ROUND  (Phase 1 & 3)
	// recursive full round function over x
	// where S-Box is applied to all elements in x
	// ARC is applied to x with round-constants c
	// and x is mixed with m
	FULL_ROUND = Program{
		Identity: "FULL_ROUND_PHASE",
		Src: `
		template _FULL_ROUND(C, M, r, t){
            input signal x[t];
      		output signal z[t];

            var sbox[t] = MULTI_SBOX(t)(x);
            var arc[t] = ARC(C, r, t)(sbox);
            z <== MIXM(M, t)(arc);
		}

		template FULL_ROUND_PHASE(C, S, M , P, t, phase, nRoundsF, nRoundsP) {
		    assert(phase == 1 || phase == 3);
      		var rounds = nRoundsF/2-1;

      		input signal x[t];
      		output signal z[t];

            var round[t] = x;
            var _r_ = 0;
      		for (var r = 0; r < rounds; r++) {
                if (phase == 1) {
                    _r_ =(r+1)*t;
                } else if (phase == 3) {
                    _r_ = (nRoundsF/2+1)*t+nRoundsP+r*t;
                }
                var _round[t] = _FULL_ROUND(C, M, _r_, t)(round);
                for (var i = 0; i < t; i++) {
                    round[i] = _round[i];
                }
            }

            z <== round;
        }
	`}

	// FINAL_ROUND (Phase 4)
	FINAL_ROUND = Program{
		Identity: "FINAL_ROUND_PHASE",
		Src: `
		template FINAL_ROUND_PHASE(C, S, M , P, t, phase, nRoundsF, nRoundsP) {
		    assert(phase == 4);

		    input signal x[t];
            output signal z[t];

            var sbox[t] = MULTI_SBOX(t)(x);
            z <==  MIXM(M, t)(sbox);
        }
	`}
)

package poseidon

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

//go:embed constants.circom
var static_constants string

// PoseidonCircuitPkg is the package containg all circuit
// blocks tu support the Poseidon hash function
// Import this package to use Poseidon hash function in your circuit
var PoseidonCircuitPkg = CircuitPkg{
	TargetVersion: "2.2.0",
	Field:         "bn128",
	Programs: []Program{
		POSEIDON_STD,
		POSEIDON_HASH,
		STATE_PERMUTATION,
		// Phases
		PRE_ROUND,
		PARTIAL_ROUND,
		FULL_ROUND,
		FINAL_ROUND,
		// rounds-function
		MULTISBOX,
		SBOX,
		ARC,
		MIXM,
		MIXS,
		// Constants
		POSEIDON_STATIC_CONSTANTS,
		STATIC_ROUND_CONSTANTS,
	},
}

var (
	POSEIDON_STATIC_CONSTANTS = Program{
		Identity: "POSEIDON_STATIC_CONSTANTS",
		Src:      static_constants,
	}

	// STATIC_ROUND_CONSTANTS: The number of partial & full rounds
	// for varying state width
	// from https://eprint.iacr.org/2019/458.pdf
	STATIC_ROUND_CONSTANTS = Program{
		Identity: "STATIC_ROUND_CONSTANTS",
		Src: `
        function N_ROUNDS_F() { return 8; }
		function N_ROUNDS_P(t) {
		    assert(t < 18);
            var N_ROUNDS_P[16] = [56, 57, 56, 60, 60, 63, 64, 63, 60, 66, 60, 65, 70, 60, 64, 68];
            return N_ROUNDS_P[t-2];
        }
	`}

	// STATE_PERMUTATION implements the HADES based
	// Poseidon permutation function/strategy
	// from: https://eprint.iacr.org/2019/458.pdf
	//
	// State is sequentially processed through the different
	// rounds (see round_function.go for more details on the rounds)
	//
	// Assume state allready has the initial state value at index 0
	STATE_PERMUTATION = Program{
		Identity: "STATE_PERMUTATION",
		Src: `
        template STATE_PERMUTATION(t) {
            var nRoundsF = N_ROUNDS_F();
            var nRoundsP = N_ROUNDS_P(t);

            var C[t*nRoundsF + nRoundsP] = POSEIDON_C(t);
            var S[  nRoundsP  *  (t*2-1)  ]  = POSEIDON_S(t);
            var M[t][t] = POSEIDON_M(t);
            var P[t][t] = POSEIDON_P(t);

            input signal state[t];
            output signal out[t];

            signal pre_round[t] <== PRE_ROUND_PHASE(C,S,M,P,t,0,nRoundsF,nRoundsP)(state);
            signal first_half_full_rounds[t] <== FULL_ROUND_PHASE(C,S,M,P,t,1,nRoundsF,nRoundsP)(pre_round);
            signal partial_rounds[t] <== PARTIAL_ROUND_PHASE(C,S,M,P,t,2,nRoundsF,nRoundsP)(first_half_full_rounds);
            signal second_half_full_rounds[t] <== FULL_ROUND_PHASE(C,S,M,P,t,3,nRoundsF,nRoundsP)(partial_rounds);

            out <== FINAL_ROUND_PHASE(C,S,M,P,t,4,nRoundsF,nRoundsP)(second_half_full_rounds);
        }
		`,
	}
	// POSEIDON_HASH implements the Poseidon hash function
	// suppring varying expected inputs and desired outputs
	POSEIDON_HASH = Program{
		Identity: "POSEIDON_HASH",
		Src: `
        template POSEIDON_HASH(nIn, nOut) {
            var t = nIn + 1;

            input signal inputs[nIn];
            input signal domain;

            output signal hash[nOut];

            component state_permutator = STATE_PERMUTATION(t);
            state_permutator.state[0] <== domain;
            for (var i = 1; i < t; i++) {
                state_permutator.state[i] <== inputs[i-1];
            }

            for (var i=0; i<nOut; i++) {
                hash[i] <== state_permutator.out[i];
            }

            for (var i=nOut; i<t; i++) {
                _ <== state_permutator.out[i] * state_permutator.out[i];
            }
        }
		`,
	}
	// POSEIDON_STD implements the Poseidon hash function
	// with multi-width inputs but 1 output
	POSEIDON_STD = Program{
		Identity: "POSEIDON_STD",
		Src: `
        template POSEIDON_STD(n) {
            input signal inputs[n];
            output signal hash;

            component hasher =  POSEIDON_HASH(n,1);
            hasher.domain <== 0;

            for (var i=0; i<n; i++) {
                hasher.inputs[i] <== inputs[i];
            }
            hash <== hasher.hash[0];
        }
		`,
	}
)

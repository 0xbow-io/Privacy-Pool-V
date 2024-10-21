package multiplexer

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// MultiplexerCircuitPkg contains circuit blocks
// for various multiplexers
var MultiplexerCircuitPkg = CircuitPkg{
	TargetVersion: "2.2.0",
	Field:         "bn128",
	Programs: []Program{
		MultiMux1,
		Mux1,
		MultiMux3,
		Mux3,
	},
}

// Below Implementations are from circomlib
// See: https://github.com/iden3/circomlib/
var (
	MultiMux1 = Program{
		Identity: "MultiMux1",
		Src: `
		 template MultiMux1(n) {
            input signal c[n][2];  // Constants
            input signal s;   // Selector
            output signal out[n];

            for (var i=0; i<n; i++) {

                out[i] <== (c[i][1] - c[i][0])*s + c[i][0];

            }
        }
	`}
	Mux1 = Program{
		Identity: "Mux1",
		Src: `
		 template Mux1() {
            var i;
            input signal c[2];  // Constants
            input signal s;   // Selector
            output signal out;

            component mux = MultiMux1(1);

            for (i=0; i<2; i++) {
                mux.c[0][i] <== c[i];
            }

            s ==> mux.s;

            mux.out[0] ==> out;
        }
	`}
	MultiMux3 = Program{
		Identity: "MultiMux3",
		Src: `
		template MultiMux3(n) {
            input signal c[n][8];  // Constants
            input signal s[3];   // Selector
            output signal out[n];

            signal a210[n];
            signal a21[n];
            signal a20[n];
            signal a2[n];

            signal a10[n];
            signal a1[n];
            signal a0[n];
            signal a[n];

            // 4 constrains for the intermediary variables
            signal  s10;
            s10 <== s[1] * s[0];

            for (var i=0; i<n; i++) {

                a210[i] <==  ( c[i][ 7]-c[i][ 6]-c[i][ 5]+c[i][ 4] - c[i][ 3]+c[i][ 2]+c[i][ 1]-c[i][ 0] ) * s10;
                a21[i] <==  ( c[i][ 6]-c[i][ 4]-c[i][ 2]+c[i][ 0] ) * s[1];
                a20[i] <==  ( c[i][ 5]-c[i][ 4]-c[i][ 1]+c[i][ 0] ) * s[0];
                a2[i] <==  ( c[i][ 4]-c[i][ 0] );

                a10[i] <==  ( c[i][ 3]-c[i][ 2]-c[i][ 1]+c[i][ 0] ) * s10;
                a1[i] <==  ( c[i][ 2]-c[i][ 0] ) * s[1];
                a0[i] <==  ( c[i][ 1]-c[i][ 0] ) * s[0];
                    a[i] <==  ( c[i][ 0] );

                out[i] <== ( a210[i] + a21[i] + a20[i] + a2[i] ) * s[2] +
                            (  a10[i] +  a1[i] +  a0[i] +  a[i] );

            }
        }
	`}
	Mux3 = Program{
		Identity: "Mux3",
		Src: `
        template Mux3() {
            var i;
            input signal c[8];  // Constants
            input signal s[3];   // Selector
            output signal out;

            component mux = MultiMux3(1);

            for (i=0; i<8; i++) {
                mux.c[0][i] <== c[i];
            }

            for (i=0; i<3; i++) {
            s[i] ==> mux.s[i];
            }

            mux.out[0] ==> out;
        }
	`}
)

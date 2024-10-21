package ecdh

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

// EcdhCircuitPkg is the package containing the ecdh template
// @NOTICE: escalarmulany pkg is required but not included here
var EcdhCircuitPkg = CircuitPkg{
	TargetVersion: "2.2.0",
	Field:         "bn128",
	Programs: []Program{
		Ecdh,
	},
}

// Below Implementations are from zk-kit.circom
// See: https://github.com/privacy-scaling-explorations/zk-kit.circom
// packages/ecdh/src/ecdh.circom
var (
	// ECDH Is a a template which allows to generate a shared secret
	// from a private key and a public key on the baby jubjub curve
	// It is important that the private key is hashed and pruned first
	Ecdh = Program{
		Identity: "Ecdh",
		Src: `
		template Ecdh() {
            // the private key must pass through deriveScalar first
            input signal privateKey;
            input signal publicKey[2];

            output signal sharedKey[2];

            // convert the private key to its bits representation
            var out[253];
            out = Num2Bits(253)(privateKey);

            // multiply the public key by the private key
            var mulFix[2];
            mulFix = EscalarMulAny(253)(out, publicKey);

            // we can then wire the output to the shared secret signal
            sharedKey <== mulFix;
        }
	`}
)

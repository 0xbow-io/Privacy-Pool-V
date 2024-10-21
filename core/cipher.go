package core

import (
	_ "embed"

	. "github.com/0xBow-io/veritas"
)

var (
	// Poseidon Decryption
	// Forked from zk-kit implementation: https://github.com/privacy-scaling-explorations/zk-kit.circom/
	// packages/poseidon-cipher/src/poseidon-cipher.circom
	// requires comparators to be included
	PoseidonDecryptWithoutCheck = Program{
		Identity: "PoseidonDecryptWithoutCheck",
		Src: `
		template PoseidonDecryptWithoutCheck(length) {
            var decryptedLength = length;
            while (decryptedLength % 3 != 0) {
                decryptedLength++;
            }

            input signal ciphertext[decryptedLength+1];
            input signal nonce;
            input signal key[2];
            output signal decrypted[decryptedLength];

            component iterations = PoseidonDecryptIterations(length);
            iterations.nonce <== nonce;
            iterations.key[0] <== key[0];
            iterations.key[1] <== key[1];
            for (var i = 0; i < decryptedLength + 1; i++) {
                iterations.ciphertext[i] <== ciphertext[i];
            }

            for (var i = 0; i < decryptedLength; i ++) {
                decrypted[i] <== iterations.decrypted[i];
            }
        }
	`}

	PoseidonDecryptIterations = Program{
		Identity: "PoseidonDecryptIterations",
		Src: `
		template PoseidonDecryptIterations(l) {
            var decryptedLength = l;

            while (decryptedLength % 3 != 0) {
                   decryptedLength++;
            }

            input signal ciphertext[decryptedLength + 1];
            input signal nonce;
            input signal key[2];

            output signal decrypted[decryptedLength];
            output signal decryptedLast;


            var two128 = 2 ** 128;

            component lt = LessThan(252);
            lt.in[0] <== nonce;
            lt.in[1] <== two128;
            lt.out === 1;

            var n = (decryptedLength + 1) \ 3;

            component strategies[n + 1];

            strategies[0] = POSEIDON_HASH(3, 4);
            // TODO: Investigate applying a different domain
            strategies[0].domain <== 0;
            strategies[0].inputs[0] <== key[0];
            strategies[0].inputs[1] <== key[1];
            strategies[0].inputs[2] <== nonce + (l * two128);

            for (var i = 0; i < n; i ++) {
                for (var j = 0; j < 3; j ++) {
                    decrypted[i * 3 + j] <== ciphertext[i * 3 + j] - strategies[i].hash[j + 1];
                }
                strategies[i + 1] = POSEIDON_HASH(3, 4);
                strategies[i + 1].domain <== strategies[i].hash[0];
                for (var j = 0; j < 3; j ++) {
                    strategies[i + 1].inputs[j] <== ciphertext[i * 3 + j];
                }
            }

            decryptedLast <== strategies[n].hash[1];
        }
	`}
)

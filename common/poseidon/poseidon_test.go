package poseidon

import (
	"testing"

	. "github.com/0xBow-io/veritas"
	"github.com/test-go/testify/require"
)

// TODO: Extend Tests
func Test_Poseidon_STD(t *testing.T) {
	var (
		lib      = NewEmptyLibrary()
		test_pkg = CircuitPkg{
			TargetVersion: "2.2.0",
			Field:         "bn128",
			Programs: []Program{
				{
					Identity: "main",
					Src:      "component main {public[in]}= Test();",
				},
				{
					Identity: "Test",
					Src: `
					// Sampled from https://github.com/iden3/go-iden3-crypto/blob/master/poseidon/poseidon_test.go
					function GetTestCase(n){
    				    if (n == 0) {
                            return [1];
                       	} else if (n == 1) {
                            return [1,2];
                        } else if (n == 2) {
                            return [1,2,0,0,0];
                        } else {
                            return [0];
                        }
					}
    				function GetExpected(n){
        				if (n == 0) {
                            return 18586133768512220936620570745912940619677854269274689475585506675881198879027;
        				} else if (n == 1) {
                            return 7853200120776062878684798364095072458815029376092732009249414926327459813530;
                        } else if (n == 2) {
                            return 1018317224307729531995786483840663576608797660851238720571059489595066344487;
                        } else if (n == 3) {
                            return 15336558801450556532856248569924170992202208561737609669134139141992924267169;
                        } else if (n == 4) {
                            return 5811595552068139067952687508729883632420015185677766880877743348592482390548;
                        } else {
                            return 0;
                        }
    				}

    				template Test(){
    				    input signal in;
    					output signal out;

    					// test case 0
    					signal tc_0_out <== POSEIDON_STD(1)(GetTestCase(0));
    					0 === tc_0_out - GetExpected(0);

    					// test case 1
                        signal tc_1_out <== POSEIDON_STD(2)(GetTestCase(1));
                        0 === tc_1_out - GetExpected(1);

                        // test case 2
                        signal tc_2_out <== POSEIDON_STD(5)(GetTestCase(2));
                        0 === tc_2_out - GetExpected(2);

                        out <== in * in;
    				}`,
				},
			},
		}
	)
	defer lib.Burn()

	reports, err := lib.Compile(test_pkg, PoseidonCircuitPkg)
	require.Nil(t, err)

	if reports != nil && len(reports) > 0 {
		println(reports.String())
		t.FailNow()
	}

	evaluation, err := lib.Evaluate([]byte(`{"in":"1"}`))

	// Execute the circuit
	require.Nil(t, err)
	require.NotNil(t, evaluation)

	//	print(evaluation.String())

	// Check for any reports
	reports, err = lib.GetReports()
	require.Nil(t, err)
	require.Len(t, reports, 0)

	//ensure that all symbols are assigned to the correct witness value
	evaluation.AssignWitToSym()

	// Check that constraints are satisfied
	require.True(t, len(evaluation.SatisfiedConstraints()) > 0)
	require.Len(t, evaluation.UnSatisfiedConstraints(), 0)
}

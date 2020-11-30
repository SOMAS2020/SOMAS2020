package rules

import (
	"gonum.org/v1/gonum/mat"
	"testing"
)

func TestRegisterNewRule(t *testing.T) {
	registerTestRule()
	if _, ok := AvailableRules["Kinda Test Rule"]; !ok {
		t.Errorf("Global rule register unable to register new rules")
	}
}

func registerTestRule() {

	//A very contrived rule//
	name := "Kinda Test Rule"
	reqVar := []string{
		"number_of_islands_contributing_to_common_pool",
		"number_of_failed_forages",
		"number_of_broken_agreements",
		"max_severity_of_sanctions",
	}

	v := []float64{1, 0, 0, 0, -4, 0, -1, -1, 0, 2, 0, 0, 0, 1, -2, 0, 0, 1, 0, -1}
	CoreMatrix := mat.NewDense(4, 5, v)
	aux := []float64{2, 2, 2, 2}
	AuxiliaryVector := mat.NewVecDense(4, aux)

	RegisterNewRule(name, reqVar, *CoreMatrix, *AuxiliaryVector)
	// Check internal/clients/team3/client.go for an implementation of a basic evaluator for this rule
}

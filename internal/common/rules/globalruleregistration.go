package rules

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
)

// Init Registers all global scoped rules
func init() {
	registerDemoRule()
	registerRulesByMass()
}

// registerDemoRule Defines and registers demo rule
func registerDemoRule() {

	//A very contrived rule//
	name := "Kinda Complicated Rule"
	reqVar := []string{
		"number_of_islands_contributing_to_common_pool",
		"number_of_failed_forages",
		"number_of_broken_agreements",
		"max_severity_of_sanctions",
	}

	v := []float64{1, 0, 0, 0, -4, 0, -1, -1, 0, 2, 0, 0, 0, 1, -2, 0, 0, 1, 0, -1}
	CoreMatrix := mat.NewDense(4, 5, v)
	aux := []float64{1, 1, 2, 0}
	AuxiliaryVector := mat.NewVecDense(4, aux)

	RegisterNewRule(name, reqVar, *CoreMatrix, *AuxiliaryVector)
	// Check internal/clients/team3/client.go for an implementation of a basic evaluator for this rule
}

func registerRulesByMass() {
	ruleSpecs := []struct {
		name   string
		reqVar []string
		v      []float64
		aux    []float64
	}{
		{
			name: "inspect_ballot_rule",
			reqVar: []string{
				"no_islands_alive",
				"no_ballots_cast",
			},
			v:   []float64{1, -1, 0},
			aux: []float64{0},
		},
		{
			name: "inspect_allocation_rule",
			reqVar: []string{
				"no_islands_alive",
				"no_allocations_sent",
			},
			v:   []float64{1, -1, 0},
			aux: []float64{0},
		},
		{
			name: "check_taxation_rule",
			reqVar: []string{
				"island_tax_contribution",
				"expected_tax_contribution",
			},
			v:   []float64{1, -1, 0},
			aux: []float64{2},
		},
		{
			name: "check_allocation_rule",
			reqVar: []string{
				"island_allocation",
				"expected_allocation",
			},
			v:   []float64{1, -1, 0},
			aux: []float64{0},
		},
	}

	for _, rs := range ruleSpecs {
		rowLength := len(rs.reqVar) + 1
		if len(rs.v)%rowLength != 0 {
			panic(fmt.Sprintf("Rule '%v' was registered without correct matrix dimensions", rs.name))
		}
		nrows := len(rs.v) / rowLength
		CoreMatrix := mat.NewDense(nrows, rowLength, rs.v)
		AuxiliaryVector := mat.NewVecDense(nrows, rs.aux)
		_, err := RegisterNewRule(rs.name, rs.reqVar, *CoreMatrix, *AuxiliaryVector)
		if err != nil {
			panic(fmt.Sprintf("%v", err.Error()))
		}
	}
}

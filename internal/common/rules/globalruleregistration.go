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
	reqVar := []VariableFieldName{
		NumberOfIslandsContributingToCommonPool,
		NumberOfFailedForages,
		NumberOfBrokenAgreements,
		MaxSeverityOfSanctions,
	}

	v := []float64{1, 0, 0, 0, -4, 0, -1, -1, 0, 2, 0, 0, 0, 1, -2, 0, 0, 1, 0, -1}
	CoreMatrix := mat.NewDense(4, 5, v)
	aux := []float64{1, 1, 2, 0}
	AuxiliaryVector := mat.NewVecDense(4, aux)

	_, ruleErr := RegisterNewRule(name, reqVar, *CoreMatrix, *AuxiliaryVector, false)
	if ruleErr != nil {
		panic(ruleErr.Error())
	}
	// Check internal/clients/team3/client.go for an implementation of a basic evaluator for this rule
}

func registerRulesByMass() {
	ruleSpecs := []struct {
		name    string
		reqVar  []VariableFieldName
		v       []float64
		aux     []float64
		mutable bool
	}{
		{
			name: "inspect_ballot_rule",
			reqVar: []VariableFieldName{
				NumberOfIslandsAlive,
				NumberOfBallotsCast,
			},
			v:       []float64{1, -1, 0},
			aux:     []float64{0},
			mutable: false,
		},
		{
			name: "inspect_allocation_rule",
			reqVar: []VariableFieldName{
				NumberOfIslandsAlive,
				NumberOfAllocationsSent,
			},
			v:       []float64{1, -1, 0},
			aux:     []float64{0},
			mutable: false,
		},
		{
			name: "check_taxation_rule",
			reqVar: []VariableFieldName{
				IslandTaxContribution,
				ExpectedTaxContribution,
			},
			v:       []float64{1, -1, 0},
			aux:     []float64{2},
			mutable: false,
		},
		{
			name: "check_allocation_rule",
			reqVar: []VariableFieldName{
				IslandAllocation,
				ExpectedAllocation,
			},
			v:       []float64{1, -1, 0},
			aux:     []float64{0},
			mutable: false,
		},
		{
			name: "vote_called_rule",
			reqVar: []VariableFieldName{
				RuleSelected,
				VoteCalled,
			},
			v:       []float64{1, -1, 0},
			aux:     []float64{0},
			mutable: false,
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
		_, ruleError := RegisterNewRule(rs.name, rs.reqVar, *CoreMatrix, *AuxiliaryVector, rs.mutable)
		if ruleError != nil {
			panic(ruleError.Error())
		}
	}
}

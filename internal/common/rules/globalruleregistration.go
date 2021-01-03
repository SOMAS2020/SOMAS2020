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

// RawRuleSpecification allows a user to use the CompileRuleCase function to build a rule matrix
type RawRuleSpecification struct {
	Name    string
	ReqVar  []VariableFieldName
	Values  []float64
	Aux     []float64
	Mutable bool
}

func registerRulesByMass() {
	ruleSpecs := []RawRuleSpecification{
		{
			Name: "inspect_ballot_rule",
			ReqVar: []VariableFieldName{
				NumberOfIslandsAlive,
				NumberOfBallotsCast,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: false,
		},
		{
			Name: "allocations_made_rule",
			ReqVar: []VariableFieldName{
				AllocationRequestsMade,
				AllocationMade,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: false,
		},
		{
			Name: "judge_inspection_rule",
			ReqVar: []VariableFieldName{
				JudgeInspectionPerformed,
			},
			Values:  []float64{1, -1},
			Aux:     []float64{0},
			Mutable: false,
		},
		{
			Name: "check_taxation_rule",
			ReqVar: []VariableFieldName{
				IslandTaxContribution,
				ExpectedTaxContribution,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{2},
			Mutable: false,
		},
		{
			Name: "check_allocation_rule",
			ReqVar: []VariableFieldName{
				IslandAllocation,
				ExpectedAllocation,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: false,
		},
		{
			Name: "vote_called_rule",
			ReqVar: []VariableFieldName{
				RuleSelected,
				VoteCalled,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: false,
		},
		{
			Name: "iigo_economic_sanction_1",
			ReqVar: []VariableFieldName{
				IslandReportedResources,
				ConstSanctionAmount,
				TurnsLeftOnSanction,
			},
			Values:  []float64{0, 0, 1, 0, 0, 0, 0, 0},
			Aux:     []float64{1, 4},
			Mutable: true,
		},
		{
			Name: "iigo_economic_sanction_2",
			ReqVar: []VariableFieldName{
				IslandReportedResources,
				ConstSanctionAmount,
				TurnsLeftOnSanction,
			},
			Values:  []float64{0, 0, 1, 0, 0.1, 1, 0, 0},
			Aux:     []float64{1, 4},
			Mutable: true,
		},
		{
			Name: "iigo_economic_sanction_3",
			ReqVar: []VariableFieldName{
				IslandReportedResources,
				ConstSanctionAmount,
				TurnsLeftOnSanction,
			},
			Values:  []float64{0, 0, 1, 0, 0.3, 1, 0, 0},
			Aux:     []float64{1, 4},
			Mutable: true,
		},
		{
			Name: "iigo_economic_sanction_4",
			ReqVar: []VariableFieldName{
				IslandReportedResources,
				ConstSanctionAmount,
				TurnsLeftOnSanction,
			},
			Values:  []float64{0, 0, 1, 0, 0.5, 1, 0, 0},
			Aux:     []float64{1, 4},
			Mutable: true,
		},
		{
			Name: "iigo_economic_sanction_5",
			ReqVar: []VariableFieldName{
				IslandReportedResources,
				ConstSanctionAmount,
				TurnsLeftOnSanction,
			},
			Values:  []float64{0, 0, 1, 0, 0.8, 1, 0, 0},
			Aux:     []float64{1, 4},
			Mutable: true,
		},
		{
			Name: "check_sanction_rule",
			ReqVar: []VariableFieldName{
				SanctionPaid,
				SanctionExpected,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: true,
		},
	}

	for _, rs := range ruleSpecs {
		rowLength := len(rs.ReqVar) + 1
		if len(rs.Values)%rowLength != 0 {
			panic(fmt.Sprintf("Rule '%v' was registered without correct matrix dimensions", rs.Name))
		}
		nrows := len(rs.Values) / rowLength
		CoreMatrix := mat.NewDense(nrows, rowLength, rs.Values)
		AuxiliaryVector := mat.NewVecDense(nrows, rs.Aux)
		_, ruleError := RegisterNewRule(rs.Name, rs.ReqVar, *CoreMatrix, *AuxiliaryVector, rs.Mutable)
		if ruleError != nil {
			panic(ruleError.Error())
		}
	}
}

// CompileRuleCase allows an agent to quickly build a RuleMatrix using the RawRuleSpecification
func CompileRuleCase(spec RawRuleSpecification) (RuleMatrix, bool) {
	rowLength := len(spec.ReqVar) + 1
	if len(spec.Values)%rowLength != 0 {
		return RuleMatrix{}, false
	}
	nrows := len(spec.Values) / rowLength
	CoreMatrix := mat.NewDense(nrows, rowLength, spec.Values)
	AuxiliaryVector := mat.NewVecDense(nrows, spec.Aux)
	finalRuleMatrix := RuleMatrix{
		RuleName:          spec.Name,
		RequiredVariables: spec.ReqVar,
		ApplicableMatrix:  *CoreMatrix,
		AuxiliaryVector:   *AuxiliaryVector,
		Mutable:           spec.Mutable,
	}
	return finalRuleMatrix, true
}

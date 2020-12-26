package baseclient

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

type BaseJudge struct {
}

// PayPresident pays the President a salary.
func (j *BaseJudge) PayPresident(presidentSalary shared.Resources) (shared.Resources, bool) {
	// TODO Implement opinion based salary payment.
	return presidentSalary, true
}

// inspectHistoryInternal is the base implementation of InspectHistory.
func (j *BaseJudge) InspectHistory() (map[shared.ClientID]roles.EvaluationReturn, bool) {
	outputMap := map[shared.ClientID]roles.EvaluationReturn{}
	for _, entry := range gamestate.TurnHistory {
		variablePairs := entry.Pairs
		clientID := entry.ClientID
		var rulesAffected []string
		for _, variable := range variablePairs {
			valuesToBeAdded, err := PickUpRulesByVariable(variable.VariableName, rules.RulesInPlay)
			if err == nil {
				rulesAffected = append(rulesAffected, valuesToBeAdded...)
			}
			err = rules.UpdateVariable(variable.VariableName, variable)
			if err != nil {
				return map[shared.ClientID]roles.EvaluationReturn{}, false
			}
		}
		if _, ok := outputMap[clientID]; !ok {
			outputMap[clientID] = roles.EvaluationReturn{
				Rules:       []rules.RuleMatrix{},
				Evaluations: []bool{},
			}
		}
		tempReturn := outputMap[clientID]
		for _, rule := range rulesAffected {
			evaluation, err := rules.BasicBooleanRuleEvaluator(rule)
			if err != nil {
				return outputMap, false
			}
			tempReturn.Rules = append(tempReturn.Rules, rules.RulesInPlay[rule])
			tempReturn.Evaluations = append(tempReturn.Evaluations, evaluation)
		}
	}
	return outputMap, true
}

// DeclareSpeakerPerformance checks how well the speaker did their job.
func (j *BaseJudge) DeclareSpeakerPerformance(inspectBallot bool, conductedRole bool) (result bool, didRole bool) {
	// TODO: Implement opinion based Speaker performance declaration.
	return inspectBallot, conductedRole
}

// DeclarePresidentPerformance checks how well the president did their job.
func (j *BaseJudge) DeclarePresidentPerformance(inspectBallot bool, conductedRole bool) (result bool, didRole bool) {
	// TODO: Implement opinion based President performance declaration.
	return inspectBallot, conductedRole
}

// PickUpRulesByVariable returns a list of rule_id's which are affected by certain variables.
func PickUpRulesByVariable(variableName rules.VariableFieldName, ruleStore map[string]rules.RuleMatrix) ([]string, error) {
	var Rules []string
	if _, ok := rules.VariableMap[variableName]; ok {
		for k, v := range ruleStore {
			_, err := searchForVariableInArray(variableName, v.RequiredVariables)
			if err != nil {
				Rules = append(Rules, k)
			}
		}
		return Rules, nil
	} else {
		return []string{}, errors.Errorf("Variable name '%v' was not found in the variable cache", variableName)
	}
}

func searchForVariableInArray(val rules.VariableFieldName, array []rules.VariableFieldName) (int, error) {
	for i, v := range array {
		if v == val {
			return i, nil
		}
	}
	return 0, errors.Errorf("Not found")
}

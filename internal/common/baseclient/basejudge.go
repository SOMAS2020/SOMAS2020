package baseclient

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type BaseJudge struct {
}

// PayPresident pays the President a salary.
func (j *BaseJudge) PayPresident(presidentSalary shared.Resources) (shared.Resources, bool) {
	// TODO Implement opinion based salary payment.
	return presidentSalary, true
}

// inspectHistoryInternal is the base implementation of InspectHistory.
func (j *BaseJudge) InspectHistory(iigoHistory []shared.Accountability) (map[shared.ClientID]roles.EvaluationReturn, bool) {
	outputMap := map[shared.ClientID]roles.EvaluationReturn{}
	for _, entry := range iigoHistory {
		variablePairs := entry.Pairs
		clientID := entry.ClientID
		var rulesAffected []string
		for _, variable := range variablePairs {
			valuesToBeAdded, foundRules := PickUpRulesByVariable(variable.VariableName, rules.RulesInPlay)
			if foundRules {
				rulesAffected = append(rulesAffected, valuesToBeAdded...)
			}
			updatedVariable := rules.UpdateVariable(variable.VariableName, variable)
			if !updatedVariable {
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
		outputMap[clientID] = tempReturn
	}
	return outputMap, true
}

// PickUpRulesByVariable returns a list of rule_id's which are affected by certain variables.
func PickUpRulesByVariable(variableName rules.VariableFieldName, ruleStore map[string]rules.RuleMatrix) ([]string, bool) {
	var Rules []string
	if _, ok := rules.VariableMap[variableName]; ok {
		for k, v := range ruleStore {
			_, found := searchForVariableInArray(variableName, v.RequiredVariables)
			if !found {
				Rules = append(Rules, k)
			}
		}
		return Rules, true
	} else {
		// fmt.Sprintf("Variable name '%v' was not found in the variable cache", variableName)
		return []string{}, false
	}
}

func searchForVariableInArray(val rules.VariableFieldName, array []rules.VariableFieldName) (int, bool) {
	for i, v := range array {
		if v == val {
			return i, true
		}
	}
	return -1, false
}

func (j *BaseJudge) CallPresidentElection(turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.Plurality,
		IslandsToVote: allIslands,
		HoldElection:  true,
	}
	return electionsettings
}

func (j *BaseJudge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	return winner
}

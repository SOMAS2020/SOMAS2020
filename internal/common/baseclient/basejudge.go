package baseclient

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type BaseJudge struct {
}

// GetRuleViolationSeverity returns a custom map of named rules and how severe the sanction should be for transgressing them
// If a rule is not named here, the default sanction value added is 1
func (j *BaseJudge) GetRuleViolationSeverity() map[string]roles.IIGOSanctionScore {
	return map[string]roles.IIGOSanctionScore{}
}

// GetSanctionThresholds returns a custom map of sanction score thresholds for different sanction tiers
// For any unfilled sanction tiers will be filled with default values (given in judiciary.go)
func (j *BaseJudge) GetSanctionThresholds() map[roles.IIGOSanctionTier]roles.IIGOSanctionScore {
	return map[roles.IIGOSanctionTier]roles.IIGOSanctionScore{}
}

// PayPresident pays the President a salary.
func (j *BaseJudge) PayPresident(presidentSalary shared.Resources) (shared.Resources, bool) {
	// TODO Implement opinion based salary payment.
	return presidentSalary, true
}

// InspectHistory is the base implementation of evaluating islands choices the last turn.
func (j *BaseJudge) InspectHistory(iigoHistory []shared.Accountability) (map[shared.ClientID]roles.EvaluationReturn, bool) {
	outputMap := map[shared.ClientID]roles.EvaluationReturn{}
	for _, entry := range iigoHistory {
		variablePairs := entry.Pairs
		clientID := entry.ClientID
		var rulesAffected []string
		for _, variable := range variablePairs {
			valuesToBeAdded, foundRules := rules.PickUpRulesByVariable(variable.VariableName, rules.RulesInPlay, rules.VariableMap)
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

// GetPardonedIslands allows a client to check all the sanctions that are currently levied against any islands.
// By returning true for a particular position a pardon will be issued for that sanction,
// false will continue the sanction as is
func (j *BaseJudge) GetPardonedIslands(currentSanctions map[int][]roles.Sanction) map[int][]bool {
	return map[int][]bool{}
}

// HistoricalRetributionEnabled Setting this to true will instruct the judiciary to give the judge role the entire
// history cache to browse through for evaluations
func (j *BaseJudge) HistoricalRetributionEnabled() bool {
	return false
}

// CallPresidentElection is called by the judiciary to decide on power-transfer
func (j *BaseJudge) CallPresidentElection(turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.Plurality,
		IslandsToVote: allIslands,
		HoldElection:  true,
	}
	return electionsettings
}

// DecideNextPresident returns the ID of chosen next President
func (j *BaseJudge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	return winner
}

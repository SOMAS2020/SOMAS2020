package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type judge struct {
	// Base implementation
	*baseclient.BaseJudge
	// Our client
	c *client
}

// Override functions here, see president.go for examples
func (j *judge) PayPresident(salary shared.Resources) (shared.Resources, bool) {
	// Use the base implementation
	return j.BaseJudge.PayPresident(salary)
}

func (j *judge) InspectHistory(iigoHistory []shared.Accountability) (map[shared.ClientID]roles.EvaluationReturn, bool) {
	outMap := map[shared.ClientID]roles.EvaluationReturn{}
	if j.c.localPool < config.IIGOConfig.InspectHistoryActionCost {
		//	dummy evaluation map
		for _, entry := range iigoHistory {
			outMap[entry.ClientID] = roles.EvaluationReturn{
				Rules:       []rules.RuleMatrix{},
				Evaluations: []bool{},
			}
		}
		return outMap, true
	}
	for _, entry := range iigoHistory {
		variablePairs := entry.Pairs
		clientID := entry.ClientID
		var rulesAffected []string
		for _, variable := range variablePairs {
			valuesToBeAdded, foundRules := rules.PickUpRulesByVariable(variable.VariableName, rules.RulesInPlay)
			if foundRules {
				rulesAffected = append(rulesAffected, valuesToBeAdded...)
			}
			updatedVariable := rules.UpdateVariable(variable.VariableName, variable)
			if !updatedVariable {
				return map[shared.ClientID]roles.EvaluationReturn{}, false
			}
		}
		if _, ok := outMap[clientID]; !ok {
			outMap[clientID] = roles.EvaluationReturn{
				Rules:       []rules.RuleMatrix{},
				Evaluations: []bool{},
			}
		}
		if j.c.trustScore[clientID] > 80 {
			tempReturn := outMap[clientID]
			for _, rule := range rulesAffected {
				tempReturn.Rules = append(tempReturn.Rules, rules.RulesInPlay[rule])
				tempReturn.Evaluations = append(tempReturn.Evaluations, true)
			}
			outMap[clientID] = tempReturn
		} else {
			tempReturn := outMap[clientID]
			for _, rule := range rulesAffected {
				evaluation, err := rules.BasicBooleanRuleEvaluator(rule)
				if err != nil {
					return outMap, false
				}
				tempReturn.Rules = append(tempReturn.Rules, rules.RulesInPlay[rule])
				tempReturn.Evaluations = append(tempReturn.Evaluations, evaluation)
			}
			outMap[clientID] = tempReturn
		}
	}
	return outMap, true
}

func (j *judge) CallPresidentElection(turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.Plurality,
		IslandsToVote: allIslands,
		HoldElection:  true,
	}
	return electionsettings
	// return j.BaseJudge.CallPresidentElection(turnsInPower, allIslands)
}

func (j *judge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	// Naively choose group 0
	return shared.ClientID(0)
}

func (j *judge) GetRuleViolationSeverity() map[string]roles.IIGOSanctionScore {
	return j.BaseJudge.GetRuleViolationSeverity()
}

func (j *judge) GetSanctionThresholds() map[roles.IIGOSanctionTier]roles.IIGOSanctionScore {
	return j.BaseJudge.GetSanctionThresholds()
}

func (j *judge) GetPardonedIslands(currentSanctions map[int][]roles.Sanction) map[int][]bool {
	return j.BaseJudge.GetPardonedIslands(currentSanctions)
}

func (j *judge) HistoricalRetributionEnabled() bool {
	return j.BaseJudge.HistoricalRetributionEnabled()
}

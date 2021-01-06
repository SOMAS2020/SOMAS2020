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

// PayPresident pays the president's salary
func (j *judge) PayPresident() (shared.Resources, bool) {

	// Strategy: Pay the president the amount they are owed, no changing amount.
	return j.BaseJudge.PayPresident()
}

// InspectHistory returns an evaluation on whether islands have adhered to the rules for that turn as a boolean.
func (j *judge) InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]roles.EvaluationReturn, bool) {
	outMap := map[shared.ClientID]roles.EvaluationReturn{}

	// If we do not have sufficient budget to conduct the inspection,
	// then we will return an empty map with true evaluations.
	// if j.c.getLocalResources() < config.IIGOConfig.InspectHistoryActionCost {
	// 	// dummy evaluation map
	// 	for _, entry := range iigoHistory {
	// 		outMap[entry.ClientID] = roles.EvaluationReturn{
	// 			Rules:       []rules.RuleMatrix{},
	// 			Evaluations: []bool{},
	// 		}
	// 	}
	// 	return outMap, true
	// }

	// Else, carry out inspectHistory as base implementation.
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
		if _, ok := outMap[clientID]; !ok {
			outMap[clientID] = roles.EvaluationReturn{
				Rules:       []rules.RuleMatrix{},
				Evaluations: []bool{},
			}
		}

		// If the island's trustScore is above 80, then return true for all rule evaluations
		// without actually evaluating if the island(s) adhered to any rules.
		if j.c.trustScore[clientID] > 80 {
			tempReturn := outMap[clientID]
			for _, rule := range rulesAffected {
				tempReturn.Rules = append(tempReturn.Rules, rules.RulesInPlay[rule])
				tempReturn.Evaluations = append(tempReturn.Evaluations, true)
			}
			outMap[clientID] = tempReturn
		} else {
			// All other islands will be evaluated fairly using base implementation.
			tempReturn := outMap[clientID]
			for _, rule := range rulesAffected {
				evaluation := rules.EvaluateRule(rule)
				if evaluation.EvalError != nil {
					return outMap, false
				}
				tempReturn.Rules = append(tempReturn.Rules, rules.RulesInPlay[rule])
				tempReturn.Evaluations = append(tempReturn.Evaluations, evaluation.RulePasses)
			}
			outMap[clientID] = tempReturn
		}
	}
	return outMap, true
}

// CallPresidentElection sets the election settings for the next president election
func (j *judge) CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.Approval,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}

	// Base implementation calls an election if monitoring was performed and the result was negative
	// or if the number of turnsInPower exceeds 3
	if monitoring.Performed && !monitoring.Result {
		electionsettings.HoldElection = true
	}
	// TODO: think if we want to change strategy here
	if turnsInPower >= 2 {
		electionsettings.HoldElection = true
	}
	return electionsettings
}

// DecideNextPresident declares who the next president will be
func (j *judge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	// If the election winner's trust score is high, we will declare them as the next President.
	// If not, we will replace it with the island who's trust score is the highest.
	if j.c.trustScore[winner] < 70 {
		// we can change this to be mailicious everytime
		for island := range j.c.trustScore {
			if j.c.trustScore[island] > j.c.trustScore[winner] {
				winner = island
			}
		}
	}
	return winner
}

// GetRuleViolationSeverity returns a custom map of named rules and
// how severe the sanction should be for transgressing them
// If a rule is not named here, the default sanction value added is 1
func (j *judge) GetRuleViolationSeverity() map[string]roles.IIGOSanctionScore {
	return map[string]roles.IIGOSanctionScore{}
}

// GetSanctionThresholds returns a custom map of sanction score thresholds for different sanction tiers
// For any unfilled sanction tiers will be filled with default values (given in judiciary.go)
func (j *judge) GetSanctionThresholds() map[roles.IIGOSanctionTier]roles.IIGOSanctionScore {
	return j.BaseJudge.GetSanctionThresholds()
}

// GetPardonedIslands decides which islands to pardon i.e. no longer impose sanctions on
// COMPULSORY: decide which islands, if any, to forgive
func (j *judge) GetPardonedIslands(currentSanctions map[int][]roles.Sanction) map[int][]bool {
	pardons := make(map[int][]bool)
	for key, sanctionList := range currentSanctions {
		lst := make([]bool, len(sanctionList))
		pardons[key] = lst
		for index, sanction := range sanctionList {
			if j.c.trustScore[sanction.ClientID] > 50 && j.c.params.friendliness > 40 {
				pardons[key][index] = true
			} else {
				pardons[key][index] = false
			}
		}
	}
	return pardons
}

// HistoricalRetributionEnabled enables historical retribution of inspection (automatically set to 3 turns ago)
func (j *judge) HistoricalRetributionEnabled() bool {
	return true
}

package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
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
	var PresidentSalary shared.Resources = 0
	PresidentSalaryRule, ok := j.GameState.RulesInfo.CurrentRulesInPlay["salary_cycle_president"]
	if ok {
		PresidentSalary = shared.Resources(PresidentSalaryRule.ApplicableMatrix.At(0, 1))
	}
	return PresidentSalary, true
}

// InspectHistory returns an evaluation on whether islands have adhered to the rules for that turn as a boolean.
func (j *judge) InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]shared.EvaluationReturn, bool) {

	if j.c.params.adv != nil {
		ret, bl, done := j.c.params.adv.InspectHistory(iigoHistory, turnsAgo)
		if done {
			return ret, bl
		}
	}

	outMap := map[shared.ClientID]shared.EvaluationReturn{}

	// Carry out base implementation of iterating through history to find rules to sanction
	for _, entry := range iigoHistory {
		variablePairs := entry.Pairs
		clientID := entry.ClientID
		copyOfGlobalVarCache := rules.CopyVariableMap(j.c.ServerReadHandle.GetGameState().RulesInfo.VariableMap)
		var rulesAffected []string
		for _, variable := range variablePairs {
			valuesToBeAdded, foundRules := rules.PickUpRulesByVariable(variable.VariableName, j.c.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay, copyOfGlobalVarCache)
			if foundRules {
				rulesAffected = append(rulesAffected, valuesToBeAdded...)
			}
			updatedVariable := rules.UpdateVariableInternal(variable.VariableName, variable, copyOfGlobalVarCache)
			if !updatedVariable {
				return map[shared.ClientID]shared.EvaluationReturn{}, false
			}
		}
		if _, ok := outMap[clientID]; !ok {
			outMap[clientID] = shared.EvaluationReturn{
				Rules:       []rules.RuleMatrix{},
				Evaluations: []bool{},
			}
		}

		// If the island's trustScore is above 80, then return true for all rule evaluations
		// without actually evaluating if the island(s) adhered to any rules.
		if j.c.trustScore[clientID] > 80 {
			tempReturn := outMap[clientID]
			for _, rule := range rulesAffected {
				tempReturn.Rules = append(tempReturn.Rules, j.c.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay[rule])
				tempReturn.Evaluations = append(tempReturn.Evaluations, true)
			}
			outMap[clientID] = tempReturn
		} else {
			// All other islands will be evaluated fairly using base implementation.
			tempReturn := outMap[clientID]
			for _, rule := range rulesAffected {
				evaluation := rules.EvaluateRuleFromCaches(rule, j.c.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay, copyOfGlobalVarCache)
				if evaluation.EvalError != nil {
					return outMap, false
				}
				tempReturn.Rules = append(tempReturn.Rules, j.c.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay[rule])
				if clientID == j.c.GetID() {
					tempReturn.Evaluations = append(tempReturn.Evaluations, true)
				} else {
					tempReturn.Evaluations = append(tempReturn.Evaluations, evaluation.RulePasses)
				}
			}
			outMap[clientID] = tempReturn
		}
	}
	return outMap, true
}

// CallPresidentElection sets the election settings for the next president election
func (j *judge) CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {

	if j.c.params.adv != nil {
		ret, done := j.c.params.adv.CallPresidentElection(monitoring, turnsInPower, allIslands)
		if done {
			return ret
		}
	}
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.InstantRunoff,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}

	// Base implementation calls an election if monitoring was performed and the result was negative
	// or if the number of turnsInPower exceeds 2
	if monitoring.Performed && !monitoring.Result {
		electionsettings.HoldElection = true
	}
	if turnsInPower >= 2 {
		electionsettings.HoldElection = true
	}
	return electionsettings
}

// DecideNextPresident declares who the next president will be
func (j *judge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	// If the election winner's trust score is high, we will declare them as the next President.
	// If not, we will replace it with the island who's trust score is the highest.
	if j.c.params.adv != nil {
		ret, done := j.c.params.adv.DecideNextPresident(winner)
		if done {
			return ret
		}
	}

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
func (j *judge) GetRuleViolationSeverity() map[string]shared.IIGOSanctionsScore {
	if j.c.params.adv != nil {
		ret, done := j.c.params.adv.GetRuleViolationSeverity()
		if done {
			return ret
		}
	}
	return map[string]shared.IIGOSanctionsScore{}
}

// GetSanctionThresholds returns a custom map of sanction score thresholds for different sanction tiers
// For any unfilled sanction tiers will be filled with default values (given in judiciary.go)
// All sanction tiers have linear scaling and hence it is easier (than default) to fall in higher
// sanction tiers. The aim is to enforce harsher penalties with hope to encourage more obedient
// behaviour from other agents.
func (j *judge) GetSanctionThresholds() map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore {

	// Linear increase of 5 per sanction tier.
	return map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
		shared.SanctionTier1: 1,
		shared.SanctionTier2: 6,
		shared.SanctionTier3: 11,
		shared.SanctionTier4: 16,
		shared.SanctionTier5: 21,
	}

}

// GetPardonedIslands decides which islands to pardon i.e. no longer impose sanctions on
// COMPULSORY: decide which islands, if any, to forgive
func (j *judge) GetPardonedIslands(currentSanctions map[int][]shared.Sanction) map[int][]bool {
	if j.c.params.adv != nil {
		ret, done := j.c.params.adv.GetPardonedIslands(currentSanctions)
		if done {
			return ret
		}
	}
	pardons := make(map[int][]bool)
	for key, sanctionList := range currentSanctions {
		lst := make([]bool, len(sanctionList))
		pardons[key] = lst
		for index, sanction := range sanctionList {
			if (j.c.trustScore[sanction.ClientID] >= 50 && j.c.params.friendliness > 0.5) || sanction.ClientID == j.c.GetID() {
				pardons[key][index] = true
			} else {
				pardons[key][index] = false
			}
		}
	}
	return pardons
}

// HistoricalRetributionEnabled enables historical retribution of inspection (automatically set to 3 turns ago)
// Strategy: If the rule is in play, we adhere to it else we will break it.
func (j *judge) HistoricalRetributionEnabled() bool {

	var ans bool = true
	res := rules.EvaluateRuleFromCaches("judge_historical_retribution_permission", j.GameState.RulesInfo.CurrentRulesInPlay, j.GameState.RulesInfo.VariableMap)
	if res.RulePasses && res.EvalError == nil {
		ans = false
	}
	return ans
}

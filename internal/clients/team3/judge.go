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
			// TODO: fix this PickUpRulesByVariable last argument to return a map
			valuesToBeAdded, foundRules := rules.PickUpRulesByVariable(variable.VariableName, rules.RulesInPlay, variablePairs)
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

func (j *judge) CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	// example implementation calls an election if monitoring was performed and the result was negative
	// or if the number of turnsInPower exceeds 3
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.Plurality,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}
	if monitoring.Performed && !monitoring.Result {
		electionsettings.HoldElection = true
	}
	// TODO: think if we want to change strategy here
	if turnsInPower >= 2 {
		electionsettings.HoldElection = true
	}
	return electionsettings
}

func (j *judge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
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

func (j *judge) GetRuleViolationSeverity() map[string]roles.IIGOSanctionScore {
	return map[string]roles.IIGOSanctionScore{}
}

func (j *judge) GetSanctionThresholds() map[roles.IIGOSanctionTier]roles.IIGOSanctionScore {
	return j.BaseJudge.GetSanctionThresholds()
}

func (j *judge) GetPardonedIslands(currentSanctions map[int][]roles.Sanction) map[int][]bool {
	return j.BaseJudge.GetPardonedIslands(currentSanctions)
}

func (j *judge) HistoricalRetributionEnabled() bool {
	return true
}

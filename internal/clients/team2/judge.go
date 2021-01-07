package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Judge struct {
	*baseclient.BaseJudge
	c *client
}

// GetRuleViolationSeverity returns a custom map of named rules and how severe the sanction should be for transgressing them
// If a rule is not named here, the default sanction value added is 1
// OPTIONAL: override to set custom sanction severities for specific rules
func (j *Judge) GetRuleViolationSeverity() map[string]roles.IIGOSanctionScore {
	return map[string]roles.IIGOSanctionScore{}
}

// GetSanctionThresholds returns a custom map of sanction score thresholds for different sanction tiers
// For any unfilled sanction tiers will be filled with default values (given in judiciary.go)
// OPTIONAL: override to set custom sanction thresholds
func (j *Judge) GetSanctionThresholds() map[roles.IIGOSanctionTier]roles.IIGOSanctionScore {
	return map[roles.IIGOSanctionTier]roles.IIGOSanctionScore{}
}

// PayPresident pays the President a salary.
// OPTIONAL: override to pay the President less than the full amount.
func (j *Judge) PayPresident() (shared.Resources, bool) {
	// TODO Implement opinion based salary payment.
	PresidentSalaryRule, ok := rules.RulesInPlay["salary_cycle_president"]
	var PresidentSalary shared.Resources = 0
	if ok {
		PresidentSalary = shared.Resources(PresidentSalaryRule.ApplicableMatrix.At(0, 1))
	}
	return PresidentSalary, true
}

// InspectHistory is the base implementation of evaluating islands choices the last turn.
// OPTIONAL: override if you want to evaluate the history log differently.
func (j *Judge) InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]roles.EvaluationReturn, bool) {
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
			ret := rules.EvaluateRule(rule)
			if ret.EvalError != nil {
				return outputMap, false
			}
			tempReturn.Rules = append(tempReturn.Rules, rules.RulesInPlay[rule])
			tempReturn.Evaluations = append(tempReturn.Evaluations, ret.RulePasses)
		}
		outputMap[clientID] = tempReturn
	}
	return outputMap, true
}

// GetPardonedIslands decides which islands to pardon i.e. no longer impose sanctions on
// Random Justice: there is a 33% chance each turn an island is forgiven
//TODO: change this strategy, an agent is forgiven after paying the sanction
//TODO: add the trustworthiness and state consideration
func (j *Judge) GetPardonedIslands(currentSanctions map[int][]roles.Sanction) map[int][]bool {
	Pardoned := make(map[int][]bool)
	for i, List := range currentSanctions {
		List2 := make([]bool, len(List))
		pardons[i] = List2
		for index, sanction := range List {
			if j.c.trustworthiness[sanction.ClientID] > 80 && MethodOfPlay() != 1 {
				e
				pardons[i][index] = true
			} else {
				pardons[i][index] = false
			}
		}
	}
	return Pardoned
}

// HistoricalRetributionEnabled allows you to punish more than the previous turns transgressions
// OPTIONAL: override if you want to punish historical transgressions
func (j *Judge) HistoricalRetributionEnabled() bool {
	return false
}

// CallPresidentElection is called by the judiciary to decide on power-transfer
// COMPULSORY: decide when to call an election following relevant rulesInPlay if you wish
func (j *Judge) CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	// example implementation calls an election if monitoring was performed and the result was negative
	// or if the number of turnsInPower exceeds 3
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.Runoff,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}
	if monitoring.Performed && !monitoring.Result {
		electionsettings.HoldElection = true
	}
	if turnsInPower >= 2 {
		electionsettings.HoldElection = true
	}
	return electionsettings
}

// DecideNextPresident returns the ID of chosen next President
// OPTIONAL: override to manipulate the result of the election
func (j *Judge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	return winner
}

package baseclient

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type BaseJudge struct {
	GameState gamestate.ClientGameState
}

// GetRuleViolationSeverity returns a custom map of named rules and how severe the sanction should be for transgressing them
// If a rule is not named here, the default sanction value added is 1
// OPTIONAL: override to set custom sanction severities for specific rules
func (j *BaseJudge) GetRuleViolationSeverity() map[string]shared.IIGOSanctionsScore {
	return map[string]shared.IIGOSanctionsScore{}
}

// GetSanctionThresholds returns a custom map of sanction score thresholds for different sanction tiers
// For any unfilled sanction tiers will be filled with default values (given in judiciary.go)
// OPTIONAL: override to set custom sanction thresholds
func (j *BaseJudge) GetSanctionThresholds() map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore {
	return map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{}
}

// PayPresident pays the President a salary.
// OPTIONAL: override to pay the President less than the full amount.
func (j *BaseJudge) PayPresident() (shared.Resources, bool) {
	// TODO Implement opinion based salary payment.
	PresidentSalaryRule, ok := j.GameState.RulesInfo.CurrentRulesInPlay["salary_cycle_president"]
	var PresidentSalary shared.Resources = 0
	if ok {
		PresidentSalary = shared.Resources(PresidentSalaryRule.ApplicableMatrix.At(0, 1))
	}
	return PresidentSalary, true
}

// InspectHistory is the base implementation of evaluating islands choices the last turn.
// OPTIONAL: override if you want to evaluate the history log differently.
func (j *BaseJudge) InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]shared.EvaluationReturn, bool) {
	outputMap := map[shared.ClientID]shared.EvaluationReturn{}
	copyOfVarCache := rules.CopyVariableMap(j.GameState.RulesInfo.VariableMap)
	for _, entry := range iigoHistory {
		variablePairs := entry.Pairs
		clientID := entry.ClientID
		var rulesAffected []string
		for _, variable := range variablePairs {
			valuesToBeAdded, foundRules := rules.PickUpRulesByVariable(variable.VariableName, j.GameState.RulesInfo.CurrentRulesInPlay, copyOfVarCache)
			if foundRules {
				rulesAffected = append(rulesAffected, valuesToBeAdded...)
			}
			updatedVariable := rules.UpdateVariableInternal(variable.VariableName, variable, copyOfVarCache)
			if !updatedVariable {
				return map[shared.ClientID]shared.EvaluationReturn{}, false
			}
		}
		if _, ok := outputMap[clientID]; !ok {
			outputMap[clientID] = shared.EvaluationReturn{
				Rules:       []rules.RuleMatrix{},
				Evaluations: []bool{},
			}
		}
		tempReturn := outputMap[clientID]
		for _, rule := range rulesAffected {
			ret := rules.EvaluateRuleFromCaches(rule, j.GameState.RulesInfo.CurrentRulesInPlay, copyOfVarCache)
			if ret.EvalError != nil {
				return outputMap, false
			}
			tempReturn.Rules = append(tempReturn.Rules, j.GameState.RulesInfo.CurrentRulesInPlay[rule])
			tempReturn.Evaluations = append(tempReturn.Evaluations, ret.RulePasses)
		}
		outputMap[clientID] = tempReturn
	}
	return outputMap, true
}

// GetPardonedIslands decides which islands to pardon i.e. no longer impose sanctions on
// COMPULSORY: decide which islands, if any, to forgive
func (j *BaseJudge) GetPardonedIslands(currentSanctions map[int][]shared.Sanction) map[int][]bool {
	return map[int][]bool{}
}

// HistoricalRetributionEnabled allows you to punish more than the previous turns transgressions
// OPTIONAL: override if you want to punish historical transgressions
func (j *BaseJudge) HistoricalRetributionEnabled() bool {
	return false
}

// CallPresidentElection is called by the judiciary to decide on power-transfer
// COMPULSORY: decide when to call an election following relevant rulesInPlay if you wish
func (j *BaseJudge) CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
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
func (j *BaseJudge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	return winner
}

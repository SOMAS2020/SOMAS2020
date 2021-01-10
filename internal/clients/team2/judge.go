package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Judge struct {
	*baseclient.BaseJudge
	c *client
}

// Pay President default amount
func (j *Judge) PayPresident() (shared.Resources, bool) {
	return j.BaseJudge.PayPresident()
}

// InspectHistory returns an evaluation on whether islands have adhered to the rules for that turn as a boolean.
func (j *Judge) InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]shared.EvaluationReturn, bool) {
	outMap := map[shared.ClientID]shared.EvaluationReturn{}
	copyOfVarCache := rules.CopyVariableMap(j.GameState.RulesInfo.VariableMap)

	j.c.updateRoleTrust(iigoHistory)

	for _, entry := range iigoHistory {
		variablePairs := entry.Pairs
		clientID := entry.ClientID
		j.c.confidenceRestrospect("RoleOpinion", clientID)
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
		if _, ok := outMap[clientID]; !ok {
			outMap[clientID] = shared.EvaluationReturn{
				Rules:       []rules.RuleMatrix{},
				Evaluations: []bool{},
			}
		}

		// If the island's trustworthiness is above the threshold, then return true for all rule evaluations
		if j.c.confidence("RoleOpinion", clientID) > 80 {
			tempReturn := outMap[clientID]
			for _, rule := range rulesAffected {
				tempReturn.Rules = append(tempReturn.Rules, j.GameState.RulesInfo.CurrentRulesInPlay[rule])
				tempReturn.Evaluations = append(tempReturn.Evaluations, true)
			}
			outMap[clientID] = tempReturn
		} else {
			// All other islands will be evaluated fairly using base implementation.
			tempReturn := outMap[clientID]
			for _, rule := range rulesAffected {
				evaluation := rules.EvaluateRuleFromCaches(rule, j.GameState.RulesInfo.CurrentRulesInPlay, copyOfVarCache)
				if evaluation.EvalError != nil {
					return outMap, false
				}
				tempReturn.Rules = append(tempReturn.Rules, j.GameState.RulesInfo.CurrentRulesInPlay[rule])
				tempReturn.Evaluations = append(tempReturn.Evaluations, evaluation.RulePasses)
			}
			outMap[clientID] = tempReturn
		}
	}

	return outMap, true
}

// GetPardonedIslands decides which islands to pardon i.e. no longer impose sanctions on
func (j *Judge) GetPardonedIslands(currentSanctions map[int][]shared.Sanction) map[int][]bool {
	Pardoned := make(map[int][]bool)
	for i, List := range currentSanctions {
		List2 := make([]bool, len(List))
		Pardoned[i] = List2
		// TODO: what are we intending to check here with agent strategy because the values were used wrong across multiple files
		for index, sanction := range List {
			if j.c.confidence("RoleOpinion", sanction.ClientID) > 80 && j.c.setAgentStrategy() != 1 {
				Pardoned[i][index] = true
			} else {
				Pardoned[i][index] = false
			}
		}
	}
	return Pardoned
}

// HistoricalRetributionEnabled allows you to punish more than the previous turns transgressions
func (j *Judge) HistoricalRetributionEnabled() bool {
	return true
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

	// TODO: This could be a config value - is there a rule on turns for power transfer?
	// TODO: could use confidence to try to keep friends in power - cheat the accountability cycle
	if turnsInPower >= 2 {
		electionsettings.HoldElection = true
	}

	return electionsettings
}

func (j *Judge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	// If the election winner's trust score is okay, we will declare them as the next President.
	// If not, we will replace it with the island who's trust score is higher
	opWinner := j.c.confidence("RoleOpinion", winner)

	if opWinner < 30 {
		aliveIslands := j.c.getAliveClients()
		for _, island := range aliveIslands {
			opIsland := j.c.confidence("RoleOpinion", island)
			//this only replaces the winner with someone with a higher trust
			if opIsland > opWinner {
				winner = island
				opWinner = opIsland
			}
		}
	}

	return winner
}

// GetRuleViolationSeverity returns a custom map of named rules and how severe the sanction should be for transgressing them
// If a rule is not named here, the default sanction value added is 1
// OPTIONAL: override to set custom sanction severities for specific rules
func (j *Judge) GetRuleViolationSeverity() map[string]shared.IIGOSanctionsScore {
	return map[string]shared.IIGOSanctionsScore{}
}

// GetSanctionThresholds returns a custom map of sanction score thresholds for different sanction tiers
// For any unfilled sanction tiers will be filled with default values (given in judiciary.go)
// OPTIONAL: override to set custom sanction thresholds
func (j *Judge) GetSanctionThresholds() map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore {
	return map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{}
}

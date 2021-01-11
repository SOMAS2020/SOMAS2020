package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type judge struct {
	*baseclient.BaseJudge
	c *client
}

func (c *client) GetClientJudgePointer() roles.Judge {
	c.Logf("Team 5 became the Judge, the Jury and Executioner.")
	return &judge{c: c, BaseJudge: &baseclient.BaseJudge{GameState: c.ServerReadHandle.GetGameState()}}
}

// GetRuleViolationSeverity returns a custom map of named rules and how severe the sanction should be for transgressing them
// If a rule is not named here, the default sanction value added is 1
// OPTIONAL: override to set custom sanction severities for specific rules
func (j *judge) GetRuleViolationSeverity() map[string]shared.IIGOSanctionsScore {
	return map[string]shared.IIGOSanctionsScore{}
}

// GetSanctionThresholds returns a custom map of sanction score thresholds for different sanction tiers
// For any unfilled sanction tiers will be filled with default values (given in judiciary.go)
// OPTIONAL: override to set custom sanction thresholds
func (j *judge) GetSanctionThresholds() map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore {
	return map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{}
}

// Pay president based on the status of our own wealth
// If we are not doing verywell, pay President less so we have more in the CP to take from
func (j *judge) PayPresident() (shared.Resources, bool) {

	PresidentSalaryRule, ok := j.GameState.RulesInfo.CurrentRulesInPlay["salary_cycle_president"]
	var salary shared.Resources = 0
	if ok {
		salary = shared.Resources(PresidentSalaryRule.ApplicableMatrix.At(0, 1))
	}
	if j.c.wealth() == jeffBezos {
		return salary, true
	} else if j.c.wealth() == middleClass {
		salary = salary * 0.8
	} else {
		salary = salary * 0.5
	}
	return salary, true

}

// InspectHistory is the base implementation of evaluating islands choices the last turn.
// OPTIONAL: override if you want to evaluate the history log differently.
func (j *judge) InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]shared.EvaluationReturn, bool) {
	return j.BaseJudge.InspectHistory(iigoHistory, turnsAgo)
}

// Pardon ourselves and homies
func (j *judge) GetPardonedIslands(currentSanctions map[int][]shared.Sanction) map[int][]bool {
	pardons := make(map[int][]bool)
	for key, sanctionList := range currentSanctions {
		lst := make([]bool, len(sanctionList))
		pardons[key] = lst
		for index, sanction := range sanctionList {
			if j.c.opinions[sanction.ClientID].getScore() > 0.5 || sanction.ClientID == shared.Team5 {
				pardons[key][index] = true
			} else {
				pardons[key][index] = false
			}

		}
	}
	return pardons
}

// HistoricalRetributionEnabled allows you to punish more than the previous turns transgressions
// OPTIONAL: override if you want to punish historical transgressions
func (j *judge) HistoricalRetributionEnabled() bool {
	return false
}

// CallPresidentElection is called by the judiciary to decide on power-transfer
// COMPULSORY: decide when to call an election following relevant rulesInPlay if you wish
func (j *judge) CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
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

// if the real winner is on our bad side, then we choose our best friend
func (j *judge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	aliveTeams := j.c.getAliveTeams(false) //not including us
	if j.c.opinions[winner].getScore() < 0 {
		ballot := j.c.VoteForElection(shared.President, aliveTeams)
		winner = ballot[0] //choose the first one in Borda Vote
	}
	return winner
}

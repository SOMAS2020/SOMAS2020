package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type judge struct {
	*baseclient.BaseJudge
	*client
}

func (j *judge) GetRuleViolationSeverity() map[string]roles.IIGOSanctionScore {
	return map[string]roles.IIGOSanctionScore{}
}

func (j *judge) GetSanctionThresholds() map[roles.IIGOSanctionTier]roles.IIGOSanctionScore {
	return j.BaseJudge.GetSanctionThresholds()
}

func (j *judge) InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]roles.EvaluationReturn, bool) {
	return j.BaseJudge.InspectHistory(iigoHistory, turnsAgo)
}

// We don't have reasons pardon any islands being sanctioned
func (j *judge) GetPardonedIslands(currentSanctions map[int][]roles.Sanction) map[int][]bool {
	return j.BaseJudge.GetPardonedIslands(currentSanctions)
}

func (j *judge) HistoricalRetributionEnabled() bool {
	return j.BaseJudge.HistoricalRetributionEnabled()
}

func (j *judge) CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {

	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.BordaCount,
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

func (j *judge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	return winner
}

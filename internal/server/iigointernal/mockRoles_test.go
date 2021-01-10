package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type mockJudge struct {
	violationSeverity     map[string]shared.IIGOSanctionsScore
	sanctionThresholds    map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore
	payPresidentVal       shared.Resources
	payPresidentChoice    bool
	inspectHistoryReturn  map[shared.ClientID]shared.EvaluationReturn
	inspectHistoryChoice  bool
	pardonedIslandsReturn map[int][]bool
	historicalRetribution bool
	callPresidentElection shared.ElectionSettings
	decidePresidentResult shared.ClientID
}

// GetRuleViolationSeverity returns a custom map of named rules and how severe the sanction should be for transgressing them
// If a rule is not named here, the default sanction value added is 1
func (j *mockJudge) GetRuleViolationSeverity() map[string]shared.IIGOSanctionsScore {
	return j.violationSeverity
}

// GetSanctionThresholds returns a custom map of sanction score thresholds for different sanction tiers
// For any unfilled sanction tiers will be filled with default values (given in judiciary.go)
func (j *mockJudge) GetSanctionThresholds() map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore {
	return j.sanctionThresholds
}

// PayPresident pays the President a salary.
func (j *mockJudge) PayPresident() (shared.Resources, bool) {
	// TODO Implement opinion based salary payment.
	return j.payPresidentVal, j.payPresidentChoice
}

// InspectHistory is the base implementation of evaluating islands choices the last turn.
func (j *mockJudge) InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]shared.EvaluationReturn, bool) {
	return j.inspectHistoryReturn, j.inspectHistoryChoice
}

// GetPardonedIslands allows a client to check all the sanctions that are currently levied against any islands.
// By returning true for a particular position a pardon will be issued for that sanction,
// false will continue the sanction as is
func (j *mockJudge) GetPardonedIslands(currentSanctions map[int][]shared.Sanction) map[int][]bool {
	return j.pardonedIslandsReturn
}

// HistoricalRetributionEnabled Setting this to true will instruct the judiciary to give the judge role the entire
// history cache to browse through for evaluations
func (j *mockJudge) HistoricalRetributionEnabled() bool {
	return j.historicalRetribution
}

// CallPresidentElection is called by the judiciary to decide on power-transfer
func (j *mockJudge) CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	return j.callPresidentElection
}

// DecideNextPresident returns the ID of chosen next President
func (j *mockJudge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	return j.decidePresidentResult
}

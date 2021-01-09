package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Judge is an interface that is implemented by baseJudge but can also be
// optionally implemented by individual islands.
// Judge is the decision interface for the judiciary branch of IIGO
type Judge interface {
	PayPresident() (shared.Resources, bool)
	InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]shared.EvaluationReturn, bool)
	CallPresidentElection(shared.MonitorResult, int, []shared.ClientID) shared.ElectionSettings
	DecideNextPresident(shared.ClientID) shared.ClientID
	GetRuleViolationSeverity() map[string]shared.IIGOSanctionsScore
	GetSanctionThresholds() map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore
	GetPardonedIslands(currentSanctions map[int][]shared.Sanction) map[int][]bool
	HistoricalRetributionEnabled() bool
}

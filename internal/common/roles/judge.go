package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Judge is an interface that is implemented by baseJudge but can also be
// optionally implemented by individual islands.

// EvaluationReturn is a data-structure allowing clients to return which rules they've evaluated and the results
type EvaluationReturn struct {
	Rules       []rules.RuleMatrix
	Evaluations []bool
}

// Sanction is a data-structure that represents a sanction on an agent, including how long it has left
type Sanction struct {
	ClientID     shared.ClientID
	SanctionTier IIGOSanctionTier
	TurnsLeft    int
}

// IIGOSanctionScore provides typed integer score for each island
type IIGOSanctionScore int

// IIGOSanctionTier provides typed integer tiers for sanctions
type IIGOSanctionTier int

// Provides enumerated tiers for IIGO sanctions
const (
	SanctionTier1 IIGOSanctionTier = iota
	SanctionTier2
	SanctionTier3
	SanctionTier4
	SanctionTier5
	NoSanction
)

// Judge is the decision interface for the judiciary branch of IIGO
type Judge interface {
	PayPresident(presidentSalary shared.Resources) (shared.Resources, bool)
	InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]EvaluationReturn, bool)
	CallPresidentElection(shared.MonitorResult, int, []shared.ClientID) shared.ElectionSettings
	DecideNextPresident(shared.ClientID) shared.ClientID
	GetRuleViolationSeverity() map[string]IIGOSanctionScore
	GetSanctionThresholds() map[IIGOSanctionTier]IIGOSanctionScore
	GetPardonedIslands(currentSanctions map[int][]Sanction) map[int][]bool
	HistoricalRetributionEnabled() bool
}

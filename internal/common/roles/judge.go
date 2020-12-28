package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Judge is an interface that is implemented by baseJudge but can also be
// optionally implemented by individual islands.

type EvaluationReturn struct {
	Rules       []rules.RuleMatrix
	Evaluations []bool
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
)

type Judge interface {
	PayPresident() (shared.Resources, bool)
	InspectHistory(iigoHistory []shared.Accountability) (map[shared.ClientID]EvaluationReturn, bool)
	DeclareSpeakerPerformance() (result bool, didRole bool)
	DeclarePresidentPerformance() (result bool, didRole bool)
	GetRuleViolationSeverity() map[string]IIGOSanctionScore
	GetSanctionThresholds() map[IIGOSanctionTier]int
}

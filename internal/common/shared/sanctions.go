package shared

import "github.com/SOMAS2020/SOMAS2020/internal/common/rules"

// IIGOSanctionsScore provides typed integer score for each island
type IIGOSanctionsScore int

// IIGOSanctionsTier provides typed integer tiers for sanctions
type IIGOSanctionsTier int

// Provides enumerated tiers for IIGO sanctions
const (
	SanctionTier1 IIGOSanctionsTier = iota
	SanctionTier2
	SanctionTier3
	SanctionTier4
	SanctionTier5
	NoSanction
)

// EvaluationReturn is a data-structure allowing clients to return which rules they've evaluated and the results
type EvaluationReturn struct {
	Rules       []rules.RuleMatrix
	Evaluations []bool
}

// Sanction is a data-structure that represents a sanction on an agent, including how long it has left
type Sanction struct {
	ClientID     ClientID
	SanctionTier IIGOSanctionsTier
	TurnsLeft    int
}

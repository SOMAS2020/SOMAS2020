package shared

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

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

func (i IIGOSanctionsTier) String() string {
	strs := [...]string{
		"SanctionTier1",
		"SanctionTier2",
		"SanctionTier3",
		"SanctionTier4",
		"SanctionTier5",
		"NoSanction",
	}
	if i >= 0 && int(i) < len(strs) {
		return strs[i]
	}
	return fmt.Sprintf("UNKNOWN IIGOSanctionsTier '%v'", int(i))
}

// GoString implements GoStringer
func (i IIGOSanctionsTier) GoString() string {
	return i.String()
}

// MarshalText implements TextMarshaler
func (i IIGOSanctionsTier) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(i.String())
}

// MarshalJSON implements RawMessage
func (i IIGOSanctionsTier) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(i.String())
}

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

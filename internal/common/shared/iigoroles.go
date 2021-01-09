package shared

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

// Role provides enumerated type for IIGO roles (President, Speaker and Judge)
type Role int

const (
	President Role = iota
	Speaker
	Judge
)

func (r Role) String() string {
	strs := [...]string{"President", "Speaker", "Judge"}
	if r >= 0 && int(r) < len(strs) {
		return strs[r]
	}
	return fmt.Sprintf("UNKNOWN Role '%v'", int(r))
}

// GoString implements GoStringer
func (r Role) GoString() string {
	return r.String()
}

// MarshalText implements TextMarshaler
func (r Role) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(r.String())
}

// MarshalJSON implements RawMessage
func (r Role) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(r.String())
}

// RuleVoteType provides enumerated values for Approving, Rejecting or Abstaining from a vote.
type RuleVoteType int

const (
	Approve = iota
	Reject
	Abstain
)

func (r RuleVoteType) String() string {
	strs := [...]string{"Approve", "Reject", "Abstain"}
	if r >= 0 && int(r) < len(strs) {
		return strs[r]
	}
	return fmt.Sprintf("UNKNOWN RuleVoteType '%v'", int(r))
}

// GoString implements GoStringer
func (r RuleVoteType) GoString() string {
	return r.String()
}

// MarshalText implements TextMarshaler
func (r RuleVoteType) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(r.String())
}

// MarshalJSON implements RawMessage
func (r RuleVoteType) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(r.String())
}

// MonitorResult is a type for communicating whether
// monitoring has been performed and the decided result
type MonitorResult struct {
	Performed bool
	Result    bool
}

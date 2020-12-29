package shared

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

type ElectionVotingMethod int

const (
	BordaCount ElectionVotingMethod = iota
	Plurality
	Majority
)

type ElectionSettings struct {
	VotingMethod  ElectionVotingMethod
	IslandsToVote []ClientID
	HoldElection  bool
}

func (e ElectionVotingMethod) String() string {
	strs := [...]string{
		"BordaCount",
		"Plurality",
		"Majority",
	}
	if e >= 0 && int(e) < len(strs) {
		return strs[e]
	}
	return fmt.Sprintf("UNKNOWN ElectionVotingMethod '%v'", int(e))
}

// GoString implements GoStringer
func (e ElectionVotingMethod) GoString() string {
	return e.String()
}

// MarshalText implements TextMarshaler
func (e ElectionVotingMethod) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(e.String())
}

// MarshalJSON implements RawMessage
func (e ElectionVotingMethod) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(e.String())
}

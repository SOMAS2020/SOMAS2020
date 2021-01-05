package shared

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

// ElectionVotingMethod provides enumerated type for selection of voting system to be used
type ElectionVotingMethod int

// Methods for winner selection in IIGO elections
const (
	BordaCount ElectionVotingMethod = iota
	Runoff
	InstantRunoff
	Approval
)

// ElectionSettings allows islands to configure elections for power transfer in IIGO
type ElectionSettings struct {
	VotingMethod  ElectionVotingMethod
	IslandsToVote []ClientID
	HoldElection  bool
}

func (e ElectionVotingMethod) String() string {
	strs := [...]string{
		"BordaCount",
		"Runoff",
		"InstantRunoff",
		"Approval",
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

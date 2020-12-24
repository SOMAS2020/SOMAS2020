package shared

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

// ForageType selects which resource the agents want to forage in
type ForageType int

const (
	// DeerForageType is hunting resource, it is described at length in foraging package
	DeerForageType ForageType = iota
	// FishForageType is another foraging resource also defined in foraging package
	FishForageType
)

// ForageDecision is used to represent a foraging decision made by agents
type ForageDecision struct {
	Type         ForageType
	Contribution Resources
}

// ForagingDecisionsDict is a map of clients' foraging decisions
type ForagingDecisionsDict = map[ClientID]ForageDecision

func (ft ForageType) String() string {
	strings := [...]string{"DeerForageType", "FishForageType"}
	if ft >= 0 && int(ft) < len(strings) {
		return strings[ft]
	}
	return fmt.Sprintf("UNKNOWN ForageType '%v'", int(ft))
}

// GoString implements GoStringer
func (ft ForageType) GoString() string {
	return ft.String()
}

// MarshalText implements TextMarshaler
func (ft ForageType) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(ft.String())
}

// MarshalJSON implements RawMessage
func (ft ForageType) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(ft.String())
}

package shared

import "fmt"

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
	switch ft {
	case DeerForageType:
		return "DeerForageType"
	case FishForageType:
		return "FishForageType"
	default:
		return fmt.Sprintf("InvalidForageType(%d)", ft)
	}
}

package shared

// ForageContribution is a number between 0 and something that
// represents how many resources the agents will allocate to foraging
type ForageContribution = float64

// ForageType selects which resource the agents want to forage in
type ForageType = int

const (
	// DeerForageType is hunting resource, it is described at length in foraging package
	DeerForageType ForageType = iota
	// FishForageType is another foraging resource also defined in foraging package
	FishForageType
)

// ForageDecision is used to represent a foraging decision made by agents
type ForageDecision struct {
	Type         ForageType
	Contribution ForageContribution
}

// ForagingDecisionsDict is a map of clients' foraging decisions
type ForagingDecisionsDict = map[ClientID]ForageDecision

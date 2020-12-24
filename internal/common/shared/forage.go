package shared

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
	Contribution Resources
}

// ForagingDecisionsDict is a map of clients' foraging decisions
type ForagingDecisionsDict = map[ClientID]ForageDecision

// ForageShareInfo is used to represent the forage informations shared by agents
type ForageShareInfo struct {
	DecisionMade     ForageDecision
	ResourceObtained Resources
	// ShareTo is used to show which agents/clients to SEND the foraging information to
	ShareTo []ClientID
	// SharedFrom is used to show where the information came from
	SharedFrom ClientID
}

// ForagingOfferDict is a map of client -> foraging decisions, their resource obtained, and which
// clients to send this information to
type ForagingOfferDict = map[ClientID]ForageShareInfo

// ForagingReceiptDict is a map of client -> array of information of other clients
// foraging decisions, their resource obtained, and which clients to sent this information
type ForagingReceiptDict = map[ClientID][]ForageShareInfo

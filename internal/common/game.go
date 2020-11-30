package common

// DefaultResources is the default number of resources at the start of the game
const DefaultResources = 100

// ClientInfo contains the client struct as well as the client's attributes
type ClientInfo struct {
	Client Client

	// Resources contains the amount of resources owned by the client.
	Resources uint

	Alive bool

	// [INFRA] add more client information here
}

// GameState represents the game's state.
type GameState struct {
	// Day represents the current (1-index) day of the game.
	Day int
	// ClientInfos map from the ClientID to ClientInfo.
	// EXTRA note: Golang maps are made to be random!
	ClientInfos map[ClientID]ClientInfo

	// 	[INFRA] add more details regarding state of game here
}

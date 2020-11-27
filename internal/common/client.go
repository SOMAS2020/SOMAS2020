package common

import "fmt"

// ClientID is an enum for client IDs
type ClientID int

const (
	// Team1 ID
	Team1 ClientID = iota
	// Team2 ID
	Team2
	// Team3 ID
	Team3
	// Team4 ID
	Team4
	// Team5 ID
	Team5
	// Team6 ID
	Team6
)

// TeamIDs contain sequential IDs of all teams
var TeamIDs = [...]ClientID{Team1, Team2, Team3, Team4, Team5, Team6}

// Client is a base interface to be implemented by each client struct.
type Client interface {
	Echo(s string) string
	GetID() ClientID
	Logf(format string, a ...interface{})

	// GetForageInvestment returns the amt of resources the team is willing to spend for foraging
	GetForageInvestment(gs GameState) uint
}

// RegisteredClients contain all registered clients, exposed for the server.
var RegisteredClients = map[ClientID]Client{}

// RegisterClient registers clients into RegisteredClients
func RegisterClient(id ClientID, c Client) {
	// prevent double registrations
	if _, ok := RegisteredClients[id]; ok {
		// OK to panic here, as this is a _crucial_ step.
		panic(fmt.Sprintf("Duplicate client ID %v in RegisterClient!", id))
	}
	RegisteredClients[id] = c
}

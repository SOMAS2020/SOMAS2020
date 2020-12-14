package common

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Client is a base interface to be implemented by each client struct.
type Client interface {
	Echo(s string) string
	GetID() shared.ClientID

	// StartOfTurnUpdate is where SOMASServer.updateIsland sends the game state over
	// at start of turn. Do whatever you like here :).
	StartOfTurnUpdate(gameState GameState)

	// EndOfTurnActions should return all end of turn actions.
	EndOfTurnActions() []Action
}

// RegisteredClients contain all registered clients, exposed for the server.
var RegisteredClients = map[shared.ClientID]Client{}

// RegisterClient registers clients into RegisteredClients
func RegisterClient(id shared.ClientID, c Client) {
	// prevent double registrations
	if _, ok := RegisteredClients[id]; ok {
		// OK to panic here, as this is a _crucial_ step.
		panic(fmt.Sprintf("Duplicate client ID %v in RegisterClient!", id))
	}
	RegisteredClients[id] = c
}

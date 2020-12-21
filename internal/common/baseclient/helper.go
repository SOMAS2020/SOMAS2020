package baseclient

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

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

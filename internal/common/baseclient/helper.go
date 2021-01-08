package baseclient

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// ClientFactory is a function that returns a new Client instance.
type ClientFactory = func() Client

// RegisteredClientFactories contain all registered clients and their respective clientFactories exposed for
// the server.
var RegisteredClientFactories = map[shared.ClientID]ClientFactory{}

// RegisterClientFactory registers client factories into RegisteredClientFactories
func RegisterClientFactory(id shared.ClientID, c ClientFactory) {
	// prevent double registrations
	if _, ok := RegisteredClientFactories[id]; ok {
		// OK to panic here, as this is a _crucial_ step.
		panic(fmt.Sprintf("Duplicate client ID %v in RegisteredClientFactories!", id))
	}
	RegisteredClientFactories[id] = c
}

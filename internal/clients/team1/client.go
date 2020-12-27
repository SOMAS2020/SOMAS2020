// Package team1 contains code for team 1's client implementation
package team1

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team1

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			Client:  baseclient.NewClient(id),
			profits: map[shared.ForageDecision][]shared.Resources{},
		},
	)
}

type client struct {
	baseclient.Client
	gameState gamestate.ClientGameState

	serverReadHandle baseclient.ServerReadHandle
	profits          map[shared.ForageDecision][]shared.Resources
}

func (c *client) Initialise(handle baseclient.ServerReadHandle) {
	c.serverReadHandle = handle
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	return shared.ForageDecision{Type: shared.DeerForageType, Contribution: 5}, nil
}

func (c *client) ForageUpdate(forageDecision shared.ForageDecision, reward shared.Resources) {
	c.profits[forageDecision] = append(c.profits[forageDecision], reward)

	c.Logf("New resources: %v", c.serverReadHandle.GetGameState().ClientInfo.Resources)
	c.Logf("Profit: %v", reward - forageDecision.Contribution)
}

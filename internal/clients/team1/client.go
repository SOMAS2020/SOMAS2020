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
			profits: map[shared.ForageDecision][]int{},
		},
	)
}

type client struct {
	baseclient.Client
	gameState gamestate.ClientGameState

	lastForagingDecision shared.ForageDecision
	profits              map[shared.ForageDecision][]int
}

func (c *client) GameStateUpdate(gameState gamestate.ClientGameState, updateData gamestate.UpdateData) {
	profit :=
		( gameState.ClientInfo.Resources - c.gameState.ClientInfo.Resources ) // Cost of living

	c.Logf("New resources: %d", gameState.ClientInfo.Resources)
	c.Logf("Profit: %d", profit)

	c.profits[c.lastForagingDecision] = append(c.profits[c.lastForagingDecision], profit)
	c.gameState = gameState
}

func (c *client) StartOfTurnUpdate(gameState gamestate.ClientGameState) {
	c.GameStateUpdate(gameState)
}



func (c *client) DecideForage() (shared.ForageDecision, error) {
	contribution := 50.0
	c.Logf("Foraging deer with %f", contribution)

	c.lastForagingDecision = shared.ForageDecision{Type: shared.DeerForageType, Contribution: contribution}
	return c.lastForagingDecision, nil
}

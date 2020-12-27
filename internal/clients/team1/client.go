// Package team1 contains code for team 1's client implementation
package team1

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team1

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			Client:        baseclient.NewClient(id),
			forageResults: map[shared.ForageType][]shared.Resources{},
		},
	)
}

type client struct {
	baseclient.Client
	gameState gamestate.ClientGameState

	serverReadHandle baseclient.ServerReadHandle
	forageResults    map[shared.ForageType][]shared.Resources
}

func (c *client) Initialise(handle baseclient.ServerReadHandle) {
	c.serverReadHandle = handle
}

func (c *client) forageHistorySize() int {
	length := 0
	for _, lst := range c.forageResults {
		length += len(lst)
	}
	return length
}

func (c *client) clientInfo() gamestate.ClientInfo {
	return c.serverReadHandle.GetGameState().ClientInfo
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	// Up to 30% of our current resources
	forageContribution := shared.Resources(0.1*rand.Float64()) * c.clientInfo().Resources

	if c.forageHistorySize() > 5 {
		// Find the forageType with the best average returns
		bestForageType := shared.ForageType(-1)
		bestForageTypeReturns := shared.Resources(0)

		for forageType, returns := range c.forageResults {
			totalReturns := shared.Resources(0)
			for _, v := range returns {
				totalReturns += v
			}
			averageReturns := totalReturns / shared.Resources(len(returns))
			if averageReturns > bestForageTypeReturns {
				bestForageTypeReturns = averageReturns
				bestForageType = forageType
			}
		}

		// Not foraging is best
		if bestForageType == shared.ForageType(-1) {
			return shared.ForageDecision{shared.FishForageType, 0}, nil
		} else {
			return shared.ForageDecision{bestForageType, forageContribution}, nil
		}
	} else {
		return shared.ForageDecision{shared.DeerForageType, forageContribution}, nil
	}
}

func (c *client) ForageUpdate(forageDecision shared.ForageDecision, reward shared.Resources) {
	c.forageResults[forageDecision.Type] = append(c.forageResults[forageDecision.Type], reward)

	c.Logf("New resources: %v", c.serverReadHandle.GetGameState().ClientInfo.Resources)
	c.Logf("Profit: %v", reward-forageDecision.Contribution)
}

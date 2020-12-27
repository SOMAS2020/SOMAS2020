// Package team1 contains code for team 1's client implementation
package team1

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team1

// ForageOutcome is how much we put in to a foraging trip and how much we got
// back
type ForageOutcome struct {
	contribution shared.Resources
	reward       shared.Resources
}
type ForageHistory map[shared.ForageType][]ForageOutcome

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			Client:        baseclient.NewClient(id),
			forageHistory: ForageHistory{},
		},
	)
}

type client struct {
	baseclient.Client
	gameState gamestate.ClientGameState

	serverReadHandle baseclient.ServerReadHandle
	forageHistory    ForageHistory
}

func (c *client) Initialise(handle baseclient.ServerReadHandle) {
	c.serverReadHandle = handle
}

func (c *client) forageHistorySize() int {
	length := 0
	for _, lst := range c.forageHistory {
		length += len(lst)
	}
	return length
}

func (c *client) clientInfo() gamestate.ClientInfo {
	return c.serverReadHandle.GetGameState().ClientInfo
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	if c.forageHistorySize() < 5 {
		// Up to 10% of our current resources
		forageContribution := shared.Resources(0.1*rand.Float64()) * c.clientInfo().Resources
		return shared.ForageDecision{shared.DeerForageType, forageContribution}, nil
	}

	// Find the forageType with the best average returns
	bestForageType := shared.ForageType(-1)
	bestROI := 0.0

	// TODO: We could also do this in one go. So just repeat the exact decision
	// that gave the best ROI, instead of picking forageType based on its
	// average ROI

	for forageType, outcomes := range c.forageHistory {
		ROIsum := 0.0
		for _, outcome := range outcomes {
			if outcome.contribution != 0 {
				ROIsum += float64(outcome.reward / outcome.contribution)
			}
		}

		averageROI := ROIsum / float64(len(outcomes))

		if averageROI > bestROI {
			bestROI = averageROI
			bestForageType = forageType
		}
	}

	// Not foraging is best
	if bestForageType == shared.ForageType(-1) {
		return shared.ForageDecision{shared.FishForageType, 0}, nil
	}

	// Find the value of resources that gave us the best return and add some
	// noise to it. Cap to 20% of our stockpile
	pastOutcomes := c.forageHistory[bestForageType]
	bestValue := shared.Resources(0)
	bestValueROI := shared.Resources(0)
	for _, outcome := range pastOutcomes {
		if outcome.contribution != 0 {
			ROI := outcome.reward / outcome.contribution
			if ROI > bestValueROI {
				bestValue = outcome.contribution
				bestValueROI = ROI
			}
		}
	}
	bestValue = shared.Resources(math.Min(
		float64(bestValue),
		float64(0.2*c.clientInfo().Resources)),
	)
	bestValue += shared.Resources(math.Min(
		rand.Float64(),
		float64(0.01*c.clientInfo().Resources),
	))

	return shared.ForageDecision{bestForageType, bestValue}, nil
}

func (c *client) ForageUpdate(forageDecision shared.ForageDecision, reward shared.Resources) {
	c.forageHistory[forageDecision.Type] = append(c.forageHistory[forageDecision.Type], ForageOutcome{
		contribution: forageDecision.Contribution,
		reward:       reward,
	})

	c.Logf("Profit: %v", reward-forageDecision.Contribution)
	c.Logf("New resources: %v", c.serverReadHandle.GetGameState().ClientInfo.Resources)
	c.Logf("History: %v", c.forageHistory)
}

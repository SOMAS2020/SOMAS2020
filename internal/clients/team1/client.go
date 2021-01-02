// Package team1 contains code for team 1's client implementation
package team1

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	// "github.com/SOMAS2020/SOMAS2020/internal/clients/team1/foraging"
)

const id = shared.Team1

type EmotionalState int

const (
	Normal EmotionalState = iota
	Anxious
	Desperate
)

func (st EmotionalState) String() string {
	strings := [...]string{"Normal", "Anxious", "Desperate"}
	if st >= 0 && int(st) < len(strings) {
		return strings[st]
	}
	return fmt.Sprintf("UNKNOWN ForageType '%v'", int(st))
}

func init() {
	baseclient.RegisterClient(id, NewClient(id))
}

type clientConfig struct {
	// At the start of the game forage randomly for this many turns. If true,
	// pay some initial tax to help get the first IIGO running
	randomForageTurns uint

	// If resources go below this limit, go into "desperation" mode
	anxietyThreshold shared.Resources

	// If true, ignore requests for taxes
	evadeTaxes bool

	// If true, pay some initial tax to help get the first IIGO running
	kickstartTaxPercent shared.Resources

	// desperateStealAmount is the amount the agent will steal from the commonPool
	desperateStealAmount shared.Resources

	// forageContributionCapPercent is the maximum percent of current resources we
	// will use for foraging
	forageContributionCapPercent float64

	// desperateStealAmount is the max randomness we will add to our foraging
	// amount as a percentage of current resources
	forageContributionNoisePercent float64

	forageDecider forageDecider
}

// client is Lucy.
type client struct {
	*baseclient.BaseClient

	forageHistory        ForageHistory
	expectedForageReward shared.Resources
	taxAmount            shared.Resources

	// allocation is the president's response to your last common pool resource request
	allocation shared.Resources

	config clientConfig
}

func NewClient(clientID shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:    baseclient.NewClient(clientID),
		forageHistory: ForageHistory{},
		taxAmount:     0,
		allocation:    0,
		config: clientConfig{
			randomForageTurns:              5,
			anxietyThreshold:               20,
			desperateStealAmount:           30,
			evadeTaxes:                     false,
			kickstartTaxPercent:            0.25,
			forageContributionCapPercent:   0.2,
			forageContributionNoisePercent: 0.01,
			forageDecider:                  defaultForageDecider,
		},
	}
}

func defaultForageDecider(c client) (shared.ForageDecision, shared.Resources) {
	if c.forageHistorySize() < c.config.randomForageTurns {
		return randomDecider(c)
	} else if c.emotionalState() == Desperate {
		return desperateDecider(c)
	} else {
		return regressionDecider(c)
	}
}

func (c client) emotionalState() EmotionalState {
	ci := c.gameState().ClientInfo
	switch {
	case ci.LifeStatus == shared.Critical:
		return Desperate
	case ci.Resources < c.config.anxietyThreshold:
		return Anxious
	default:
		return Normal
	}
}

func (c *client) StartOfTurn() {
	c.Logf("Emotional state: %v", c.emotionalState())
	c.Logf("Resources: %v", c.gameState().ClientInfo.Resources)

	for clientID, status := range c.gameState().ClientLifeStatuses {
		if status != shared.Dead && clientID != c.GetID() {
			return
		}
	}
	c.Logf("I'm all alone :c")

}

/********************/
/***    Helpers     */
/********************/

func (c *client) forageHistorySize() uint {
	length := uint(0)
	for _, lst := range c.forageHistory {
		length += uint(len(lst))
	}
	return length
}

func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

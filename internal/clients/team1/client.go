// Package team1 contains code for team 1's client implementation
package team1

import (
	"fmt"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
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
	return fmt.Sprintf("UNKNOWN EmotionalState '%v'", int(st))
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

	// maxOpinion is the boundary where we either give resources without questioning
	// or we refuse to give them resources.
	maxOpinion int

	// flipForageScale scales the amount contributed by flipForage
	flipForageScale float64
}

// client is Lucy.
type client struct {
	*baseclient.BaseClient

	forageHistory        ForageHistory
	expectedForageReward shared.Resources

	// IITO/Gifts
	opinionTeams  []opinionOnTeam
	receivedOffer map[shared.ClientID]shared.Resources

	// allocation is the president's response to your last common pool resource request
	allocation shared.Resources
	// Disaster
	// A map of Turns -> DisasterInfo
	disasterInfo disaster

	// aliveClients
	aliveClients []shared.ClientID

	// For foraging
	switchType bool

	config clientConfig
}

func NewClient(clientID shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:    baseclient.NewClient(clientID),
		forageHistory: ForageHistory{},
		config: clientConfig{
			randomForageTurns:              0,
			anxietyThreshold:               20,
			desperateStealAmount:           30,
			evadeTaxes:                     false,
			kickstartTaxPercent:            0,
			forageContributionCapPercent:   0.2,
			forageContributionNoisePercent: 0.01,
			maxOpinion:                     10,
			flipForageScale:                0.3,
		},
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

	// This should only happen at the start of the game.
	if c.gameState().Turn == 1 {
		c.disasterInfo.meanDisaster = disasters.DisasterReport{}
		c.switchType = false
		if c.ServerReadHandle.GetGameConfig().DisasterConfig.DisasterPeriod.Valid == true {
			c.disasterInfo.estimatedDDay = c.ServerReadHandle.GetGameConfig().DisasterConfig.DisasterPeriod.Value
		} else {
			c.disasterInfo.estimatedDDay = uint(rand.Intn(10))
		}
	}

	// if opinionTeams is empty. Initialise it.
	if len(c.opinionTeams) <= 0 {
		for _, clientID := range shared.TeamIDs {
			c.opinionTeams = append(c.opinionTeams, opinionOnTeam{clientID: clientID, opinion: 0})
		}
	}

	// Reset the map
	c.receivedOffer = nil
	if c.receivedOffer == nil {
		c.receivedOffer = make(map[shared.ClientID]shared.Resources)
	}

	// Reset alive list
	c.aliveClients = nil
	for clientID, status := range c.gameState().ClientLifeStatuses {
		if status != shared.Dead && clientID != c.GetID() {
			c.aliveClients = append(c.aliveClients, clientID)
		}
	}

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
		for _, outcome := range lst {
			if outcome.participant == c.GetID() {
				length++
			}
		}
	}
	return length
}

func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

func (c *client) livingClients() (livingClients []shared.ClientID) {
	ids := []shared.ClientID{}
	for id, livingState := range c.gameState().ClientLifeStatuses {
		if livingState == shared.Alive {
			ids = append(ids, id)
		}
	}
	return ids
}

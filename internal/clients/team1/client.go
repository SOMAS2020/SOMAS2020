// Package team1 contains code for team 1's client implementation
package team1

import (
	"fmt"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// DefaultClient creates the client that will be used for most simulations. All
// other personalities are considered alternatives. To give a different
// personality for your agent simply create another (exported) function with the
// same signature as "DefaultClient" that creates a different agent, and inform
// someone on the simulation team that you would like it to be included in
// testing
func DefaultClient(id shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:    baseclient.NewClient(id),
		BasePresident: &baseclient.BasePresident{},
		config: team1Config{
			anxietyThreshold:               50,
			randomForageTurns:              0,
			flipForageScale:                0.3,
			forageContributionCapPercent:   0.2,
			forageContributionNoisePercent: 0.01,
			evadeTaxes:                     false,
			kickstartTaxPercent:            0,
			desperateStealAmount:           30,
			maxOpinion:                     10,
			soloDeerHuntContribution:       40,
		},

		forageHistory:     ForageHistory{},
		reportedResources: map[shared.ClientID]bool{},
		teamOpinions:      map[shared.ClientID]Opinion{},
		receivedOffer:     map[shared.ClientID]shared.Resources{},
		trustTeams:        map[shared.ClientID]float64{},
	}
}

type EmotionalState int
type Opinion int

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

type team1Config struct {
	// If resources go below this limit, go into "desperation" mode
	anxietyThreshold shared.Resources

	//** Foraging **//

	// At the start of the game forage randomly for this many turns. If true,
	// pay some initial tax to help get the first IIGO running
	randomForageTurns uint

	// flipForageScale scales the amount contributed by flipForage
	flipForageScale float64

	// forageContributionCapPercent is the maximum percent of current resources we
	// will use for foraging
	forageContributionCapPercent float64

	// desperateStealAmount is the max randomness we will add to our foraging
	// amount as a percentage of current resources
	forageContributionNoisePercent float64

	// soloDeerHuntContribution is the amount of resources that will be used to
	// hunt if no one hunted last turn
	soloDeerHuntContribution shared.Resources

	//** Taxes **//

	// If true, ignore requests for taxes
	evadeTaxes bool

	// If true, pay some initial tax to help get the first IIGO running
	kickstartTaxPercent shared.Resources

	//** Allocation **//

	// commonPoolResourceRequestScale scales the request made to the commonPool.
	// The basis is the cost of living
	resourceRequestScale float64

	// desperateStealAmount is the amount the agent will steal from the commonPool
	desperateStealAmount shared.Resources

	//** Opinions **//

	// maxOpinion is the boundary where we either give resources without questioning
	// or we refuse to give them resources.
	maxOpinion Opinion
}

// client is Lucy.
type client struct {
	*baseclient.BaseClient
	*baseclient.BasePresident

	forageHistory        ForageHistory
	expectedForageReward shared.Resources

	reportedResources map[shared.ClientID]bool

	// IITO/Gifts
	teamOpinions  map[shared.ClientID]Opinion
	receivedOffer map[shared.ClientID]shared.Resources

	// allocation is the president's response to your last common pool resource request
	allocation shared.Resources

	// Disaster
	disasterInfo             disaster
	othersDisasterPrediction shared.ReceivedDisasterPredictionsDict
	// The higher the score, the more trustworthy they are
	trustTeams map[shared.ClientID]float64

	// Foraging
	forageType shared.ForageType

	config team1Config
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

	// Initialise President with gamestate
	c.BasePresident.GameState = c.gameState()
	c.LocalVariableCache = map[rules.VariableFieldName]rules.VariableValuePair{}
	// This should only happen at the start of the game.
	if c.gameState().Turn == 1 {
		c.disasterInfo.meanDisaster = disasters.DisasterReport{}
		c.forageType = shared.DeerForageType
		if c.gameConfig().DisasterConfig.DisasterPeriod.Valid {
			c.disasterInfo.estimatedDDay = c.gameConfig().DisasterConfig.DisasterPeriod.Value
		} else {
			c.disasterInfo.estimatedDDay = uint(rand.Intn(10))
		}

		c.trustTeams = make(map[shared.ClientID]float64)
		c.othersDisasterPrediction = make(shared.ReceivedDisasterPredictionsDict)
	}

	// if opinionTeams is empty. Initialise it.
	if len(c.teamOpinions) <= 0 {
		for _, clientID := range shared.TeamIDs {
			c.teamOpinions[clientID] = 0
		}
	}

	// Reset the map
	c.receivedOffer = nil
	if c.receivedOffer == nil {
		c.receivedOffer = make(map[shared.ClientID]shared.Resources)
	}

	// Reset alive list
	if len(c.aliveClients()) == 0 {
		c.Logf("I'm all alone :c")
	}
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

func (c *client) gameConfig() config.ClientConfig {
	return c.BaseClient.ServerReadHandle.GetGameConfig()
}

func (c *client) aliveClients() []shared.ClientID {
	result := []shared.ClientID{}
	for clientID, status := range c.gameState().ClientLifeStatuses {
		if status != shared.Dead && clientID != c.GetID() {
			result = append(result, clientID)
		}
	}
	return result
}

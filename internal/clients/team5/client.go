// Package team5 contains code for team 5's client implementation
package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
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
	return createClient(id)
}

func createClient(ourClientID shared.ClientID) *client {
	return &client{
		BaseClient:              baseclient.NewClient(ourClientID),
		cpRequestHistory:        cpRequestHistory{},
		cpAllocationHistory:     cpAllocationHistory{},
		forageHistory:           forageHistory{},
		resourceHistory:         resourceHistory{},
		team5President:          president{},
		giftHistory:             map[shared.ClientID]giftExchange{},
		forecastHistory:         forecastHistory{},
		receivedForecastHistory: receivedForecastHistory{},
		disasterHistory:         disasterHistory{},
		cpResourceHistory:       cpResourceHistory{0: 0},

		disasterModel: disasterModel{},

		config: getClientConfig(),
	}
}

func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {

	c.ServerReadHandle = serverReadHandle // don't change this
	c.LocalVariableCache = rules.CopyVariableMap(c.gameState().RulesInfo.VariableMap)
	c.initOpinions()
	c.initGiftHist()
	// Assign the thresholds according to the amount of resouces in the first turn
	c.config.jbThreshold = (c.gameState().ClientInfo.Resources) * c.config.jbThreshold
	c.config.middleThreshold = (c.gameState().ClientInfo.Resources) * c.config.middleThreshold
	c.config.imperialThreshold = (c.gameState().ClientInfo.Resources) * c.config.imperialThreshold

	// Gift Requests
	c.config.dyingGiftRequestAmount = float64(c.getGameConfig().CostOfLiving) * c.config.dyingGiftRequestAmount
	c.config.imperialGiftRequestAmount = float64(c.getGameConfig().CostOfLiving) * c.config.imperialGiftRequestAmount
	c.config.middleGiftRequestAmount = float64(c.getGameConfig().CostOfLiving) * c.config.middleGiftRequestAmount

	// Gift offers
	c.config.offertoDyingIslands = (float64(c.getGameConfig().CostOfLiving)) * c.config.offertoDyingIslands
	// Print the Thresholds
	c.Logf("[Initialise [%v]] JB TH %v | Middle TH %v | Imperial TH %v",
		c.config.agentMentality, c.config.jbThreshold, c.config.middleThreshold, c.config.imperialThreshold)
}

// StartOfTurn functions that are needed when our agent starts its turn
func (c *client) StartOfTurn() {
	c.Logf("[StartOfTurn][%v]: Wealth class: %v | Money In the Bank: %v | Teams still alive: %v ", c.getTurn(), c.wealth(), c.gameState().ClientInfo.Resources, c.gameState().ClientLifeStatuses)

	c.updateResourceHistory(c.resourceHistory) // First update the history of our resources
	c.opinionHistory[c.getTurn()] = c.opinions // assign last turn's opinions as default for this turn
	c.cpResourceHistory[c.getTurn()] = c.getCP()

	//update cpResourceHistory
	turn := c.getTurn()
	c.cpResourceHistory[turn] = c.getCP()

	for clientID, status := range c.gameState().ClientLifeStatuses { //if not dead then can start the turn, else no return
		if status != shared.Dead && clientID != c.GetID() {
			return
		}
	}
	c.Logf("Team 5 Died-ed lol")
}

//================================================================
/*	Wealth class
	Calculates the class of wealth we are in according
	to thresholds */
//=================================================================
func (c client) wealth() wealthTier {
	cData := c.gameState().ClientInfo
	switch {
	case cData.LifeStatus == shared.Critical: // We dying
		return dying
	case cData.Resources > c.config.jbThreshold:
		return jeffBezos // Rich
	case cData.Resources >= c.config.imperialThreshold && cData.Resources <= c.config.jbThreshold:
		return middleClass // Middle
	default:
		return imperialStudent // Middle class
	}
}

//================================================================
/*	Resource History
	Stores the level of resources we have at each turn */
//=================================================================

func (c *client) updateResourceHistory(resourceHistory resourceHistory) {
	currentResources := c.gameState().ClientInfo.Resources
	c.resourceHistory[c.gameState().Turn] = currentResources
	if c.gameState().Turn >= 2 {
		amount := c.resourceHistory[c.gameState().Turn-1]
		c.Logf("[updateResourceHistory]: Previous round: (%v) | Amount: %v", c.getTurn(), amount)
	}
	c.Logf("[updateResourceHistory]: Current round: (%v) | Amount: %v", c.getTurn(), currentResources)
}

func (c client) gameState() gamestate.ClientGameState {
	return c.ServerReadHandle.GetGameState()
}

// Package team5 contains code for team 5's client implementation
package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func init() {
	baseclient.RegisterClientFactory(ourClientID, func() baseclient.Client { return createClient() })
}

func createClient() *client {
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

		taxAmount:      0,
		allocation:     0,
		sanctionAmount: 0,

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

	// Print the Thresholds
	c.Logf("[Debug] - [Start of Turn] JB TH %v | Middle TH %v | Imperial TH %v",
		c.config.jbThreshold, c.config.middleThreshold, c.config.imperialThreshold)

}

// StartOfTurn functions that are needed when our agent starts its turn
func (c *client) StartOfTurn() {
	c.Logf("[Start of Turn] Wealth: %v | Money In the Bank: %v", c.wealth(), c.gameState().ClientInfo.Resources)
	c.Logf("[Teams still alive]: Teams %v", c.gameState().ClientLifeStatuses)

	c.updateResourceHistory(c.resourceHistory) // First update the history of our resources
	c.opinionHistory[c.getTurn()] = c.opinions // assign last turn's opinions as default for this turn
	c.cpResourceHistory[c.getTurn()] = c.getCP()

	for clientID, status := range c.gameState().ClientLifeStatuses { //if not dead then can start the turn, else no return
		if status != shared.Dead && clientID != c.GetID() {
			return
		}
	}
	c.Logf("Died-ed lol")
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
		c.Logf("[Resource History] Previous round (%v) amount: %v", c.getTurn(), amount)
	}
	c.Logf("[Resource History] Current round (%v) amount: %v", c.getTurn(), currentResources)
}

func (c client) gameState() gamestate.ClientGameState {
	return c.ServerReadHandle.GetGameState()
}

// Package team5 contains code for team 5's client implementation
package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func init() {
	c := client{
		BaseClient:          baseclient.NewClient(ourClientID),
		cpRequestHistory:    cpRequestHistory{},
		cpAllocationHistory: cpAllocationHistory{},
		forageHistory:       forageHistory{},
		resourceHistory:     resourceHistory{},
		team5president:      president{},
		giftHistory:         map[shared.ClientID]giftExchange{},

		taxAmount:  0,
		allocation: 0,
		config:     getClientConfig(),
	}

	baseclient.RegisterClient(ourClientID, &c)
}

// StartOfTurn functions that are needed when our agent starts its turn
func (c *client) StartOfTurn() {
	c.updateResourceHistory(c.resourceHistory) // First update the history of our resources
	c.wealth()

	c.opinionHistory[c.getTurn()] = c.opinions // assign last turn's opinions as default for this turn

	// Print the level of wealth we are at
	c.Logf("[Debug] - [Start of Turn] Class: %v | Money In the Bank: %v", c.wealth(), c.gameState().ClientInfo.Resources)

	for clientID, status := range c.gameState().ClientLifeStatuses { //if not dead then can start the turn, else no return
		if status != shared.Dead && clientID != c.GetID() {
			return
		}
	}

	//update cpResourceHistory
	turn := c.getTurn()
	c.cpResourceHistory[turn] = c.getCP()

	// update opinion score of current roles. Ideally at startOfTurns but for some reasons it doesn't update opinion score here
	// move to GetVoteForElection() for now
	//c.evaluateRoles()

}

func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle // don't change this
	c.LocalVariableCache = rules.CopyVariableMap()
	c.initOpinions()

	// Assign the thresholds according to the amount of resouces in the first turn
	c.config.jbThreshold = c.resourceHistory[1] * 2
	c.config.middleThreshold = c.resourceHistory[1] * 0.95
	c.config.imperialThreshold = c.resourceHistory[1] * 0.5

	// Print the Thresholds
	c.Logf("[Debug] - [Start of Turn] JB TH %v | Middle TH %v | Imperial TH %v",
		c.config.jbThreshold, c.config.middleThreshold, c.config.imperialThreshold)
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
		// c.Logf("[Team 5][Wealth:%v][Class:%v]", cData.Resources,c.config.JBThreshold)      // Debugging
		return jeffBezos // Rich
	case cData.Resources > c.config.imperialThreshold && cData.Resources <= c.config.jbThreshold:
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
		c.Logf("[Debug] - Previous round amount: %v", amount)
	}
	c.Logf("[Debug] - Current round amount: %v", currentResources)
}

func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

//================================================================
/*	Comunication
	Gets information on minimum tax amount and cp allocation */
//=================================================================
func (c *client) ReceiveCommunication(
	sender shared.ClientID,
	data map[shared.CommunicationFieldName]shared.CommunicationContent,
) {
	for field, content := range data {
		switch field {
		case shared.TaxAmount:
			c.taxAmount = shared.Resources(content.IntegerData)
		case shared.AllocationAmount:
			c.allocation = shared.Resources(content.IntegerData)
		}
	}
}

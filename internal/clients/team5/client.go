// Package team5 contains code for team 5's client implementation
package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			// Old config
			// forageHistory: ForageHistory{},
			BaseClient:          baseclient.NewClient(id),
			cpRequestHistory:    cpRequestHistory{},
			cpAllocationHistory: cpAllocationHistory{},
			forageHistory:       forageHistory{},
			resourceHistory:     resourceHistory{},
			giftHistory:         giftHistory{},
			taxAmount:           0,
			allocation:          0,
			config: clientConfig{
				//Variables for Intial forage
				InitialForageTurns:      3,
				MinimumForagePercentage: 0.01,
				NormalForagePercentage:  0.05,
				JBForagePercentage:      0.10, // % of our resources when JB is Normal< X < JB

				// Variables for Normal forage
				SkipForage:           1,
				NormalRandomIncrease: 0.05,
				MaxForagePercentage:  0.20,

				// Threshold for wealth
				JBThreshold:       100,
				MiddleThreshold:   60.0,
				ImperialThreshold: 30.0, // surely should be - 100e6? (your right we are so far indebt)
				//  Dying threshold is 0 < Dying < Imperial

				// Gifts Config
				DyingGiftRequestAmount:    10,
				ImperialGiftRequestAmount: 5,
				MiddleGiftRequestAmount:   2,
			},
		},
	)
}

// StartOfTurn functions that are needed when our agent starts its turn
func (c *client) StartOfTurn() {
	c.updateResourceHistory(c.resourceHistory) // First update the history of our resources
	c.wealth()
	// Assign the thresholds according to the amount of resouces in the first turn
	c.config.JBThreshold = c.resourceHistory[1] * 2
	c.config.MiddleThreshold = c.resourceHistory[1] * 0.95
	c.config.ImperialThreshold = c.resourceHistory[1] * 0.5

	// Print the Thresholds
	c.Logf("[Debug] - [Start of Turn] JB TH %v | Middle TH %v | Imperial TH %v",
		c.config.JBThreshold, c.config.MiddleThreshold, c.config.ImperialThreshold)

	// Print the level of wealth we are at
	c.Logf("[Debug] - [Start of Turn] Class: %v | Money In the Bank: %v", c.wealth(), c.gameState().ClientInfo.Resources)
	for clientID, status := range c.gameState().ClientLifeStatuses { //if not dead then can start the turn, else no return
		if status != shared.Dead && clientID != c.GetID() {
			return
		}
	}

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
	case cData.Resources > c.config.JBThreshold:
		// c.Logf("[Team 5][Wealth:%v][Class:%v]", cData.Resources,c.config.JBThreshold)      // Debugging
		return jeffBezos // Rich
	case cData.Resources > c.config.ImperialThreshold && cData.Resources <= c.config.JBThreshold:
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

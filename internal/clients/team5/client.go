// Package team5 contains code for team 5's client implementation
package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

//  Old config doesn't work for some reason?
/*
	func NewClient(clientID shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:    baseclient.NewClient(clientID),
		forageHistory: ForageHistory{},
		taxAmount:     0,
		allocation:    0,
		config: clientConfig{
			InitialForageTurns: 10,
			SkipForage:         5,

			JBThreshold:       1000.0,
			MiddleThreshold:   60.0,
			ImperialThreshold: 30.0, // surely should be - 100e6? (your right we are so far indebt)
		},
	}
}
*/

//================================================================
/*  Init */
//================================================================
func init() {
	baseclient.RegisterClient(
		id,
		&client{
			// BaseClient:    baseclient.NewClient(id),
			// forageHistory: ForageHistory{},
			BaseClient:    baseclient.NewClient(id),
			forageHistory: ForageHistory{},
			taxAmount:     0,
			allocation:    0,
			config: clientConfig{
				InitialForageTurns: 3,
				SkipForage:         3,

				JBThreshold:       100.0,
				MiddleThreshold:   60.0,
				ImperialThreshold: 30.0, // surely should be - 100e6? (your right we are so far indebt)
			},
		},
	)
}

func (c *client) StartOfTurn() {
	c.Logf("[Debug] - [Start of Turn] Class: %v | Money In the Bank: %v", c.wealth(), c.gameState().ClientInfo.Resources)
	// c.Logf("[The Pitts]: %v", c.gameState().ClientInfo.Resources)
	for clientID, status := range c.gameState().ClientLifeStatuses { //if not dead then can start the turn, else no return
		if status != shared.Dead && clientID != c.GetID() {
			return
		}
	}

}

//================================================================
/*  Wealth class */
//================================================================

func (c client) wealth() WealthTier {
	cData := c.gameState().ClientInfo
	switch {
	case cData.LifeStatus == shared.Critical: // We dying
		return Dying
	case cData.Resources > c.config.ImperialThreshold && cData.Resources < c.config.MiddleThreshold:
		return ImperialStudent // Poor
	case cData.Resources > c.config.JBThreshold:
		// c.Logf("[Team 5][Wealth:%v][Class:%v]", cData.Resources,c.config.JBThreshold)      // Debugging
		return JeffBezos // Rich
	default:
		return MiddleClass // Middle class
	}
}

/********************/
/***    IIFO        */
/********************/

func (c *client) MakeForageInfo() shared.ForageShareInfo {
	var shareTo []shared.ClientID

	for id, status := range c.gameState().ClientLifeStatuses {
		if status != shared.Dead {
			shareTo = append(shareTo, id)
		}
	}

	lastDecisionTurn := -1
	var lastDecision shared.ForageDecision
	var lastRevenue shared.Resources

	for forageType, outcomes := range c.forageHistory {
		for _, outcome := range outcomes {
			if int(outcome.turn) > lastDecisionTurn {
				lastDecisionTurn = int(outcome.turn)
				lastDecision = shared.ForageDecision{
					Type:         forageType,
					Contribution: outcome.input,
				}
				lastRevenue = outcome.output
			}
		}
	}

	if lastDecisionTurn < 0 {
		shareTo = []shared.ClientID{}
	}

	forageInfo := shared.ForageShareInfo{
		ShareTo:          shareTo,
		ResourceObtained: lastRevenue,
		DecisionMade:     lastDecision,
	}

	c.Logf("Sharing forage info: %v", forageInfo)
	return forageInfo
}

func (c *client) ReceiveForageInfo(forageInfos []shared.ForageShareInfo) {
	for _, forageInfo := range forageInfos {
		c.forageHistory[forageInfo.DecisionMade.Type] =
			append(
				c.forageHistory[forageInfo.DecisionMade.Type],
				ForageOutcome{
					input:  forageInfo.DecisionMade.Contribution,
					output: forageInfo.ResourceObtained,
				},
			)
	}
}

//================================================================
/*  Wealth class */
//================================================================

// gameState() gets the data from the server about our island
func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

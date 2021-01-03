// Package team5 contains code for team 5's client implementation
package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team5

// Client is the island number
type client struct {
	*baseclient.BaseClient

	forageHistory ForageHistory // Stores our previous foraging data

	taxAmount shared.Resources

	// allocation is the president's response to your last common pool resource request
	allocation shared.Resources

	config clientConfig
}

// NewClient get the base client
func NewClient(clientID shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:    baseclient.NewClient(clientID),
		forageHistory: ForageHistory{},
		taxAmount:     0,
		allocation:    0,
		config: clientConfig{
			InitialForageTurns: 10,
			// SkipForage:         5,

			JBThreshold:       1000.0,
			MiddleThreshold:   60.0,
			ImperialThreshold: 30.0, // surely should be - 100e6?
		},
	}
}

//================================================================
/*  Wealth class */
//================================================================

func (c client) wealth() WealthTier {
	cData := c.gameState().ClientInfo
	switch {
	case cData.LifeStatus == shared.Critical:
		c.Logf("Critical") // Debugging
		return Dying
	case cData.Resources > c.config.ImperialThreshold && cData.Resources < c.config.MiddleThreshold:
		c.Logf("Student") // Debugging
		return ImperialStudent
	case cData.Resources > c.config.JBThreshold:
		c.Logf("Threshold %f", c.config.JBThreshold) // Debugging
		c.Logf("Resources %f", cData.Resources)      // Debugging
		return JeffBezos
	default:
		c.Logf("Middle class") // Debugging
		return MiddleClass
	}
}

//================================================================
/*  Init */
//================================================================
func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient:    baseclient.NewClient(id),
			forageHistory: ForageHistory{},
		},
	)
}

func (c *client) StartOfTurn() {
	c.Logf("Wealth: %v", c.wealth())
	c.Logf("Resources: %v", c.gameState().ClientInfo.Resources)

	for clientID, status := range c.gameState().ClientLifeStatuses { //if not dead then can start the turn, else no return
		if status != shared.Dead && clientID != c.GetID() {
			return
		}
	}

}

func (c *client) CommonPoolResourceRequest() shared.Resources {
	switch c.wealth() {
	case Dying:
		c.Logf("Common pool request: 20")
		return 20
	default:
		return 0
	}
}

// func (c *client) RequestAllocation() shared.Resources {
// 	var allocation shared.Resources

// 	if c.wealth() == Dying {
// 		allocation = c.config.desperateStealAmount
// 	} else if c.allocation != 0 {
// 		allocation = c.allocation
// 		c.allocation = 0
// 	}

// 	if allocation != 0 {
// 		c.Logf("Taking %v from common pool", allocation)
// 	}
// 	return allocation
// }

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

/*Foraging History*/
func (c *client) forageHistorySize() uint {
	length := uint(0)
	for _, lst := range c.forageHistory {
		length += uint(len(lst))
	}
	return length // Return how many turns of foraging we have been on depending on the History
}

func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

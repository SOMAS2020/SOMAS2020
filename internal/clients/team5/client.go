// Package team5 contains code for team 5's client implementation
package team5

import (
	//"fmt"
	//"math"
	//"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team5

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient:    baseclient.NewClient(id),
			cpRequestHistory: CPRequestHistory{},
			cpAllocationHistory: CPAllocationHistory{},
			taxAmount:     0,
			allocation:    0,
			config: clientConfig {  
				JBThreshold:       100,      
				MiddleThreshold:   60,     
				ImperialThreshold: 30,	
			},
		},
	)
}

type client struct {
	*baseclient.BaseClient

	cpRequestHistory  	CPRequestHistory
	cpAllocationHistory	CPAllocationHistory

	taxAmount     		shared.Resources

	// allocation is the president's response to your last common pool resource request
	allocation 			shared.Resources

	config 				clientConfig
}

// func NewClient(clientID shared.ClientID) baseclient.Client {
// 	return &client{
// 		BaseClient:    baseclient.NewClient(clientID),
// 		forageHistory: ForageHistory{},
// 		cpRequestHistory: CPRequestHistory{},
// 		cpAllocationHistory: CPAllocationHistory{},
// 		taxAmount:     0,
// 		allocation:    0,
// 		config: clientConfig{
// 			randomForageTurns: 5,      // Richness Thresholds      
// 			JBThreshold:       100,      
// 			MiddleThreshold:   60,     
// 			ImperialThreshold: 30,	
// 		},
// 	}
// }

func (c client) wealth() WealthTier {  
	c_data := c.gameState().ClientInfo  
	switch {  
		case c_data.LifeStatus == shared.Critical:    
			return Dying  
		case c_data.Resources > c.config.ImperialThreshold && c_data.Resources < c.config.MiddleThreshold :    
			return Imperial_Student  
		case c_data.Resources > c.config.JBThreshold:    
			return Jeff_Bezos  
		default:    
			return Middle_Class  
		}
}

func (c *client) StartOfTurn() {
	c.Logf("Wealth state: %v", c.wealth())
	c.Logf("Resources: %v", c.gameState().ClientInfo.Resources)

	for clientID, status := range c.gameState().ClientLifeStatuses {
		if status != shared.Dead && clientID != c.GetID() {
			return
		}
	}
	c.Logf("I'm all alone :c")

}

func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

//------------------------------------Comunication--------------------------------------------------------//
// to get information on minimum tax amount and cp allocation
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



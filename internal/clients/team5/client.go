// Package team5 contains code for team 5's client implementation
package team5

import (
	"fmt"
	//"math"
	//"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team5

// ForageOutcome is how much we put in to a foraging trip and how much we got
// back
type ForageOutcome struct {
	turn         uint
	contribution shared.Resources
	revenue      shared.Resources
}

type ForageHistory map[shared.ForageType][]ForageOutcome

type CPRequestHistory []shared.Resources
type CPAllocationHistory []shared.Resources 

type Wealth int

const (  
	Dying            Wealth = iota // Sets values = 0  
	Imperial_Student               // iota sets the folloing values =1  
	Middle_Class                   // = 2  
	Jeff_Bezos                     // = 3
)

func (st Wealth) String() string {  
	strings := [...]string{"Dying", "Imperial_Student", "Middle_Class", "Jeff_Bezos"}  
	if st >= 0 && int(st) < len(strings) {    
		return strings[st]  
	}  
	return fmt.Sprintf("Unkown internal state '%v'", int(st))
}

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient:    baseclient.NewClient(id),
			forageHistory: ForageHistory{},
		},
	)
}

type clientConfig struct {
	randomForageTurns uint

	JBThreshold			shared.Resources
	MiddleThreshold		shared.Resources
	ImperialThreshold	shared.Resources
}

// client is Lucy.
type client struct {
	*baseclient.BaseClient

	forageHistory 		ForageHistory

	cpRequestHistory  	CPRequestHistory
	cpAllocationHistory	CPAllocationHistory

	taxAmount     		shared.Resources

	// allocation is the president's response to your last common pool resource request
	allocation 			shared.Resources

	config 				clientConfig
}

func NewClient(clientID shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:    baseclient.NewClient(clientID),
		forageHistory: ForageHistory{},
		cpRequestHistory: CPRequestHistory{},
		cpAllocationHistory: CPAllocationHistory{},
		taxAmount:     0,
		allocation:    0,
		config: clientConfig{
			randomForageTurns: 5,      // Richness Thresholds      
			JBThreshold:       100,      
			MiddleThreshold:   60,     
			ImperialThreshold: 30,	
		},
	}
}

func (c client) wealth() Wealth {  
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

//------------------------------------Comunication--------------------------------------------------------//
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

//--------------------------------------------------COMMONPOOL---------------------------------------------//
// Request resource from the President
func (c *client) CommonPoolResourceRequest() shared.Resources {
	// Initially, request the minimum
	newRequest := c.config.ImperialThreshold
	if c.gameState().Turn == 1 && c.gameState().Season == 1 {
		newRequest = c.config.ImperialThreshold
		c.Logf("Common Pool Request: %v", newRequest)
	} else if c.wealth() == Imperial_Student || c.wealth() == Dying {
		// If we are as poor as imperial student, request more resource from cp
		if c.config.ImperialThreshold <  (c.gameState().CommonPool / 6) {
			newRequest = c.gameState().CommonPool / 6
		} else {
			newRequest = c.config.ImperialThreshold
		}
	} else {
		// For other scenarios, look at the history and make decisions based on that 
		lastAllocation := c.cpAllocationHistory[len(c.cpAllocationHistory)-1]
		if lastAllocation == 0 {
			newRequest = c.config.ImperialThreshold
			c.Logf("Common Pool Request: %v", newRequest)
		} else {
			newRequest = lastAllocation + 10
			c.Logf("Common Pool Request: %v", newRequest)
		}
	}
	//Update request history
	c.cpRequestHistory = append(c.cpRequestHistory, newRequest)
	return newRequest
}

//Request resource from the common pool
func (c *client) RequestAllocation() shared.Resources {
	var allocation shared.Resources
	c.Logf("Current cp allocation amount: %v", c.allocation)
	
	// if we are poor we get this amount no matter the approval
	if c.wealth() == Imperial_Student || c.wealth() == Dying  {
		allocation = c.cpRequestHistory[len(c.cpRequestHistory)-1]
	} else {
		allocation = c.allocation
	}
	c.Logf("Taking %v from common pool", allocation)
	c.cpAllocationHistory = append(c.cpAllocationHistory, c.allocation)
	c.Logf("cpAllocationHistory: ", c.cpAllocationHistory)
	return allocation
}

//----------------------------------------------TAX----------------------------------------------//
func (c *client) GetTaxContribution() shared.Resources {
	c.Logf("Current tax amount: %v", c.taxAmount)
	if c.wealth() == Imperial_Student || c.wealth() == Dying {
		return 0
	}
	if c.gameState().Turn == 1 {
		return 0.5* c.taxAmount
	}
	return c.taxAmount
}


func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

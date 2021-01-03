package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Request resource from the President
func (c *client) CommonPoolResourceRequest() shared.Resources {
	// Initially, request the minimum
	newRequest := c.config.ImperialThreshold
	if c.gameState().Turn == 1 && c.gameState().Season == 1 {
		newRequest = c.config.ImperialThreshold
		c.Logf("Common Pool Request: %v", newRequest)
	} else if c.wealth() == Imperial_Student || c.wealth() == Dying {
		// If we are as poor as imperial student, request more resource from cp (whichever number is higher)
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
		//allocation = c.cpRequestHistory[len(c.cpRequestHistory)-1] //this one not working rn but it should be the same as the longer code. 
		// will be fixed by preet's team
		if c.config.ImperialThreshold <  (c.gameState().CommonPool / 6) {
			allocation = c.gameState().CommonPool / 6
		} else {
			allocation = c.config.ImperialThreshold
		}
	} else {
		allocation = c.allocation
	}
	c.Logf("Taking %v from common pool", allocation)
	c.cpAllocationHistory = append(c.cpAllocationHistory, c.allocation)
	c.Logf("cpAllocationHistory: ", c.cpAllocationHistory)
	return allocation
}
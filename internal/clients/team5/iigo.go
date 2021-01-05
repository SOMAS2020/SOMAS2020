package team5

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Request resource from the President
// This function asking permission from the President to take resource from the commonpool legally
// The President will reply with an allocation amount
func (c *client) CommonPoolResourceRequest() shared.Resources {
	// Initially, request the minimum
	reqAmount := c.config.imperialThreshold
	turn := c.getTurn()
	currentCP := c.getCP()

	if turn == 1 && c.getSeason() == 1 {
		reqAmount = c.config.imperialThreshold
	} else if c.wealth() == imperialStudent || c.wealth() == dying {
		// If we are as poor as imperial student, request more resource from cp (whichever number is higher)
		reqAmount = shared.Resources(math.Max(float64(c.config.imperialThreshold), float64(currentCP/6)))

	} else if turn > 1 {
		// For other scenarios, look at the history and make decisions based on that
		lastAllocation := c.cpAllocationHistory[turn-1]
		if lastAllocation == 0 {
			reqAmount = c.config.imperialThreshold
		} else {
			reqAmount = lastAllocation + 10
		}
	}
	c.Logf("Submitting CP resource request of %v resources", reqAmount)
	//Update request history
	c.cpRequestHistory[turn] = reqAmount
	return reqAmount
}

// Take resource from the common pool ideally from the allocation given by President
// If we are in imperial state, we may take resources regardless what the Presdient say (steal mode)
func (c *client) RequestAllocation() shared.Resources {
	var allocation shared.Resources
	c.Logf("Current cp allocation amount: %v", c.allocation)

	// if we are poor we get this amount no matter the approval
	if c.wealth() == imperialStudent || c.wealth() == dying {
		//allocation = c.cpRequestHistory[len(c.cpRequestHistory)-1] //this one not working rn but it should be the same as the longer code.
		// will be fixed by preet's team
		if c.config.imperialThreshold < (c.gameState().CommonPool / 6) {
			allocation = c.gameState().CommonPool / 6
		} else {
			allocation = c.config.imperialThreshold
		}
	} else {
		allocation = c.allocation
	}
	c.Logf("Taking %v from common pool", allocation)
	c.cpAllocationHistory[c.getTurn()] = c.allocation
	c.Logf("cpAllocationHistory: %v", c.cpAllocationHistory)
	return allocation
}

func (c *client) commonPoolContribution() shared.Resources {
	turn := c.getTurn()
	var contribution shared.Resources
	// Day 1 we don't contribute anything
	if turn == 1 {
		contribution = 0
	} else {
		difference := c.cpResourceHistory[turn] - c.cpResourceHistory[turn-1]
		if difference < 0 {
			contribution = 0
		} else {
			contribution = difference / 6
		}
	}
	return contribution
}

// Pay tax and contribution to the common pool
func (c *client) GetTaxContribution() shared.Resources {
	c.Logf("Current tax amount: %v", c.taxAmount)
	if c.wealth() == imperialStudent || c.wealth() == dying {
		return 0
	}
	if c.gameState().Turn == 1 && c.gameState().Season == 1 {
		return 0.5 * c.taxAmount
	}
	contribution := c.commonPoolContribution()
	return c.taxAmount + contribution
}

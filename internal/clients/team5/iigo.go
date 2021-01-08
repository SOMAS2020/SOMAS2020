package team5

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Request resource from the President
// This function asking permission from the President to take resource from the commonpool legally
// The President will reply with an allocation amount
func (c *client) CommonPoolResourceRequest() shared.Resources {
	c.evaluateRoles()
	// Initially, request the minimum
	turn := c.getTurn()
	season := c.getSeason()
	currentCP := c.getCP()
	currentTier := c.wealth()

	reqAmount := c.calculateRequestAllocation(turn, season, currentTier, currentCP)

	c.Logf("Submitting CP resource request of %v resources", reqAmount)
	//Update request history
	c.cpRequestHistory[turn] = reqAmount
	return reqAmount
}

// Take resource from the common pool ideally from the allocation given by President
// If we are in imperial state, we may take resources regardless what the Presdient say (steal mode)
func (c *client) RequestAllocation() shared.Resources {
	currentCP := c.getCP()
	currentTier := c.wealth()
	c.Logf("Current cp allocation amount: %v", c.allocation)
	allocation := c.calculateRequestCommonPool(currentTier, currentCP)

	c.Logf("Taking %v from common pool", allocation)
	c.cpAllocationHistory[c.getTurn()] = c.allocation
	c.Logf("cpAllocationHistory: %v", c.cpAllocationHistory)
	return allocation
}

// Pay tax and contribution to the common pool
func (c *client) GetTaxContribution() shared.Resources {
	c.Logf("Current tax amount: %v", c.taxAmount)
	turn := c.getTurn()
	season := c.getSeason()
	currentTier := c.wealth()

	taxContribution := c.calculateTaxContribution(turn, season, currentTier)
	contribution := c.calculateCPContribution(turn, season)
	total := taxContribution + contribution
	return total
}

//MonitorIIGORole decides whether to perform monitoring on a role
//COMPULOSRY: must be implemented
//always monitor a role
func (c *client) MonitorIIGORole(roleName shared.Role) bool {
	return true
}

//DecideIIGOMonitoringAnnouncement decides whether to share the result of monitoring a role and what result to share
//COMPULSORY: must be implemented
// always broadcast monitoring result
func (c *client) DecideIIGOMonitoringAnnouncement(monitoringResult bool) (resultToShare bool, announce bool) {
	resultToShare = monitoringResult
	announce = true
	return
}

// Cannot evade sanction :c
func (c *client) GetSanctionPayment() shared.Resources {
	return c.sanctionAmount
}

func (c *client) calculateRequestCommonPool(currentTier wealthTier, currentCP shared.Resources) (allocation shared.Resources) {
	if currentTier == imperialStudent || currentTier == dying {
		if c.config.imperialThreshold < (currentCP / 6) {
			if (currentCP / 6) > c.allocation {
				allocation = currentCP / 6
			} else {
				allocation = c.allocation
			}
		} else {
			if c.config.imperialThreshold > c.allocation {
				allocation = c.config.imperialThreshold
			} else {
				allocation = c.allocation
			}
		}
	} else {
		allocation = c.allocation
	}
	return allocation
}

func (c *client) calculateRequestAllocation(turn uint, season uint, currentTier wealthTier, currentCP shared.Resources) (reqAmount shared.Resources) {
	if turn == 1 && season == 1 {
		reqAmount = c.config.imperialThreshold
	} else if currentTier == imperialStudent || currentTier == dying {
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
	return
}

func (c *client) calculateTaxContribution(turn uint, season uint, currentTier wealthTier) (contribution shared.Resources) {
	if turn == 1 && season == 1 {
		contribution = 0.5 * c.taxAmount
	} else if currentTier == imperialStudent || currentTier == dying {
		contribution = 0
		return contribution
	} else {
		contribution = c.taxAmount
	}
	return
}

// Calculate our contribution to common pool
func (c *client) calculateCPContribution(turn uint, season uint) (contribution shared.Resources) {
	// Day 1 we don't contribute anything
	if turn == 1 && season == 1 {
		contribution = 0
	} else { // other days we contribute based on cashflow of commonpool
		difference := c.cpResourceHistory[turn] - c.cpResourceHistory[turn-1]
		if difference < 0 {
			contribution = 0
		} else {
			contribution = difference / 6
		}
	}
	return contribution
}
package team1

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) GetTaxContribution() shared.Resources {
	taxDecisionMade := c.BaseClient.LocalVariableCache[rules.TaxDecisionMade].Values[0] != 0

	if !taxDecisionMade && c.gameState().Turn == 1 {
		// Put in some initial resources
		contribution := c.config.kickstartTaxPercent * c.gameState().ClientInfo.Resources
		c.Logf("Paying initial contribution: %v", contribution)
		return contribution
	}

	if !taxDecisionMade || c.config.evadeTaxes {
		c.Logf("Tax decision not made")
		return 0
	}

	// taxDecisionMade && !c.config.evadeTaxes && c.gameState().Turn > 1
	expectedTaxContribution := shared.Resources(c.BaseClient.LocalVariableCache[rules.ExpectedTaxContribution].Values[0])
	c.Logf("Paying tax: %v", expectedTaxContribution)
	return expectedTaxContribution
}

func (c *client) CommonPoolResourceRequest() shared.Resources {
	switch c.emotionalState() {
	case Anxious:
		c.Logf("Common pool request: 20")
		return 20
	default:
		return 0
	}
}

func (c *client) RequestAllocation() shared.Resources {
	allocationMade := c.BaseClient.LocalVariableCache[rules.AllocationMade].Values[0] != 0
	allocationAmount := shared.Resources(0)

	if allocationMade {
		allocationAmount = shared.Resources(c.BaseClient.LocalVariableCache[rules.ExpectedAllocation].Values[0])
	} else if c.emotionalState() == Desperate {
		allocationAmount = c.config.desperateStealAmount
	}

	if allocationAmount != 0 {
		c.Logf("Taking %v from common pool", allocationAmount)
	}
	return allocationAmount
}

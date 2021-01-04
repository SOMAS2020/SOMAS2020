package team1

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

func (c *client) GetTaxContribution() shared.Resources {
	if c.config.evadeTaxes {
		return 0
	}
	if c.gameState().Turn == 1 {
		// Put in some initial resources to get IIGO to run
		return c.config.kickstartTaxPercent * c.gameState().ClientInfo.Resources
	}
	return c.taxAmount
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
	var allocation shared.Resources

	if c.emotionalState() == Desperate {
		allocation = c.config.desperateStealAmount
	} else if c.allocation != 0 {
		allocation = c.allocation
		c.allocation = 0
	}

	if allocation != 0 {
		c.Logf("Taking %v from common pool", allocation)
	}
	return allocation
}

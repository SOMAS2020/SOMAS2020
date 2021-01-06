package team1

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) GetTaxContribution() shared.Resources {
	contribution, success := c.GetRecommendation(rules.IslandTaxContribution)
	if !success {
		c.Logf("Cannot determine correct tax, paying 0")
		return 0
	}
	if c.config.evadeTaxes {
		c.Logf("Evading tax")
		return 0
	}
	c.Logf("Paying tax: %v", contribution)
	return shared.Resources(contribution.Values[0])
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

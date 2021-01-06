package team1

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
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
		amount := 2 * c.ServerReadHandle.GetGameConfig().CostOfLiving
		c.Logf("Common pool request: %v", amount)
		return amount
	default:
		return c.gameState().CommonPool
	}
}

func (c *client) RequestAllocation() shared.Resources {
	c.LocalVariableCache[rules.IslandAllocation] = rules.VariableValuePair{
		VariableName: rules.IslandAllocation,
		Values:       []float64{float64(c.gameState().CommonPool)},
	}
	allocationPair, success := c.GetRecommendation(rules.IslandAllocation)
	if !success {
		c.Logf("Cannot determine allocation, trying to get all resources in CP.")
		return c.gameState().CommonPool
	}

	if c.emotionalState() == Desperate {
		allocationAmount := c.config.desperateStealAmount
		c.Logf("Desperate for %v to stay alive.", allocationAmount)
		return allocationAmount
	}

	allocation := allocationPair.Values[0]
	if allocation != 0 {
		c.Logf("Taking %v from common pool", allocation)
	}
	return shared.Resources(allocation)
}

	if allocationAmount != 0 {
		c.Logf("Taking %v from common pool", allocationAmount)
	}
	return allocationAmount
}

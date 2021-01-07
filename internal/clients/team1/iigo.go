package team1

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

/**********************/
/*    Presidency      */
/**********************/

func (c *client) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {
	taxAmountMap := make(map[shared.ClientID]shared.Resources)

	for clientID, clientReport := range islandsResources {
		if !clientReport.Reported {
			// If they are not reporting their wealth, they will probably also
			// not pay tax.
			taxAmountMap[clientID] = 9999999
			c.addToOpinion(clientID, -1)
			continue
		}

		taxRate := 0.0 // Percentage of their resources the client will be taxed
		resources := clientReport.ReportedAmount
		livingCost := c.BaseClient.ServerReadHandle.GetGameConfig().CostOfLiving

		switch {
		case resources < 2*livingCost:
			taxRate = 0
		case resources < 10*livingCost:
			taxRate = 0.05
		case resources < 50*livingCost:
			taxRate = 0.3
		case resources < 100*livingCost:
			taxRate = 0.4
		default:
			taxRate = 0.6
		}

		// https://bit.ly/3s7dRXt
		taxAmountMap[clientID] = shared.Resources(float64(clientReport.ReportedAmount) * taxRate)
	}
	return shared.PresidentReturnContent{
		ContentType: shared.PresidentTaxation,
		ResourceMap: taxAmountMap,
		ActionTaken: true,
	}
}

/*************************/
/* Taxes and allocations */
/*************************/

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
		amount := 2 * c.gameConfig().CostOfLiving
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

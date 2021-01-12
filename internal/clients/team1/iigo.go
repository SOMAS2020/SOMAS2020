package team1

import (
	"math"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) GetClientPresidentPointer() roles.President {
	return c
}

/**********************/
/*    Presidency      */
/**********************/

func (c *client) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {
	taxAmountMap := make(map[shared.ClientID]shared.Resources)

	for clientID, clientReport := range islandsResources {
		c.reportedResources[clientID] = clientReport.Reported

		if !clientReport.Reported {
			// If they are not reporting their wealth, they will probably also
			// not pay tax.
			taxAmountMap[clientID] = 20 * c.gameConfig().CostOfLiving
			c.teamOpinions[clientID]--
			continue
		}

		taxRate := 0.0 // Percentage of their resources the client will be taxed
		resources := clientReport.ReportedAmount
		livingCost := c.BaseClient.ServerReadHandle.GetGameConfig().CostOfLiving

		switch {
		case resources < 2*livingCost:
			taxRate = 0
		case resources < 10*livingCost:
			taxRate = 0.1
		case resources < 50*livingCost || clientID == c.GetID():
			// Common pool feels empty :/
			taxRate = 0.3
		case resources < 100*livingCost:
			taxRate = 0.4
		default:
			taxRate = 0.6
		}

		// https://bit.ly/3s7dRXt
		taxAmountMap[clientID] = shared.Resources(float64(clientReport.ReportedAmount) * taxRate)
	}
	c.Logf("[IIGO] taxAmountMap: %v", taxAmountMap)
	return shared.PresidentReturnContent{
		ContentType: shared.PresidentTaxation,
		ResourceMap: taxAmountMap,
		ActionTaken: true,
	}
}

// Make allocations as the base president would, but not to agents that are not
// reporting their wealth.
func (c *client) EvaluateAllocationRequests(
	resourceRequests map[shared.ClientID]shared.Resources,
	commonPool shared.Resources) shared.PresidentReturnContent {

	chosenRequests := map[shared.ClientID]shared.Resources{}

	for clientID, request := range resourceRequests {
		reportedResources, noData := c.reportedResources[clientID]
		if reportedResources || noData {
			switch {
			case noData:
				c.Logf("Granting request to %v, no data on tax evasion", clientID)
			case reportedResources:
				c.Logf("Granting request to %v, they reported resources", clientID)
			}
			chosenRequests[clientID] = request
		}
	}

	return c.BasePresident.EvaluateAllocationRequests(chosenRequests, commonPool)

}

/*************/
/*  Voting   */
/*************/

func (c *client) VoteForElection(
	roleToElect shared.Role,
	candidateList []shared.ClientID) []shared.ClientID {

	opinionRank := sortByOpinion{}

	for _, candidate := range candidateList {
		if candidate == c.GetID() {
			opinionRank = append(opinionRank, opinionOnTeam{
				clientID: candidate,
				opinion:  math.MaxInt64,
			})
		} else {
			opinionRank = append(opinionRank, opinionOnTeam{
				clientID: candidate,
				opinion:  c.teamOpinions[candidate],
			})
		}
	}
	sort.Sort(opinionRank)

	ballot := []shared.ClientID{}

	for _, opinionOnTeam := range opinionRank {
		ballot = append(ballot, opinionOnTeam.clientID)
	}

	return ballot
}

/*************************/
/* Taxes and allocations */
/*************************/

func (c *client) GetTaxContribution() shared.Resources {
	valToBeReturned := shared.Resources(0)
	c.LocalVariableCache[rules.IslandTaxContribution] = rules.VariableValuePair{
		VariableName: rules.IslandTaxContribution,
		Values:       []float64{float64(valToBeReturned)},
	}
	c.Logf("[IIGO]: Expected Tax Contribution: %v", c.LocalVariableCache[rules.ExpectedTaxContribution].Values)
	contribution, success := c.GetRecommendation(rules.IslandTaxContribution)
	if !success {
		c.Logf("Cannot determine correct tax, paying 0")
		return 0
	}
	if c.config.evadeTaxes {
		c.Logf("Evading tax")
		return 0
	}
	c.Logf("[IIGO]: Paying tax: %v", contribution)
	return shared.Resources(contribution.Values[0])
}

func (c *client) CommonPoolResourceRequest() shared.Resources {
	switch c.emotionalState() {
	case Normal:
		return shared.Resources(2 * float64(c.gameConfig().CostOfLiving))
	case Desperate, Anxious:
		amount := shared.Resources(c.config.resourceRequestScale) * c.gameConfig().CostOfLiving
		c.Logf("Common pool request: %v", amount)
		return amount
	default:
		return c.gameState().CommonPool
	}
}

// Gets called at the end of IIGO
func (c *client) RequestAllocation() shared.Resources {
	if c.emotionalState() == Desperate && c.config.desperateStealAmount != 0 {
		allocation := c.config.desperateStealAmount
		c.Logf("Desperate for %v to stay alive.", allocation)
		return shared.Resources(
			math.Min(float64(allocation), float64(c.gameState().CommonPool)),
		)
	}

	c.LocalVariableCache[rules.IslandAllocation] = rules.VariableValuePair{
		VariableName: rules.IslandAllocation,
		Values:       []float64{float64(c.ServerReadHandle.GetGameState().CommonPool)},
	}

	allocationPair, success := c.GetRecommendation(rules.IslandAllocation)
	if !success {
		c.Logf("Cannot determine allocation, trying to get all resources in CP.")
		return c.gameState().CommonPool
	}

	// Unintentionally nicking from commonPool so limiting amount. GetRecommendation is too powerful.
	allocation := allocationPair.Values[0]
	if allocation != 0 {
		c.Logf("Taking %v from common pool", allocation)
	}
	return shared.Resources(
		math.Min(allocation, float64(c.gameState().CommonPool)),
	)
}

// ResourceReport is an island's self-report of its own resources. This is called by
// the President to help work out how many resources to allocate each island.
func (c *client) ResourceReport() shared.ResourcesReport {
	amountReported := c.gameState().ClientInfo.Resources
	c.Logf("[IIGO]: amountReported %v", amountReported)
	c.LocalVariableCache[rules.IslandReportedResources] = rules.MakeVariableValuePair(rules.IslandReportedResources, []float64{float64(amountReported)})
	return shared.ResourcesReport{
		ReportedAmount: amountReported,
		Reported:       true,
	}
}

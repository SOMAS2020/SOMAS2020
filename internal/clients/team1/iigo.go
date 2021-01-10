package team1

import (
	"fmt"
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
			taxAmountMap[clientID] = 9999999
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
		opinionRank = append(opinionRank, opinionOnTeam{
			clientID: candidate,
			opinion:  c.teamOpinions[candidate],
		})
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
	contribution, success := c.BaseClient.GetRecommendation(rules.IslandTaxContribution)
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
		fmt.Printf("Amount %v", amount)
		c.Logf("Common pool request: %v", amount)
		return amount
	default:
		return c.gameState().CommonPool
	}
}

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
		Values:       []float64{float64(c.gameState().CommonPool)},
	}

	allocationPair, success := c.GetRecommendation(rules.IslandAllocation)
	if !success {
		c.Logf("Cannot determine allocation, trying to get all resources in CP.")
		return c.gameState().CommonPool
	}

	allocation := allocationPair.Values[0]
	if allocation != 0 {
		c.Logf("Taking %v from common pool", allocation)
	}
	return shared.Resources(
		math.Min(allocation, float64(c.gameState().CommonPool)),
	)
}

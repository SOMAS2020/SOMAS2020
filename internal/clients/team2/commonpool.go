package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Updates Common Pool History with current Common Pool Level
func (c *client) commonPoolUpdate() {
	c.commonPoolHistory[c.gameState().Turn] = c.gameState().CommonPool
	c.Logf("Common Pool History updated: ", c.commonPoolHistory)
}

// Updates Resource Level History with our current resource Level
func (c *client) resourceHistoryUpdate(resourceLevelHistory ResourcesLevelHistory) {
	c.resourceLevelHistory[c.gameState().Turn] = c.gameState().ClientInfo.Resources
	c.Logf("Resource Level History updated: ", resourceLevelHistory)
}

// Updates Pres Common Pool History with current resource request
func (c *client) presCommonPoolUpdate(request shared.Resources) {
	// Initialises a complete commonPool update
	c.Logf("We request ", request)
	commonPool := CommonPoolInfo{
		tax:             0,
		requestedToPres: request,
		allocatedByPres: 0,
		takenFromCP:     0,
	}
	// Checks if the map is not initialised
	if _, ok := c.presCommonPoolHist[c.gameState().PresidentID]; !ok {
		c.presCommonPoolHist[c.gameState().PresidentID] = make(map[uint]CommonPoolInfo)
		c.presCommonPoolHist[c.gameState().PresidentID][c.gameState().Turn] = commonPool

		c.Logf("Initialised presCommonPoolHist and added request", c.presCommonPoolHist)
	} else if pastCommonPool, ok := c.presCommonPoolHist[c.gameState().PresidentID][c.gameState().Turn]; ok {
		// If we have a previous entry, update the requestedToPres
		commonPool = pastCommonPool
		commonPool.requestedToPres = request
		c.Logf("President Common Pool History updated", c.presCommonPoolHist)

	}
	c.presCommonPoolHist[c.gameState().PresidentID][c.gameState().Turn] = commonPool
}

// Returns a resource request to ask the President for from the common pool
// of type shared.Resources and updates presCommonPoolHist
func (c *client) CommonPoolResourceRequest() shared.Resources {
	// DVP - if we are heavily affected - worry more about ours - 0.25 - 1 (1 highly vulnerable)
	agentDVP := c.islandDVPs[c.GetID()]

	// 1 + agentDVP - 1.25x if vuln 1.8x
	// Base request is amount needed to survive 1 more turn
	request := c.taxAmount + c.gameConfig().CostOfLiving

	if c.criticalStatus() {
		request = request*3 + c.gameConfig().MinimumResourceThreshold
		c.Logf("Critical status! Set Common Pool request to: ", request)
	} else if c.getAgentExcessResources() == shared.Resources(0) || agentDVP > 0.6 {
		request = request*2 + c.gameConfig().MinimumResourceThreshold
		c.Logf("Resources getting LOW or at risk! Gather up resources for one turn to stay safe")
	} else if c.getAgentExcessResources() > shared.Resources(0) {
		multiplier := shared.Resources(0)
		switch c.setAgentStrategy() {
		case Altruist:
			multiplier = 0.4
		case FairSharer:
			multiplier = 0.5
		case Selfish:
			multiplier = 1
		}
		request = request*multiplier + c.gameConfig().MinimumResourceThreshold
		c.Logf("Resources in Excess! Request based on agent Strategy: ", c.taxAmount)
	} else {
		request = c.taxAmount
	}

	c.presCommonPoolUpdate(request)

	return request
}

// Determines how many resources you actually take
func (c *client) RequestAllocation() shared.Resources {
	// this ignores the presidents allocation
	request := c.CommonPoolResourceRequest()

	switch c.getAgentStrategy() {
	case Selfish:
		spiteFactor := shared.Resources(1.2)
		if request > c.commonPoolAllocation {
			return request * spiteFactor
		}

		request = c.commonPoolAllocation
	case Altruist:
		request = Min(request, c.commonPoolAllocation)
	default:
		// Leave request as is if Fair Strategy
	}

	commonPool := CommonPoolInfo{
		tax:             0,
		requestedToPres: 0,
		allocatedByPres: 0,
		takenFromCP:     request,
	}

	// Checks if the map is not initialised
	if _, ok := c.presCommonPoolHist[c.gameState().PresidentID]; !ok {
		c.presCommonPoolHist[c.gameState().PresidentID] = make(map[uint]CommonPoolInfo)
		c.presCommonPoolHist[c.gameState().PresidentID][c.gameState().Turn] = commonPool
		c.Logf("Initialised presCommonPoolHist and added request", c.presCommonPoolHist)
	} else if pastCommonPool, ok := c.presCommonPoolHist[c.gameState().PresidentID][c.gameState().Turn]; ok {
		// If we have a previous entry, update the requestedToPres
		commonPool = pastCommonPool
		commonPool.takenFromCP = request
		c.Logf("President Common Pool History updated", c.presCommonPoolHist)
	}

	c.presCommonPoolHist[c.gameState().PresidentID][c.gameState().Turn] = commonPool

	return request
}

// GetTaxContribution determines how much we put into pool
func (c *client) GetTaxContribution() shared.Resources {
	contribution := shared.Resources(0)
	strategicContribution := c.getStrategicContribution()

	if c.getAgentExcessResources() == 0 || c.criticalStatus() {
		contribution = shared.Resources(0)
	} else if c.getAgentExcessResources() < strategicContribution {
		// Not enough excess to mitigate disasters oh no...
		// hoard some for ourselves
		contribution = c.getAgentExcessResources() / 2
	} else if c.getAgentExcessResources() < strategicContribution+c.taxAmount {
		// enough to mitigate disasters but not pay tax so screw tax
		contribution = strategicContribution
	} else {
		// enough to mitigate disasters and pay tax
		contribution = strategicContribution + c.taxAmount
	}

	c.updatePresidentTrust()
	c.confidenceRestrospect("President", c.gameState().PresidentID)
	c.Logf("Common Pool Contribution: ", contribution)

	return contribution
}

// Returns the number of turns until the next disaster (could be prediction/known value)
func (c *client) getTurnsUntilDisaster() uint {
	currentTurn := c.gameState().Turn

	// add extra when disaster is soon
	if c.gameConfig().DisasterConfig.DisasterPeriod.Valid {
		period := c.gameConfig().DisasterConfig.DisasterPeriod.Value
		return period - (currentTurn % period)
	}

	// Not able to predict disasters in first season as no prev known data
	if c.gameState().Season == 1 {
		return c.config.InitialDisasterTurnGuess - currentTurn
	} else {
		return c.CombinedDisasterPred.TimeLeft
	}
}

// Determines resources we need to be above critical, pay tax and cost of living, put resources aside proportional to incoming disaster
// TODO: could also be used to return negatives and pass that on to gifts and common pool request to know how much we need
// How much we can spend
func (c *client) getAgentExcessResources() shared.Resources {
	// At a minimum we should be able to pay cost of living
	excess := c.gameState().ClientInfo.Resources
	excess -= c.gameConfig().CostOfLiving - c.gameConfig().MinimumResourceThreshold
	return Min(0, excess)
}

// getStrategicContribution finds the best contribution to the common pool based on the method of play we
// are in and whether the threshold is known. If so we should use our disaster prediction also.
// Returns  proposed contribution as shared.Resources
func (c *client) getStrategicContribution() shared.Resources {
	strategicContribution := c.config.DefaultContribution
	disasterContribution := c.config.DefaultContribution

	if c.gameState().Turn == 1 {
		return strategicContribution
	}

	if c.gameConfig().DisasterConfig.CommonpoolThreshold.Valid {
		missingResources := Max(0, c.gameConfig().DisasterConfig.CommonpoolThreshold.Value-c.gameState().CommonPool)
		if missingResources != 0 {
			// If we know the threshold get a disaster contribution, if not always use strategic
			// trust no one - not even the server to stop the game
			if c.getNumAliveClients() != 0 && c.getTurnsUntilDisaster() != 0 {
				disasterContribution = missingResources / shared.Resources(c.getTurnsUntilDisaster()) / shared.Resources(c.getNumAliveClients())
			} else if c.getNumAliveClients() != 0 {
				disasterContribution = missingResources / shared.Resources(c.getNumAliveClients())
			}
		}
	}

	switch c.setAgentStrategy() {
	case FairSharer:
		// contribute the weighted average contribution
		ResourceHistory := c.commonPoolHistory
		runningAverageCPChange := shared.Resources(0)

		// compute running average
		for j := c.gameState().Turn; j > 0; j-- {
			runningAverageCPChange = runningAverageCPChange + (ResourceHistory[j]-ResourceHistory[j-1]-runningAverageCPChange)/shared.Resources(j)
		}

		// trust no one - not even the server to stop the game
		if shared.Resources(c.getNumAliveClients()) != 0 {
			strategicContribution = runningAverageCPChange / shared.Resources(c.getNumAliveClients())
		}

		// If we know the threshold take average between strategic and disaster contributions
		if c.gameConfig().DisasterConfig.CommonpoolThreshold.Valid {
			return (strategicContribution + disasterContribution) / shared.Resources(2)
		}

		// otherwise just return strategic
		return strategicContribution
	case Altruist:
		// contribute weighted average contribution multiplied by a factor
		ResourceHistory := c.commonPoolHistory
		runningAverageCPChange := shared.Resources(0)

		// compute running average
		for j := c.gameState().Turn; j > 0; j-- {
			runningAverageCPChange = runningAverageCPChange + (ResourceHistory[j]-ResourceHistory[j-1]-runningAverageCPChange)/shared.Resources(j)
		}

		// TODO: setting this multiplier in the config
		// do not trust anyone - even the server - check for divide by 0
		if shared.Resources(c.getNumAliveClients()) == shared.Resources(0) {
			strategicContribution = runningAverageCPChange * 1.2 / shared.Resources(c.getNumAliveClients())
		}

		return Max(strategicContribution, disasterContribution)
	default:
		// if we are Selfish contribute nothing shared.Resources(0)
		return shared.Resources(0)
	}
}

// Pays sanction unless we are being selfish or if we are critical
func (c *client) GetSanctionPayment() shared.Resources {
	if value, ok := c.LocalVariableCache[rules.SanctionExpected]; ok {
		if c.criticalStatus() || c.getAgentStrategy() == Selfish {
			c.Logf("Yeah I don't know about those sanctions...not feeling like it :P")
			return 0

		} else {
			c.Logf("Not happy about it but okay...we'll pay your sanction")
			return shared.Resources(value.Values[rules.SingleValueVariableEntry])
		}
	}
	return 0
}

// Returns the intended Contribution to the teams selected to share it with
func (c *client) ShareIntendedContribution() shared.IntendedContribution {
	shareWith := make([]shared.ClientID, 0)
	aliveClients := c.getAliveClients()

	// Share contributions with our friends
	for _, island := range aliveClients {
		if c.confidence("Gifts", island) > 30 {
			shareWith = append(shareWith, island)
		}
	}

	return shared.IntendedContribution{
		Contribution:   c.GetTaxContribution(),
		TeamsOfferedTo: shareWith,
	}
}

// TODO: this is completely empty
// func (c *client) ReceiveIntendedContribution(receivedIntendedContributions shared.ReceivedIntendedContributionDict) {
// we check how much each island intends to contribute
// Compute the average amount needed for the common pool threshold
// form an opinion based on how far their contribution is from the average
// could help us determine empathy level? (ie altruist, etc)
// }

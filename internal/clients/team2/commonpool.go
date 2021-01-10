package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Updates Common Pool History with current Common Pool Level
// TODO: this is never used
func CommonPoolUpdate(c *client, commonPoolHistory CommonPoolHistory) {
	c.commonPoolHistory[c.gameState().Turn] = c.gameState().CommonPool
	c.Logf("Common Pool History updated: ", commonPoolHistory)
}

// Updates Resource Level History with our current resource Level
// TODO: this is never used
func resourceHistoryUpdate(c *client, resourceLevelHistory ResourcesLevelHistory) {
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

// If we are critical request the full threshold to shift us back to security
// If our current resources are below the threshold request enough to reach the threshold
// If our resources are above the threshold request the tax
// So we can pay our tax for at least one turn (we may be granted less)
func (c *client) determineBaseCommonPoolRequest() shared.Resources {
	currResources := c.gameState().ClientInfo.Resources
	if c.criticalStatus() {
		c.Logf("Critical status! Set base request to agent threshold: ", c.agentThreshold())
		return c.agentThreshold()
	} else if currResources < c.agentThreshold() {
		c.Logf("Resources Low! Make up resources to reach agent threshold")
		return (c.agentThreshold() - currResources)
	} else {
		c.Logf("Resources Ok! Request tax amount from Common Pool: ", c.taxAmount)
		return c.taxAmount
	}
}

// Returns a resource request to ask the President for from the common pool
// of type shared.Resources and updates presCommonPoolHist
func (c *client) CommonPoolResourceRequest() shared.Resources {
	request := c.determineBaseCommonPoolRequest() * c.commonPoolMultiplier()
	// TODO: code is logging both requests in the common pool but one is the request
	// TODO: and the other is when it actually happens (in RequestAllocation)
	// TODO: should probably only log one
	c.presCommonPoolUpdate(request)

	return request
}

// type CommonPoolInfo struct {
// 	turn            uint
// 	tax             shared.Resources
// 	requestedToPres shared.Resources
// 	allocatedByPres shared.Resources
// 	takenFromCP     shared.Resources
// }

// Determines how many resources you actually take
func (c *client) RequestAllocation() shared.Resources {
	request := c.determineBaseCommonPoolRequest() * c.commonPoolMultiplier()
	// TODO: check if it's fine to just always take the biggest (before we also checked for status critical)
	request = Max(request, c.commonPoolAllocation)

	c.presCommonPoolUpdate(request)

	// This uses outdated logic without error handling
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

func (c *client) calculateContribution() shared.Resources {
	ourResources := c.gameState().ClientInfo.Resources
	// This is our default allocation, this determines how much to give based off of previous common pool level
	defaultAllocation := AverageCommonPoolDilemma(c) + c.taxAmount
	agentThreshold := c.agentThreshold()
	surplus := ourResources - agentThreshold

	if c.criticalStatus() || agentThreshold == 0 {
		// tax evasion by necessity
		return shared.Resources(0)
	} else if ourResources <= agentThreshold {
		// Below the threshold, pay proportion of taxAmount
		return c.taxAmount * (ourResources / agentThreshold)
	} else if checkOthersCrit(c) {
		// Others are in a critical state (Long term survival)
		return Min(surplus/c.config.HelpCritOthersDivisor, c.taxAmount)
	} else {
		// Give the smallest contribution
		return Min(surplus, defaultAllocation)
	}
}

// GetTaxContribution determines how much we put into pool
func (c *client) GetTaxContribution() shared.Resources {
	contribution := c.calculateContribution()

	c.updatePresidentTrust()
	c.confidenceRestrospect("President", c.gameState().PresidentID)

	return contribution
}

func (c *client) calculateDisasterMagPred() float64 {
	turn := c.gameState().Turn

	// If we know the common pool threshold
	if c.gameConfig().DisasterConfig.CommonpoolThreshold.Valid {
		return float64(c.gameConfig().DisasterConfig.CommonpoolThreshold.Value) / float64(c.getNumAliveClients())
	} else if len(c.disasterHistory) == 0 {
		// If we don't know the common pool threshold
		return float64(c.gameState().ClientInfo.Resources / c.config.BaseDisasterProtectionDivisor) //initial disaster threshold guess when we start playing
	} else {
		sampleMeanMag, magnitudePrediction := getMagnitudePrediction(c, float64(turn))

		// TODO: why are we accessing the first value here???
		baseThreshold := float64(c.resourceLevelHistory[1] / c.config.BaseResourcesToGiveDivisor)
		disasterBasedAdjustment := 0.0
		if c.checkForDisaster() {
			if c.resourceLevelHistory[turn] >= c.resourceLevelHistory[turn-1] { //no resources taken by disaster
				if disasterBasedAdjustment > 5 {
					disasterBasedAdjustment -= 5
				}
			}
			//disaster took our resources
			disasterBasedAdjustment += 5
		}
		//change factor by if next mag > or < prev mag
		return baseThreshold*(magnitudePrediction/sampleMeanMag) + disasterBasedAdjustment
	}
}

func (c *client) calculateTimeRemaining() float64 {
	turn := c.gameState().Turn
	// add extra when disaster is soon
	if c.gameConfig().DisasterConfig.DisasterPeriod.Valid {
		period := c.gameConfig().DisasterConfig.DisasterPeriod.Value
		return float64(period - (turn % period))
	}
	if c.gameState().Season == 1 { //not able to predict disasters in first season as no prev known data
		return c.config.InitialDisasterTurnGuess - float64(turn)
	} else {
		sampleMeanX, timeRemainingPrediction := getTimeRemainingPrediction(c, float64(turn))
		turnsLeftConfidence := getTimeRemainingConfidence(c, float64(turn), sampleMeanX)
		return float64(timeRemainingPrediction) * (turnsLeftConfidence / 100)
	}

}

// Determines esources we need to be above critical, pay tax and cost of living, put resources aside proportional to incoming disaster
func (c *client) agentThreshold() shared.Resources {
	criticalThreshold := c.gameConfig().MinimumResourceThreshold
	costOfLiving := c.gameConfig().CostOfLiving
	basicCosts := criticalThreshold + c.taxAmount + costOfLiving
	vulnerabilityMultiplier := 0.75 + getIslandDVPs(c.gameState().Geography)[c.GetID()] //1 to 1.75 (1.75 being the most vulnerable)

	// Add resources based on expected/predicted disaster magnitude
	disasterMagProtection := c.calculateDisasterMagPred()

	timeRemaining := c.calculateTimeRemaining()
	disasterTimeProtectionMultiplier := 1.0
	if timeRemaining < c.config.TimeLeftIncreaseDisProtection {
		disasterTimeProtectionMultiplier = c.config.DisasterSoonProtectionMultiplier
	}

	return basicCosts + shared.Resources(disasterTimeProtectionMultiplier*disasterMagProtection*vulnerabilityMultiplier)
}

// Checks if there was a disaster in the previous turn
func (c *client) checkForDisaster() bool {
	var prevSeason uint
	if c.gameState().Turn == 1 {
		prevSeason = 1
		return false
	}
	if prevSeason != c.gameState().Season {
		prevSeason++
		return true
	}
	return false
}

// AverageCommonPoolDilemma determines how much to contribute to the common pool depending on whether other agents are altruists,fair sharers or free riders
// TODO: improve comment description of what the function does
func AverageCommonPoolDilemma(c *client) shared.Resources {
	turn := c.gameState().Turn
	altruistContribution := c.determineAltruistContribution(turn)
	fairContribution := c.determineFairContribution(turn)

	if turn == 1 {
		return c.config.DefaultFirstTurnContribution
	}

	switch c.setAgentStrategy() {
	case 0:
		return shared.Resources(altruistContribution)
	case 1:
		return shared.Resources(fairContribution)
	default:
		// Use Selfish approach if neither case is matched
		return shared.Resources(0)
	}
}

func (c *client) determineAltruistContribution(turn uint) shared.Resources { //identical to fair sharing but a larger factor to multiple the average contribution by
	ResourceHistory := c.commonPoolHistory
	tuneAlt := shared.Resources(c.config.AltruistFactorOfAvToGive) //what factor of the average to contribute when being altruistic, will be much higher than fair sharing
	for j := turn; j > 0; j-- {                                    //we are trying to find the most recent instance of the common pool increasing and then use that value
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn] > 0 {
			if shared.Resources(c.getNumAliveClients())*tuneAlt != shared.Resources(0) {
				return ((ResourceHistory[j] - ResourceHistory[prevTurn]) / shared.Resources(c.getNumAliveClients())) * tuneAlt
			}
		}
	}
	return 0
}

// Can make more sophisticated! Right now just contribute the average
func (c *client) determineFairContribution(turn uint) shared.Resources {
	ResourceHistory := c.commonPoolHistory
	// What factor of the average to contribute when fair sharing, default is 1 to give the average
	tuneAverage := shared.Resources(c.config.FairShareFactorOfAvToGive)
	// We are trying to find the most recent instance of the common pool increasing and then use that value
	for j := turn; j > 0; j-- {
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn] > 0 {
			if shared.Resources(c.getNumAliveClients())*tuneAverage != shared.Resources(0) {
				return ((ResourceHistory[j] - ResourceHistory[prevTurn]) / shared.Resources(c.getNumAliveClients())) * tuneAverage
			}
		}
	}
	return 0
}

// TODO: RENAME MethodOfPlay to Agent Mode
func (c *client) commonPoolMultiplier() shared.Resources {
	var multiplier float64

	switch c.setAgentStrategy() {
	case Altruist:
		// when the pool is struggling, we will forage less to hav emo
		multiplier = 0.4
	case FairSharer:
		multiplier = 0.6
	case Selfish:
		multiplier = 1.2
	}

	return shared.Resources(multiplier)
}

// TODO: make switch case on agent mode.
func (c *client) SanctionHopeful() shared.Resources {
	switch c.setAgentStrategy() {

	}
	return 0
}

// Checks the sanction amount against what we expect
// TODO: this function is not implementing any logic
func (c *client) GetSanctionPayment() shared.Resources {
	if value, ok := c.LocalVariableCache[rules.SanctionExpected]; ok {
		if c.gameState().ClientLifeStatuses[c.GetID()] != shared.Critical {
			if shared.Resources(value.Values[0]) <= c.SanctionHopeful() {
				return shared.Resources(value.Values[0])
			} else {
				return c.SanctionHopeful()
			}
		} else {
			return 0
		}
	}
	return 0
}

// Returns the intended Contribution to the teams selected to share it with
func (c *client) ShareIntendedContribution() shared.IntendedContribution {
	shareWith := make([]shared.ClientID, 0)
	aliveClients := c.getAliveClients()

	for _, island := range aliveClients {
		if c.confidence("Gifts", island) > 30 {
			shareWith = append(shareWith, island)
		}
	}

	return shared.IntendedContribution{
		Contribution:   c.calculateContribution(),
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

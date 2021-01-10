package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Updates Common Pool History with current Common Pool Level
// TODO: this is never used
func commonPoolUpdate(c *client, commonPoolHistory CommonPoolHistory) {
	c.commonPoolHistory[c.gameState().Turn] = c.gameState().CommonPool
	c.Logf("Common Pool History updated: ", commonPoolHistory)
}

// Updates Resource Level History with our current resource Level
// TODO: this is never used
func resourceHistoryUpdate(c *client, resourceLevelHistory ResourcesLevelHistory) {
	c.resourceLevelHistory[c.gameState().Turn] = c.gameState().ClientInfo.Resources
	c.Logf("Resource Level History updated: ", resourceLevelHistory)
}

// type CommonPoolInfo struct {
// 	turn            uint
// 	tax             shared.Resources
// 	requestedToPres shared.Resources
// 	allocatedByPres shared.Resources
// 	takenFromCP     shared.Resources
// }
// TODO: this does not match - you are leaving half the values empty each time you access it
// Updates Pres Common Pool History with current resource request
func (c *client) presCommonPoolUpdate(request shared.Resources) {
	var commonPool CommonPoolInfo
	if _, ok := c.presCommonPoolHist[c.gameState().PresidentID]; !ok {
		c.presCommonPoolHist[c.gameState().PresidentID] = make([]CommonPoolInfo, 0)
		c.Logf("Initialised presCommonPoolHist", c.presCommonPoolHist)
	} else {
		presHist := c.presCommonPoolHist[c.gameState().PresidentID]
		// TODO: this logic doesn't make sense - we get the most recent item from the hist
		// TODO: set the turn and requestedToPres to the current turn (which we already know is the same)
		// TODO: and then append an item that already exists to the presCommonPoolHist
		if presHist[len(presHist)-1].turn == c.gameState().Turn {
			commonPool = presHist[len(presHist)-1]
			commonPool.turn = c.gameState().Turn
			commonPool.requestedToPres = request
			c.presCommonPoolHist[c.gameState().PresidentID] = append(presHist, commonPool)
			c.Logf("President Common Pool History updated", c.presCommonPoolHist)
		}
	}
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
// TODO: this does not match - you are leaving half the values empty each time you access it
// TODO: so it seems if you leave the push empty it makes the other values zero
// TODO: then you just have to update the values if they're already there instead of pushing another item onto the end
// Determines how many resources you actually take
func (c *client) RequestAllocation() shared.Resources {
	request := c.determineBaseCommonPoolRequest() * c.commonPoolMultiplier()
	c.presCommonPoolUpdate(request)

	// This uses outdated logic without error handling
	var commonPool CommonPoolInfo
	presHist := c.presCommonPoolHist[c.gameState().PresidentID]

	if len(presHist) != 0 {
		presHist[len(presHist)-1].takenFromCP = request
		presHist[len(presHist)-1].turn = c.gameState().Turn
	} else {
		commonPool = CommonPoolInfo{
			takenFromCP: request,
			turn:        c.gameState().Turn,
		}
	}

	// TODO: same bug here -> also appending incomplete objects
	// TODO: also, commonpool history is only assigned a value in the else so since this is outside we often append NaN
	c.presCommonPoolHist[c.gameState().PresidentID] = append(presHist, commonPool)

	if c.criticalStatus() && c.commonPoolAllocation < request {
		return request
	}

	return c.commonPoolAllocation
}

// // 1 to 1.75 (1.75 being the most vulnerable)
// vulnerabilityMultiplier := 0.75 + c.islandDVPs[c.GetID()]

// // Add resources based on expected/predicted disaster magnitude
// c.commonPoolThreshold = c.updateCommonPoolThreshold()
// timeRemaining := c.getTurnsUntilDisaster()

// //commonpool/timeRemaining -> give to pool each turn

// disasterTimeProtectionMultiplier := 1.0
// if timeRemaining < c.config.TimeLeftIncreaseDisProtection {
// 	disasterTimeProtectionMultiplier = c.config.DisasterSoonProtectionMultiplier
// }

// return basicCosts + shared.Resources(disasterTimeProtectionMultiplier*c.commonPoolThreshold*vulnerabilityMultiplier)

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

// if we know what threshold we need that's fine
// guess if we don't know
// once we have some knowledge on disasters, if it is going to be bigger - we increase the amount we're giving
// if our predictions were not enough before we want to increase it
func (c *client) updateCommonPoolThreshold() shared.Resources {
	// If we know the common pool threshold
	if c.gameConfig().DisasterConfig.CommonpoolThreshold.Valid {
		return c.gameConfig().DisasterConfig.CommonpoolThreshold.Value
	} else if len(c.disasterHistory) == 0 {
		// If we don't know the common pool threshold make a guess
		return c.gameState().ClientInfo.Resources * shared.Resources(c.config.InitialThresholdProportionGuess) * shared.Resources(c.getNumAliveClients())
	} else {
		magnitudePrediction := c.CombinedDisasterPred.Magnitude
		disasterOccurred, lastDisasterReport := c.getLastDisasterInfo()

		// TODO: this relies heavily on an accurate first guess
		// if disaster happened in the previous turn, update our prediction, otherwise
		if disasterOccurred && (lastDisasterReport.Turn == c.gameState().Turn-1) {
			return shared.Resources(magnitudePrediction/lastDisasterReport.Report.Magnitude) * c.commonPoolThreshold
		}
	}
}

// Returns the number of turns until the next disaster (could be prediction/known value)
func (c *client) getTurnsUntilDisaster() uint {
	currTurn := c.gameState().Turn

	// add extra when disaster is soon
	if c.gameConfig().DisasterConfig.DisasterPeriod.Valid {
		period := c.gameConfig().DisasterConfig.DisasterPeriod.Value
		return period - (currTurn % period)
	}

	// Not able to predict disasters in first season as no prev known data
	if c.gameState().Season == 1 {
		return c.config.InitialDisasterTurnGuess - currTurn
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
	excess -= c.gameConfig().CostOfLiving - c.gameConfig().MinimumResourceThreshold - c.taxAmount

	return Min(0, excess)
}

// Gets turn and report from last disaster. Returns true and DisasterOccurrence{} if a disaster occurred
// Else returns false and an empty DisasterOccurrence{}
func (c *client) getLastDisasterInfo() (bool, DisasterOccurrence) {
	lastDisaster := DisasterOccurrence{}
	numDisasters := len(c.disasterHistory)
	if numDisasters == 0 {
		return false, lastDisaster
	} else {
		// get the most recent disaster turn and report
		lastDisaster.Turn = c.disasterHistory[numDisasters-1].Turn
		lastDisaster.Report = c.disasterHistory[numDisasters-1].Report

		return true, lastDisaster
	}
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

// Identical to fair sharing but a larger factor to multiple the average contribution by
func (c *client) determineAltruistContribution(turn uint) shared.Resources {
	ResourceHistory := c.commonPoolHistory
	//what factor of the average to contribute when being altruistic, will be much higher than fair sharing
	tuneAlt := shared.Resources(c.config.AltruistFactorOfAvToGive)
	//we are trying to find the most recent instance of the common pool increasing and then use that value
	for j := turn; j > 0; j-- {
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

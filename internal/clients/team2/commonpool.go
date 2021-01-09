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
	request := c.determineBaseCommonPoolRequest() * methodConfPool(c)
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

// Determines how many resources you actually take
func (c *client) RequestAllocation() shared.Resources {
	request := c.determineBaseCommonPoolRequest() * methodConfPool(c)
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
	c.presCommonPoolHist[c.gameState().PresidentID] = append(presHist, commonPool)

	if c.criticalStatus() && c.commonPoolAllocation < request {
		return request
	}

	return c.commonPoolAllocation
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
		return Min(surplus/HelpCritOthersDivisor, c.taxAmount)
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
		return float64(c.gameState().ClientInfo.Resources / BaseDisasterProtectionDivisor) //initial disaster threshold guess when we start playing
	} else {
		sampleMeanMag, magnitudePrediction := GetMagnitudePrediction(c, float64(turn))

		// TODO: why are we accessing the first value here???
		baseThreshold := float64(c.resourceLevelHistory[1] / BaseResourcesToGiveDivisor)
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
		return InitialDisasterTurnGuess - float64(turn)
	} else {
		sampleMeanX, timeRemainingPrediction := GetTimeRemainingPrediction(c, float64(turn))
		turnsLeftConfidence := GetTimeRemainingConfidence(float64(turn), sampleMeanX)
		return float64(timeRemainingPrediction) * (turnsLeftConfidence / 100)
	}

}

// Determines esources we need to be above critical, pay tax and cost of living, put resources aside proportional to incoming disaster
func (c *client) agentThreshold() shared.Resources {
	criticalThreshold := c.gameConfig().MinimumResourceThreshold
	costOfLiving := c.gameConfig().CostOfLiving
	basicCosts := criticalThreshold + c.taxAmount + costOfLiving
	vulnerabilityMultiplier := 0.75 + GetIslandDVPs(c.gameState().Geography)[c.GetID()] //1 to 1.75 (1.75 being the most vulnerable)

	// Add resources based on expected/predicted disaster magnitude
	disasterMagProtection := c.calculateDisasterMagPred()

	timeRemaining := c.calculateTimeRemaining()
	disasterTimeProtectionMultiplier := 1.0
	if timeRemaining < TimeLeftIncreaseDisProtection {
		disasterTimeProtectionMultiplier = DisasterSoonProtectionMultiplier
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
		return DefaultFirstTurnContribution
	}

	switch c.MethodOfPlay() {
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
	tuneAlt := shared.Resources(AltruistFactorOfAvToGive) //what factor of the average to contribute when being altruistic, will be much higher than fair sharing
	for j := turn; j > 0; j-- {                           //we are trying to find the most recent instance of the common pool increasing and then use that value
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn] > 0 {
			if shared.Resources(c.getNumAliveClients())*tuneAlt != shared.Resources(0) {
				return ((ResourceHistory[j] - ResourceHistory[prevTurn]) / shared.Resources(c.getNumAliveClients())) * tuneAlt
			}
		}
	}
	return 0
}

func (c *client) determineFairContribution(turn uint) shared.Resources { //can make more sophisticated! Right now just contribute the average, default matters the most
	ResourceHistory := c.commonPoolHistory
	tuneAverage := shared.Resources(FairShareFactorOfAvToGive) //what factor of the average to contribute when fair sharing, default is 1 to give the average
	for j := turn; j > 0; j-- {                                //we are trying to find the most recent instance of the common pool increasing and then use that value
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn] > 0 {
			if shared.Resources(c.getNumAliveClients())*tuneAverage != shared.Resources(0) {
				return ((ResourceHistory[j] - ResourceHistory[prevTurn]) / shared.Resources(c.getNumAliveClients())) * tuneAverage
			}
		}
	}
	return 0
}

func methodConfPool(c *client) shared.Resources {
	var modeMult float64
	switch c.MethodOfPlay() {
	case 0:
		modeMult = 0.4 //when the pool is struggling, we will forage less to hav emo
	case 1:
		modeMult = 0.6 //default
	case 2:
		modeMult = 1.2 //when free riding we mostly take from the pool
	}
	return shared.Resources(modeMult)
}

func (c *client) SanctionHopeful() shared.Resources {
	return 0
}

//Checks the sanction amount aginst what we expect
func (c *client) GetSanctionPayment() shared.Resources {
	if value, ok := c.LocalVariableCache[rules.SanctionExpected]; ok {
		if c.gameState().ClientLifeStatuses[c.GetID()] != shared.Critical {
			if shared.Resources(value.Values[0]) <= c.SanctionHopeful() {
				return shared.Resources(value.Values[0])
			} else {
				// TODO: make switch case on agent mode.
				return c.SanctionHopeful()
			}
		} else {
			return 0
		}
	}
	return 0
}

func (c *client) ShareIntendedContribution() shared.IntendedContribution {
	shareWith := make([]shared.ClientID, 0)
	aliveClients := c.getAliveClients()
	for _, island := range aliveClients {
		if c.confidence("Gifts", island) > 30 {
			shareWith = append(shareWith, island)
		}
	}
	intendedContribution := shared.IntendedContribution{
		Contribution:   c.calculateContribution(),
		TeamsOfferedTo: shareWith,
	}
	return intendedContribution
}

// TODO: this is completely empty
// func (c *client) ReceiveIntendedContribution(receivedIntendedContributions shared.ReceivedIntendedContributionDict) {
// we check how much each island intends to contribute
// Compute the average amount needed for the common pool threshold
// form an opinion based on how far their contribution is from the average
// could help us determine empathy level? (ie altruist, etc)
// }

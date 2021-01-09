package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// CommonPoolUpdate Records history of common pool levels
func CommonPoolUpdate(c *client, commonPoolHistory CommonPoolHistory) {
	currentPool := c.gameState().CommonPool
	c.commonPoolHistory[c.gameState().Turn] = currentPool
}

// Records our resources level each turn
func ourResourcesHistoryUpdate(c *client, resourceLevelHistory ResourcesLevelHistory) {
	currentLevel := c.gameState().ClientInfo.Resources
	c.resourceLevelHistory[c.gameState().Turn] = currentLevel
}

// How much we ask the President for from pool
func (c *client) CommonPoolResourceRequest() shared.Resources {
	request := determineAllocation(c) * methodConfPool(c)
	var commonPool CommonPoolInfo

	if _, ok := c.presCommonPoolHist[c.gameState().PresidentID]; !ok {
		c.presCommonPoolHist[c.gameState().PresidentID] = make([]CommonPoolInfo, 0)
	} else {
		presHist := c.presCommonPoolHist[c.gameState().PresidentID]
		if presHist[len(presHist)-1].turn == c.gameState().Turn {
			commonPool = presHist[len(presHist)-1]
		}
	}
	presHist := c.presCommonPoolHist[c.gameState().PresidentID]
	commonPool.requestedToPres = request
	commonPool.turn = c.gameState().Turn
	c.presCommonPoolHist[c.gameState().PresidentID] = append(presHist, commonPool)

	return request
}

//Determines how many resources we want to obtain this round through the pool
func determineAllocation(c *client) shared.Resources {
	ourResources := c.gameState().ClientInfo.Resources
	if c.criticalStatus() {
		return c.agentThreshold()
	}
	if c.gameState().ClientInfo.Resources < c.agentThreshold() {
		return (c.agentThreshold() - ourResources)
	}
	return 0
}

//determines how many resources you actually take
func (c *client) RequestAllocation() shared.Resources {
	request := determineAllocation(c) * shared.Resources(methodConfPool(c))
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
	// TODO: same bug here
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
	//add extra when disaster is soon
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

//Determines esources we need to be above critical, pay tax and cost of living, put resources aside proportional to incoming disaster
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

//Checks if there was a disaster in the previous turn
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

//AverageCommonPoolDilemma determines how much to contribute to the common pool depending on whether other agents are altruists,fair sharers or free riders
func AverageCommonPoolDilemma(c *client) shared.Resources {
	turn := c.gameState().Turn
	if turn == 1 {
		return DefaultFirstTurnContribution
	}
	altruist := c.determineAltruist(turn) //determines altruist amount
	fairSharer := c.determineFair(turn)   //determines fair sharer amount
	method := c.MethodOfPlay()
	switch method {
	case 0:
		return shared.Resources(altruist)
	case 1:
		return shared.Resources(fairSharer)
	case 2:
		return shared.Resources(0)
	}
	return shared.Resources(fairSharer)
}

func (c *client) determineAltruist(turn uint) shared.Resources { //identical to fair sharing but a larger factor to multiple the average contribution by
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

func (c *client) determineFair(turn uint) shared.Resources { //can make more sophisticated! Right now just contribute the average, default matters the most
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
func (c *client) ReceiveIntendedContribution(receivedIntendedContributions shared.ReceivedIntendedContributionDict) {
	// we check how much each island intends to contribute
	// Compute the average amount needed for the common pool threshold
	// form an opinion based on how far their contribution is from the average
	// could help us determine empathy level? (ie altruist, etc)
}

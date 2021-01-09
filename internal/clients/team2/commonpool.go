package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

//CommonPoolUpdate Records history of common pool levels
func CommonPoolUpdate(c *client, commonPoolHistory CommonPoolHistory) {
	currentPool := c.gameState().CommonPool
	c.commonPoolHistory[c.gameState().Turn] = float64(currentPool)
}

//Records our resources level each turn
func ourResourcesHistoryUpdate(c *client, resourceLevelHistory ResourcesLevelHistory) {
	currentLevel := c.gameState().ClientInfo.Resources
	c.resourceLevelHistory[c.gameState().Turn] = currentLevel
}

// How much we ask the President for from pool
func (c *client) CommonPoolResourceRequest() shared.Resources {
	request := determineAllocation(c) * shared.Resources(methodConfPool(c))
	var commonPool CommonPoolInfo
	presHist := c.commonPoolHist[c.gameState().PresidentID]
	if len(presHist) != 0 {
		presHist[len(presHist)-1].requestedToPres = request
		presHist[len(presHist)-1].turn = c.gameState().Turn
	} else {
		commonPool = CommonPoolInfo{
			requestedToPres: request,
			turn:            c.gameState().Turn,
		}
	}
	c.commonPoolHist[c.gameState().PresidentID] = append(presHist, commonPool)
	return request
}

//Determines how many resources we want to obtain this round through the pool
func determineAllocation(c *client) shared.Resources {
	ourResources := c.gameState().ClientInfo.Resources
	if criticalStatus(c) {
		return c.agentThreshold()
	}
	if c.gameState().ClientInfo.Resources < c.agentThreshold() {
		return (c.agentThreshold() - ourResources)
	}
	return 0
}

//determines how many resources you actually take - currrently going to take however much we say (playing nicely)
// How much we request the server (we're given as much as there is in the CP)
func (c *client) RequestAllocation() shared.Resources {
	request := determineAllocation(c) * shared.Resources(methodConfPool(c))
	var commonPool CommonPoolInfo
	presHist := c.commonPoolHist[c.gameState().PresidentID]
	if len(presHist) != 0 {
		presHist[len(presHist)-1].takenFromCP = request
		presHist[len(presHist)-1].turn = c.gameState().Turn
	} else {
		commonPool = CommonPoolInfo{
			takenFromCP: request,
			turn:        c.gameState().Turn,
		}
	}
	c.commonPoolHist[c.gameState().PresidentID] = append(presHist, commonPool)
	if criticalStatus(c) && c.commonPoolAllocation < request {
		return request
	}
	return c.commonPoolAllocation
}

func (c *client) calculateContribution() shared.Resources {
	ourResources := c.gameState().ClientInfo.Resources
	Taxmin := determineTax(c)
	allocation := AverageCommonPoolDilemma(c) + Taxmin //This is our default allocation, this determines how much to give based off of previous common pool level
	if criticalStatus(c) {
		return 0 //tax evasion by necessity
	}
	if ourResources < c.agentThreshold() {
		return Taxmin
	}
	if checkOthersCrit(c) {
		return (ourResources - c.agentThreshold() - Taxmin) / 2
	}

	allocation = AverageCommonPoolDilemma(c) + Taxmin
	return allocation
}

//GetTaxContribution determines how much we put into pool
func (c *client) GetTaxContribution() shared.Resources {
	allocation := c.calculateContribution()
	c.updatePresidentTrust()
	c.confidenceRestrospect("President", c.gameState().PresidentID)
	return allocation
}

//determineTax returns how much tax we have to pay
func determineTax(c *client) shared.Resources {
	return c.taxAmount
}

//Determines esources we need to be above critical, pay tax and cost of living, put resources aside proportional to incoming disaster
func (c *client) agentThreshold() shared.Resources {
	criticaThreshold := c.gameConfig().MinimumResourceThreshold
	tax := c.taxAmount
	costOfLiving := c.gameConfig().CostOfLiving
	basicCosts := criticaThreshold + tax + costOfLiving
	vulnerability := GetIslandDVPs(c.gameState().Geography)[c.GetID()] //0.25 to 1 (1 being the most vulnerable)
	vulnerabilityMultiplier := 0.75 + vulnerability

	//add resources based on expected/predicted disaster magnitude
	turn := c.gameState().Turn
	var disasterMagProtection float64
	if c.gameConfig().DisasterConfig.CommonpoolThreshold.Valid { //resources for disaster needed known
		disasterMagProtection = float64(c.gameConfig().DisasterConfig.CommonpoolThreshold.Value) / float64(c.getNumAliveClients())
	} else { //resources for disaster needed not known
		sampleMeanM, magnitudePrediction := GetMagnitudePrediction(c, float64(turn))

		if turn == 1 {
			disasterMagProtection = float64(c.gameState().ClientInfo.Resources / BaseDisasterProtectionDivisor) //initial disaster threshold guess when we start playing
		}
		baseThreshold := float64(c.resourceLevelHistory[1] / BaseResourcesToGiveDivisor)
		if c.gameState().Season == 1 { //keep threshold from first turn
			disasterMagProtection = baseThreshold
		}
		disasterBasedAdjustment := 0
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
		disasterMagProtection = (baseThreshold)*(magnitudePrediction/sampleMeanM) + float64(disasterBasedAdjustment)
	}

	var timeRemaining float64
	//add extra when disaster is soon
	if c.gameConfig().DisasterConfig.DisasterPeriod.Valid {
		period := c.gameConfig().DisasterConfig.DisasterPeriod.Value
		timeRemaining = float64(period - (c.gameState().Turn % period))
	} else {
		sampleMeanX, timeRemainingPrediction := GetTimeRemainingPrediction(c, float64(turn))
		turnsLeftConfidence := GetTimeRemainingConfidence(float64(turn), sampleMeanX)
		timeRemaining = float64(timeRemainingPrediction) * (turnsLeftConfidence / 100)
	}
	disasterTimeProtectionMultiplier := float64(1)
	if timeRemaining < TimeLeftIncreaseDisProtection {
		disasterTimeProtectionMultiplier = disasterSoonProtectionMultiplier
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

//Determines how much to contribute to the common pool depending on whether other agents are altruists,fair sharers or free riders
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

func (c *client) determineAltruist(turn uint) float64 { //identical to fair sharing but a larger factor to multiple the average contribution by
	ResourceHistory := c.commonPoolHistory
	tuneAlt := AltruistFactorOfAvToGive //what factor of the average to contribute when being altruistic, will be much higher than fair sharing
	for j := turn; j > 0; j-- {         //we are trying to find the most recent instance of the common pool increasing and then use that value
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn] > 0 {
			return ((ResourceHistory[j] - ResourceHistory[prevTurn]) / float64(c.getNumAliveClients())) * tuneAlt
		}
	}
	return 0
}

func (c *client) determineFair(turn uint) float64 { //can make more sophisticated! Right now just contribute the average, default matters the most
	ResourceHistory := c.commonPoolHistory
	tuneAverage := FairShareFactorOfAvToGive //what factor of the average to contribute when fair sharing, default is 1 to give the average
	for j := turn; j > 0; j-- {              //we are trying to find the most recent instance of the common pool increasing and then use that value
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn] > 0 {
			return ((ResourceHistory[j] - ResourceHistory[prevTurn]) / float64(c.getNumAliveClients())) * tuneAverage
		}
	}
	return 0
}

func methodConfPool(c *client) float64 {
	var modeMult float64
	switch c.MethodOfPlay() {
	case 0:
		modeMult = 0.4 //when the pool is struggling, we will forage less to hav emo
	case 1:
		modeMult = 0.6 //default
	case 2:
		modeMult = 1.2 //when free riding we mostly take from the pool
	}
	return modeMult
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
func (c *client) ReceiveIntendedContribution(receivedIntendedContributions shared.ReceivedIntendedContributionDict) {

}

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

//determines how much to request from the pool
// How much we ask the President for
func (c *client) CommonPoolResourceRequest() shared.Resources {
	return determineAllocation(c) * shared.Resources(methodConfPool(c))
}

//TODO:Update to consider IIGO allocated amount, opinion on pres -> whether we take the allocated amount
func determineAllocation(c *client) shared.Resources {
	ourResources := c.gameState().ClientInfo.Resources
	if criticalStatus(c) {
		return c.agentThreshold() //not sure about this amount
	}
	if c.gameState().ClientInfo.Resources < c.agentThreshold() {
		return (c.agentThreshold() - ourResources)
	}
	//TODO: maybe separate standard gameplay when no-one is critical vs when others are critical
	return 0
}

//determines how many resources you actually take - currrently going to take however much we say (playing nicely)
// How much we request the server (we're given as much as there is in the CP)
func (c *client) RequestAllocation() shared.Resources {
	return determineAllocation(c) * shared.Resources(methodConfPool(c))
}

//GetTaxContribution determines how much we put into pool
func (c *client) GetTaxContribution() shared.Resources {
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

//determineTax returns how much tax we have to pay
func determineTax(c *client) shared.Resources {
	return c.taxAmount //TODO: not sure if this is correct tax amount to use
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
			disasterMagProtection = float64(c.gameState().ClientInfo.Resources / 4) //initial disaster threshold guess when we start playing
			//TODO: tune
		}
		baseThreshold := float64(c.resourceLevelHistory[1] / 4) //TODO: tune
		if c.gameState().Season == 1 {                          //keep threshold from first turn
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
	if timeRemaining < 3 {
		disasterTimeProtectionMultiplier = float64(1.2) //TODO:tune
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

//this function determines how much to contribute to the common pool depending on whether other agents are altruists,fair sharers etc
//it only needs the current resource level and the current turn as inputs
//the output will be an integer which is a recommendation on how much to add to the pool, with this recommendation there will be a weighting of how important it is we contribute that exact amount
//this will be a part of other decision making functions which will have their own weights

//tunable parameters:
//how much to give the pool on our first turn: default_strat
//after how many rounds of struggling pool to intervene and become altruist: intervene
//the number of turns at the beginning we cannot free ride: no_freeride
//the factor in which the common pool increases by to decide if we should free ride: freeride
//the factor which we multiply the fair_sharer average by: tune_average
//the factor which we multiply the altruist value by: tune_alt

//Extra Functionality TODO: The bigger the drop in the common pool, the more we give in the altruistic mode
//TODO: inside determine_fair and determine_altruist make the 6 how many alive agents there are
func AverageCommonPoolDilemma(c *client) shared.Resources {
	ResourceHistory := c.commonPoolHistory
	turn := c.gameState().Turn
	var defaultStrat float64 = 50 //this parameter will determine how much we contribute on the first turn when there is no data to make a decision

	var fairSharer float64 //this is how much we contribute when we are a fair sharer and altruist
	var altruist float64

	//var decreasing_pool float64 //records for how many turns the common pool is decreasing
	var noFreeride float64 = 3 //how many turns at the beginning we cannot free ride for
	var freeride float64 = 5   //what factor the common pool must increase by for us to considered free riding
	var altfactor float64 = 5  //what factor the common pool must drop by for us to consider altruist

	if turn == 1 { //if there is no historical data then use default strategy
		return shared.Resources(defaultStrat)
	}

	altruist = c.determineAltruist(turn) //determines altruist amount
	fairSharer = c.determineFair(turn)   //determines fair sharer amount

	prevTurn := turn - 1
	prevTurn2 := turn - 2
	if ResourceHistory[prevTurn] > (ResourceHistory[turn] * altfactor) { //decreasing common pool means consider altruist
		if ResourceHistory[prevTurn2] > (ResourceHistory[prevTurn] * altfactor) {
			return shared.Resources(altruist)
		}
	}

	if float64(turn) > noFreeride { //we will not allow ourselves to use free riding at the start of the game
		if (ResourceHistory[prevTurn] * freeride) < ResourceHistory[turn] {
			if (ResourceHistory[prevTurn2] * freeride) < ResourceHistory[prevTurn] { //two large jumps then we free ride
				return 0
			}
		}
	}
	return shared.Resources(fairSharer) //by default we contribute a fair share
}

func (c *client) determineAltruist(turn uint) float64 { //identical to fair sharing but a larger factor to multiple the average contribution by
	ResourceHistory := c.commonPoolHistory
	var tuneAlt float64 = 2     //what factor of the average to contribute when being altruistic, will be much higher than fair sharing
	for j := turn; j > 0; j-- { //we are trying to find the most recent instance of the common pool increasing and then use that value
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn] > 0 {
			return ((ResourceHistory[j] - ResourceHistory[prevTurn]) / float64(c.getNumAliveClients())) * tuneAlt
		}
	}
	return 0
}

func (c *client) determineFair(turn uint) float64 { //can make more sophisticated! Right now just contribute the average, default matters the most
	ResourceHistory := c.commonPoolHistory
	var tuneAverage float64 = 1 //what factor of the average to contribute when fair sharing, default is 1 to give the average
	for j := turn; j > 0; j-- { //we are trying to find the most recent instance of the common pool increasing and then use that value
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn] > 0 {
			return ((ResourceHistory[j] - ResourceHistory[prevTurn]) / float64(c.getNumAliveClients())) * tuneAverage //make 6 variable for no of agents
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

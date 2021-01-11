package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"math"
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
func (c *client) CommonPoolResourceRequest() shared.Resources {
	return determineAllocation(c) * 0.6
}

//determines how many resources you actually take - currrently going to take however much we say (playing nicely)
func determineAllocation(c *client) shared.Resources {
	ourResources := c.gameState().ClientInfo.Resources
	if criticalStatus(c) {
		return c.determineThreshold() //not sure about this amount
	}
	if determineTax(c)+c.determineThreshold() > ourResources {
		return (determineTax(c) - ourResources)
	}
	if c.gameState().ClientInfo.Resources < c.internalThreshold() {
		return (c.internalThreshold() - ourResources)
	}
	//TODO: maybe separate standard gameplay when no-one is critical vs when others are critical
	return 0
}

func (c *client) RequestAllocation() shared.Resources {
	return shared.Resources(math.Max(float64(determineAllocation(c)*0.6), 0))
}

//GetTaxContribution determines how much we put into pool
func (c *client) GetTaxContribution() shared.Resources {
	ourResources := c.gameState().ClientInfo.Resources
	Taxmin := determineTax(c)
	allocation := AverageCommonPoolDilemma(c) + Taxmin //This is our default allocation, this determines how much to give based off of previous common pool level
	if criticalStatus(c) {
		return 0 //tax evasion
	}
	if determineTax(c)+c.determineThreshold() > ourResources {
		return 0 //tax evasion
	}
	if ourResources < c.internalThreshold() {
		return Taxmin
	}
	if checkOthersCrit(c) {
		return (ourResources - c.internalThreshold() - Taxmin) / 2
	}

	allocation = AverageCommonPoolDilemma(c) + Taxmin
	return allocation
}

//determineTax returns how much tax we have to pay
func determineTax(c *client) shared.Resources {
	return shared.Resources(c.BaseClient.LocalVariableCache[rules.ExpectedTaxContribution].Values[0]) //TODO: not sure if this is correct tax amount to use
}

//internalThreshold determines our internal threshold for survival, allocationrec is the output of the function AverageCommonPool which determines which role we will be
func (c *client) internalThreshold() shared.Resources {
	gameThreshold := c.determineThreshold()
	allocationrec := AverageCommonPoolDilemma(c)
	ourVulnerability := GetIslandDVPs(c.gameState().Geography)[c.GetID()] //TODO: get value from init function
	//turnsLeftUntilDisaster := 3           //TODO: get this value when known
	totalTurns := float64(c.gameState().Turn)
	sampleMeanX, timeRemainingPrediction := GetTimeRemainingPrediction(c, totalTurns)
	turnsLeftConfidence := GetTimeRemainingConfidence(totalTurns, sampleMeanX)
	//TODO: Update these to be functions more specific to our island rather than general mag
	sampleMeanM, magnitudePrediction := GetMagnitudePrediction(c, totalTurns)
	confidenceMagnitude := GetMagnitudeConfidence(totalTurns, sampleMeanM)

	if magnitudePrediction > sampleMeanM { //larger mag than average expected
		return (gameThreshold + allocationrec) * shared.Resources((confidenceMagnitude/10)*(1+ourVulnerability)) //tune
	}
	if timeRemainingPrediction < 3 { //tune
		return (gameThreshold + allocationrec) * shared.Resources((turnsLeftConfidence/10)*(1+ourVulnerability)) //tune
	}
	return gameThreshold + allocationrec
}

//Checks if there was a disaster in the previous turn
func checkForDisaster(c *client) bool {
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

//Finds the game threshold if this information is available or works it out based on history of turns
func (c *client) determineThreshold() shared.Resources {
	costOfLiving := c.gameConfig().CostOfLiving
	if c.gameConfig().DisasterConfig.CommonpoolThreshold.Valid {
		var threshold = (c.gameConfig().DisasterConfig.CommonpoolThreshold.Value) / 6
		return threshold*1.1 + costOfLiving //*1.1 as threshold may not always be enough
	}
	ourResources := c.gameState().ClientInfo.Resources
	turn := c.gameState().Turn
	season := c.gameState().Season
	var disasterBasedAdjustment float64
	ourVulnerability := GetIslandDVPs(c.gameState().Geography)[c.GetID()] //TODO: get value from init function
	//turnsLeftUntilDisaster := 3           //TODO: get this value when available in client config
	totalTurns := float64(c.gameState().Turn)
	sampleMeanM, magnitudePrediction := GetMagnitudePrediction(c, totalTurns)

	if turn == 1 {
		return ourResources / 4 //TODO: tune initial threshold guess when we start playing
	}
	baseThreshold := c.resourceLevelHistory[1] / 4
	if season == 1 || sampleMeanM == 0.0 { //keep threshold from first turn
		return baseThreshold
	}
	if checkForDisaster(c) {
		if c.resourceLevelHistory[turn] >= c.resourceLevelHistory[turn-1] { //no resources taken by disaster
			if disasterBasedAdjustment > 5 {
				disasterBasedAdjustment -= 5
			}
		}
		//disaster took our resources
		disasterBasedAdjustment += 5
	}
	//change factor by if next mag > or < prev mag
	return shared.Resources(float64(baseThreshold)*(magnitudePrediction/sampleMeanM)*(1+ourVulnerability)+disasterBasedAdjustment) + costOfLiving
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
	var default_strat float64 = 50 //this parameter will determine how much we contribute on the first turn when there is no data to make a decision

	var fair_sharer float64 //this is how much we contribute when we are a fair sharer and altruist
	var altruist float64

	//var decreasing_pool float64 //records for how many turns the common pool is decreasing
	var no_freeride float64 = 3 //how many turns at the beginning we cannot free ride for
	var freeride float64 = 5    //what factor the common pool must increase by for us to considered free riding

	if turn == 1 { //if there is no historical data then use default strategy
		return shared.Resources(default_strat)
	}

	altruist = c.determine_altruist(turn) //determines altruist amount
	fair_sharer = c.determine_fair(turn)  //determines fair sharer amount

	prevTurn := turn - 1
	prevTurn2 := turn - 2
	if ResourceHistory[prevTurn] > ResourceHistory[turn] { //decreasing common pool means consider altruist
		if ResourceHistory[prevTurn2] > ResourceHistory[prevTurn] {
			return shared.Resources(altruist)
		}
	}

	if float64(turn) > no_freeride { //we will not allow ourselves to use free riding at the start of the game
		if ResourceHistory[prevTurn] < (ResourceHistory[turn] * freeride) {
			if ResourceHistory[prevTurn2] < (ResourceHistory[prevTurn] * freeride) { //two large jumps then we free ride
				return 0
			}
		}
	}
	return shared.Resources(fair_sharer) //by default we contribute a fair share
}

func (c *client) determine_altruist(turn uint) float64 { //identical to fair sharing but a larger factor to multiple the average contribution by
	ResourceHistory := c.commonPoolHistory
	var tune_alt float64 = 2    //what factor of the average to contribute when being altruistic, will be much higher than fair sharing
	for j := turn; j > 0; j-- { //we are trying to find the most recent instance of the common pool increasing and then use that value
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn] > 0 {
			return ((ResourceHistory[j] - ResourceHistory[prevTurn]) / 6) * tune_alt
		}
	}
	return 0
}

func (c *client) determine_fair(turn uint) float64 { //can make more sophisticated! Right now just contribute the average, default matters the most
	ResourceHistory := c.commonPoolHistory
	var tune_average float64 = 1 //what factor of the average to contribute when fair sharing, default is 1 to give the average
	for j := turn; j > 0; j-- {  //we are trying to find the most recent instance of the common pool increasing and then use that value
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn] > 0 {
			return ((ResourceHistory[j] - ResourceHistory[prevTurn]) / 6) * tune_average //make 6 variable for no of agents
		}
	}
	return 0
}

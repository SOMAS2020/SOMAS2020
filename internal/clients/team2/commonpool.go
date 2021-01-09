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
	totalAlloc := determineAllocation(c)

	request := determineAllocation(c) * shared.Resources(methodConfPool(c))

	//HISTORY STUFF
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

	//RULES STUFF
	c.LocalVariableCache[rules.IslandAllocation] = rules.VariableValuePair{
		VariableName: rules.IslandAllocation,
		Values:       []float64{float64(c.gameState().CommonPool)},
	}
	allocationPair, success := c.GetRecommendation(rules.IslandAllocation)

	// when not sure we try to get all the resources from commonpool
	if !success {
		return c.gameState().CommonPool
	}

	//in critical we want, the resources to survive and more
	if c.gameState().ClientLifeStatuses[c.GetID()] == shared.Critical {
		return totalAlloc * shared.Resources(2)
	}

	allocation := allocationPair.Values[0]

	return determineAllocation(c)*shared.Resources(methodConfPool(c)) + shared.Resources(allocation)
}

func (c *client) calculateContribution() shared.Resources {
	ourResources := c.gameState().ClientInfo.Resources
	Taxmin := determineTax(c)
	contribution, success := c.GetRecommendation(rules.IslandTaxContribution)
	allocation := AverageCommonPoolDilemma(c) + Taxmin //This is our default allocation, this determines how much to give based off of previous common pool level
	if !success {
		return 0 //tax not determined correctly by our check
	}
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

	//if asked tax less than our calculation
	if allocation < shared.Resources(contribution.Values[0]) {
		return allocation
	}

	return shared.Resources(contribution.Values[0])

}

//GetTaxContribution determines how much we put into pool
func (c *client) GetTaxContribution() shared.Resources {
	allocation := c.calculateContribution()
	c.updatePresidentTrust()
	c.confidenceRestrospect("President", c.gameState().PresidentID)
	return allocation
}

//determineTax returns how much tax we have to pay
//TODO: not sure if this is correct tax amount to use
func determineTax(c *client) shared.Resources {
	return shared.Resources(c.BaseClient.LocalVariableCache[rules.ExpectedTaxContribution].Values[0])
}

//Determines esources we need to be above critical, pay tax and cost of living, put resources aside proportional to incoming disaster
func (c *client) agentThreshold() shared.Resources {
	criticaThreshold := c.gameConfig().MinimumResourceThreshold
	tax := shared.Resources(0)
	costOfLiving := c.gameConfig().CostOfLiving
	basicCosts := criticaThreshold + tax + costOfLiving
	vulnerability := GetIslandDVPs(c.gameState().Geography)[c.GetID()] //0.25 to 1 (1 being the most vulnerable)
	vulnerabilityMultiplier := 0.75 + vulnerability                    //1 to 1.75
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
			}
			// TODO: make switch case on agent mode.
			return c.SanctionHopeful()

		}
		return 0
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

//***********HELPER FUNCTIONS************************
// DELETE FROM FINAL CODE BUT NOT YET

// func (c *client) UpdateCache(oldCache map[rules.VariableFieldName]rules.VariableValuePair, newValues map[rules.VariableFieldName]rules.VariableValuePair) map[rules.VariableFieldName]rules.VariableValuePair {
// 	for key, val := range newValues {
// 		oldCache[key] = val
// 	}
// 	return oldCache
// }

// func (c *client) dynamicAssistedResult(variablesChanged map[rules.VariableFieldName]rules.VariableValuePair) (newVals map[rules.VariableFieldName]rules.VariableValuePair) {
// 	if c.LocalVariableCache != nil {
// 		c.LocalVariableCache = c.UpdateCache(c.LocalVariableCache, variablesChanged)
// 		// For testing using available rules
// 		return c.BaseClient.GetRecommendation(rules.IslandAllocation)
// 	}
// 	return variablesChanged
// }

// func (c *client) GetVoteForRule(matrix rules.RuleMatrix) bool {

// 	newRulesInPlay := make(map[string]rules.RuleMatrix)

// 	for key, value := range rules.RulesInPlay {
// 		if key == matrix.RuleName {
// 			newRulesInPlay[key] = matrix
// 		} else {
// 			newRulesInPlay[key] = value
// 		}
// 	}

// 	if _, ok := rules.RulesInPlay[matrix.RuleName]; ok {
// 		delete(newRulesInPlay, matrix.RuleName)
// 	} else {
// 		newRulesInPlay[matrix.RuleName] = rules.AvailableRules[matrix.RuleName]
// 	}

// 	// TODO: define postion -> list of variables and values associated with the rule (obtained from IIGO communications)

// 	// distancetoRulesInPlay = CalculateDistanceFromRuleSpace(rules.RulesInPlay, position)
// 	// distancetoNewRulesInPlay = CalculateDistanceFromRuleSpace(newRulesInPlay, position)

// 	// if distancetoRulesInPlay < distancetoNewRulesInPlay {
// 	//  return false
// 	// } else {
// 	//  return true
// 	// }

// 	return true
// }

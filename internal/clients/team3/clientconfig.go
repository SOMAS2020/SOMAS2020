package team3

func getislandParams() islandParams {
	return islandParams{
		equity:                  1,    // 0-1 // do we need this or can it be replaced?
		complianceLevel:         0.5,  // 0-1 //this seems good
		resourcesSkew:           2,    // >1 //same as equity
		saveCriticalIsland:      true, //seems like this will always be true
		selfishness:             1,    //0-1 //higher is better for us
		riskFactor:              0.8,  // 0-1 // increasing this has mixed results
		friendliness:            1,    // 0-1 // agent performs better when this matches selfishness
		adv:                     nil,  // keep it off pls
		giftInflationPercentage: 1,    // 0-1 //this doesn't have a noticeable effect
		sensitivity:             0.5,  // 0-1
		controlLoop:             true,
	}
}

// TODOS:
// 1. finalize the list of parameters.
// 2. pick 3 sets of values that correspond to different behaviours & submit to naim. play mainly with selfishness, risk, friendliness.
// 3. Improve foraging if we can
// 4. WE ARE NOT VOTING WELL! Add a sensitivity parameter and use lists changed by evalPerformance to vote

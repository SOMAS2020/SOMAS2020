package team3

import "github.com/SOMAS2020/SOMAS2020/internal/clients/team3/adv"

func getislandParams() islandParams {
	return ardernLeeMacron()
}

func evilPeople() islandParams {
	return islandParams{
		equity:                  1,     // 0-1 // do we need this or can it be replaced?
		complianceLevel:         0.5,   // 0-1 //this seems good
		resourcesSkew:           2,     // >1 //same as equity
		saveCriticalIsland:      false, //seems like this will always be true
		selfishness:             1.0,   //0-1 //higher is better for us
		riskFactor:              0.8,   // 0-1 // increasing this has mixed results
		friendliness:            1,
		advType:                 adv.MaliceAdv, // 0-1 // agent performs better when this matches selfishness
		giftInflationPercentage: 1,             // 0-1 //this doesn't have a noticeable effect
		sensitivity:             0.5,           // 0-1
		controlLoop:             true,
	}
}

func gandhi() islandParams {
	return islandParams{
		equity:                  1.2,  // 0-1 // do we need this or can it be replaced?
		complianceLevel:         1,    // 0-1 //this seems good
		resourcesSkew:           1,    // >1 //same as equity
		saveCriticalIsland:      true, //seems like this will always be true
		selfishness:             0.5,  //0-1 //higher is better for us
		riskFactor:              0.05, // 0-1 // increasing this has mixed results
		friendliness:            1,
		advType:                 adv.NoAdv, // 0-1 // agent performs better when this matches selfishness
		giftInflationPercentage: 1,         // 0-1 //this doesn't have a noticeable effect
		sensitivity:             0.5,       // 0-1
		controlLoop:             false,
	}
}

func ardernLeeMacron() islandParams {
	return islandParams{
		equity:                  1.1,   // 0-1 // do we need this or can it be replaced?
		complianceLevel:         0.75,  // 0-1 //this seems good
		resourcesSkew:           1.5,   // >1 //same as equity
		saveCriticalIsland:      true,  //seems like this will always be true
		selfishness:             0.75,  //0-1 //higher is better for us
		riskFactor:              0.125, // 0-1 // increasing this has mixed results
		friendliness:            1,
		advType:                 adv.NoAdv, // 0-1 // agent performs better when this matches selfishness
		giftInflationPercentage: 1,         // 0-1 //this doesn't have a noticeable effect
		sensitivity:             0.5,       // 0-1
		controlLoop:             true,
	}
}

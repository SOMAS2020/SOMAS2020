package team3

func getislandParams() islandParams {
	return islandParams{
		equity:                      0.1, // 0-1
		complianceLevel:             0.1, // 0-1
		resourcesSkew:               1.3, // >1
		saveCriticalIsland:          true,
		escapeCritcaIsland:          true,
		selfishness:                 0.3,  //0-1
		disasterPredictionWeighting: 0.4,  // 0-1
		recidivism:                  0.3,  // real number
		riskFactor:                  0.2,  // 0-1
		friendliness:                0.3,  // 0-1
		anger:                       0.2,  // 0-1
		aggression:                  0.4,  // 0-1
		localPoolThreshold:          100,  // % of common pool
		giftInflationPercentage:     0.1,  // 0-1
		trustConstantAdjustor:       1,    // arbitrary
		trustParameter:              0.5,  // 0-1
		NoRequestGiftParam:          0.25, // 0-1
	}
}

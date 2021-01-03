package team3

func getislandParams() islandParams {
	return islandParams{
		giftingThreshold:            45,  // real number
		equity:                      0.1, // 0-1
		complianceLevel:             0.1, // 0-1
		resourcesSkew:               1.3, // >1
		saveCriticalIsland:          true,
		escapeCritcaIsland:          true,
		selfishness:                 0.3, //0-1
		minimumRequest:              30,  // real number
		disasterPredictionWeighting: 0.4, // 0-1
		recidivism:                  0.3, // real number
		riskFactor:                  0.2, // 0-1
		friendliness:                0.3, // 0-1
		anger:                       0.2, // 0-1
		aggression:                  0.4, // 0-1
	}
}

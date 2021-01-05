package team5

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// define config structure here
type clientConfig struct {

	// Initial non planned foraging
	InitialForageTurns      uint
	MinimumForagePercentage float64
	NormalForagePercentage  float64
	JBForagePercentage      float64

	// Normal foraging
	NormalRandomIncrease float64
	MaxForagePercentage  float64
	SkipForage           uint // Skip for X turns if no positive RoI

	// If resources go above this limit we are balling with money
	jbThreshold shared.Resources
	// Middle class:  Middle < Jeff bezos
	middleThreshold shared.Resources
	// Poor: Imperial student < Middle
	imperialThreshold shared.Resources

	// How much to request when we are dying
	dyingGiftRequestAmount shared.Resources
	// How much to request when we are at Imperial
	imperialGiftRequestAmount shared.Resources
	// How much to request when we are dying
	middleGiftRequestAmount shared.Resources

	// Disasters and IIFO
	forecastTrustTreshold opinionScore // min opinion score of another team to consider their forecast in creating ours
	maxForecastVariance   float64      // maximum tolerable variance in historical forecast values
}

// set param values here. In order to add a new value, you need to add a definition in struct above.
func getClientConfig() clientConfig {
	return clientConfig{
		//Variables for Intial forage
		InitialForageTurns:      5,
		MinimumForagePercentage: 0.01,
		NormalForagePercentage:  0.05,
		JBForagePercentage:      0.10, // % of our resources when JB is Normal< X < JB

		// Variables for Normal forage
		SkipForage:           1,
		NormalRandomIncrease: 0.05,
		MaxForagePercentage:  0.20,

		// Threshold for wealth
		jbThreshold:       95,
		middleThreshold:   60.0,
		imperialThreshold: 30.0, // surely should be - 100e6? (your right we are so far indebt)
		//  Dying threshold is 0 < Dying < Imperial

		// Gifts Config
		dyingGiftRequestAmount:    10,
		imperialGiftRequestAmount: 5,
		middleGiftRequestAmount:   2,

		// Disasters and IIFO
		forecastTrustTreshold: 0.0, // neutral opinion
		maxForecastVariance:   100.0,
	}
}

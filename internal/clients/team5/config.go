package team5

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// define config structure here
type clientConfig struct {

	// Initial non planned foraging
	InitialForageTurns      uint
	MinimumForagePercentage float64
	NormalForagePercentage  float64
	JBForagePercentage      float64

	// Deciding Foraging Method
	RandomChanceToFish           float64
	RandomChanceToHunt           float64
	IncreasePerHunterLastTurn    float64 // % increase per Hunter last turn
	IncreasePerFisherMenLastTurn float64 // % increase per Fisherman last turn
	DeerTurnsToLookBack          uint    // Find hunters in LB turns (not includding previous)
	DecreasePerHunterInLookBack  float64 // % decrease per hunter in past LB turns

	// NormalForage
	SkipForage           uint // Skip for X turns if no positive RoI
	NormalRandomIncrease float64
	MaxForagePercentage  float64

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

	//How much to boost gifts based on opinion
	giftBoosting shared.Resources
	giftReduct   shared.Resources

	// Disasters and IIFO
	forecastTrustTreshold opinionScore // min opinion score of another team to consider their forecast in creating ours
	maxForecastVariance   float64      // maximum tolerable variance in historical forecast values
	periodConfidenceBias  float64      // value between 0 and 1. Greater than 0.5 => weight period confidence more than mag. and spatial in overall conf
}

// set param values here. In order to add a new value, you need to add a definition in struct above.
func getClientConfig() clientConfig {
	return clientConfig{
		//Variables for initial forage
		InitialForageTurns:      5,
		MinimumForagePercentage: 0.01,
		NormalForagePercentage:  0.05,
		JBForagePercentage:      0.10, // % of our resources when JB is Normal< X < JB

		// Deciding foraging type
		RandomChanceToFish:           0.2, // Chacne to switch to Hunting/Fishing
		RandomChanceToHunt:           0.2,
		IncreasePerHunterLastTurn:    0.05, // % increase for each Hunter
		IncreasePerFisherMenLastTurn: 0.05, // % incrase for each Fisher
		DeerTurnsToLookBack:          5,    // Number of turns to look back at for deer (not including last)
		DecreasePerHunterInLookBack:  0.03, // lower for less emphasis on looking at previous turn hunters (MAX 0.07 will skip if 6 hunters in 5 turns)

		// Normal Forage
		SkipForage:           1,
		NormalRandomIncrease: 0.05,
		MaxForagePercentage:  0.20,

		// Threshold for wealth as multiplier
		jbThreshold:       2.0,
		middleThreshold:   1.0,
		imperialThreshold: 0.6, // surely should be - 100e6? (your right we are so far indebt)
		//  Dying threshold is 0 < Dying < Imperial

		// Gifts Config
		dyingGiftRequestAmount:    10,
		imperialGiftRequestAmount: 5,
		middleGiftRequestAmount:   2,
		giftBoosting:              1.4,
		giftReduct:                0.5,

		// Disasters and IIFO
		forecastTrustTreshold: 0.0, // neutral opinion
		maxForecastVariance:   100.0,
	}
}

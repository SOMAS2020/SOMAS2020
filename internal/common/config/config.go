// Package config contains types for the configuration of the game.
// DO NOT depend on other packages outside this folder!
// Add default values etc. in <root>/params.go
package config

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// Config is the type for the game configuration.
type Config struct {
	// MaxSeasons is the maximum number of 1-indexed seasons to run the game.
	MaxSeasons uint

	// MaxTurns is the maximum numbers of 1-indexed turns to run the game.
	MaxTurns uint

	// InitialResources is the default number of resources at the start of the game.
	InitialResources shared.Resources

	// InitialCommonPool is the default number of resources in the common pool at the start of the game.
	InitialCommonPool shared.Resources

	// CostOfLiving is subtracted from an islands pool before
	// the next turn. This is the simulation-level equivalent to using resources to stay
	// alive (e.g. food consumed). These resources are permanently consumed and do
	// NOT go into the common pool. Note: this is NOT the same as the tax.
	CostOfLiving shared.Resources

	// MinimumResourceThreshold is the minimum resources required for an island to not be
	// in Critical state.
	MinimumResourceThreshold shared.Resources

	// MaxCriticalConsecutiveTurns is the maximum consecutive turns an island can be in the critical state.
	MaxCriticalConsecutiveTurns uint

	// Wrapped foraging config
	ForagingConfig ForagingConfig

	// Wrapped disaster config
	DisasterConfig DisasterConfig
}

// DeerHuntConfig is a subset of foraging config
type DeerHuntConfig struct {
	// Deer Hunting
	MaxDeerPerHunt        uint                                // Max possible number of deer on a single hunt (regardless of number of participants)
	IncrementalInputDecay float64                             // Determines decay of incremental input cost of hunting more deer
	BernoulliProb         float64                             // `p` param in D variable (see README). Controls prob of catching a deer or not
	ExponentialRate       float64                             // `lambda` param in W variable (see README). Controls distribution of deer sizes.
	InputScaler           float64                             // scalar value that adjusts input resources to be in a range that is commensurate with cost of living, salaries etc.
	OutputScaler          float64                             // scalar value that adjusts returns to be in a range that is commensurate with cost of living, salaries etc.
	DistributionStrategy  shared.ResourceDistributionStrategy // basis on which returns are split amongst hunters
	ThetaCritical         float64                             // Bernoulli prob of catching deer when running population = max deer per hunt
	ThetaMax              float64                             // Bernoulli prob of catching deer when running population = carrying capacity (max population)

	// Deer Population
	MaxDeerPopulation     uint    // Max possible deer population.
	DeerGrowthCoefficient float64 // Scaling parameter used in the population model. Larger coeff => deer pop. regenerates faster
}

// FishingConfig is a subset of foraging config
type FishingConfig struct {
	// Fishing
	MaxFishPerHunt        uint                                // Max possible number of fish on a single fishing expedition
	IncrementalInputDecay float64                             // Determines decay of incremental input cost of catching additional fish
	Mean                  float64                             // mean of normally distributed fish size
	Variance              float64                             // variance of normally distributed fish size
	InputScaler           float64                             // scalar value that adjusts input resources to be in a range that is commensurate with cost of living, salaries etc.
	OutputScaler          float64                             // scalar value that adjusts returns to be in a range that is commensurate with cost of living, salaries etc.
	DistributionStrategy  shared.ResourceDistributionStrategy // basis on which returns are split amongst fishermen
}

// DisasterConfig captures disaster-specific config
type DisasterConfig struct {
	XMin, XMax, YMin, YMax      shared.Coordinate     // [min, max] x,y bounds of archipelago (bounds for possible disaster)
	GlobalProb                  float64               // Bernoulli 'p' param. Chance of a disaster occurring
	SpatialPDFType              shared.SpatialPDFType // Set x,y prob. distribution of the disaster's epicentre (more post MVP)
	MagnitudeLambda             float64               // Exponential rate param for disaster magnitude
	MagnitudeResourceMultiplier float64               // multiplier to map disaster magnitude to CP resource deductions
	CommonpoolThreshold         shared.Resources      // threshold for min CP resources for disaster mitigation
}

// ForagingConfig captures foraging-specific config
type ForagingConfig struct {
	DeerHuntConfig DeerHuntConfig
	FishingConfig  FishingConfig
}

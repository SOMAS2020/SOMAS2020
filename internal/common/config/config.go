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

// ForagingConfig captures foraging-specific config
type ForagingConfig struct {
	// Deer Hunting
	MaxDeerPerHunt        uint    // Max possible number of deer on a single hunt (regardless of number of participants)
	IncrementalInputDecay float64 // Determines decay of incremental input cost of hunting more deer
	BernoulliProb         float64 // `p` param in D variable (see README). Controls prob of catching a deer or not
	ExponentialRate       float64 // `lambda` param in W variable (see README). Controls distribution of deer sizes.

	// Deer Population
	MaxDeerPopulation     uint    // Max possible deer population. Reserved for post-MVP functionality
	DeerGrowthCoefficient float64 // Scaling parameter used in the population model. Larger coeff => deer pop. regenerates faster

	// TODO: add other pertinent params here (for fishing etc)
}

// DisasterConfig captures disaster-specific config
type DisasterConfig struct {
	XMin, XMax, YMin, YMax shared.Coordinate     // [min, max] x,y bounds of archipelago (bounds for possible disaster)
	GlobalProb             float64               // Bernoulli 'p' param. Chance of a disaster occurring
	SpatialPDFType         shared.SpatialPDFType // Set x,y prob. distribution of the disaster's epicentre (more post MVP)
	MagnitudeLambda        float64               // Exponential rate param for disaster magnitude
}

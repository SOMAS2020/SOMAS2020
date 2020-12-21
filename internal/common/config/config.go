// Package config contains configurations of the game.
// DO NOT depend on other packages outside this folder!
package config

// Config is the type for the game configuration.
type Config struct {
	// MaxSeasons is the maximum number of 1-indexed seasons to run the game.
	MaxSeasons uint

	// MaxTurns is the maximum numbers of 1-indexed turns to run the game.
	MaxTurns uint

	// InitialResources is the default number of resources at the start of the game
	InitialResources int

	// CostOfLiving is subtracted from an islands pool before
	// the next term. This is the simulation-level equivalent to using resources to stay
	// alive (e.g. food consumed). These resources are permanently consumed and do
	// NOT go into the common pool. Note: this is NOT the same as the tax
	CostOfLiving int

	// MinimumResourceThreshold is the minimum resources required for an island to not be
	// in Critical state.
	MinimumResourceThreshold int

	// MaxCriticalConsecutiveTurns is the maximum consecutive turns an island can be in the critical state.
	MaxCriticalConsecutiveTurns uint

	// Wrapped foraging config
	ForagingConfig ForagingConfig
}

// ForagingConfig captures foraging-specific config
type ForagingConfig struct {
	// Deer Hunting
	MaxDeerPerHunt        uint    // Maximimum possible number of deer on a single hunt (regardless of number of participants)
	IncrementalInputDecay float64 // Determines decay of incremental input cost of hunting more deer
	BernoulliProb         float64 // `p` param in D variable (see README). Controls prob of catching a deer or not
	ExponentialRate       float64 // `lambda` param in W variable (see README). Controls distribution of deer sizes.

	// Deer Population
	MaxDeerPopulation     uint    // Maximimum possible deer population. Reserved for post-MVP functionality
	DeerGrowthCoefficient float64 // Scaling parameter used in the population model. Larger coeff => deer pop. regenerates faster

	// TOOD: add other pertinent params here (for fishing etc)
}

// GameConfig returns the configuration of the game.
// (Made a function so it cannot be altered mid-game).
func GameConfig() Config {
	foragingConf := ForagingConfig{
		MaxDeerPerHunt:        4,
		IncrementalInputDecay: 0.8,
		BernoulliProb:         0.95,
		ExponentialRate:       1,

		MaxDeerPopulation:     12,
		DeerGrowthCoefficient: 0.4,
	}

	return Config{
		MaxSeasons:                  100,
		MaxTurns:                    2,
		InitialResources:            100,
		CostOfLiving:                10,
		MinimumResourceThreshold:    5,
		MaxCriticalConsecutiveTurns: 3,
		ForagingConfig:              foragingConf,
	}
}

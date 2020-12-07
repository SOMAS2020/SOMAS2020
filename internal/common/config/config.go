// Package config contains configurations of the game.
// DO NOT depend on other packages outside this folder!
package config

// InitialSeason is the 1-indexed season the game starts with.
const InitialSeason = 1 // default: 1

// InitialTurn is the 1-indexed turn the game starts with.
const InitialTurn = 1 // default: 1

// MaxSeasons is the maximum number of 1-indexed seasons to run the game.
const MaxSeasons = 100

// MaxTurns is the maximum numbers of 1-indexed turns to run the game.
const MaxTurns = 2

// InitialResources is the default number of resources at the start of the game
const InitialResources = 100

// CostOfLiving is subtracted from an islands pool before
// the next term. This is the simulation-level equivalent to using resources to stay
// alive (e.g. food consumed). These resources are permanently consumed and do
// NOT go into the common pool. Note: this is NOT the same as the tax
const CostOfLiving = 10

// MinimumResourceThreshold is the minimum resources required for an island to not be
// in Critical state.
const MinimumResourceThreshold = 5

// MaxCriticalConsecutiveTurns is the maximum consecutive turns an island can be in the critical state.
const MaxCriticalConsecutiveTurns = 3

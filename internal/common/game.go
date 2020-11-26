package common

import "github.com/SOMAS2020/SOMAS2020/internal/common/forage"

// GameState represents the game's state.
type GameState struct {
	// Day represents the 1-index day number of the game
	Day int

	// ForageRules contains representations of the current set rules for foraging.
	ForageRules ForageRules
}

// ForageRules contain the state of the current set rules for foraging.
type ForageRules struct {
	SplitRuleKey  forage.SplitRuleKey
	PayoffRuleKey forage.PayoffRuleKey
}

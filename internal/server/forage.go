package server

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
)

func getTeamForageInvestments(logger func(format string, a ...interface{}), gameState common.GameState) (teamForageInvestments map[common.ClientID]uint, err error) {
	teamForageInvestments = map[common.ClientID]uint{}
	totalInvestments := uint(0)

	for _, ci := range gameState.ClientInfos {
		c := ci.Client
		if ci.Alive {
			inv := c.GetForageInvestment(gameState)

			if ci.Resources < inv {
				logger("Team %v is overinvesting resources, defaulting to maximum available: %v", c.GetID(), ci.Resources)
				inv = ci.Resources
			}

			teamForageInvestments[c.GetID()] = inv
			totalInvestments += inv

			logger("Team %v invests %v resources into foraging", c.GetID(), inv)
		} else {
			logger("Team %v cannot forage, it is dead!", c.GetID())
		}
	}
	return
}

// StartForageRound starts a forage round.
func (s *SOMASServer) StartForageRound() error {
	s.Logf("Starting forage round!")

	teamForageInvestments, err := getTeamForageInvestments(s.Logf, s.gameState)
	if err != nil {
		return fmt.Errorf("Can't get team forage investments: %v", err)
	}

	// we just treat that everyone who wants to forage will collaborate for now...
	totalInvestments := uint(0)
	for _, inv := range teamForageInvestments {
		totalInvestments += inv
	}
	s.Logf("Total investments: %v", totalInvestments)

	totalPayoff := common.ComputeForagePayoff(s.gameState.ForageRules.PayoffRuleKey, totalInvestments)

	s.Logf("Total payoff from forage: %v (using rule: %v)", totalPayoff, s.gameState.ForageRules.PayoffRuleKey)

	split := common.ComputeForageSplit(s.gameState.ForageRules.SplitRuleKey, teamForageInvestments, totalPayoff)

	s.Logf("Split of payoff: %v", split)

	// we should be sending a message back to the clients now updating their state :) - we should also update GameState

	return nil
}

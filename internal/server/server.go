// Package server contains server-side code
package server

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/forage"
)

// Server represents the primary server interface exposed to the simulation.
type Server interface {
	GetEcho(s string) error
	Logf(format string, a ...interface{})

	StartForageRound() error
}

// SOMASServer implements Server.
type SOMASServer struct {
	clients   []common.Client
	gameState common.GameState
}

// SOMASServerFactory returns an instance of the main server we use.
func SOMASServerFactory() Server {
	return &SOMASServer{
		clients: common.RegisteredClients,
		gameState: common.GameState{
			Day: 1,
			ForageRules: common.ForageRules{
				SplitRuleKey:  forage.DefaultSplitRuleKey,
				PayoffRuleKey: forage.DefaultPayoffRuleKey,
			},
		},
	}
}

// GetEcho retrieves an echo from all the clients and make sure they are the same.
func (s *SOMASServer) GetEcho(str string) error {
	for _, c := range s.clients {
		got := c.Echo(str)
		if str != got {
			return fmt.Errorf("Echo error: want '%v' got '%v' from client %v",
				str, got, c.GetID())
		}
		s.Logf("Received echo `%v` from client %v", str, c.GetID())
	}
	return nil
}

// Logf is the server's default logger.
func (s *SOMASServer) Logf(format string, a ...interface{}) {
	log.Printf("[SERVER]: %v", fmt.Sprintf(format, a...))
}

func (s *SOMASServer) StartForageRound() error {
	s.Logf("Starting forage round!")
	teamForageInvestments := map[int]int{}
	totalInvestments := 0

	for _, c := range s.clients {
		inv := c.GetForageInvestment(s.gameState)

		// the GameState should contain the number of available resources, which we need to check whether a team is "overinvesting".

		teamForageInvestments[c.GetID()] = inv
		totalInvestments += inv
		s.Logf("Team %v invests %v resources into foraging", c.GetID(), inv)
	}

	s.Logf("Total forage investments: %v", totalInvestments)

	// we just treat that everyone who wants to forage will collaborate for now...

	totalPayoff := forage.ComputePayoff(s.gameState.ForageRules.PayoffRuleKey, totalInvestments)

	s.Logf("Total payoff from forage: %v (using rule: %v)", totalPayoff, s.gameState.ForageRules.PayoffRuleKey)

	split := forage.ComputeSplit(s.gameState.ForageRules.SplitRuleKey, teamForageInvestments, totalPayoff)

	s.Logf("Split of payoff: %v", split)

	// we should be sending a message back to the clients now updating their state :) - we should also update GameState

	return nil
}

package server

import (
	"sync"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

// runTurn runs a turn
func (s *SOMASServer) runTurn() error {
	s.logf("start runTurn")
	defer s.logf("finish runTurn")

	s.logf("TURN: %v, Season: %v", s.gameState.Turn, s.gameState.Season)

	s.startOfTurnUpdate()

	// run all orgs
	err := s.runOrgs()
	if err != nil {
		return errors.Errorf("Error running orgs: %v", err)
	}

	if err := s.endOfTurn(); err != nil {
		return errors.Errorf("Error running end of turn procedures: %v", err)
	}

	return nil
}

// runOrgs runs all the orgs
func (s *SOMASServer) runOrgs() error {
	s.logf("start runOrgs")
	defer s.logf("finish runOrgs")

	if err := s.runIIGO(); err != nil {
		return errors.Errorf("IIGO error: %v", err)
	}

	if err := s.runIIFO(); err != nil {
		return errors.Errorf("IIFO error: %v", err)
	}

	if err := s.runIITO(); err != nil {
		return errors.Errorf("IITO error: %v", err)
	}

	return nil
}

// startOfTurnUpdate sends the gameState at the start of the turn to all non-Dead clients.
func (s *SOMASServer) startOfTurnUpdate() {
	s.logf("start startOfTurnUpdate")
	defer s.logf("finish startOfTurnUpdate")

	// send update of entire gameState to alive clients
	N := len(s.gameState.ClientInfos)
	var wg sync.WaitGroup
	wg.Add(N)
	for id, ci := range s.gameState.ClientInfos {
		go func(id shared.ClientID, ci gamestate.ClientInfo) {
			defer wg.Done()
			if ci.LifeStatus != shared.Dead {
				c := s.clientMap[id]
				c.StartOfTurnUpdate(s.gameState.Copy())
			}
		}(id, ci)
	}
	wg.Wait()
}

// gameStateUpdate sends the gameState mid-turn to all non-Dead clients.
// For use by orgs to update game state after dispatching actions.
func (s *SOMASServer) gameStateUpdate() {
	s.logf("start gameStateUpdate")
	defer s.logf("finish gameStateUpdate")

	N := len(s.gameState.ClientInfos)
	var wg sync.WaitGroup
	wg.Add(N)
	for id, ci := range s.gameState.ClientInfos {
		go func(id shared.ClientID, ci gamestate.ClientInfo) {
			defer wg.Done()
			if ci.LifeStatus != shared.Dead {
				c := s.clientMap[id]
				c.GameStateUpdate(s.gameState.Copy())
			}
		}(id, ci)
	}
	wg.Wait()
}

// endOfTurn performs end of turn updates
func (s *SOMASServer) endOfTurn() error {
	s.logf("start endOfTurn")
	defer s.logf("finish endOfTurn")

	err := s.runOrgsEndOfTurn()
	if err != nil {
		return errors.Errorf("Failed to run orgs end of turn: %v", err)
	}

	// probe for disaster
	disasterHappened, err := s.probeDisaster()
	if err != nil {
		return errors.Errorf("Failed to probe disaster: %v", err)
	}
	// increment turn & season if needed
	s.incrementTurnAndSeason(disasterHappened)

	// deduct cost of living
	s.deductCostOfLiving(config.GameConfig().CostOfLiving)

	err = s.updateIslandLivingStatus()
	if err != nil {
		return errors.Errorf("Failed to update island living status: %v", err)
	}

	return nil
}

// runOrgsEndOfTurn runs all the end of turn variants of the orgs.
func (s *SOMASServer) runOrgsEndOfTurn() error {
	s.logf("start runOrgsEndOfTurn")
	defer s.logf("finish runOrgsEndOfTurn")

	if err := s.runIIGOEndOfTurn(); err != nil {
		return errors.Errorf("IIGO EndOfTurn error: %v", err)
	}

	if err := s.runIIFOEndOfTurn(); err != nil {
		return errors.Errorf("IIFO EndOfTurn error: %v", err)
	}

	if err := s.runIITOEndOfTurn(); err != nil {
		return errors.Errorf("IITO EndOfTurn error: %v", err)
	}

	return nil
}

// incrementTurnAndSeason increments turn, and season if a disaster happened.
func (s *SOMASServer) incrementTurnAndSeason(disasterHappened bool) {
	s.logf("start incrementTurnAndSeason")
	defer s.logf("finish incrementTurnAndSeason")

	s.gameState.Turn++
	if disasterHappened {
		s.gameState.Season++
	}
}

// deductCostOfLiving deducts CoL for all living islands, including critical ones
func (s *SOMASServer) deductCostOfLiving(costOfLiving int) {
	s.logf("start deductCostOfLiving")
	defer s.logf("finish deductCostOfLiving")

	for id, ci := range s.gameState.ClientInfos {
		if ci.LifeStatus != shared.Dead {
			ci.Resources -= costOfLiving
			s.gameState.ClientInfos[id] = ci
		}
	}
}

// updateIslandLivingStatus changes the islands Alive and Critical state depending
// on the island's resource state.
// Dead islands are not resurrected.
func (s *SOMASServer) updateIslandLivingStatus() error {
	s.logf("start updateIslandLivingStatus")
	defer s.logf("finish updateIslandLivingStatus")
	for id, ci := range s.gameState.ClientInfos {
		if ci.LifeStatus != shared.Dead {
			ci, err := updateIslandLivingStatusForClient(ci,
				config.GameConfig().MinimumResourceThreshold, config.GameConfig().MaxCriticalConsecutiveTurns)
			if err != nil {
				return errors.Errorf("Unable to update island living status for %v: %v",
					id, err)
			}
			s.gameState.ClientInfos[id] = ci
		}
	}
	return nil
}

func (s *SOMASServer) gameOver(maxTurns uint, maxSeasons uint) bool {
	st := s.gameState

	if !anyClientsAlive(st.ClientInfos) {
		s.logf("All clients are dead!")
		return true
	}

	// +1 due to 1-indexing
	if st.Turn >= maxTurns+1 {
		s.logf("Max turns '%v' reached or exceeded", maxTurns)
		return true
	}

	// +1 due to 1-indexing
	if st.Season >= maxSeasons+1 {
		s.logf("Max seasons '%v' reached or exceeded", maxSeasons)
		return true
	}

	return false
}

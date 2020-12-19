package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
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

	// get all end of turn actions
	err = s.getAndDispatchEndOfTurnActions()
	if err != nil {
		return errors.Errorf("Error running end of turn actions: %v", err)
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

	for id, ci := range s.gameState.ClientInfos {
		if ci.LifeStatus != shared.Dead {
			c := s.clientMap[id]
			c.StartOfTurnUpdate(s.gameState)
		}
	}
}

// startOfTurnUpdate sends the gameState mid-turn to all non-Dead clients.
// For use by orgs to update game state after dispatching actions.
func (s *SOMASServer) gameStateUpdate() {
	s.logf("start gameStateUpdate")
	defer s.logf("finish gameStateUpdate")

	for id, ci := range s.gameState.ClientInfos {
		if ci.LifeStatus != shared.Dead {
			c := s.clientMap[id]
			c.GameStateUpdate(s.gameState)
		}
	}
}

// getAndDispatchEndOfTurnActions gets all end of turn actions from the clients
func (s *SOMASServer) getAndDispatchEndOfTurnActions() error {
	s.logf("start getAndDispatchEndOfTurnActions")
	defer s.logf("finish getAndDispatchEndOfTurnActions")
	allActions := []common.Action{}

	for id, ci := range s.gameState.ClientInfos {
		if ci.LifeStatus != shared.Dead {
			c := s.clientMap[id]
			actions := c.EndOfTurnActions()
			allActions = append(allActions, actions...)
		}
	}

	// dispatch actions
	err := s.gameState.DispatchActions(allActions)
	if err != nil {
		return errors.Errorf("Error dispatching end of turn actions: %v", err)
	}

	// send updates
	s.gameStateUpdate()

	return nil
}

// endOfTurn performs end of turn updates
func (s *SOMASServer) endOfTurn() error {
	s.logf("start endOfTurn")
	defer s.logf("finish endOfTurn")
	// increment turn
	s.gameState.Turn++

	// increment season if disaster happened
	disasterHappened, err := s.probeDisaster()
	if err != nil {
		return err
	}
	if disasterHappened {
		s.gameState.Season++
	}

	err = s.deductCostOfLiving()
	if err != nil {
		return err
	}

	err = s.updateIslandLivingStatus()
	if err != nil {
		return err
	}

	return nil
}

// deductCostOfLiving deducts CoL for all living islands, including critical ones
func (s *SOMASServer) deductCostOfLiving() error {
	s.logf("start deductCostOfLiving")
	defer s.logf("finish deductCostOfLiving")
	for id, ci := range s.gameState.ClientInfos {
		if ci.LifeStatus != shared.Dead {
			ci.Resources -= config.CostOfLiving
			s.gameState.ClientInfos[id] = ci
		}
	}
	return nil
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
				config.MinimumResourceThreshold, config.MaxCriticalConsecutiveTurns)
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

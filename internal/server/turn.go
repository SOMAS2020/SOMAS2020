package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/action"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
)

// runTurn runs a turn
func (s *SOMASServer) runTurn() error {
	s.logf("start runTurn")
	defer s.logf("finish runTurn")

	s.logf("TURN: %v, Season: %v", s.gameState.Turn, s.gameState.Season)

	if err := s.updateIslands(); err != nil {
		return err
	}

	allActions := []action.Action{}

	// run all orgs and get all actions
	if actions, err := s.runOrgs(); err != nil {
		return err
	} else {
		allActions = append(allActions, actions...)
	}

	// get all end of turn actions
	if actions, err := s.getEndOfTurnActions(); err != nil {
		return err
	} else {
		allActions = append(allActions, actions...)
	}

	// dispatch actions
	if err := s.dispatchActions(allActions); err != nil {
		return err
	}

	if err := s.endOfTurn(); err != nil {
		return err
	}

	return nil
}

// runOrgs runs all the orgs
func (s *SOMASServer) runOrgs() ([]action.Action, error) {
	s.logf("start runOrgs")
	defer s.logf("finish runOrgs")

	allActions := []action.Action{}

	if actions, err := s.runIITO(); err != nil {
		return nil, err
	} else {
		allActions = append(allActions, actions...)
	}

	if actions, err := s.runIIFO(); err != nil {
		return nil, err
	} else {
		allActions = append(allActions, actions...)
	}

	if actions, err := s.runIIGO(); err != nil {
		return nil, err
	} else {
		allActions = append(allActions, actions...)
	}

	return allActions, nil
}

// updateIsland sends all the island the gameState at the start of the turn.
func (s *SOMASServer) updateIslands() error {
	s.logf("start updateIsland")
	defer s.logf("finish updateIsland")

	// send update of entire gameState to alive clients
	for id, ci := range s.gameState.ClientInfos {
		if ci.Alive {
			c := s.clientMap[id]
			c.ReceiveGameStateUpdate(s.gameState)
		}
	}

	return nil
}

func (s *SOMASServer) dispatchActions([]action.Action) error {
	s.logf("start dispatchActions")
	defer s.logf("finish dispatchActions")
	// TODO
	return nil
}

func (s *SOMASServer) getEndOfTurnActions() ([]action.Action, error) {
	s.logf("start getEndOfTurnActions")
	defer s.logf("finish getEndOfTurnActions")
	allActions := []action.Action{}

	for id, ci := range s.gameState.ClientInfos {
		if ci.Alive {
			c := s.clientMap[id]
			if actions, err := c.EndOfTurnActions(); err != nil {
				s.logf("EndOfTurnActions error for client '%v' ignored, treated as no action: %v",
					id, err)
			} else {
				allActions = append(allActions, actions...)
			}
		}
	}

	return allActions, nil
}

// endOfTurn performs end of turn actions
func (s *SOMASServer) endOfTurn() error {
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
	for id, ci := range s.gameState.ClientInfos {
		if ci.Alive {
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
	for id, ci := range s.gameState.ClientInfos {
		if ci.Alive {
			s.gameState.ClientInfos[id] = updateIslandLivingStatusForClient(ci,
				config.MinimumResourceThreshold, config.MaxCriticalConsecutiveTurns)
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

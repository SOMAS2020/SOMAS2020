package server

import "github.com/SOMAS2020/SOMAS2020/internal/common/config"

// runTurn runs a turn
func (s *SOMASServer) runTurn() error {
	s.logf("start runTurn")
	defer s.logf("finish runTurn")

	s.logf("TURN: %v, Season: %v", s.gameState.Turn, s.gameState.Season)

	if err := s.updateIslands(); err != nil {
		return err
	}

	if err := s.runIIFO(); err != nil {
		return err
	}

	if err := s.runIITO(); err != nil {
		return err
	}

	if err := s.runIIGO(); err != nil {
		return err
	}

	if err := s.getIslandsActions(); err != nil {
		return err
	}

	if err := s.endOfTurn(); err != nil {
		return err
	}

	return nil
}

// updateIsland updates resources for each island and on global game state - what
// information is included in the global game state and a resource update will depend
// heavily on the infrastructure team, and the design implementations of the environment
// and the IIGO and for this reason cannot be specified without further discussion (this is
// the primary source of information that can be used in any strategic algorithms for the
// iterative game). An island update could, for example, include the amount of resources
// an island currently has, the amount of resources currently in the common pool, whether
// or not an agent cheated in the last turn, whether there were any proposed rule changes
// in the previous turn, whether any rules were broken in the previous turn, etc. (these are
// all examples and should be agreed on collectively potentially in a dedicated meeting on
// this subject once the rest of the specification is established - from an implementation
// point of view adding more/less parameters does not seem to be a major concern at this
// stage).
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

// runIIFO : IIFO makes recommendations about the optimal (and fairest) contributions this term.
// to mitigate common risk dilemma.
func (s *SOMASServer) runIIFO() error {
	s.logf("start runIIFO")
	defer s.logf("finish runIIFO")
	// TODO:- IIFO team
	return nil
}

// runIITO : IITO makes recommendations about the optimal (and fairest) contributions this term
// to mitigate the common pool dilemma
func (s *SOMASServer) runIITO() error {
	s.logf("start runIITO")
	defer s.logf("finish runIITO")
	// TOOD:- IITO team
	return nil
}

// runIIGO : IIGO decides rule changes, elections, sanctions
func (s *SOMASServer) runIIGO() error {
	s.logf("start runIITO")
	defer s.logf("finish runIITO")
	// TOOD:- IIGO team
	return nil
}

// getIslandsAction obtains islands' decisions on their actions to the server to
// formally end the turn.
func (s *SOMASServer) getIslandsActions() error {
	s.logf("start getIslandsActions")
	defer s.logf("finish getIslandsActions")
	// TODO:- ?
	return nil
}

// probeDisaster checks if a disaster occurs this turn
func (s *SOMASServer) probeDisaster() (bool, error) {
	s.logf("start probeDisaster")
	defer s.logf("finish probeDisaster")
	// TOOD:- env team
	return false, nil
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

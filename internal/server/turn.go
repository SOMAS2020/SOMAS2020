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
	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)
	N := len(nonDeadClients)
	var wg sync.WaitGroup
	wg.Add(N)

	for _, id := range nonDeadClients {
		go func(id shared.ClientID, ci gamestate.ClientInfo) {
			defer wg.Done()
			c := s.clientMap[id]
			c.StartOfTurnUpdate(s.gameState.GetClientGameStateCopy(id))
		}(id, s.gameState.ClientInfos[id])
	}
	wg.Wait()
}

// gameStateUpdate sends the gameState mid-turn to all non-Dead clients.
// For use by orgs to update game state after dispatching actions.
func (s *SOMASServer) gameStateUpdate() {
	s.logf("start gameStateUpdate")
	defer s.logf("finish gameStateUpdate")

	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)
	N := len(nonDeadClients)
	var wg sync.WaitGroup
	wg.Add(N)
	for _, id := range nonDeadClients {
		go func(id shared.ClientID, ci gamestate.ClientInfo) {
			defer wg.Done()
			c := s.clientMap[id]
			c.GameStateUpdate(s.gameState.GetClientGameStateCopy(id))
		}(id, s.gameState.ClientInfos[id])
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

	err = s.runForage()
	if err != nil {
		return errors.Errorf("Failed to run hunt at end of turn: %v", err)
	}

	// Run Fish hunt test
	err = s.runDummyFishingExpedition()
	if err != nil {
		return errors.Errorf("Failed to run Fish hunt at end of turn: %v", err)
	}

	// probe for disaster
	updatedEnv, err := s.probeDisaster()
	if err != nil {
		return errors.Errorf("Failed to probe disaster: %v", err)
	}
	s.gameState.Environment = updatedEnv
	// increment turn & season if needed
	disasterHappened := updatedEnv.LastDisasterReport.Magnitude > 0
	s.incrementTurnAndSeason(disasterHappened)

	s.updateDeerPopulation([]int{2})

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

	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)
	N := len(nonDeadClients)
	resChan := make(chan clientInfoUpdateResult, N)
	wg := &sync.WaitGroup{}
	wg.Add(N)

	for _, id := range nonDeadClients {
		go func(id shared.ClientID, ci gamestate.ClientInfo) {
			defer wg.Done()
			ci.Resources -= costOfLiving
			resChan <- clientInfoUpdateResult{
				ID:  id,
				Ci:  ci,
				Err: nil,
			}
		}(id, s.gameState.ClientInfos[id])
	}

	wg.Wait()
	close(resChan)

	for res := range resChan {
		id, ci := res.ID, res.Ci
		// fine to ignore error, always nil
		s.gameState.ClientInfos[id] = ci
	}
}

// updateIslandLivingStatus changes the islands Alive and Critical state depending
// on the island's resource state.
// Dead islands are not resurrected.
func (s *SOMASServer) updateIslandLivingStatus() error {
	s.logf("start updateIslandLivingStatus")
	defer s.logf("finish updateIslandLivingStatus")

	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)
	N := len(nonDeadClients)
	resChan := make(chan clientInfoUpdateResult, N)
	wg := &sync.WaitGroup{}
	wg.Add(N)

	for _, id := range nonDeadClients {
		go func(id shared.ClientID, ci gamestate.ClientInfo) {
			defer wg.Done()
			ciNew, err := updateIslandLivingStatusForClient(ci,
				config.GameConfig().MinimumResourceThreshold, config.GameConfig().MaxCriticalConsecutiveTurns)
			resChan <- clientInfoUpdateResult{
				ID:  id,
				Ci:  ciNew,
				Err: err,
			}
		}(id, s.gameState.ClientInfos[id])
	}

	wg.Wait()
	close(resChan)

	for res := range resChan {
		id, ci, err := res.ID, res.Ci, res.Err
		if err != nil {
			return errors.Errorf("Failed to update island living status for '%v': %v", id, err)
		}
		s.gameState.ClientInfos[id] = ci
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

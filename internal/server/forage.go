package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/foraging"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

func (s *SOMASServer) runDeerHunt(participants map[shared.ClientID]float64) error {
	s.logf("start runDeerHunt")
	defer s.logf("finish runDeerHunt")

	hunt, err := foraging.CreateDeerHunt(participants)
	if err != nil {
		return errors.Errorf("Error running deer hunt: %v", err)
	}
	totalReturn := hunt.Hunt()
	s.logf("Hunt generated a return of %.3f from input of %.3f", totalReturn, hunt.TotalInput())
	return nil
}

func (s *SOMASServer) runFishingExpedition(participants map[shared.ClientID]float64) error {
	s.logf("start runFishHunt")
	defer s.logf("finish runFishHunt")

	huntF, err := foraging.CreateFishingExpedition(participants)
	if err != nil {
		return errors.Errorf("Error running fish hunt: %v", err)
	}
	totalReturn := huntF.Fish()
	s.logf("Fish Hunt generated a return of %.3f from input of %.3f", totalReturn, huntF.TotalInput())
	return nil
}

func (s *SOMASServer) runDummyHunt() error {
	huntParticipants := map[shared.ClientID]float64{shared.Team1: 1.0, shared.Team2: 0.9} // just to test for now
	return s.runDeerHunt(huntParticipants)
}

func (s *SOMASServer) runDummyFishingExpedition() error {
	huntParticipants := map[shared.ClientID]float64{shared.Team1: 1.0, shared.Team2: 0.9} // just to test for now
	return s.runFishingExpedition(huntParticipants)
}

// updateDeerPopulation adjusts deer pop. based on consumption of deer after hunt. Note that len(consumption) implies the number of
// days/turns that are to be simulated. If the intention is just to update after one turn, len(consumption) should be 1 and should
// containt the number of deer removed (hunted) from the env in the last turn
func (s *SOMASServer) updateDeerPopulation(consumption []int) {
	s.gameState.DeerPopulation.Simulate(consumption) // updates pop. according to DE definition
}

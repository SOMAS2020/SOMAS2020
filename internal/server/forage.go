package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/foraging"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

// runForage runs the foraging session, interacting with the alive agents and the environment
func (s *SOMASServer) runForage() error {
	s.logf("start runForage")
	defer s.logf("finish runForage")

	foragingParticipants, err := s.getForagingDecisions()
	if err != nil {
		return errors.Errorf("Something went wrong getting the foraging decision:%v", err)
	}

	deerHunters := make(map[shared.ClientID]shared.Resources)
	for id, decision := range foragingParticipants {
		if decision.Type == shared.DeerForageType {
			err := s.takeResources(id, decision.Contribution, "deer hunt participation")
			if err == nil {
				deerHunters[id] = decision.Contribution
			} else {
				s.logf("%v did not have enough resources to participate in foraging", id)
			}
		}
	}
	err = s.runDeerHunt(deerHunters)
	if err != nil {
		return errors.Errorf("Deer hunt returned with an error:%v", err)
	}
	return nil
}

func (s *SOMASServer) getForagingDecisions() (shared.ForagingDecisionsDict, error) {
	participants := shared.ForagingDecisionsDict{}
	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)
	var err error
	for _, id := range nonDeadClients {
		c := s.clientMap[id]
		participants[id], err = c.DecideForage()
		if err != nil {
			return participants, errors.Errorf("Failed to get foraging decision from %v: %v", id, err)
		}
	}
	return participants, nil
}

func (s *SOMASServer) runDeerHunt(participants map[shared.ClientID]shared.Resources) error {
	s.logf("start runDeerHunt")
	defer s.logf("finish runDeerHunt")
	hunt, err := foraging.CreateDeerHunt(participants)
	if err != nil {
		return errors.Errorf("Error running deer hunt: %v", err)
	}
	totalReturn := hunt.Hunt()

	totalContributions := shared.Resources(0)
	for _, contribution := range participants { totalContributions += contribution }

	for participantID, contribution := range participants {
		s.giveResources(participantID, (contribution/totalContributions) * totalReturn, "Deer hunt return")
	}

	s.logf("Hunt generated a return of %.3f from input of %.3f", totalReturn, hunt.TotalInput())
	return nil
}

/*func (s *SOMASServer) runDummyHunt() error {
	huntParticipants := map[shared.ClientID]float64{shared.Team1: 1.0, shared.Team2: 0.9} // just to test for now
	return s.runDeerHunt(huntParticipants)
}*/

// updateDeerPopulation adjusts deer pop. based on consumption of deer after hunt. Note that len(consumption) implies the number of
// days/turns that are to be simulated. If the intention is just to update after one turn, len(consumption) should be 1 and should
// containt the number of deer removed (hunted) from the env in the last turn
func (s *SOMASServer) updateDeerPopulation(consumption []int) {
	s.gameState.DeerPopulation.Simulate(consumption) // updates pop. according to DE definition
}

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
		// (*target[decision.Type])[id] = decision.Contribution
	}
	errD := s.runDeerHunt(deerHunters)

	if errD != nil {
		return errors.Errorf("Deer hunt returned with an error:%v", errD)
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

func (s *SOMASServer) runDeerHunt(contributions map[shared.ClientID]shared.Resources) error {
	s.logf("start runDeerHunt")
	defer s.logf("finish runDeerHunt")

	dhConf := s.gameConfig.ForagingConfig.DeerHuntConfig

	hunt, err := foraging.CreateDeerHunt(
		contributions,
		dhConf,
	)
	if err != nil {
		return errors.Errorf("Error running deer hunt: %v", err)
	}

	huntReport := hunt.Hunt(dhConf)

	totalContributions := shared.Resources(0)
	for _, contribution := range contributions {
		totalContributions += contribution
	}

	for participantID, contribution := range contributions {
		participantReturn := shared.Resources(0)
		if totalContributions != 0 {
			participantReturn = (contribution / totalContributions) * huntReport.TotalUtility
		}
		s.giveResources(participantID, participantReturn, "Deer hunt return")
		s.clientMap[participantID].ForageUpdate(shared.ForageDecision{
			Type:         shared.DeerForageType,
			Contribution: contribution,
		}, participantReturn)
	}

	s.logf("Hunt generated a return of %.3f from input of %.3f", huntReport.TotalUtility, hunt.TotalInput())

	// update deer population // TODO: decide if there is a better place to do this
	s.logf("Updating deer population after %v deer hunted", huntReport.NumberDeerCaught)
	s.updateDeerPopulation(huntReport.NumberDeerCaught) // update deer population based on hunt
	return nil
}

func (s *SOMASServer) runFishingExpedition(participants map[shared.ClientID]shared.Resources) error {
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
	huntParticipants := map[shared.ClientID]shared.Resources{shared.Team1: 1.0, shared.Team2: 0.9} // just to test for now
	return s.runDeerHunt(huntParticipants)
}

func (s *SOMASServer) runDummyFishingExpedition() error {
	huntParticipants := map[shared.ClientID]shared.Resources{shared.Team1: 1.0, shared.Team2: 0.9} // just to test for now
	return s.runFishingExpedition(huntParticipants)
}

// updateDeerPopulation adjusts deer pop. based on consumption of deer after hunt
func (s *SOMASServer) updateDeerPopulation(consumption uint) {
	s.gameState.DeerPopulation.Simulate([]int{int(consumption)}) // updates pop. according to DE definition
}

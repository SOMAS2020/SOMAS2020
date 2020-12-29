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
	fishers := make(map[shared.ClientID]shared.Resources)

	// used to keep track of groups choosing different foraging types
	// this allows different foraging types to be added in future without having
	// to change assignment logic below
	forageGroups := map[shared.ForageType]forageContributionType{
		shared.DeerForageType: {
			partyContributions: &deerHunters,
			takeResourceReason: "deer hunt participation",
		},
		shared.FishForageType: {
			partyContributions: &fishers,
			takeResourceReason: "fishing participation",
		},
	}

	for id, decision := range foragingParticipants {
		forageGroup := forageGroups[decision.Type]
		err := s.takeResources(id, decision.Contribution, forageGroup.takeResourceReason)
		if err == nil {
			// assign contribution to client ID within appropriate forage group
			(*forageGroup.partyContributions)[id] = decision.Contribution
		} else {
			s.logf("%v did not have enough resources to participate in foraging", id)
		}
	}
	errD := s.runDeerHunt(deerHunters)
	errF := s.runFishingExpedition(fishers)

	if errD != nil {
		return errors.Errorf("Deer hunt returned with an error:%v", errD)
	}

	if errF != nil {
		return errors.Errorf("Fishing expedition returned with an error:%v", errD)
	}

	return nil
}

// forageContributionType is a helper type to store forage participants and their contributions
// and the corresponding reason for taking resources from concerned clients
type forageContributionType struct {
	partyContributions *map[shared.ClientID]shared.Resources
	takeResourceReason string
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
	huntReport.Turn = s.gameState.Turn // update report's Turn with actual turn value
	// update foraging history
	s.gameState.ForagingHistory[shared.DeerForageType] = append(s.gameState.ForagingHistory[shared.DeerForageType], huntReport)

	totalContributions := hunt.TotalInput()

	// distribute return amongst participants
	for participantID, contribution := range contributions {
		participantReturn := shared.Resources(0.0)
		if totalContributions != 0.0 {
			participantReturn = (contribution / totalContributions) * huntReport.TotalUtility
		}
		s.giveResources(participantID, participantReturn, "Deer hunt return")
		s.clientMap[participantID].ForageUpdate(shared.ForageDecision{
			Type:         shared.DeerForageType,
			Contribution: contribution,
		}, participantReturn)
	}

	s.logf("Deer hunt report: %v", huntReport.Display())

	// update deer population // TODO: decide if there is a better place to do this
	s.logf("Updating deer population after %v deer hunted", huntReport.NumberCaught)
	s.updateDeerPopulation(huntReport.NumberCaught) // update deer population based on hunt

	return nil
}

func (s *SOMASServer) runFishingExpedition(participants map[shared.ClientID]shared.Resources) error {
	s.logf("start runFishHunt")
	defer s.logf("finish runFishHunt")

	fConf := s.gameConfig.ForagingConfig.FishingConfig

	huntF, err := foraging.CreateFishingExpedition(participants)
	if err != nil {
		return errors.Errorf("Error running fish hunt: %v", err)
	}
	fishingReport := huntF.Fish(fConf)

	fishingReport.Turn = s.gameState.Turn // update report's Turn with actual turn value
	// update foraging history
	s.gameState.ForagingHistory[shared.FishForageType] = append(s.gameState.ForagingHistory[shared.FishForageType], fishingReport)

	s.logf("Fishing expedition report: %v", fishingReport.Display())

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

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

		if !shared.IsValidForageType(decision.Type) {
			s.logf("%v client selected invalid forag type in foraging decision: ", decision.Type)
		} else {
			forageGroup := forageGroups[decision.Type]
			err := s.takeResources(id, decision.Contribution, forageGroup.takeResourceReason)

			if err == nil {
				if decision.Contribution > 0.0 {
					(*forageGroup.partyContributions)[id] = decision.Contribution // assign contribution to client ID within appropriate forage group
				} else {
					s.logf("%v did not contribute resources and will not participate in this foraging round.", id)
				}
			} else {
				s.logf("%v did not have enough resources to participate in foraging", id)
			}
		}
	}

	if len(deerHunters) > 0 {
		errD := s.runDeerHunt(deerHunters)

		if errD != nil {
			return errors.Errorf("Deer hunt returned with an error: %v", errD)
		}
	}

	if len(fishers) > 0 {
		errF := s.runFishingExpedition(fishers)

		if errF != nil {
			return errors.Errorf("Fishing expedition returned with an error: %v", errF)
		}
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
		s.logf,
	)
	if err != nil {
		return errors.Errorf("Error running deer hunt: %v", err)
	}

	huntReport := hunt.Hunt(dhConf, uint(s.gameState.DeerPopulation.Population))
	huntReport.Turn = s.gameState.Turn // update report's Turn with actual turn value
	// update foraging history
	if s.gameState.ForagingHistory[shared.DeerForageType] == nil {
		return errors.Errorf("Foraging history not initialised properly: %v", err)
	}
	s.gameState.ForagingHistory[shared.DeerForageType] = append(s.gameState.ForagingHistory[shared.DeerForageType], huntReport)

	s.distributeForageReturn(contributions, huntReport)

	s.logf("Deer hunt report: %v", huntReport.Display())

	// update deer population // TODO: decide if there is a better place to do this
	s.logf("Updating deer population after %v deer hunted", huntReport.NumberCaught)
	s.updateDeerPopulation(huntReport.NumberCaught) // update deer population based on hunt

	return nil
}

func (s *SOMASServer) distributeForageReturn(contributions map[shared.ClientID]shared.Resources, huntReport foraging.ForagingReport) {
	// distribute return amongst participants

	totalContributions := shared.Resources(0)
	for _, c := range huntReport.ParticipantContributions {
		totalContributions += c
	}

	if len(huntReport.ParticipantContributions) == 0 {
		return // to prevent div0 below. Also, no need to evaluate further
	}

	// auxiliary return info. Just to avoid repetition.
	type auxReturnInfo struct {
		distrStrategy        shared.ResourceDistributionStrategy
		resourceReturnReason string
	}

	for participantID, contribution := range contributions {
		deerReturnStrat := s.gameConfig.ForagingConfig.DeerHuntConfig.DistributionStrategy
		fishReturnStrat := s.gameConfig.ForagingConfig.FishingConfig.DistributionStrategy

		returnInfoPerType := map[shared.ForageType]auxReturnInfo{
			shared.DeerForageType: {distrStrategy: deerReturnStrat, resourceReturnReason: "Deer hunt return"},
			shared.FishForageType: {distrStrategy: fishReturnStrat, resourceReturnReason: "Fishing return"},
		}
		participantReturn := shared.Resources(0.0)      // default to zero return
		retReason := "Unspecified foraging return type" // default if f. type not found
		// check if the foraging type has been specified above
		if r, ok := returnInfoPerType[huntReport.ForageType]; ok {
			retReason = r.resourceReturnReason

			if totalContributions > 0.0 {
				switch r.distrStrategy {
				case shared.InputProportionalSplit:
					participantReturn = (contribution / totalContributions) * huntReport.TotalUtility
				case shared.EqualSplit, shared.RankProportionalSplit: // RankProportional is just same as equal split for now
					participantReturn = huntReport.TotalUtility / shared.Resources(len(huntReport.ParticipantContributions)) // this casting is a bit lazy
				}
			}
		}

		err := s.giveResources(participantID, participantReturn, retReason)
		if err != nil {
			s.logf("Ignoring failure to give resources in distributeForageReturn: %v", err)
		}
		s.clientMap[participantID].ForageUpdate(shared.ForageDecision{
			Type:         huntReport.ForageType,
			Contribution: contribution,
		}, participantReturn, huntReport.NumberCaught)
	}
}

func (s *SOMASServer) runFishingExpedition(contributions map[shared.ClientID]shared.Resources) error {
	s.logf("start runFishHunt")
	defer s.logf("finish runFishHunt")

	fConf := s.gameConfig.ForagingConfig.FishingConfig

	huntF, err := foraging.CreateFishingExpedition(contributions, fConf, s.logf)
	if err != nil {
		return errors.Errorf("Error running fish hunt: %v", err)
	}
	fishingReport := huntF.Fish(fConf)

	fishingReport.Turn = s.gameState.Turn // update report's Turn with actual turn value
	if s.gameState.ForagingHistory[shared.DeerForageType] == nil {
		return errors.Errorf("Foraging history not initialised properly: %v", err)
	}
	// update foraging history
	s.gameState.ForagingHistory[shared.FishForageType] = append(s.gameState.ForagingHistory[shared.FishForageType], fishingReport)

	s.distributeForageReturn(contributions, fishingReport)

	s.logf("Fishing expedition report: %v", fishingReport.Display())

	return nil
}

// updateDeerPopulation adjusts deer pop. based on consumption of deer after hunt
func (s *SOMASServer) updateDeerPopulation(consumption uint) {
	updatedModel := s.gameState.DeerPopulation.Simulate([]int{int(consumption)}) // updates pop. according to DE definition
	s.gameState.DeerPopulation = updatedModel
}

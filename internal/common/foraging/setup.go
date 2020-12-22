package foraging

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// CreateDeerHunt receives hunt participants and their contributions and returns a DeerHunt
func CreateDeerHunt(teamResourceInputs map[shared.ClientID]float64) (DeerHunt, error) {
	fConf := config.GameConfig().ForagingConfig
	params := deerHuntParams{p: fConf.BernoulliProb, lam: fConf.ExponentialRate}
	return DeerHunt{ParticipantContributions: teamResourceInputs, params: params}, nil // returning error too for future use
}

// CreateFishHunt sees the participants and their contributions and returns the value of FishHunt
func CreateFishHunt(teamResourceInputs map[shared.ClientID]float64) (FishHunt, error) {
	fConf := config.GameConfig().ForagingConfig
	params := fishHuntParams{Mu: fConf.FishingMean, Sigma: fConf.FishingVariance}
	return FishHunt{Participants: teamResourceInputs, Params: params}, nil // returning error too for future use
}

// CreateDeerPopulationModel returns the target population model. The formulation of this model should be changed here before runtime
func CreateDeerPopulationModel() DeerPopulationModel {
	return createBasicDeerPopulationModel()
}

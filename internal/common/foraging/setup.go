package foraging

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

// CreateDeerHunt receives hunt participants and their contributions and returns a DeerHunt
func CreateDeerHunt(teamResourceInputs map[shared.ClientID]shared.Resources, dhConf config.DeerHuntConfig, logger shared.Logger) (DeerHunt, error) {
	if len(teamResourceInputs) == 0 {
		return DeerHunt{}, errors.Errorf("No deer hunt resource contributions specified!")
	}
	params := deerHuntParams{p: dhConf.BernoulliProb, lam: dhConf.ExponentialRate}
	return DeerHunt{ParticipantContributions: teamResourceInputs, params: params, logger: logger}, nil // returning error too for future use
}

// CreateFishingExpedition sees the participants and their contributions and returns the value of FishHunt
func CreateFishingExpedition(teamResourceInputs map[shared.ClientID]shared.Resources, fConf config.FishingConfig) (FishingExpedition, error) {
	if len(teamResourceInputs) == 0 {
		return FishingExpedition{}, errors.Errorf("No fishing resource contributions specified!")
	}
	params := fishingParams{Mu: fConf.Mean, Sigma: fConf.Variance}
	return FishingExpedition{ParticipantContributions: teamResourceInputs, params: params}, nil // returning error too for future use
}

// CreateDeerPopulationModel returns the target population model. The formulation of this model should be changed here before runtime
func CreateDeerPopulationModel(dhConf config.DeerHuntConfig) DeerPopulationModel {
	return createBasicDeerPopulationModel(dhConf)
}

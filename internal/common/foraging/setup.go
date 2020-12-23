package foraging

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// CreateDeerHunt receives hunt participants and their contributions and returns a DeerHunt
func CreateDeerHunt(teamResourceInputs map[shared.ClientID]shared.ForageContribution) (DeerHunt, error) {
	dhConf := config.GameConfig().ForagingConfig.DeerHuntConfig
	params := deerHuntParams{p: dhConf.BernoulliProb, lam: dhConf.ExponentialRate}
	return DeerHunt{ParticipantContributions: teamResourceInputs, params: params}, nil // returning error too for future use
}

// CreateFishingExpedition sees the participants and their contributions and returns the value of FishHunt
func CreateFishingExpedition(teamResourceInputs map[shared.ClientID]shared.ForageContribution) (FishingExpedition, error) {
	fConf := config.GameConfig().ForagingConfig.FishingConfig
	params := fishingParams{Mu: fConf.Mean, Sigma: fConf.Variance}
	return FishingExpedition{ParticipantContributions: teamResourceInputs, Params: params}, nil // returning error too for future use
}

// CreateDeerPopulationModel returns the target population model. The formulation of this model should be changed here before runtime
func CreateDeerPopulationModel() DeerPopulationModel {
	return createBasicDeerPopulationModel()
}

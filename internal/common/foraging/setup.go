package foraging

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

var gameConf = config.GameConfig()

// CreateDeerHunt receives hunt participants and their contributions and returns a DeerHunt
func CreateDeerHunt(teamResourceInputs map[shared.ClientID]float64) DeerHunt {
	params := DeerHuntParams{P: gameConf.ForagingConfig.BernoulliProb, Lam: gameConf.ForagingConfig.ExponentialRate}
	return DeerHunt{Participants: teamResourceInputs, Params: params}
}

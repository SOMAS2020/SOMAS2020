package foraging

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// CreateDeerHunt receives hunt participants and their contributions and returns a DeerHunt
func CreateDeerHunt(teamResourceInputs map[shared.ClientID]float64) (DeerHunt, error) {
	fConf := config.GameConfig().ForagingConfig
	params := deerHuntParams{p: fConf.BernoulliProb, lam: fConf.ExponentialRate}
	return DeerHunt{Participants: teamResourceInputs, params: params}, nil // returning error too for future use
}

package team2

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type ForagingResults struct {
	Hunters int //no. teams that we know (they have shared the info with us)
	Fishers int
	Result  shared.Resources //total fron all the teams we know
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	ft := int(math.Round(rand.Float64())) // implement the normal distribution which shifts closer to hunt or fish
	var Threshold float64=decideThreshold()
	//base threshold - when we know nothing
	//a threshold when we know stuff

	if rand.Float64() > Threshold { //we fish when above the threshold
		ft := 1
	} else {                       // we hunt when below the threshold
		ft := 0
	}
	return shared.ForageDecision{
		Type:         shared.ForageType(ft),
		Contribution: shared.Resources(20), //contribute fixed amount for now
	}, nil
}

// ReceiveForageInfo lets clients know what other clients have obtained from their most recent foraging attempt.
// Most recent foraging attempt includes information about: foraging DecisionMade and ResourceObtained as well
// as where this information came from.
// OPTIONAL.
func (c *client) ReceiveForageInfo(neighbourForaging []shared.ForageShareInfo) {
	// Return on Investment

	for _, val := range neighbourForaging {
		decisionMade := val.DecisionMade
		resourcesObtained := val.ResourceObtained

		newResult := ForageInfo{
			DecisionMade:      decisionMade,
			ResourcesObtained: resourcesObtained,
		}
		hist := append(c.foragingReturnsHist[val.SharedFrom], newResult)
		c.foragingReturnsHist[val.SharedFrom] = hist
	}
}

func (c *client) decideThreshold() () {
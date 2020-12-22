package foraging

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/stat/distuv"
)

// FishHunt The teams that are involved and the resources they put in.
type FishHunt struct {
	ParticipantContributions map[shared.ClientID]float64
	Params                   fishHuntParams
}

// fishHuntParams : Defines the parameters for the normal distibution for the fishing returns
type fishHuntParams struct {
	Mu    float64
	Sigma float64
}

// TotalInput provides a sum of the total inputs
func (f FishHunt) TotalInput() float64 {
	i := 0.0                                       // Sum of the inputs INIT
	for _, x := range f.ParticipantContributions { // Check for every particpant ("_" ignore the index)
		i += x // Sum the values within each index
	}
	return i //return
}

func fishUtilityTier(input float64, maxFishPerHunt uint, decay float64) uint {
	sum := 0.0                                  // Sum = cumulative sum of the tier values
	for i := uint(0); i < maxFishPerHunt; i++ { // Checks the value of the input in comparision to the minimum needed for each tier
		sum += math.Pow(decay, float64(i+1)) // incrementING sum of tier values
		if input < sum {                     // Check condition that the input is less than the sum of  tier values
			return i // if so return the tier we are located in
		}
	}
	return maxFishPerHunt
}

// fishReturn is the normal distibtuion
func fishReturn(params fishHuntParams) float64 {
	F := distuv.Normal{
		Mu:    params.Mu,    // mean of the normal dist
		Sigma: params.Sigma, // Var of the normal dist
	}
	return F.Rand()
}

// HuntFish computes the return from a fishing expedition
func (f FishHunt) HuntFish() float64 {
	fConf := config.GameConfig().ForagingConfig.FishingConfig
	input := f.TotalInput()
	decay := fConf.IncrementalInputDecay
	maxFish := fConf.MaxFishPerHunt
	nFishFromInput := fishUtilityTier(input, maxFish, decay) // get max number of fish allowed for given resource input
	utility := 0.0
	for i := uint(1); i < nFishFromInput; i++ {
		utility += fishReturn(f.Params)
	}
	return utility
}

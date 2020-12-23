package foraging

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/stat/distuv"
)

// FishingExpedition The teams that are involved and the resources they put in.
type FishingExpedition struct {
	ParticipantContributions map[shared.ClientID]shared.ForageContribution
	Params                   fishingParams
}

// fishingParams : Defines the parameters for the normal distibution for the fishing returns
type fishingParams struct {
	Mu    float64
	Sigma float64
}

// TotalInput provides a sum of the total inputs
func (f FishingExpedition) TotalInput() shared.ForageContribution {
	i := 0.0                                       // Sum of the inputs INIT
	for _, x := range f.ParticipantContributions { // Check for every particpant ("_" ignore the index)
		i += x // Sum the values within each index
	}
	return i //return
}

func fishUtilityTier(input shared.ForageContribution, maxFishPerHunt uint, decay float64) uint {
	sum := 0.0                                  // Sum = cumulative sum of the tier values
	for i := uint(0); i < maxFishPerHunt; i++ { // Checks the value of the input in comparision to the minimum needed for each tier
		sum += math.Pow(decay, float64(i+1)) // incrementING sum of tier values
		if input < sum {                     // Check condition that the input is less than the sum of  tier values
			return i // if so return the tier we are located in
		}
	}
	return maxFishPerHunt
}

// fishingReturn is the normal distibtuion
func fishingReturn(params fishingParams) shared.ForageReturn {
	F := distuv.Normal{
		Mu:    params.Mu,    // mean of the normal dist
		Sigma: params.Sigma, // Var of the normal dist
	}
	return F.Rand()
}

// Fish computes the return from a fishing expedition
func (f FishingExpedition) Fish() shared.ForageReturn {
	fConf := config.GameConfig().ForagingConfig.FishingConfig
	input := f.TotalInput()
	decay := fConf.IncrementalInputDecay
	maxFish := fConf.MaxFishPerHunt
	nFishFromInput := fishUtilityTier(input, maxFish, decay) // get max number of fish allowed for given resource input
	utility := 0.0
	for i := uint(1); i < nFishFromInput; i++ {
		utility += fishingReturn(f.Params)
	}
	return utility
}

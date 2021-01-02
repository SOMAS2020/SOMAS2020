package foraging

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/stat/distuv"
)

// FishingExpedition The teams that are involved and the resources they put in.
type FishingExpedition struct {
	ParticipantContributions map[shared.ClientID]shared.Resources
	params                   fishingParams
}

// fishingParams : Defines the parameters for the normal distibution for the fishing returns
type fishingParams struct {
	Mu    float64
	Sigma float64
}

// FishingReport holds information about the result of a fishing expedition
type FishingReport struct {
	InputResources   shared.Resources
	NumberFishermen  uint
	NumberFishCaught uint
	TotalUtility     shared.Resources
	FishWeights      []float64
}

// TotalInput provides a sum of the total inputs
// TotalInput simply sums the total group resource input of hunt participants
func (f FishingExpedition) TotalInput() shared.Resources {
	return getTotalInput(f.ParticipantContributions)
}

// fishingReturn is the normal distibtuion
func fishingReturn(params fishingParams) shared.Resources {
	F := distuv.Normal{
		Mu:    params.Mu,    // mean of the normal dist
		Sigma: params.Sigma, // Var of the normal dist
	}
	return shared.Resources(F.Rand())
}

// Fish computes the return from a fishing expedition
func (f FishingExpedition) Fish(fConf config.FishingConfig) ForagingReport {
	input := f.TotalInput()
	// get max number of deer allowed for given resource input
	nFishFromInput := utilityTier(input, fConf.MaxFishPerHunt, fConf.IncrementalInputDecay)
	returns := []shared.Resources{} // store return for each potential fish we could catch

	for i := uint(0); i < nFishFromInput; i++ {
		utility := fishingReturn(f.params) * shared.Resources(fConf.ResourceMultiplier) // scale return by resource multiplier
		returns = append(returns, utility)
	}
	return compileForagingReport(shared.FishForageType, f.ParticipantContributions, returns)
}

package foraging

// see https://colab.research.google.com/drive/1g1tiX27Ds7FGjj4_WjFB3OLj8Fat_Ur5?usp=sharing for experiments + simulations

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/stat/distuv"
)

type deerHuntParams struct {
	p   float64 // Bernoulli p variable (whether or not a deer is caught)
	lam float64 // Exponential lambda (scale) param for W (weight variable)
}

// DeerHunt captures the hunt participants (teams) and their resource contributions, as well as hunt params
type DeerHunt struct {
	ParticipantContributions map[shared.ClientID]float64
	params                   deerHuntParams
}

// TotalInput simply sums the total group resource input of hunt participants
func (d DeerHunt) TotalInput() float64 {
	i := 0.0
	for _, x := range d.ParticipantContributions {
		i += x
	}
	return i
}

// Hunt returns the utility from a deer hunt
func (d DeerHunt) Hunt() float64 {
	fConf := config.GameConfig().ForagingConfig

	input := d.TotalInput()
	decay := fConf.IncrementalInputDecay
	maxDeer := fConf.MaxDeerPerHunt
	nDeerFromInput := deerUtilityTier(input, maxDeer, decay) // get max number of deer allowed for given resource input
	utility := 0.0
	for i := uint(0); i < nDeerFromInput; i++ {
		utility += deerReturn(d.params)
	}
	return utility
}

// deerUtilityTier gets the discrete utility tier (i.e. max number of deer) for given scalar input
func deerUtilityTier(input float64, maxDeerPerHunt uint, decay float64) uint {
	sum := 0.0
	for i := uint(0); i < maxDeerPerHunt; i++ {
		sum += math.Pow(decay, float64(i))
		if input < sum {
			return i
		}
	}
	return maxDeerPerHunt
}

// deerReturn() is effectively the combination of two other RVs:
// - D: Bernoulli RV that represents the probaility of catching a deer at all (binary). Usually p - i.e. P(D=1) = p - will be fairly
// close to 1 (fairly high chance of catching a deer if you invest the resources)
// - W: A continuous RV that adds some variance to the return. This could be interpreted as the weight of the deer that is caught. W is
// exponentially distributed such that the prevalence of deer of certain size is inversely prop. to the size.
// returns H, where H = D*(1+W) is an other random variable
func deerReturn(params deerHuntParams) float64 {
	W := distuv.Exponential{Rate: params.lam} // Rate = lambda
	D := distuv.Bernoulli{P: params.p}        // Bernoulli RV where `P` = P(X=1)
	return D.Rand() * (1 + W.Rand())
}

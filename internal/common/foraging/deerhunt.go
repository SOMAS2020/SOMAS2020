package foraging

// see https://colab.research.google.com/drive/1g1tiX27Ds7FGjj4_WjFB3OLj8Fat_Ur5?usp=sharing for experiments + simulations

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/stat/distuv"
)

// defines the incremental increase in input resources required to move up a utility tier (to be able to hunt another deer)
var deerUtilityIncrements = []float64{1.0, 0.75, 0.5, 0.25} //TODO: move this to central config store

type DeerHuntParams struct {
	P   float64 // Bernoulli p variable (whether or not a deer is caught)
	Lam float64 // Exponential lambda (scale) param for W (weight variable)
}

// DeerHunt captures the hunt participants (teams) and their resource contributions, as well as hunt params
type DeerHunt struct {
	Participants map[shared.ClientID]float64
	Params       DeerHuntParams
}

// TotalInput simply sums the total group resource input of hunt participants
func (d DeerHunt) TotalInput() float64 {
	i := 0.0
	for _, x := range d.Participants {
		i += x
	}
	return i
}

// Hunt returns the utility from a deer hunt
func (d DeerHunt) Hunt() float64 {
	input := d.TotalInput()
	maxDeer := deerUtilityTier(input, deerUtilityIncrements) // get max number of deer allowed for given resource input
	utility := 0.0
	for i := 1; i < maxDeer; i++ {
		utility += deerReturn(d.Params)
	}
	return utility
}

// deerUtilityTier gets the discrete utility tier (i.e. max number of deer) for given scalar input
func deerUtilityTier(input float64, increments []float64) int {
	if len(increments) == 0 || input < increments[0] {
		return 0
	}
	sum := 0.0
	for i := 0; i < len(increments); i++ {
		sum += increments[i]
		fmt.Println(sum)
		if input < sum {
			return i
		}
	}
	return len(increments)
}

// deerReturn() is effectively the combination of two other RVs:
// - D: Bernoulli RV that represents the probaility of catching a deer at all (binary). Usually p - i.e. P(D=1) = p - will be fairly
// close to 1 (fairly high chance of catching a deer if you invest the resources)
// - W: A continuous RV that adds some variance to the return. This could be interpreted as the weight of the deer that is caught. W is
// exponentially distributed such that the prevalence of deer of certain size is inversely prop. to the size.
// returns H, where H = D*(1+W) is an other random variable
func deerReturn(params DeerHuntParams) float64 {
	W := distuv.Exponential{Rate: params.Lam} // Rate = lambda
	D := distuv.Bernoulli{P: params.P}        // Bernoulli RV where `P` = P(X=1)
	return D.Rand() * (1 + W.Rand())
}

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
	ParticipantContributions map[shared.ClientID]shared.Resources
	params                   deerHuntParams
}

// TotalInput simply sums the total group resource input of hunt participants
func (d DeerHunt) TotalInput() shared.Resources {
	return getTotalInput(d.ParticipantContributions)
}

// Hunt returns the utility from a deer hunt
func (d DeerHunt) Hunt(dhConf config.DeerHuntConfig, deerPopulation uint) ForagingReport {
	input := d.TotalInput()
	// get max number of deer allowed for given resource input
	nDeerFromInput := utilityTier(input, dhConf.MaxDeerPerHunt, dhConf.IncrementalInputDecay)
	returns := []shared.Resources{}

	for i := uint(0); i < nDeerFromInput; i++ {
		d.params.p = getPopulationLinkedProbability(dhConf, deerPopulation)
		utility := deerReturn(d.params) * shared.Resources(dhConf.ResourceMultiplier) // scale return by resource multiplier
		returns = append(returns, utility)
		if utility > 0 { // a deer was caught and so should be removed from population
			deerPopulation = uint(math.Max(0, float64(deerPopulation)-1)) // min pop is zero. Assume no population growth (from DE) effects during short hunt
		}
	}
	return compileForagingReport(shared.DeerForageType, d.ParticipantContributions, returns)
}

// deerReturn() is effectively the combination of two other RVs:
// - D: Bernoulli RV that represents the probaility of catching a deer at all (binary). Usually p - i.e. P(D=1) = p - will be fairly
// close to 1 (fairly high chance of catching a deer if you invest the resources)
// - W: A continuous RV that adds some variance to the return. This could be interpreted as the weight of the deer that is caught. W is
// exponentially distributed such that the prevalence of deer of certain size is inversely prop. to the size.
// returns H, where H = D*(1+W) is an other random variable
func deerReturn(params deerHuntParams) shared.Resources {
	W := distuv.Exponential{Rate: params.lam} // Rate = lambda
	D := distuv.Bernoulli{P: params.p}        // Bernoulli RV where `P` = P(X=1)
	return shared.Resources(D.Rand() * (1 + W.Rand()))
}

func getPopulationLinkedProbability(dhConf config.DeerHuntConfig, population uint) float64 {
	pCritical := uint(1) // TODO: move to config
	pMax := dhConf.MaxDeerPopulation / dhConf.MaxDeerPerHunt
	thetaCritical := 0.85 // TODO: move to config
	thetaMax := 0.95      // TODO: move to config
	p := population / dhConf.MaxDeerPerHunt
	alpha := (thetaMax - thetaCritical) / float64((pMax - pCritical))

	f1 := func(p uint) float64 { return thetaCritical * float64(p) }
	f2 := func(p uint) float64 { return alpha*float64((p-pCritical)) + thetaCritical }

	if p < pCritical {
		return f1(p)
	}
	return f2(p)
}

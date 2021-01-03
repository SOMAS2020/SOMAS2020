package foraging

// see https://colab.research.google.com/drive/1g1tiX27Ds7FGjj4_WjFB3OLj8Fat_Ur5?usp=sharing for experiments + simulations

import (
	"fmt"
	"log"
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
	nDeerFromInput := utilityTier(input, dhConf.MaxDeerPerHunt, dhConf.IncrementalInputDecay, dhConf.InputScaler)
	returns := []shared.Resources{}

	for i := uint(0); i < nDeerFromInput; i++ {
		d.params.p = d.getPopulationLinkedProbability(dhConf, deerPopulation)
		utility := deerReturn(d.params) * shared.Resources(dhConf.OutputScaler) // scale raw deerReturn to be in range with other resource quantities
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

// getPopulationLinkedProbability returns the Bernoulli probability of catching a deer given the current running deer population.
// The dynamics and variables implemented in this function are documented in the README in this package.
func (d DeerHunt) getPopulationLinkedProbability(dhConf config.DeerHuntConfig, population uint) float64 {
	pCritical := 1.0
	pMax := float64(dhConf.MaxDeerPopulation) / float64(dhConf.MaxDeerPerHunt)
	thetaCritical := dhConf.ThetaCritical
	thetaMax := dhConf.ThetaMax
	p := float64(population) / float64(dhConf.MaxDeerPerHunt)
	alpha := (thetaMax - thetaCritical) / (pMax - float64(pCritical))

	f1 := func(p float64) float64 { return thetaCritical * float64(p) }
	f2 := func(p float64) float64 { return alpha*float64((p-pCritical)) + thetaCritical }

	if p < pCritical {
		theta := f1(p)
		d.Logf("Deer population ratio below critical level: P(t)=%v, p=%v, p_crit=%v, theta=%.2f", population, p, pCritical, theta)
		return theta
	}
	return f2(p)
}

// Logf is a this type's custom logger
func (d DeerHunt) Logf(format string, a ...interface{}) {
	log.Printf("[SERVER][DeerHunt]: %v", fmt.Sprintf(format, a...))
}

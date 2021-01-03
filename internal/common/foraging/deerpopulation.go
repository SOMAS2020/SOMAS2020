package foraging

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/simulation"
)

// DeerPopulationModel encapsulates a deer population over time (governed by a predefined DE)
type DeerPopulationModel struct {
	deProblem  simulation.ODEProblem        // definition of DE governing rate of change of deeer pop.
	Population float64                      // current number of deer in env
	T          float64                      // temporal parameter. Time, turn or whatever other incarnation
	deState    func(float64) (t, y float64) // holds internal state of DE and allows to call next solution step
}

// Logf is the client's logger that prepends logs with your ID. This makes
// it easier to read logs. DO NOT use other loggers that will mess logs up!
// BASE: Do not overwrite in team client.
func (dp DeerPopulationModel) Logf(format string, a ...interface{}) {
	log.Printf("[SERVER]: DeerPop [t=%.1f]: %v", dp.T, fmt.Sprintf(format, a...))
}

// CreateBasicDeerPopulationModel returns a basic population model based on dP/dt = k(N-y) model. k = growth coeff., N = max deer (constants).
func createBasicDeerPopulationModel(dhConf config.DeerHuntConfig) DeerPopulationModel {
	maxDeer := dhConf.MaxDeerPopulation
	deerPopulationGrowth := func(t, y float64) float64 {
		return dhConf.DeerGrowthCoefficient * (float64(maxDeer) - y) // DE of form dy/dt = k(N-y) where k, N are constants
	}
	dp := DeerPopulationModel{
		deProblem:  simulation.ODEProblem{YPrime: deerPopulationGrowth, Y0: float64(maxDeer), T0: 0, DtStep: 0.1},
		Population: float64(maxDeer),
		T:          .0,
	}
	dp.deState = dp.deProblem.StepDeltaY() // initialise DE state
	return dp
}

// Simulate method simulates the reaction of a deer pop. over i=len(deerConsumption) days where [0, maxDeer] are hunted each day i.
// Note: if only simulating for one turn ('step'), len(deerConsumption) = 1
func (dp DeerPopulationModel) Simulate(deerConsumption []int) DeerPopulationModel {
	deStep := dp.deState // this is gets a function from deState closure containing latest state. Allows us to step next solution.

	for i := 0; i < len(deerConsumption); i++ { // note: can use DE.SolveUntilT(10) but in this case we want access to y, t at each iteration
		y0 := dp.Population - float64(deerConsumption[i])
		t, y := deStep(float64(-deerConsumption[i]))

		dp.Population, dp.T = y, t
		dp.Logf("P(t): %.2f. \tPopulation after %v deer hunted: %v, \tpopulation at end of turn (after regeneration): %v\n", y, deerConsumption[i], int(y0), int(y))
	}
	return dp
}

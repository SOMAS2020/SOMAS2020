package foraging

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/simulation"
)

// DeerPopulationModel encapsulates a deer population over time (governed by a predefined DE)
type DeerPopulationModel struct {
	deProblem  simulation.ODEProblem // defintion of DE governing rate of change of deeer pop.
	population float64               // current number of deer in env
	t          float64               // temporal parameter. Time, turn or whatever other incarnation
}

// Logf is the client's logger that prepends logs with your ID. This makes
// it easier to read logs. DO NOT use other loggers that will mess logs up!
// BASE: Do not overwrite in team client.
func (dp *DeerPopulationModel) Logf(format string, a ...interface{}) {
	log.Printf("DeerPop [t=%.1f]: %v", dp.t, fmt.Sprintf(format, a...))
}

// CreateBasicDeerPopulationModel returns a basic population model based on dP/dt = k(N-y) model. k = growth coeff., N = max deer (constants).
func createBasicDeerPopulationModel() DeerPopulationModel {
	fConf := config.GameConfig().ForagingConfig
	// definition of deer pop. gradient. Provides dy/dt given y, t.
	maxDeer := fConf.MaxDeerPopulation
	deerPopulationGrowth := func(t, y float64) float64 {
		return fConf.DeerGrowthCoefficient * (float64(maxDeer) - y) // DE of form dy/dt = k(N-y) where k, N are constants
	}
	return DeerPopulationModel{
		deProblem:  simulation.ODEProblem{YPrime: deerPopulationGrowth, Y0: float64(maxDeer), T0: 0, DtStep: 0.1},
		population: float64(maxDeer),
		t:          .0,
	}
}

// Simulate method simulates the reaction of a deer pop. over i=len(deerConsumption) days where [0, maxDeer] are hunted each day i.
// Note: if only simulating for one turn ('step'), len(deerConsumption) = 1
func (dp *DeerPopulationModel) Simulate(deerConsumption []int) {
	t := dp.t
	y := dp.population
	deStep := dp.deProblem.StepDeltaY()
	for i := 0; i < len(deerConsumption); i++ { // note: can use DE.SolveUntilT(10) but in this case we want access to y, t at each iteration
		y0 := y - float64(deerConsumption[i])
		t, y = deStep(float64(-deerConsumption[i])) // this will update population, t in receiver
		dp.Logf("Day %v - P(t): %.2f, \tdeer after hunt %v, \tdeer end of day: %v, \tdeer hunted: %v\n", t, y, int(y0), int(y), deerConsumption[i])
	}
}

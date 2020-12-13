package foraging

import "github.com/SOMAS2020/SOMAS2020/internal/common/simulation"

var maxDeer = 4.0 // TODO: move to central config

// Population DE. Deer population (y) grows at a rate inverse prop. to pop size. Returns dy/dt given y, t.
func deerPopulationGrad(t, y float64) float64 {
	return 0.1 * (maxDeer - y)
}

// SimulateDeerPopulation simulates the natural growth of a deer population according to a predefined DE
func SimulateDeerPopulation() {
	prob := simulation.ODEProblem{YPrime: deerPopulationGrad, Y0: 0, T0: 0, DtStep: 0.1}
	prob.SolveUntilT(20) // solve for x days
}

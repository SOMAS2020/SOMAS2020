package simulation

import (
	"fmt"
	"math"
	"testing"
)

// helper type to group f(x) and f'(x) for various f(x)
type testFuncPair struct {
	y     yActual
	dydt  ypFunc
	label string
}

var testPairs = []testFuncPair{
	{
		y:     func(t float64) float64 { return 2 * t },
		dydt:  func(t, y float64) float64 { return 2 },
		label: "linear",
	},
	{
		y:     func(t float64) float64 { return 3 * t * t },
		dydt:  func(t, y float64) float64 { return 6 * t },
		label: "quadratic",
	},
	{
		y:     func(t float64) float64 { return 0.5 * math.Pow(t, 3) },
		dydt:  func(t, y float64) float64 { return 1.5 * math.Pow(t, 2) },
		label: "cubic",
	},
	{
		y:     func(t float64) float64 { return math.Sin(t) },
		dydt:  func(t, y float64) float64 { return math.Cos(t) },
		label: "sinusoidal",
	},
	{
		y:     func(t float64) float64 { return math.Atan(t) },
		dydt:  func(t, y float64) float64 { return 1 / (1 + t*t) },
		label: "arctan",
	},
	{
		y:     func(t float64) float64 { return math.Tanh(t) },
		dydt:  func(t, y float64) float64 { return 1 - math.Pow((math.Tanh(t)), 2) },
		label: "hyperbolic tan",
	},
}

func TestSolveUntilT(t *testing.T) {

	maxErrorPerc := 0.001 // max error tolerance between calculated and actual values
	nTrialsPerPair := 10

	for _, tp := range testPairs {
		prob := ODEProblem{
			YPrime: tp.dydt,
			T0:     0,
			Y0:     0,
			DtStep: 0.1,
		}

		var testCases = []struct {
			yActual, yCalc float64
		}{}

		for i, ycalc := range prob.SolveUntilT(nTrialsPerPair) {
			testCases = append(testCases, struct {
				yActual float64
				yCalc   float64
			}{tp.y(float64(i + 1)), ycalc})
		}

		for _, tc := range testCases {
			testname := fmt.Sprintf("%v test", tp.label)
			t.Run(testname, func(t *testing.T) {
				t.Logf("Testing f(x)/f'(x) pair with label %v", tp.label)

				if (1 - tc.yActual/tc.yCalc) > maxErrorPerc/100 {
					t.Errorf("got %.4f, want %.4f", tc.yCalc, tc.yActual)
				}
			})
		}

	}
}

// test stepwise solver
func TestStepDeltaY(t *testing.T) {

	maxErrorPerc := 0.001 // max error tolerance between calculated and actual values
	nTrialsPerPair := 10

	for _, tp := range testPairs {
		prob := ODEProblem{
			YPrime: tp.dydt,
			T0:     0,
			Y0:     0,
			DtStep: 0.1,
		}

		var testCases = []struct {
			yActual, yCalc float64
		}{}

		deStep := prob.StepDeltaY()

		for i := 0; i < nTrialsPerPair; i++ {
			_, yCalc := deStep(0.0)
			testCases = append(testCases, struct {
				yActual float64
				yCalc   float64
			}{tp.y(float64(i + 1)), yCalc})

		}
		for _, tc := range testCases {
			testname := fmt.Sprintf("%v test", tp.label)
			t.Run(testname, func(t *testing.T) {
				t.Logf("Testing f(x)/f'(x) pair with label %v", tp.label)

				if (1 - tc.yActual/tc.yCalc) > maxErrorPerc/100 {
					t.Errorf("got %.4f, want %.4f", tc.yCalc, tc.yActual)
				}
			})
		}

	}
}

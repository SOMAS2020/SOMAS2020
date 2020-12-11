package simulation

import (
	"fmt"
	"math"
)

type ypFunc func(t, y float64) float64
type ypStepFunc func(t, y, dt float64) float64

// newRKStep takes a function representing an ODE
// and returns a function that performs a single step of the 4th order
// Runge-Kutta method. See https://en.wikipedia.org/wiki/Runge–Kutta_methods
func newRK4Step(yp ypFunc) ypStepFunc {
	return func(t, y, dt float64) float64 {
		dy1 := dt * yp(t, y)
		dy2 := dt * yp(t+dt/2, y+dy1/2)
		dy3 := dt * yp(t+dt/2, y+dy2/2)
		dy4 := dt * yp(t+dt, y+dy3)
		return y + (dy1+2*(dy2+dy3)+dy4)/6
	}
}

// differential equation func
func yprime(t, y float64) float64 {
	return 2 * t //t * math.Sqrt(y)
}

// exact solution for comparison
func actual(t float64) float64 {
	return t * t
}
func printErr(t, y float64) {
	fmt.Printf("y(%.1f) = %f Error: %e\n", t, y, math.Abs(actual(t)-y))
}

// SolveRK4 provides an approximate ODE solution using the 4th order Runge-Kutta method
func SolveRK4(yprime ypFunc, t0 int, tFinal int, y0 float64) {
	dtPrint := 1 // and to print at whole numbers.
	dtStep := .1 // step value.

	t, y := float64(t0), y0
	ypStep := newRK4Step(yprime)
	for t1 := t0 + dtPrint; t1 <= tFinal; t1 += dtPrint {
		printErr(t, y) // print intermediate result
		for steps := int(float64(dtPrint)/dtStep + .5); steps > 1; steps-- {
			y = ypStep(t, y, dtStep)
			t += dtStep
		}
		y = ypStep(t, y, float64(t1)-t) // adjust step to integer time
		t = float64(t1)
	}
	printErr(t, y) // print final result
}

// TestSolve solves a sample DE
func TestSolve() {
	SolveRK4(yprime, 0, 10, 0)
}

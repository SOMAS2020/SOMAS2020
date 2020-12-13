package simulation

import (
	"fmt"
	"math"
)

type ypFunc func(t, y float64) float64         // dy/dt definition
type ypStepFunc func(t, y, dt float64) float64 // step function for computing dy/dt

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

// SolveStep performs an update towards solving a DE using RK4. Performs approx. 1/dtStep iterations in a step.
func solveStep(t, y, dtStep float64, yPrime ypFunc) (t2, y2 float64) {
	ypStep := newRK4Step(yPrime)

	for steps := int(1 / dtStep); steps > 1; steps-- {
		y = ypStep(t, y, dtStep)
		t += dtStep
	}
	tUp := math.Ceil(t)
	dtFinal := float64(tUp) - t // perform last step over remaining dt to finish integer time increment
	y = ypStep(t, y, dtFinal)
	return tUp, y
}

// ODEProblem is simply a struct to initialise the parts of and ODE IVP
type ODEProblem struct {
	YPrime ypFunc  // func that returns dy/dt
	T0     int     // initial integer timestep
	Y0     float64 // initial y value for IVP
	DtStep float64 // solver step size. Typically 0.1
}

// Step is a closure that initialises t, y and returns a function that allows you to perform a solution
// step in the DE solution. Performs approx. 1/dtStep iterations in a step.
func (de ODEProblem) Step() func() (t2, y2 float64) {
	t, y := float64(de.T0), de.Y0

	return func() (t2, y2 float64) {
		t, y = solveStep(t, y, de.DtStep, de.YPrime)
		return t, y
	}
}

// StepDeltaY is the same as Step but allows for manipulation of y by providing a deltaY arg. Internal y
// will then be modified: y -> y+deltaY. Useful for modifiying y between steps to model external disturbances (e.g. resource consumption)
func (de ODEProblem) StepDeltaY() func(float64) (t2, y2 float64) {
	t, y := float64(de.T0), de.Y0

	return func(deltaY float64) (t2, y2 float64) {
		t, y = solveStep(t, y+deltaY, de.DtStep, de.YPrime) // adjust y by deltaY
		return t, y
	}
}

// SolveUntilT solves a DE from T0 (from initialisation) to tFinal
func (de ODEProblem) SolveUntilT(tFinal int) {
	dtPrint := 1 // and to print at whole numbers.
	t, y := float64(de.T0), de.Y0
	for t1 := de.T0 + dtPrint; t1 <= tFinal; t1 += dtPrint {
		t, y = solveStep(t, y, de.DtStep, de.YPrime)
		fmt.Println(t, y)
	}
}

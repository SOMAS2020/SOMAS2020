package simulation

import (
	"fmt"
	"math"
)

type yActual func(t float64) float64 // func to provide actual solution y(t) if known for comparisn

// Utility func to compare error of numerical solution to actual solution at t
func printErr(t, y float64, yActual yActual) {
	fmt.Printf("y(%.1f) = %f Error: %e\n", t, y, math.Abs(yActual(t)-y))
}

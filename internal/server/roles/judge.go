package roles

import (
    "fmt"
)

type Judge Interface {
	withdrawPresidentSalary()
	payPresident()
	inspectBallot()
	inspectAllocation()
	declareSpeakerPerformance()
	declarePresidentPerformance()
}
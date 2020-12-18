package roles

// Judge is an interface that is implemented by BaseJudge but can also be
// optionally implemented by individual islands.
type Judge interface {
	withdrawPresidentSalary() error
	payPresident() error
	inspectBallot() error
	inspectAllocation() error
	declareSpeakerPerformance() error
	declarePresidentPerformance() error
}
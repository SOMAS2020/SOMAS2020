package roles

// Judge is an interface that is implemented by baseJudge but can also be
// optionally implemented by individual islands.
type Judge interface {
	payPresident()
	inspectHistory() (map[int]EvaluationReturn, error)
	declareSpeakerPerformance() (int, bool, int, bool, error)
	declarePresidentPerformance() (int, bool, int, bool, error)
}

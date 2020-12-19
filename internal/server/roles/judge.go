package roles

// Judge is an interface that is implemented by baseJudge but can also be
// optionally implemented by individual islands.
type Judge interface {
	payPresident() error
	inspectHistory()
	declareSpeakerPerformance() (int, bool, int, bool)
	declarePresidentPerformance() (int, bool, int, bool)
}

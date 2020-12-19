package roles

// Speaker is an interface that is implemented by baseSpeaker but can also be
// optionally implemented by individual islands.
type Speaker interface {
	PayJudge() error
	RunVote() error
	UpdateRules() error
}

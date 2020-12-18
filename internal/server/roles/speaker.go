package roles

// Speaker is an interface that is implemented by BaseSpeaker but can also be
// optionally implemented by individual islands.
type Speaker Interface{
	PayJudge() error
	RunVote() error
	UpdateRules() error
}
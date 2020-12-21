package roles

// Speaker is an interface that is implemented by baseSpeaker but can also be
// optionally implemented by individual islands.
type Speaker interface {
	payJudge() (int, error)
	RunVote(string) (bool, error)
	DecideAnnouncement(string, bool) (string, bool, error)
}

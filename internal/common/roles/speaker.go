package roles

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// Speaker is an interface that is implemented by baseSpeaker but can also be
// optionally implemented by individual islands.
type Speaker interface {
	PayJudge() (shared.Resources, error)
	RunVote(string) (string, error)
	DecideAnnouncement(string, bool) (string, bool, error)
}

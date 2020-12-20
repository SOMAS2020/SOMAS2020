package roles

import "github.com/SOMAS2020/SOMAS2020/internal/common"

// Speaker is an interface that is implemented by baseSpeaker but can also be
// optionally implemented by individual islands.
type Speaker interface {
	payJudge(common.GameState)
	RunVote() error
	UpdateRules() error
	DeclareResult() error
}

package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Speaker is an interface that is implemented by baseSpeaker but can also be
// optionally implemented by individual islands.
type Speaker interface {
	PayJudge(judgeSalary shared.Resources) shared.SpeakerReturnContent
	DecideAgenda(rules.RuleMatrix) shared.SpeakerReturnContent
	DecideVote(rules.RuleMatrix, []shared.ClientID) shared.SpeakerReturnContent
	DecideAnnouncement(rules.RuleMatrix, bool) shared.SpeakerReturnContent
	CallJudgeElection(shared.MonitorResult, int, []shared.ClientID) shared.ElectionSettings
	DecideNextJudge(shared.ClientID) shared.ClientID
}

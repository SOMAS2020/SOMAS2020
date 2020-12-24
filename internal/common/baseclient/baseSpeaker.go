package baseclient

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type baseSpeaker struct {
}

func (s *baseSpeaker) PayJudge(salary shared.Resources) (shared.Resources, error) {
	return salary, nil
}

func (s *baseSpeaker) RunVote(ruleID string) (string, error) {
	return ruleID, nil
}

//Return "", _ for no announcement to occur, TODO: Change to error type for no announcement
func (s *baseSpeaker) DecideAnnouncement(ruleId string, result bool) (string, bool, error) {
	return ruleId, result, nil
}

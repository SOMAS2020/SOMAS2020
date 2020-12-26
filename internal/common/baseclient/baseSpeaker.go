package baseclient

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type BaseSpeaker struct {
}

func (s *BaseSpeaker) PayJudge(salary shared.Resources) (shared.Resources, error) {
	return salary, nil
}

//DecideAgenda the interface implementation and example of a well behaved Speaker
//who sets the vote to be voted on to be the rule the President provided
func (s *BaseSpeaker) DecideAgenda(ruleID string) (string, bool) {
	return ruleID, true
}

//DecideVote is the interface implementation and example of a well behaved Speaker
//who calls a vote on the proposed rule and asks all available islands to vote.
//Return an empty string or empty []shared.ClientID for no vote to occur
func (s *BaseSpeaker) DecideVote(ruleID string, aliveClients []shared.ClientID) (string, []shared.ClientID, error) {
	//TODO: disregard islands with sanctions
	return ruleID, aliveClients, nil
}

//DecideAnnouncement is the interface implementation and example of a well behaved Speaker
//A well behaved speaker announces what had been voted on and the corresponding result
//Return "", _ for no announcement to occur
func (s *BaseSpeaker) DecideAnnouncement(ruleId string, result bool) (string, bool, error) {
	return ruleId, result, nil
}

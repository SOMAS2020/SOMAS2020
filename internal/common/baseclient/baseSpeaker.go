package baseclient

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type BaseSpeaker struct {
}

// PayJudge is used for paying judge for his service
func (s *BaseSpeaker) PayJudge(salary shared.Resources) shared.SpeakerReturnContent {
	return shared.SpeakerReturnContent{
		ContentType: shared.SpeakerJudgeSalary,
		JudgeSalary: salary,
		ActionTaken: true,
	}
}

//DecideAgenda the interface implementation and example of a well behaved Speaker
//who sets the vote to be voted on to be the rule the President provided
func (s *BaseSpeaker) DecideAgenda(ruleID string) shared.SpeakerReturnContent {
	return shared.SpeakerReturnContent{
		ContentType: shared.SpeakerAgenda,
		RuleID:      ruleID,
		ActionTaken: true,
	}
}

//DecideVote is the interface implementation and example of a well behaved Speaker
//who calls a vote on the proposed rule and asks all available islands to vote.
//Return an empty string or empty []shared.ClientID for no vote to occur
func (s *BaseSpeaker) DecideVote(ruleID string, aliveClients []shared.ClientID) shared.SpeakerReturnContent {
	//TODO: disregard islands with sanctions
	return shared.SpeakerReturnContent{
		ContentType:          shared.SpeakerVote,
		ParticipatingIslands: aliveClients,
		RuleID:               ruleID,
		ActionTaken:          true,
	}
}

//DecideAnnouncement is the interface implementation and example of a well behaved Speaker
//A well behaved speaker announces what had been voted on and the corresponding result
//Return "", _ for no announcement to occur
func (s *BaseSpeaker) DecideAnnouncement(ruleID string, result bool) shared.SpeakerReturnContent {
	return shared.SpeakerReturnContent{
		ContentType:  shared.SpeakerAnnouncement,
		RuleID:       ruleID,
		VotingResult: result,
		ActionTaken:  true,
	}
}

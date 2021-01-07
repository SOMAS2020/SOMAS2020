package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type speaker struct {
	*baseclient.BaseSpeaker
	*client
}

/*
func (s *speaker) PayJudge(salary shared.Resources) shared.SpeakerReturnContent {
	return s.BaseSpeaker.PayJudge(salary)
}
*/
func (s *speaker) DecideAgenda(ruleID rules.RuleMatrix) shared.SpeakerReturnContent {
	return s.BaseSpeaker.DecideAgenda(ruleID)
}

func (s *speaker) DecideVote(ruleID rules.RuleMatrix, aliveClients []shared.ClientID) shared.SpeakerReturnContent {
	return s.BaseSpeaker.DecideVote(ruleID, aliveClients)
}

func (s *speaker) DecideAnnouncement(ruleID rules.RuleMatrix, result bool) shared.SpeakerReturnContent {
	return s.BaseSpeaker.DecideAnnouncement(ruleID, result)
}

func (s *speaker) CallJudgeElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.BordaCount,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}
	if monitoring.Performed && !monitoring.Result {
		electionsettings.HoldElection = true
	}
	if turnsInPower >= 2 {
		electionsettings.HoldElection = true
	}
	return electionsettings
}

func (s *speaker) DecideNextJudge(winner shared.ClientID) shared.ClientID {
	return winner
}

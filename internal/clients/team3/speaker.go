package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type speaker struct {
	// Base implementation
	*baseclient.BaseSpeaker
	// Our client
	c *client
}

// Override functions here, see president.go for examples

func (s *speaker) PayJudge(salary shared.Resources) shared.SpeakerReturnContent {
	// Use the base implementation
	return s.BaseSpeaker.PayJudge(salary)
}

func (s *speaker) DecideAgenda(ruleID string) shared.SpeakerReturnContent {
	return s.BaseSpeaker.DecideAgenda(ruleID)
}

func (s *speaker) DecideVote(ruleID string, aliveClients []shared.ClientID) shared.SpeakerReturnContent {
	var chosenClients []shared.ClientID
	for _, islandID := range aliveClients {
		if s.c.iigoInfo.sanctions.islandSanctions[islandID] != roles.NoSanction {
			chosenClients = append(chosenClients, islandID)
		}
	}
	if s.c.shouldICheat() {
		for _, islandID := range aliveClients {
			if s.c.trustScore[islandID] > 0.5 {
				chosenClients = append(chosenClients, islandID)
			}
		}
	}

	return shared.SpeakerReturnContent{
		ContentType:          shared.SpeakerVote,
		ParticipatingIslands: chosenClients,
		RuleID:               ruleID,
		ActionTaken:          true,
	}

}

func (s *speaker) DecideAnnouncement(ruleID string, result bool) shared.SpeakerReturnContent {
	if s.c.shouldICheat() {
		result = s.c.GetVoteForRule(ruleID)
	}

	return shared.SpeakerReturnContent{
		ContentType:  shared.SpeakerAnnouncement,
		RuleID:       ruleID,
		VotingResult: result,
		ActionTaken:  true,
	}

}

func (s *speaker) CallJudgeElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	return s.BaseSpeaker.CallJudgeElection(monitoring, turnsInPower, allIslands)
}

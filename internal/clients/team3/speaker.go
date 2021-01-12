package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type speaker struct {
	// Base implementation
	*baseclient.BaseSpeaker
	// Our client
	c *client
}

// Override functions here, see president.go for examples

func (s *speaker) PayJudge() shared.SpeakerReturnContent {
	// Use the base implementation
	return s.BaseSpeaker.PayJudge()
}

func (s *speaker) DecideAgenda(ruleMat rules.RuleMatrix) shared.SpeakerReturnContent {
	return s.BaseSpeaker.DecideAgenda(ruleMat)
}

func (s *speaker) DecideVote(ruleMatrix rules.RuleMatrix, aliveClients []shared.ClientID) shared.SpeakerReturnContent {
	var chosenClients []shared.ClientID
	for _, islandID := range aliveClients {
		if s.c.iigoInfo.sanctions.islandSanctions[islandID] != shared.NoSanction {
			chosenClients = append(chosenClients, islandID)
		}
	}
	/*if s.c.shouldICheat() {
		for _, islandID := range aliveClients {
			if s.c.trustScore[islandID] > 35 {
				chosenClients = append(chosenClients, islandID)
			}
		}
	}*/

	return shared.SpeakerReturnContent{
		ContentType:          shared.SpeakerVote,
		ParticipatingIslands: chosenClients,
		RuleMatrix:           ruleMatrix,
		ActionTaken:          true,
	}

}

func (s *speaker) DecideAnnouncement(ruleMatrix rules.RuleMatrix, result bool) shared.SpeakerReturnContent {

	/*if s.c.shouldICheat() {
		res := s.c.iigoInfo.ruleVotingResults[ruleMatrix.RuleName].ourVote
		if res == shared.Approve {
			result = true
		} else {
			result = false
		}
	}*/

	return shared.SpeakerReturnContent{
		ContentType:  shared.SpeakerAnnouncement,
		RuleMatrix:   ruleMatrix,
		VotingResult: result,
		ActionTaken:  true,
	}

}

func (s *speaker) CallJudgeElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	if s.c.params.adv != nil {
		ret, done := s.c.params.adv.CallJudgeElection(monitoring, turnsInPower, allIslands)
		if done {
			return ret
		}
	}

	return s.BaseSpeaker.CallJudgeElection(monitoring, turnsInPower, allIslands)
}

// DecideNextJudge returns the ID of chosen next Judge
// OPTIONAL: override to manipulate the result of the election
func (s *speaker) DecideNextJudge(winner shared.ClientID) shared.ClientID {
	if s.c.params.adv != nil {
		ret, done := s.c.params.adv.DecideNextJudge(winner)
		if done {
			return ret
		}
	}
	return winner
}

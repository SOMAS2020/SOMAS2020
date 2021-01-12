package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type speaker struct {
	*baseclient.BaseSpeaker
	*client
}

func (s *speaker) CallJudgeElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	currJudgeID := s.client.ServerReadHandle.GetGameState().JudgeID
	otherIslands := []shared.ClientID{}

	for _, team := range allIslands {
		if currJudgeID == team {
			continue
		}

		otherIslands = append(otherIslands, team)
	}

	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.Runoff,
		IslandsToVote: otherIslands,
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
	if winner == s.client.ServerReadHandle.GetGameState().JudgeID {
		return s.client.GetID()
	}

	if s.client.friendship[winner] < s.clientConfig.maxFriendship/1.5 {
		return s.client.GetID()
	}

	return winner
}

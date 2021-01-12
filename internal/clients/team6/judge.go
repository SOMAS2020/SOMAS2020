package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type judge struct {
	*baseclient.BaseJudge
	*client
}

func (j *judge) GetPardonedIslands(currentSanctions map[int][]shared.Sanction) map[int][]bool {
	pardons := map[int][]bool{}
	maxSanctionTime := int(j.client.ServerReadHandle.GetGameConfig().IIGOClientConfig.SanctionCacheDepth - 1)

	for timeStep, sanctions := range currentSanctions {
		pardons[timeStep] = make([]bool, len(sanctions))
		for who, sanction := range sanctions {
			if timeStep == maxSanctionTime && j.client.friendship[sanction.ClientID] == j.client.clientConfig.maxFriendship {
				// we can pardon certain islands having maximum friendship with us
				pardons[timeStep][who] = true
			}
		}
	}

	return pardons
}

func (j *judge) CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	currPresidentID := j.client.ServerReadHandle.GetGameState().PresidentID
	otherIslands := []shared.ClientID{}

	for _, team := range allIslands {
		if currPresidentID == team {
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

func (j *judge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	if winner == j.client.ServerReadHandle.GetGameState().PresidentID {
		return j.client.GetID()
	}

	if j.client.friendship[winner] < j.clientConfig.maxFriendship/1.5 {
		return j.client.GetID()
	}

	return winner
}

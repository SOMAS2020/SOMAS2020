package team2

import (
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	return &c.currSpeaker
}

func (c *client) GetClientJudgePointer() roles.Judge {
	return &c.currJudge
}

func (c *client) GetClientPresidentPointer() roles.President {
	return &c.currPresident
}

func (c *client) VoteForElection(roleToElect shared.Role, candidateList []shared.ClientID) []shared.ClientID {
	var situation Situation
	switch roleToElect {
	case shared.President:
		situation = "President"
	case shared.Judge:
		situation = "Judge"
	default:
		situation = "Gifts"
	}

	var trustRank IslandTrustList
	for _, candidate := range candidateList {
		islandConf := IslandTrust{
			island: candidate,
			trust:  c.confidence(situation, candidate),
		}
		trustRank = append(trustRank, islandConf)
	}

	sort.Sort(trustRank)
	bordaList := make([]shared.ClientID, 0)

	for _, islandTrust := range trustRank {
		bordaList = append(bordaList, islandTrust.island)
	}

	return bordaList
}

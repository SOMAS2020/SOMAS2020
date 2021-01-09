package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Speaker struct {
	*baseclient.BaseSpeaker
	c *client
}

// DecideVote calls a vote on the rule decided on during DecideAgenda
// As speaker we will only count the votes of islands who have not had a snaction in the last 10 turns
func (s *Speaker) DecideVote(ruleMatrix rules.RuleMatrix, aliveClients []shared.ClientID) shared.SpeakerReturnContent {
	var chosenClients []shared.ClientID
	for _, islandID := range aliveClients {
		if s.c.islandSanctions[islandID][len(s.c.islandSanctions[islandID])-1].Turn <= s.c.gameState().Turn-10 {
			chosenClients = append(chosenClients, islandID)
		}
	}

	return shared.SpeakerReturnContent{
		ContentType:          shared.SpeakerVote,
		ParticipatingIslands: chosenClients,
		RuleMatrix:           ruleMatrix,
		ActionTaken:          true,
	}
}

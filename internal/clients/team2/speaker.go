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
// In the first Num turns or if the last sanction is far away, count all the votes
func (s *Speaker) DecideVote(ruleMatrix rules.RuleMatrix, aliveClients []shared.ClientID) shared.SpeakerReturnContent {
	var chosenClients []shared.ClientID
	chosenClients = append(chosenClients, s.c.GetID())

	for _, islandID := range aliveClients {
		// do not add our own island twice
		if islandID == s.c.GetID() {
			continue
		}

		if lastSanctionTurn, ok := s.c.islandSanctions[islandID]; ok {
			if s.c.gameState().Turn <= 10 || lastSanctionTurn.Turn <= s.c.gameState().Turn-10 {
				chosenClients = append(chosenClients, islandID)
			}
		} else {
			chosenClients = append(chosenClients, islandID)
		}

	}

	// chosen client is never null - sneaky fix
	return shared.SpeakerReturnContent{
		ContentType:          shared.SpeakerVote,
		ParticipatingIslands: chosenClients,
		RuleMatrix:           ruleMatrix,
		ActionTaken:          true,
	}
}

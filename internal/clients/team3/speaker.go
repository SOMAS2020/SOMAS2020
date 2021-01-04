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

func (s *speaker) PayJudge(salary shared.Resources) (shared.Resources, bool) {
	// Use the base implementation
	return s.BaseSpeaker.PayJudge(salary)
}

func (s *speaker) DecideVote(ruleID string, aliveClients []shared.ClientID) (string, []shared.ClientID, bool) {
	var chosenClients []shared.ClientID
	for _, islandID := range aliveClients {
		if s.c.iigoInfo.sanctions.islandSanctions[islandID] != roles.NoSanction {
			chosenClients = append(chosenClients, islandID)
		}
	}
	return ruleID, chosenClients, true
}

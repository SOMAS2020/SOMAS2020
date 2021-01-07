// Package team2 contains code for team 2's client implementation
package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	return &c.ourSpeaker
}

func (c *client) GetClientJudgePointer() roles.Judge {
	return &c.ourJudge
}

func (c *client) GetClientPresidentPointer() roles.President {
	return &c.ourPresident
}

//resetIIGOInfo clears the island's information regarding IIGO at start of turn
func (c *client) resetIIGOInfo() {
	c.iigoInfo.ourRole = nil // TODO unused, remove
	c.iigoInfo.commonPoolAllocation = 0
	c.iigoInfo.taxationAmount = 0
	c.iigoInfo.monitoringOutcomes = make(map[shared.Role]bool)
	c.iigoInfo.monitoringDeclared = make(map[shared.Role]bool)
	c.iigoInfo.startOfTurnJudgeID = c.ServerReadHandle.GetGameState().JudgeID
	c.iigoInfo.startOfTurnPresidentID = c.ServerReadHandle.GetGameState().PresidentID
	c.iigoInfo.startOfTurnSpeakerID = c.ServerReadHandle.GetGameState().SpeakerID
	c.iigoInfo.sanctions = &sanctionInfo{
		tierInfo:        make(map[roles.IIGOSanctionTier]roles.IIGOSanctionScore),
		rulePenalties:   make(map[string]roles.IIGOSanctionScore),
		islandSanctions: make(map[shared.ClientID]roles.IIGOSanctionTier),
		ourSanction:     roles.IIGOSanctionScore(0),
	}
	c.iigoInfo.ruleVotingResults = make(map[string]*ruleVoteInfo)
	c.iigoInfo.ourRequest = 0
	c.iigoInfo.ourDeclaredResources = 0
}

func (c *client) getOurRole() string {
	if c.iigoInfo.startOfTurnJudgeID == c.BaseClient.GetID() {
		return "Judge"
	}
	if c.iigoInfo.startOfTurnPresidentID == c.BaseClient.GetID() {
		return "President"
	}
	if c.iigoInfo.startOfTurnSpeakerID == c.BaseClient.GetID() {
		return "Speaker"
	}
	return "None"
}

func (c *client) GetSanctionPayment() shared.Resources {
	if value, ok := c.LocalVariableCache[rules.SanctionExpected]; ok {
		idealVal, available := c.locationService.switchDetermineFunction(rules.SanctionPaid, value.Values)
		if available {
			variablesChanged := map[rules.VariableFieldName]rules.VariableValuePair{
				rules.SanctionPaid: {
					rules.SanctionPaid,
					idealVal,
				},
				rules.SanctionExpected: {
					rules.SanctionExpected,
					c.LocalVariableCache[rules.SanctionExpected].Values,
				},
			}

			recommendedValues := c.dynamicAssistedResult(variablesChanged)
			if c.params.complianceLevel > 80 {
				return shared.Resources(recommendedValues[rules.SanctionPaid].Values[rules.SingleValueVariableEntry])
			}
			return shared.Resources(idealVal[rules.SingleValueVariableEntry])
		}
		return shared.Resources(value.Values[rules.SingleValueVariableEntry])
	}
	return 0
}

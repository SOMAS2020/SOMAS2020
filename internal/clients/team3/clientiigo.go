package team3

import (
	// "github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	// "github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

/*
	//IIGO: COMPULSORY
	MonitorIIGORole(shared.Role) bool
	DecideIIGOMonitoringAnnouncement(bool) (bool, bool)

	GetVoteForRule(ruleName string) bool
	GetVoteForElection(roleToElect shared.Role) []shared.ClientID

	CommonPoolResourceRequest() shared.Resources
	ResourceReport() shared.Resources
	RuleProposal() string
	GetClientPresidentPointer() roles.President
	GetClientJudgePointer() roles.Judge
	GetClientSpeakerPointer() roles.Speaker
	TaxTaken(shared.Resources)
	GetTaxContribution() shared.Resources
	RequestAllocation() shared.Resources
*/

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	// c.Logf("became speaker")
	role := shared.Speaker
	c.iigoInfo.ourRole = &role
	return &speaker{c: c}
}

func (c *client) GetClientJudgePointer() roles.Judge {
	// c.Logf("became judge")
	role := shared.Judge
	c.iigoInfo.ourRole = &role
	return &judge{c: c}
}

func (c *client) GetClientPresidentPointer() roles.President {
	// c.Logf("became president")
	role := shared.President
	c.iigoInfo.ourRole = &role
	return &president{c: c}
}

//resetIIGOInfo clears the island's information regarding IIGO at start of turn
func (c *client) resetIIGOInfo() {
	c.iigoInfo.ourRole = nil
	c.iigoInfo.commonPoolAllocation = 0
	c.iigoInfo.taxationAmount = 0
	c.iigoInfo.ruleVotingResults = make(map[string]bool)
	c.iigoInfo.ruleVotingResultAnnounced = make(map[string]bool)
	c.iigoInfo.monitoringOutcomes = make(map[shared.Role]bool)
	c.iigoInfo.monitoringDeclared = make(map[shared.Role]bool)
	c.iigoInfo.startOfTurnJudgeID = c.ServerReadHandle.GetGameState().JudgeID
	c.iigoInfo.startOfTurnPresidentID = c.ServerReadHandle.GetGameState().PresidentID
	c.iigoInfo.startOfTurnSpeakerID = c.ServerReadHandle.GetGameState().SpeakerID
}

func (c *client) getOurRole() string {
	if c.iigoInfo.startOfTurnJudgeID == shared.ClientID(3) {
		return "Judge"
	}
	if c.iigoInfo.startOfTurnPresidentID == shared.ClientID(3) {
		return "President"
	}
	if c.iigoInfo.startOfTurnSpeakerID == shared.ClientID(3) {
		return "Speaker"
	}
	return "None"
}

// ReceiveCommunication is a function called by IIGO to pass the communication sent to the client.
// This function is overridden to receive information and update local info accordingly.
func (c *client) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	c.Communications[sender] = append(c.Communications[sender], data)

	for contentType, content := range data {
		switch contentType {
		case shared.TaxAmount:
			c.iigoInfo.taxationAmount = shared.Resources(content.IntegerData)
		case shared.AllocationAmount:
			c.iigoInfo.commonPoolAllocation = shared.Resources(content.IntegerData)
		case shared.RuleName:
			currentRuleID := content.TextData
			c.iigoInfo.ruleVotingResultAnnounced[currentRuleID] = true
			c.iigoInfo.ruleVotingResults[currentRuleID] = data[shared.RuleVoteResult].BooleanData
		case shared.RoleMonitored:
			c.iigoInfo.monitoringDeclared[content.IIGORole] = true
			c.iigoInfo.monitoringOutcomes[content.IIGORole] = data[shared.MonitoringResult].BooleanData
		}
	}
}

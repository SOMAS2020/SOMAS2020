package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) GetClientPresidentPointer() roles.President {
	return &president{client: c}
}

func (c *client) GetClientJudgePointer() roles.Judge {
	return &judge{client: c}
}

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	return &speaker{client: c}
}

// ------ TODO: COMPULSORY ------
func (c *client) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	c.BaseClient.ReceiveCommunication(sender, data)
}

// ------ TODO: COMPULSORY -----
func (c *client) MonitorIIGORole(roleName shared.Role) bool {
	return c.BaseClient.MonitorIIGORole(roleName)
}

// ------ TODO: COMPULSORY -----
func (c *client) DecideIIGOMonitoringAnnouncement(monitoringResult bool) (resultToShare bool, announce bool) {
	return c.BaseClient.DecideIIGOMonitoringAnnouncement(monitoringResult)
}

func (c *client) CommonPoolResourceRequest() shared.Resources {

}

func (c *client) ResourceReport() shared.ResourcesReport {
	// if we are selfish, will report 1/2 of the actual resources
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	fakeReport := shared.ResourcesReport{
		ReportedAmount: (1 / 2) * ourResources,
		Reported:       true,
	}

	if c.getPersonality() == Selfish {
		return fakeReport
	}

	return c.BaseClient.ResourceReport()
}

// ------ TODO: COMPULSORY -----
func (c *client) RuleProposal() rules.RuleMatrix {
	return c.BaseClient.RuleProposal()
}

// ------ TODO: COMPULSORY -----
func (c *client) GetTaxContribution() shared.Resources {
	return c.BaseClient.GetTaxContribution()
}

func (c *client) GetSanctionPayment() shared.Resources {
	return c.BaseClient.GetSanctionPayment()
}

func (c *client) RequestAllocation() shared.Resources {
	return c.BaseClient.RequestAllocation()
}

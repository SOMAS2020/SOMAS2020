package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
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

// ------ TODO: COMPULSORY -----
func (c *client) CommonPoolResourceRequest() shared.Resources {
	return c.BaseClient.CommonPoolResourceRequest()
}

// ------ TODO: COMPULSORY -----
func (c *client) ResourceReport() shared.ResourcesReport {
	return c.BaseClient.ResourceReport()
}

// ------ TODO: COMPULSORY -----
func (c *client) RuleProposal() string {
	return c.BaseClient.RuleProposal()
}

// ------ TODO: COMPULSORY -----
func (c *client) GetTaxContribution() shared.Resources {
	return c.BaseClient.GetTaxContribution()
}

// ------ TODO: COMPULSORY -----
func (c *client) GetSanctionPayment() shared.Resources {
	return c.BaseClient.GetSanctionPayment()
}

func (c *client) RequestAllocation() shared.Resources {
	return c.BaseClient.RequestAllocation()
}

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

// func (c *client) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
// 	for fieldName, content := range data {
// 		switch fieldName {
// 		case shared.TaxAmount:
// 			c.config.payingTax = shared.Resources(content.IntegerData)
// 		} //add sth else
// 	}
// }

func (c *client) MonitorIIGORole(roleName shared.Role) bool {
	return false
}

func (c *client) DecideIIGOMonitoringAnnouncement(monitoringResult bool) (resultToShare bool, announce bool) {
	resultToShare = monitoringResult
	announce = true
	return
}

func (c *client) CommonPoolResourceRequest() shared.Resources {
	minThreshold := c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold
	ownResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	if ownResources > minThreshold { //if current resource > threshold, our agent skip to request resource from common pool
		return 0
	}
	return minThreshold - ownResources
}

func (c *client) RequestAllocation() shared.Resources {
	//we will take 10% of the common pool when we are critical or dying
	ourStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	if ourStatus == shared.Critical || ourStatus == shared.Dead {
		return c.ServerReadHandle.GetGameState().CommonPool / 10
	}
	return 0
}

//Optional
func (c *client) ResourceReport() shared.ResourcesReport {
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	return shared.ResourcesReport{
		ReportedAmount: ourResources,
		Reported:       true,
	}
}

// ------ TODO: COMPULSORY -----
// func (c *client) RuleProposal() string {
// 	return c.BaseClient.RuleProposal()
// }

func (c *client) GetTaxContribution() shared.Resources {
	ourPersonality := c.getPersonality()
	if ourPersonality == Selfish { //evade tax when we are selfish
		return 0
	}
	return c.clientConfig.payingTax
}

// ------ TODO: COMPULSORY -----
func (c *client) GetSanctionPayment() shared.Resources {
	return 0
}

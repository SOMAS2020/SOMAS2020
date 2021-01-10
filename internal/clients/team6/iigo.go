package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) GetClientPresidentPointer() roles.President {
	return &president{client: c, BasePresident: &baseclient.BasePresident{}}
}

func (c *client) GetClientJudgePointer() roles.Judge {
	return &judge{client: c, BaseJudge: &baseclient.BaseJudge{GameState: c.ServerReadHandle.GetGameState()}}
}

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	return &speaker{client: c, BaseSpeaker: &baseclient.BaseSpeaker{}}
}

func (c *client) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	for fieldName, content := range data {
		switch fieldName {
		case shared.IIGOTaxDecision:
			c.payingTax = shared.Resources(content.IIGOValueData.Amount)
		case shared.IIGOAllocationDecision:
			c.payingSanction = shared.Resources(content.IIGOValueData.Amount)
			//add sth else
		default:
		}

	}
}

func (c *client) MonitorIIGORole(roleName shared.Role) bool {
	return true
}

func (c *client) DecideIIGOMonitoringAnnouncement(monitoringResult bool) (resultToShare bool, announce bool) {
	resultToShare = monitoringResult
	announce = true
	return
}

func (c *client) CommonPoolResourceRequest() shared.Resources {
	var reqResource shared.Resources = 0
	minThreshold := c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold
	ownResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	if ownResources > minThreshold { //if current resource > threshold, our agent skip to request resource from common pool
		reqResource = minThreshold - ownResources
	}
	c.Logf("Request %v from common pool", reqResource)
	return reqResource
}

func (c *client) RequestAllocation() shared.Resources {
	numberAlive := shared.Resources(c.getNumOfAliveIslands())
	ourStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	if ourStatus == shared.Critical || ourStatus == shared.Dead {
		if numberAlive > 0 {
			takenResource := c.ServerReadHandle.GetGameState().CommonPool / numberAlive
			c.Logf("Taken %v from common pool", takenResource)
			return takenResource
		}
	}
	return 0
}

func (c *client) ResourceReport() shared.ResourcesReport {
	// if we are selfish, will report 1/2 of the actual resources
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	ourPersonality := c.getPersonality()
	fakeReport := shared.ResourcesReport{
		ReportedAmount: shared.Resources(float64(1/2) * float64(ourResources)),
		Reported:       true,
	}

	if ourPersonality == Selfish {
		return fakeReport
	}

	return shared.ResourcesReport{
		ReportedAmount: ourResources,
		Reported:       true,
	}
}

func (c *client) GetTaxContribution() shared.Resources {
	var total shared.Resources = 0
	prediction, ok := c.disasterPredictions[c.GetID()]
	if ok && prediction.TimeLeft == 1 {
		total = shared.Resources(prediction.Magnitude * prediction.Confidence / 100)
	}

	ourPersonality := c.getPersonality()
	if ourPersonality == Selfish {
		if total > c.payingTax {
			return c.payingTax
		}
		return total
	}
	if total > c.payingTax {
		return total
	}
	return c.payingTax
}

func (c *client) GetSanctionPayment() shared.Resources {
	return c.payingSanction
}

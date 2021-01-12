package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) GetClientPresidentPointer() roles.President {
	return &president{client: c, BasePresident: &baseclient.BasePresident{GameState: c.ServerReadHandle.GetGameState()}}
}

func (c *client) GetClientJudgePointer() roles.Judge {
	return &judge{client: c, BaseJudge: &baseclient.BaseJudge{GameState: c.ServerReadHandle.GetGameState()}}
}

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	return &speaker{client: c, BaseSpeaker: &baseclient.BaseSpeaker{GameState: c.ServerReadHandle.GetGameState()}}
}

func (c *client) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	for fieldName, content := range data {
		switch fieldName {
		case shared.IIGOTaxDecision:
			c.payingTax = shared.Resources(content.IIGOValueData.Amount)
		case shared.SanctionAmount:
			c.payingSanction = shared.Resources(content.IntegerData)
		case shared.IIGOAllocationDecision:
			c.allocation = shared.Resources(content.IIGOValueData.Amount)
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
	ourPersonality := c.getPersonality()
	numberAlive := shared.Resources(c.getNumOfAliveIslands())
	livingCost := c.ServerReadHandle.GetGameConfig().CostOfLiving
	minThreshold := c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold
	ownResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources

	cprLeft := c.ServerReadHandle.GetGameState().CommonPool
	//when common pool does not have enough resource, will not request
	if cprLeft < numberAlive*livingCost {
		return 0
	}
	if ownResources < minThreshold {
		reqResource = minThreshold - ownResources + livingCost
	} else {
		if ourPersonality == Selfish {
			reqResource = 2 * livingCost
		} else {
			reqResource = livingCost
		}
	}
	c.Logf("Request %v from common pool", reqResource)
	return reqResource
}

func (c *client) RequestAllocation() shared.Resources {
	numberAlive := shared.Resources(c.getNumOfAliveIslands())
	ourStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	//if we are critical or dying
	if ourStatus == shared.Critical || ourStatus == shared.Dead {
		if numberAlive > 0 {
			takenResource := c.ServerReadHandle.GetGameState().CommonPool / numberAlive
			c.Logf("Taken %v from common pool", takenResource)
			return takenResource
		}
	}
	return c.allocation
}

func (c *client) ResourceReport() shared.ResourcesReport {
	// if we are selfish, will report 1/2 of the actual resources
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	ourPersonality := c.getPersonality()
	fakeReport := shared.ResourcesReport{
		ReportedAmount: shared.Resources(float64(1/2) * float64(ourResources)),
		Reported:       true,
	}

	if ourPersonality == Generous {
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
	if ok && prediction.TimeLeft < 1 {
		total = shared.Resources(prediction.Magnitude*prediction.Confidence/3) + c.payingTax
	}

	ourPersonality := c.getPersonality()
	if ourPersonality == Selfish {
		return c.payingTax
	}
	return total
}

func (c *client) GetSanctionPayment() shared.Resources {
	return c.payingSanction
}

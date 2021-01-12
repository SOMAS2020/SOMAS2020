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
			c.taxDemanded = shared.Resources(content.IIGOValueData.Amount)
		case shared.SanctionAmount:
			c.sanctionDemanded = shared.Resources(content.IntegerData)
		case shared.IIGOAllocationDecision:
			c.allocationAllowed = shared.Resources(content.IIGOValueData.Amount)
		default:
		}
	}
}

func (c *client) MonitorIIGORole(roleName shared.Role) bool {
	currPresident := c.ServerReadHandle.GetGameState().PresidentID

	if currPresident == c.GetID() {
		return true
	}

	return false
}

func (c *client) DecideIIGOMonitoringAnnouncement(monitoringResult bool) (resultToShare bool, announce bool) {
	resultToShare = monitoringResult
	announce = true
	return
}

func (c *client) CommonPoolResourceRequest() shared.Resources {
	ourPersonality := c.getPersonality()
	numberAlive := shared.Resources(c.getNumOfAliveIslands())
	livingCost := c.ServerReadHandle.GetGameConfig().CostOfLiving
	minThreshold := c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	ourStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	cprLeft := c.ServerReadHandle.GetGameState().CommonPool
	reqResource := shared.Resources(3.0 * livingCost)

	if ourStatus == shared.Critical && ourResources < minThreshold {
		return minThreshold - ourResources
	}

	//when common pool does not have enough resource, will not request
	if ourPersonality == Generous && cprLeft/numberAlive < livingCost {
		reqResource = livingCost
	}

	c.Logf("Request %v from common pool", reqResource)

	return reqResource
}

func (c *client) RequestAllocation() shared.Resources {
	resourceTaken := c.allocationAllowed
	numberAlive := shared.Resources(c.getNumOfAliveIslands())
	ourStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	ourPersonality := c.getPersonality()
	minThreshold := c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold
	livingCost := c.ServerReadHandle.GetGameConfig().CostOfLiving
	commonPool := c.ServerReadHandle.GetGameState().CommonPool
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources

	//if we are critical or dying
	if c.ServerReadHandle.GetGameState().ClientInfo.CriticalConsecutiveTurnsCounter == 2 {
		if ourResources < minThreshold {
			return minThreshold - ourResources
		}
	}

	if ourStatus == shared.Critical {
		return minThreshold - ourResources + livingCost
	}

	if numberAlive <= 1 {
		return commonPool
	} else if ourPersonality == Selfish && c.allocationAllowed < livingCost && commonPool >= livingCost*numberAlive {
		resourceTaken = livingCost
	}

	return resourceTaken
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
	payTax := c.taxDemanded
	prediction, ok := c.disasterPredictions[c.GetID()]
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	ourPersonality := c.getPersonality()
	ourStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	minThreshold := c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold

	if ourStatus == shared.Critical {
		return 0
	} else if ourResources-c.taxDemanded < minThreshold {
		payTax = ourResources - minThreshold
	}

	if ok && prediction.TimeLeft < 1 && ourPersonality == Generous {
		payTax = shared.Resources(prediction.Magnitude*prediction.Confidence/3) + c.taxDemanded
	}

	return payTax
}

func (c *client) GetSanctionPayment() shared.Resources {
	paySanction := c.sanctionDemanded
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	ourStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	minThreshold := c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold

	if ourStatus == shared.Critical {
		return 0
	} else if ourResources-c.sanctionDemanded < minThreshold {
		paySanction = ourResources - minThreshold
	}

	return paySanction
}

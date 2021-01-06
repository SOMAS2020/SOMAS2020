package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

func (c *client) GetClientJudgePointer() roles.Judge {
	return &c.clientJudge
}

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	return &c.clientSpeaker
}

func (c *client) RequestAllocation() shared.Resources {
	allocationGranted := c.obs.iigoObs.allocationGranted
	uncomplianceThreshold := 5.0
	importance := mat.NewVecDense(4, []float64{
		5.0, 1.0, 0.0, -1.0, 5.0, 0.0,
		// TODO: add multiplier for the 0.0 ones
	})

	parameters := mat.NewVecDense(4, []float64{
		c.internalParam.greediness,
		c.internalParam.selfishness,
		c.internalParam.fairness,
		c.internalParam.collaboration,
		c.internalParam.riskTaking,
		c.internalParam.agentsTrust[0], // TODO: index properly based on president in that turn
	})

	uncomplianceLevel := mat.Dot(importance, parameters) - uncomplianceThreshold
	// TODO: if we're in critical state, take the resource needed to get us out of it. Maybe take a protion of what we need
	// until the very last turns in which we're about to die and take all we need to get out of critical state.
	// c.ServerReadHandle.GetGameConfig().maxCriticalConsecutiveTurns
	allocDemanded := allocationGranted
	if uncomplianceLevel > 0 {
		allocDemanded = allocationGranted * shared.Resources((uncomplianceLevel + 1))
	}

	return allocDemanded
}

// this function is used to receive tax amount, allocation amount rule name etc from the server. Use this to receive information
// and store it inside our agent's observation
func (c *client) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	c.Communications[sender] = append(c.Communications[sender], data)
	// TODO parse sanction info
	for contentType, content := range data {
		switch contentType {
		case shared.IIGOTaxDecision:
			c.obs.iigoObs.taxDemanded = shared.Resources(content.IntegerData)
		case shared.IIGOAllocationDecision:
			c.obs.iigoObs.allocationGranted = shared.Resources(content.IntegerData)
		case shared.RuleName:
			// currentRuleID := content.TextData
			// if _, ok := c.iigoInfo.ruleVotingResults[currentRuleID]; ok {
			// 	c.iigoInfo.ruleVotingResults[currentRuleID].resultAnnounced = true
			// 	c.iigoInfo.ruleVotingResults[currentRuleID].result = data[shared.RuleVoteResult].BooleanData
			// } else {
			// 	c.iigoInfo.ruleVotingResults[currentRuleID] = &ruleVoteInfo{resultAnnounced: true, result: data[shared.RuleVoteResult].BooleanData}
			// }
		case shared.RoleMonitored:
			// c.iigoInfo.monitoringDeclared[content.IIGORoleData] = true
			// c.iigoInfo.monitoringOutcomes[content.IIGORoleData] = data[shared.MonitoringResult].BooleanData
		}
	}
}

func (c *client) CommonPoolResourceRequest() shared.Resources {
	// TODO: Implement needs based resource request.

	// available observations
	commonPoolLevel := c.ServerReadHandle.GetGameState().CommonPool
	ourResource := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	ourLifeStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	otherAgentsLifeStatuses := c.ServerReadHandle.GetGameState().ClientLifeStatuses

	numClientAlive := 0
	for _, status := range otherAgentsLifeStatuses {
		if status == shared.Alive || status == shared.Critical {
			numClientAlive++
		}
	}

	resNeeded := commonPoolLevel / shared.Resources(numClientAlive) //tempcomment: Allocation is taken before taxation before disaster

	if ourLifeStatus == shared.Critical {
		// turnsInCriticalState := c.ServerReadHandle.GetGameState().ClientInfo.CriticalConsecutiveTurnsCounter //TODO: probably don't need this, only need this in RequestAllocation()
		resNeeded = c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold - ourResource
	}

	// TODO: define how much we want -> resNeeded
	greedyThreshold := 2.5

	importance := mat.NewVecDense(4, []float64{
		5.0, 1.0, -1.0, -1.0, 1.0, 0.0,
		// TODO: add multiplier for the 0.0 ones
	})

	parameters := mat.NewVecDense(4, []float64{
		c.internalParam.greediness,
		c.internalParam.selfishness,
		c.internalParam.fairness,
		c.internalParam.collaboration,
		c.internalParam.riskTaking,
		c.internalParam.agentsTrust[0], // TODO: index properly based on president in that turn
	})
	greedyLevel := mat.Dot(importance, parameters) - greedyThreshold

	allocRequested := resNeeded // if we're selfless, still request and take resNeeded, but gift the extra to critical islands
	if greedyLevel > 0 {
		allocRequested = resNeeded * shared.Resources((greedyLevel + 1))
	}

	return allocRequested
}

func (c *client) ResourceReport() shared.ResourcesReport {
	// Parameters initialisation.
	currentResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	lyingThreshold := 3.0
	reporting := true

	// If collaboration and trust are above average chose to report, otherwise abstain!
	if (c.internalParam.collaboration + c.internalParam.agentsTrust[0]) < 1 { // agent trust towards the president, TODO: change to president index
		reporting = false
	}

	// Initialise importance vector and parameters vector.
	importance := mat.NewVecDense(4, []float64{
		5.0, 5.0, -5.0, -5.0, 1.0, 5.0,
		// TODO: add multiplier for the 0.0 ones.
	})

	parameters := mat.NewVecDense(4, []float64{
		c.internalParam.greediness,
		c.internalParam.selfishness,
		c.internalParam.fairness,
		c.internalParam.collaboration,
		c.internalParam.riskTaking,
		c.internalParam.agentsTrust[0], // TODO: index properly based on president and judge: respectively
		// to measure your trust on the fairness of the tax you will get/how
		// much you trust that agent with this info and how much you think the
		// judge is likely to inspect you.
	})

	// lyingLevel will be positive when agent is inclined to lie.
	lyingLevel := mat.Dot(importance, parameters) - lyingThreshold

	// Construct output struct.
	resReported := currentResources

	// Agent lies linearly based on lyingLevel.
	if lyingLevel > 0 {
		resReported = currentResources * (1 / shared.Resources((lyingLevel + 1)))
	}

	resReportStruct := shared.ResourcesReport{
		ReportedAmount: resReported,
		Reported:       reporting,
	}

	return resReportStruct
}

////////////// TODO: FUNCTION WAITING ON BASECLIENT PR /////////////
// GetTaxContribution()
// GetSanctionPayment()

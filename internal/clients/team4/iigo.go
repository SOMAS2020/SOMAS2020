package team4

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

func (c *client) GetClientJudgePointer() roles.Judge {
	c.clientJudge.GameState = c.ServerReadHandle.GetGameState()
	return &c.clientJudge
}

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	c.clientSpeaker.GameState = c.ServerReadHandle.GetGameState()
	return &c.clientSpeaker
}

// EvaluateParamVector returns the dot product of the decision matrix and the internal parameters
func (c *client) evaluateParamVector(decisionVector *mat.VecDense, agent shared.ClientID, threshold float64) float64 {
	parameters := mat.NewVecDense(5, []float64{
		c.internalParam.greediness,
		c.internalParam.selfishness,
		c.internalParam.fairness,
		c.internalParam.collaboration,
		c.internalParam.riskTaking,
		c.trustMatrix.GetClientTrust(agent),
	})
	return mat.Dot(decisionVector, parameters) - threshold
}

func (c *client) RequestAllocation() shared.Resources {
	ourLifeStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	allocationGranted := c.BaseClient.RequestAllocation() //c.obs.iigoObs.allocationGranted
	uncomplianceThreshold := 5.0
	importance := c.importances.requestAllocationImportance
	commonPool := c.ServerReadHandle.GetGameState().CommonPool

	parameters := mat.NewVecDense(6, []float64{
		c.internalParam.greediness,
		c.internalParam.selfishness,
		c.internalParam.fairness,
		c.internalParam.collaboration,
		c.internalParam.riskTaking,
		c.getTrust(c.getPresident()),
	})

	uncomplianceLevel := mat.Dot(importance, parameters) - uncomplianceThreshold

	// if alive and compliant then take nothing if granted nothing
	allocDemanded := allocationGranted
	if allocationGranted == 0 {
		c.internalParam.giftExtra = false
	}
	resNeeded := c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold + c.ServerReadHandle.GetGameConfig().CostOfLiving - c.getOurResources()
	if resNeeded < 0 {
		resNeeded = (c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold + c.ServerReadHandle.GetGameConfig().CostOfLiving) * shared.Resources(1+c.internalParam.greediness)
	}
	if ourLifeStatus == shared.Critical {
		c.internalParam.giftExtra = false
		maxTurnsInCritical := c.ServerReadHandle.GetGameConfig().MaxCriticalConsecutiveTurns
		turnsInCritical := c.ServerReadHandle.GetGameState().ClientInfo.CriticalConsecutiveTurnsCounter
		if turnsInCritical == maxTurnsInCritical {
			allocDemanded = shared.Resources(math.Min(float64(resNeeded*5), float64(commonPool)))
		} else if turnsInCritical == maxTurnsInCritical-1 {
			allocDemanded = shared.Resources(math.Min(float64(resNeeded*3), float64(commonPool)))
		} else if turnsInCritical == maxTurnsInCritical-2 {
			allocDemanded = shared.Resources(math.Min(float64(resNeeded*2), float64(commonPool)))

		}
	} else if ourLifeStatus == shared.Dead {
		c.internalParam.giftExtra = false
		allocDemanded = shared.Resources(0)
	} else if ourLifeStatus == shared.Alive {
		if uncomplianceLevel > 0 {
			if allocationGranted == 0 {
				c.internalParam.giftExtra = false
				allocDemanded = shared.Resources(math.Min(float64(resNeeded*2)*(uncomplianceLevel+1), float64(commonPool)))
			} else {
				allocDemanded = allocationGranted * shared.Resources((uncomplianceLevel + 1))
			}
		}
	}
	c.Logf("Allocation granted: %v", allocationGranted)

	c.Logf("Allocation demanded: %v", allocDemanded)
	return allocDemanded
}

// this function is used to receive tax amount, allocation amount rule name etc from the server. Use this to receive information
// and store it inside our agent's observation
func (c *client) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	c.BaseClient.ReceiveCommunication(sender, data)
	// TODO parse sanction info
	c.updateTrustMonitoring(data)

	for contentType, content := range data {
		switch contentType {
		case shared.IIGOTaxDecision:
			c.obs.iigoObs.taxDemanded = content.IIGOValueData.Amount //shared.Resources(content.IntegerData)
		case shared.IIGOAllocationDecision:
			c.obs.iigoObs.allocationGranted = content.IIGOValueData.Amount //shared.Resources(content.IntegerData)
		case shared.RuleName:
		// currentRuleID := content.TextData
		// if _, ok := c.iigoInfo.ruleVotingResults[currentRuleID]; ok {
		// 	c.iigoInfo.ruleVotingResults[currentRuleID].resultAnnounced = true
		// 	c.iigoInfo.ruleVotingResults[currentRuleID].result = data[shared.RuleVoteResult].BooleanData
		// } else {
		// 	c.iigoInfo.ruleVotingResults[currentRuleID] = &ruleVoteInfo{resultAnnounced: true, result: data[shared.RuleVoteResult].BooleanData}
		// }
		case shared.SanctionClientID:
			if sanctionTier, ok := data[shared.IIGOSanctionTier]; ok {
				sanctionedClient := shared.ClientID(content.IntegerData)
				sanctionTierData := shared.IIGOSanctionsTier(sanctionTier.IntegerData)
				c.obs.iigoObs.sanctionTiers[sanctionedClient] = sanctionTierData
			}
		default: //[exhaustive] reported by reviewdog 🐶
			return
			//missing cases in switch of type shared.CommunicationFieldName: BallotID, IIGOSanctionScore, IIGOSanctionTier, MonitoringResult, PardonClientID, PardonTier, PresidentID, ResAllocID, RoleConducted, RuleVoteResult, SanctionAmount, SanctionClientID, SpeakerID (exhaustive)

		}
	}
}

func (c *client) CommonPoolResourceRequest() shared.Resources {
	// TODO: Implement needs based on resource request.

	// available observations
	c.internalParam.giftExtra = false
	commonPoolLevel := c.ServerReadHandle.GetGameState().CommonPool
	ourResource := c.getOurResources()
	ourLifeStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	otherAgentsLifeStatuses := c.ServerReadHandle.GetGameState().ClientLifeStatuses

	numClientAlive := 0
	for _, status := range otherAgentsLifeStatuses {
		if status == shared.Alive || status == shared.Critical {
			numClientAlive++
		}
	}

	resNeeded := shared.Resources(0)
	if ourLifeStatus != shared.Dead {
		if numClientAlive != 0 {
			eachClient := commonPoolLevel / shared.Resources(numClientAlive) //tempcomment: Allocation is taken before taxation before disaster

			if ourLifeStatus == shared.Alive {
				resNeeded = shared.Resources(2) * c.getSafeResourceLevel()
				resNeeded += shared.Resources(numClientAlive * 10)
				c.internalParam.giftExtra = true

			} else if ourLifeStatus == shared.Critical {
				resNeeded = c.getSafeResourceLevel() - ourResource
				resNeeded *= shared.Resources(3)
			}
			resNeeded = shared.Resources(math.Min(float64(eachClient), float64(resNeeded)))
		} else {
			resNeeded = commonPoolLevel * shared.Resources(rand.Float64())
		}
	}

	greedyThreshold := 2.5
	importance := c.importances.commonPoolResourceRequestImportance

	parameters := mat.NewVecDense(6, []float64{
		c.internalParam.greediness,
		c.internalParam.selfishness,
		c.internalParam.fairness,
		c.internalParam.collaboration,
		c.internalParam.riskTaking,
		c.getTrust(c.getPresident()),
	})
	greedyLevel := mat.Dot(importance, parameters) - greedyThreshold

	allocRequested := resNeeded // if we're selfless, still request and take resNeeded, but gift the extra to critical islands.
	if greedyLevel > 0 {
		allocRequested = resNeeded * shared.Resources((greedyLevel + 1))
	}
	c.Logf("Allocation requested: %v", allocRequested)

	return allocRequested
}

func (c *client) ResourceReport() shared.ResourcesReport {
	// Parameters initialisation.
	currentResources := c.getOurResources()
	lyingThreshold := 3.0
	reporting := true

	presidentID := c.ServerReadHandle.GetGameState().PresidentID
	// If collaboration and trust are above average chose to report, otherwise abstain!
	if (c.internalParam.collaboration + c.trustMatrix.GetClientTrust(presidentID)) < 1 { // agent trust towards the president
		reporting = false
	}

	// Initialise importance vector and parameters vector.
	importance := c.importances.resourceReportImportance

	parameters := mat.NewVecDense(6, []float64{
		c.internalParam.greediness,
		c.internalParam.selfishness,
		c.internalParam.fairness,
		c.internalParam.collaboration,
		c.internalParam.riskTaking,
		c.getTrust(c.getPresident()),
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

// GetTaxContribution gives value of how much the island wants to pay in taxes
// The tax is the minimum contribution, you can pay more if you want to
// COMPULSORY

func (c *client) GetTaxContribution() shared.Resources {
	valToBeReturned := c.BaseClient.GetTaxContribution() // c.obs.iigoObs.taxDemanded

	c.Logf("Team4 tax expected: %v", valToBeReturned)

	currentWealth := c.getOurResources()

	collaborationThreshold := 3.0
	wealthThreshold := 5 * valToBeReturned

	// Initialise importance vector and parameters vector.
	importance := c.importances.getTaxContributionImportance

	parameters := mat.NewVecDense(4, []float64{
		c.internalParam.greediness,
		c.internalParam.selfishness,
		c.internalParam.collaboration,
		c.getTrust(c.getPresident()),
	})

	collaborationLevel := mat.Dot(importance, parameters)

	if collaborationLevel > collaborationThreshold &&
		currentWealth > wealthThreshold {
		// Deliberately pay more (collaborationLevel is larger than 1)
		valToBeReturned = valToBeReturned * shared.Resources(0.2*(collaborationLevel-collaborationThreshold))

	}

	c.Logf("Tax paid: %v", valToBeReturned)
	return valToBeReturned

}

// GetSanctionPayment()

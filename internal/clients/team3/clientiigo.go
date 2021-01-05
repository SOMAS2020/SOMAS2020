package team3

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/clients/team3/dynamics"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
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
   RequestAllocation() shared.Resources
*/

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	c.clientPrint("became speaker")
	return &c.ourSpeaker
}

func (c *client) GetClientJudgePointer() roles.Judge {
	c.clientPrint("became judge")
	return &c.ourJudge
}

func (c *client) GetClientPresidentPointer() roles.President {
	c.clientPrint("became president")
	return &c.ourPresident
}

//resetIIGOInfo clears the island's information regarding IIGO at start of turn
func (c *client) resetIIGOInfo() {
	c.iigoInfo.ourRole = nil
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

func (c *client) GetTaxContribution() shared.Resources {
	commonPool := c.BaseClient.ServerReadHandle.GetGameState().CommonPool
	totalToPay := 100 - commonPool
	if len(c.disasterPredictions) > int(c.ServerReadHandle.GetGameState().Turn) {
		if disaster, ok := c.disasterPredictions[int(c.BaseClient.ServerReadHandle.GetGameState().Turn)][c.BaseClient.GetID()]; ok {
			totalToPay = (shared.Resources(disaster.Magnitude) - commonPool) / shared.Resources(disaster.TimeLeft)
		}
	}
	sumTrust := 0.0
	for id, trust := range c.trustScore {
		if id != c.BaseClient.GetID() {
			sumTrust += trust
		} else {
			sumTrust += (1 - c.params.selfishness) * 100
		}
	}
	toPay := (totalToPay / shared.Resources(sumTrust)) * (1 - shared.Resources(c.params.selfishness)) * 100
	targetResources := shared.Resources(2-c.params.riskFactor) * (c.criticalStatePrediction.upperBound)
	if c.getLocalResources()-toPay <= targetResources {
		toPay = shared.Resources(math.Max(float64(c.getLocalResources()-targetResources), 0.0))
	}
	if (c.iigoInfo.taxationAmount > toPay) && !c.shouldICheat() {
		return c.iigoInfo.taxationAmount
	}
	c.clientPrint("Paying %v in tax", toPay)
	return toPay

}

// ReceiveCommunication is a function called by IIGO to pass the communication sent to the client.
// This function is overridden to receive information and update local info accordingly.
func (c *client) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	c.Communications[sender] = append(c.Communications[sender], data)
	// TODO parse sanction info
	for contentType, content := range data {
		switch contentType {
		case shared.TaxAmount:
			c.iigoInfo.taxationAmount = shared.Resources(content.IntegerData)
		case shared.AllocationAmount:
			c.iigoInfo.commonPoolAllocation = shared.Resources(content.IntegerData)
		case shared.RuleName:
			currentRuleID := content.TextData
			if _, ok := c.iigoInfo.ruleVotingResults[currentRuleID]; ok {
				c.iigoInfo.ruleVotingResults[currentRuleID].resultAnnounced = true
				c.iigoInfo.ruleVotingResults[currentRuleID].result = data[shared.RuleVoteResult].BooleanData
			} else {
				c.iigoInfo.ruleVotingResults[currentRuleID] = &ruleVoteInfo{resultAnnounced: true, result: data[shared.RuleVoteResult].BooleanData}
			}
		case shared.RoleMonitored:
			c.iigoInfo.monitoringDeclared[content.IIGORoleData] = true
			c.iigoInfo.monitoringOutcomes[content.IIGORoleData] = data[shared.MonitoringResult].BooleanData
		}
	}
}

func (c *client) RuleProposal() string {
	c.locationService.syncGameState(c.ServerReadHandle.GetGameState())
	c.locationService.syncTrustScore(c.trustScore)
	internalMap := copyRulesMap(rules.RulesInPlay)
	inputMap := c.locationService.TranslateCommunications(c.localVariableCache)
	c.localInputsCache = inputMap
	shortestSoFar := -2.0
	selectedRule := ""
	for key, rule := range rules.AvailableRules {
		if _, ok := rules.RulesInPlay[key]; !ok {
			reqInputs := dynamics.SourceRequiredInputs(rule, inputMap)
			idealLoc, valid := c.locationService.checkIfIdealLocationAvailable(rule, reqInputs)
			if valid {
				ruleDynamics := dynamics.BuildAllDynamics(rule, rule.AuxiliaryVector)
				distance := dynamics.GetDistanceToSubspace(ruleDynamics, idealLoc)
				if distance == -1 {
					return key
				}
				if shortestSoFar == -2.0 || shortestSoFar > distance {
					shortestSoFar = distance
					selectedRule = rule.RuleName
				}

			}
		} else {
			lstRules := dynamics.RemoveFromMap(internalMap, key)
			dist := dynamics.CalculateDistanceFromRuleSpace(lstRules, inputMap)
			if shortestSoFar == -2.0 || dist < shortestSoFar {
				selectedRule = rule.RuleName
				shortestSoFar = dist
			}
		}
	}
	if selectedRule == "" {
		return "inspect_ballot_rule"
	}
	return selectedRule
}

func (c *client) GetVoteForRule(ruleName string) bool {

	newRulesInPlay := make(map[string]rules.RuleMatrix)

	for key, value := range rules.RulesInPlay {
		newRulesInPlay[key] = value
	}

	if _, ok := rules.RulesInPlay[ruleName]; ok {
		delete(newRulesInPlay, ruleName)
	} else {
		newRulesInPlay[ruleName] = rules.AvailableRules[ruleName]
	}

	// TODO: define postion -> list of variables and values associated with the rule (obtained from IIGO communications)

	// distancetoRulesInPlay = CalculateDistanceFromRuleSpace(rules.RulesInPlay, position)
	// distancetoNewRulesInPlay = CalculateDistanceFromRuleSpace(newRulesInPlay, position)

	// if distancetoRulesInPlay < distancetoNewRulesInPlay {
	//  return false
	// } else {
	//  return true
	// }

	return true
}

// RequestAllocation gives how much island is taking from common pool
func (c *client) RequestAllocation() shared.Resources {
	ourAllocation := c.iigoInfo.commonPoolAllocation
	currentState := c.BaseClient.ServerReadHandle.GetGameState()
	escapeCritical := c.params.escapeCritcaIsland && currentState.ClientInfo.LifeStatus == shared.Critical
	distCriticalThreshold := ((c.criticalStatePrediction.upperBound + c.criticalStatePrediction.lowerBound) / 2) - ourAllocation

	if escapeCritical && (ourAllocation < distCriticalThreshold) {
		// Get enough to save ourselves
		return distCriticalThreshold
	}

	if c.shouldICheat() {
		// Scale up allocation a bit
		return ourAllocation + shared.Resources(float64(ourAllocation)*c.params.selfishness)
	}

	// Base return - take what we are allocated, but make sure we are stolen from!
	if ourAllocation < shared.Resources(0) {
		ourAllocation = shared.Resources(0)
	}
	c.clientPrint("Taking %f from common pool", ourAllocation)
	return ourAllocation

}

// CommonPoolResourceRequest is called by the President in IIGO to
// request an allocation of resources from the common pool.
func (c *client) CommonPoolResourceRequest() shared.Resources {
	var request shared.Resources

	currentState := c.BaseClient.ServerReadHandle.GetGameState()
	ourResources := currentState.ClientInfo.Resources
	escapeCritical := c.params.escapeCritcaIsland && currentState.ClientInfo.LifeStatus == shared.Critical
	distCriticalThreshold := ((c.criticalStatePrediction.upperBound + c.criticalStatePrediction.lowerBound) / 2) - ourResources

	request = shared.Resources(c.params.minimumRequest)
	if escapeCritical {
		if request < distCriticalThreshold {
			request = distCriticalThreshold
		}
	}
	if c.shouldICheat() {
		request += shared.Resources(float64(request) * c.params.selfishness)
	}
	// TODO request based on disaster prediction
	c.clientPrint("Our Request: %f", request)
	return request
}

func copyRulesMap(inp map[string]rules.RuleMatrix) map[string]rules.RuleMatrix {
	newMap := make(map[string]rules.RuleMatrix)
	for key, val := range inp {
		newMap[key] = val
	}
	return newMap
}

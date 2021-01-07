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

func (c *client) VoteForElection(roleToElect shared.Role, candidateList []shared.ClientID) []shared.ClientID {
	candidates := map[int]shared.ClientID{}
	for i := 0; i < len(candidateList); i++ {
		candidates[i] = candidateList[i]
	}
	// Recombine map, in shuffled order
	var returnList []shared.ClientID
	for _, v := range candidates {
		returnList = append(returnList, v)
	}
	return returnList
}

//resetIIGOInfo clears the island's information regarding IIGO at start of turn
func (c *client) resetIIGOInfo() {
	c.iigoInfo.ourRole = nil // TODO unused, remove
	c.clientPrint("IIGO cache from previous turn: %+v", c.iigoInfo)
	c.clientPrint("IIGO sanction info from previous turn: %+v", c.iigoInfo.sanctions)
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
	variablesChanged := map[rules.VariableFieldName]rules.VariableValuePair{
		rules.IslandTaxContribution: {
			rules.IslandTaxContribution,
			[]float64{float64(toPay)},
		},
		rules.ExpectedTaxContribution: {
			rules.ExpectedTaxContribution,
			c.LocalVariableCache[rules.ExpectedTaxContribution].Values,
		},
	}
	recommendedValues := c.dynamicAssistedResult(variablesChanged)
	resolve := shared.Resources(recommendedValues[rules.IslandTaxContribution].Values[rules.SingleValueVariableEntry])
	if c.params.complianceLevel > 80 {
		return resolve
	}
	if toPay != resolve {
		affectedRules, success := rules.PickUpRulesByVariable(rules.IslandTaxContribution, rules.RulesInPlay, c.LocalVariableCache)
		if success {
			c.oldBrokenRules = append(c.oldBrokenRules, affectedRules...)
		}
	}
	return toPay

}

func (c *client) dynamicAssistedResult(variablesChanged map[rules.VariableFieldName]rules.VariableValuePair) (newVals map[rules.VariableFieldName]rules.VariableValuePair) {
	if c.LocalVariableCache != nil {
		c.LocalVariableCache = c.locationService.UpdateCache(c.LocalVariableCache, variablesChanged)
		// For testing using available rules
		return c.locationService.GetRecommendations(variablesChanged, rules.AvailableRules, c.LocalVariableCache)
	}
	return variablesChanged
}

// ReceiveCommunication is a function called by IIGO to pass the communication sent to the client.
// This function is overridden to receive information and update local info accordingly.
func (c *client) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	c.Communications[sender] = append(c.Communications[sender], data)
	// c.clientPrint("Received communication: %+v", data)
	for contentType, content := range data {
		switch contentType {
		case shared.IIGOTaxDecision:
			c.iigoInfo.taxationAmount = shared.Resources(content.IntegerData)
		case shared.IIGOAllocationDecision:
			c.iigoInfo.commonPoolAllocation = shared.Resources(content.IntegerData)
		case shared.RuleName:
			currentRuleID := content.TextData
			// Rule voting
			if _, ok := data[shared.RuleVoteResult]; ok {
				if _, ok := c.iigoInfo.ruleVotingResults[currentRuleID]; ok {
					c.iigoInfo.ruleVotingResults[currentRuleID].resultAnnounced = true
					c.iigoInfo.ruleVotingResults[currentRuleID].result = data[shared.RuleVoteResult].BooleanData
				} else {
					c.iigoInfo.ruleVotingResults[currentRuleID] = &ruleVoteInfo{resultAnnounced: true, result: data[shared.RuleVoteResult].BooleanData}
				}
			}
			// Rule sanctions
			if _, ok := data[shared.IIGOSanctionScore]; ok {
				// c.clientPrint("Received sanction info: %+v", data)
				c.iigoInfo.sanctions.rulePenalties[currentRuleID] = roles.IIGOSanctionScore(data[shared.IIGOSanctionScore].IntegerData)
			}
		case shared.RoleMonitored:
			c.iigoInfo.monitoringDeclared[content.IIGORoleData] = true
			c.iigoInfo.monitoringOutcomes[content.IIGORoleData] = data[shared.MonitoringResult].BooleanData
		case shared.SanctionClientID:
			c.iigoInfo.sanctions.islandSanctions[shared.ClientID(content.IntegerData)] = roles.IIGOSanctionTier(data[shared.IIGOSanctionTier].IntegerData)
		case shared.IIGOSanctionTier:
			c.iigoInfo.sanctions.tierInfo[roles.IIGOSanctionTier(content.IntegerData)] = roles.IIGOSanctionScore(data[shared.IIGOSanctionScore].IntegerData)
		case shared.SanctionAmount:
			c.iigoInfo.sanctions.ourSanction = roles.IIGOSanctionScore(content.IntegerData)
		}
	}
}

func (c *client) VoteForRule(matrix rules.RuleMatrix) shared.RuleVoteType {

	newRulesInPlay := make(map[string]rules.RuleMatrix)

	for key, value := range rules.RulesInPlay {
		if key == matrix.RuleName {
			newRulesInPlay[key] = matrix
		} else {
			newRulesInPlay[key] = value
		}
	}

	if _, ok := rules.RulesInPlay[matrix.RuleName]; ok {
		delete(newRulesInPlay, matrix.RuleName)
	} else {
		newRulesInPlay[matrix.RuleName] = rules.AvailableRules[matrix.RuleName]
	}

	// TODO: define postion -> list of variables and values associated with the rule (obtained from IIGO communications)

	distancetoRulesInPlay := dynamics.CalculateDistanceFromRuleSpace(dynamics.CollapseRuleMap(rules.RulesInPlay), c.locationService.TranslateToInputs(c.LocalVariableCache))
	distancetoNewRulesInPlay := dynamics.CalculateDistanceFromRuleSpace(dynamics.CollapseRuleMap(newRulesInPlay), c.locationService.TranslateToInputs(c.LocalVariableCache))

	if distancetoRulesInPlay < distancetoNewRulesInPlay {
		return shared.Reject
	} else {
		return shared.Approve
	}
}

func (c *client) RuleProposal() rules.RuleMatrix {
	c.locationService.syncGameState(c.ServerReadHandle.GetGameState())
	c.locationService.syncTrustScore(c.trustScore)
	internalMap := copyRulesMap(rules.RulesInPlay)
	inputMap := c.locationService.TranslateToInputs(c.LocalVariableCache)
	c.localInputsCache = inputMap
	shortestSoFar := -2.0
	selectedRule := ""
	if c.params.intelligence {
		newMat, success := c.intelligentShift()
		if success {
			return newMat
		}
	}
	for key, rule := range rules.AvailableRules {
		if _, ok := rules.RulesInPlay[key]; !ok {
			reqInputs := dynamics.SourceRequiredInputs(rule, inputMap)
			idealLoc, valid := c.locationService.checkIfIdealLocationAvailable(rule, reqInputs)
			if valid {
				ruleDynamics := dynamics.BuildAllDynamics(rule, rule.AuxiliaryVector)
				distance := dynamics.GetDistanceToSubspace(ruleDynamics, idealLoc)
				if distance == -1 {
					return rules.AvailableRules[key]
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
		selectedRule = "inspect_ballot_rule"
	}
	return rules.AvailableRules[selectedRule]
}

func (c *client) intelligentShift() (rules.RuleMatrix, bool) {
	if len(c.oldBrokenRules) == 0 {
		return rules.RuleMatrix{}, false
	}
	luckyRule := c.oldBrokenRules[0]
	inputMap := c.locationService.TranslateToInputs(c.LocalVariableCache)
	return dynamics.Shift(rules.RulesInPlay[luckyRule], inputMap)
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

	variablesChanged := map[rules.VariableFieldName]rules.VariableValuePair{
		rules.IslandAllocation: {
			rules.IslandAllocation,
			[]float64{float64(ourAllocation)},
		},
		rules.ExpectedAllocation: {
			rules.ExpectedAllocation,
			c.LocalVariableCache[rules.ExpectedAllocation].Values,
		},
	}

	recommendedValues := c.dynamicAssistedResult(variablesChanged)
	resolve := shared.Resources(recommendedValues[rules.IslandAllocation].Values[rules.SingleValueVariableEntry])
	if c.params.complianceLevel > 80 {
		return resolve
	}
	if ourAllocation != resolve {
		affectedRules, success := rules.PickUpRulesByVariable(rules.IslandAllocation, rules.RulesInPlay, c.LocalVariableCache)
		if success {
			c.oldBrokenRules = append(c.oldBrokenRules, affectedRules...)
		}
	}
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

func (c *client) GetSanctionPayment() shared.Resources {
	if value, ok := c.LocalVariableCache[rules.SanctionExpected]; ok {
		idealVal, available := c.locationService.switchDetermineFunction(rules.SanctionPaid, value.Values)
		if available {
			variablesChanged := map[rules.VariableFieldName]rules.VariableValuePair{
				rules.SanctionPaid: {
					rules.SanctionPaid,
					idealVal,
				},
				rules.SanctionExpected: {
					rules.SanctionExpected,
					c.LocalVariableCache[rules.SanctionExpected].Values,
				},
			}

			recommendedValues := c.dynamicAssistedResult(variablesChanged)
			resolve := shared.Resources(recommendedValues[rules.SanctionPaid].Values[rules.SingleValueVariableEntry])
			if c.params.complianceLevel > 80 {
				return resolve
			}
			if shared.Resources(idealVal[rules.SingleValueVariableEntry]) != resolve {
				affectedRules, success := rules.PickUpRulesByVariable(rules.SanctionPaid, rules.RulesInPlay, c.LocalVariableCache)
				if success {
					c.oldBrokenRules = append(c.oldBrokenRules, affectedRules...)
				}
			}
			return shared.Resources(idealVal[rules.SingleValueVariableEntry])
		}
		return shared.Resources(value.Values[rules.SingleValueVariableEntry])
	}
	return 0
}

func copyRulesMap(inp map[string]rules.RuleMatrix) map[string]rules.RuleMatrix {
	newMap := make(map[string]rules.RuleMatrix)
	for key, val := range inp {
		newMap[key] = val
	}
	return newMap
}

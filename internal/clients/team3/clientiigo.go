package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team3/dynamics"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"math"
	"math/rand"
	"sort"
)

/*
	//IIGO: COMPULSORY
	MonitorIIGORole(shared.Role) bool
	DecideIIGOMonitoringAnnouncement(bool) (bool, bool)

	VoteForRule(ruleName string) bool
	VoteForElection(roleToElect shared.Role) []shared.ClientID

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

// Vote based on island's past performance in the role and trust score if they have not previously held that role
func (c *client) VoteForElection(roleToElect shared.Role, candidateList []shared.ClientID) []shared.ClientID {

	// Get relevant map of past performance
	var pastRolePerformance = make(map[shared.ClientID]float64)
	if roleToElect == shared.President {
		pastRolePerformance = c.presidentPerformance
	}
	if roleToElect == shared.Judge {
		pastRolePerformance = c.judgePerformance
	} else {
		pastRolePerformance = c.speakerPerformance
	}

	returnList := []shared.ClientID{c.GetID()}

	// Calculate combined trust and past performance metric
	var trustPerformanceScore = make(map[shared.ClientID]float64)
	for _, island := range candidateList {
		if island != c.GetID() {
			trustPerformanceScore[island] = c.trustScore[island]
			if val, ok := pastRolePerformance[island]; ok {
				trustPerformanceScore[island] += val
			}
		}
	}

	return append(returnList, sortTrustPerformanceScore(trustPerformanceScore)...)
}

func sortTrustPerformanceScore(trustPerformanceScore map[shared.ClientID]float64) []shared.ClientID {
	final := []shared.ClientID{}
	for k := range trustPerformanceScore {
		final = append(final, k)
	}
	sort.Slice(final, func(i, j int) bool {
		return trustPerformanceScore[final[i]] > trustPerformanceScore[final[j]]
	})
	return final
}

//resetIIGOInfo clears the island's information regarding IIGO at start of turn
func (c *client) resetIIGOInfo() {
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
		tierInfo:        make(map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore),
		rulePenalties:   make(map[string]shared.IIGOSanctionsScore),
		islandSanctions: make(map[shared.ClientID]shared.IIGOSanctionsTier),
		ourSanction:     shared.IIGOSanctionsScore(0),
	}
	c.iigoInfo.ruleVotingResults = make(map[string]*ruleVoteInfo)
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
	totalToPay := shared.Resources(math.Max(float64(c.getIIGOCost()-commonPool), 0))
	if len(c.globalDisasterPredictions) > int(c.ServerReadHandle.GetGameState().Turn) {
		disaster := c.globalDisasterPredictions[int(c.ServerReadHandle.GetGameState().Turn)]
		totalToPay += safeDivResources(shared.Resources(disaster.Magnitude)-commonPool, shared.Resources(disaster.TimeLeft+1))
	}
	sumTrust := 0.0
	for id, trust := range c.trustScore {
		if id != c.BaseClient.GetID() {
			sumTrust += trust
		} else {
			sumTrust += (1 - c.params.selfishness) * 100
		}
	}
	toPay := shared.Resources(0)
	if sumTrust != 0 {
		toPay = safeDivResources(totalToPay, shared.Resources(sumTrust)) * (1 - shared.Resources(c.params.selfishness)) * 100
	}
	targetResources := shared.Resources(2-c.params.riskFactor) * (c.criticalThreshold)
	if c.getLocalResources()-toPay <= targetResources {
		toPay = shared.Resources(math.Max(float64(c.getLocalResources()-targetResources), 0.0))
	}
	if (c.iigoInfo.taxationAmount > toPay) && !c.shouldICheat() {
		return c.iigoInfo.taxationAmount
	}
	c.clientPrint("Paying %v in tax", toPay)
	variablesChanged := map[rules.VariableFieldName]rules.VariableValuePair{
		rules.IslandTaxContribution: {
			VariableName: rules.IslandTaxContribution,
			Values:       []float64{float64(toPay)},
		},
		rules.ExpectedTaxContribution: {
			VariableName: rules.ExpectedTaxContribution,
			Values:       c.LocalVariableCache[rules.ExpectedTaxContribution].Values,
		},
	}
	recommendedValues := c.dynamicAssistedResult(variablesChanged)
	resolve := shared.Resources(recommendedValues[rules.IslandTaxContribution].Values[rules.SingleValueVariableEntry])
	if c.params.complianceLevel > 80 {
		return resolve
	}
	if toPay != resolve {
		rulesInPlay := c.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay
		affectedRules, success := rules.PickUpRulesByVariable(rules.IslandTaxContribution, rulesInPlay, c.LocalVariableCache)
		if success {
			c.oldBrokenRules = append(c.oldBrokenRules, affectedRules...)
		}
	}

	c.account.LoadTaxation(toPay)

	return toPay

}

func (c *client) dynamicAssistedResult(variablesChanged map[rules.VariableFieldName]rules.VariableValuePair) (newVals map[rules.VariableFieldName]rules.VariableValuePair) {
	if c.LocalVariableCache != nil {
		c.LocalVariableCache = c.locationService.UpdateCache(c.LocalVariableCache, variablesChanged)
		// For testing using available rules
		return c.locationService.GetRecommendations(variablesChanged, c.ServerReadHandle.GetGameState().RulesInfo.AvailableRules, c.LocalVariableCache)
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
			c.iigoInfo.taxationAmount = content.IIGOValueData.Amount
		case shared.IIGOAllocationDecision:
			c.iigoInfo.commonPoolAllocation = content.IIGOValueData.Amount
		case shared.RuleName:

			// Rule voting
			if _, ok := data[shared.RuleVoteResult]; ok {
				currentRuleID := content.RuleMatrixData.RuleName
				if _, ok := c.iigoInfo.ruleVotingResults[currentRuleID]; ok {
					c.iigoInfo.ruleVotingResults[currentRuleID].resultAnnounced = true
					c.iigoInfo.ruleVotingResults[currentRuleID].result = data[shared.RuleVoteResult].BooleanData
				} else {
					c.iigoInfo.ruleVotingResults[currentRuleID] = &ruleVoteInfo{resultAnnounced: true, result: data[shared.RuleVoteResult].BooleanData}
				}
			}
			// Rule sanctions
			if _, ok := data[shared.RuleSanctionPenalty]; ok {
				currentRuleID := content.TextData
				// c.clientPrint("Received sanction info: %+v", data)
				c.iigoInfo.sanctions.rulePenalties[currentRuleID] = shared.IIGOSanctionsScore(data[shared.RuleSanctionPenalty].IntegerData)
			}
		case shared.RoleMonitored:
			c.iigoInfo.monitoringDeclared[content.IIGORoleData] = true
			c.iigoInfo.monitoringOutcomes[content.IIGORoleData] = data[shared.MonitoringResult].BooleanData
		case shared.SanctionClientID:
			c.iigoInfo.sanctions.islandSanctions[shared.ClientID(content.IntegerData)] = shared.IIGOSanctionsTier(data[shared.IIGOSanctionTier].IntegerData)
		case shared.IIGOSanctionTier:
			c.iigoInfo.sanctions.tierInfo[shared.IIGOSanctionsTier(content.IntegerData)] = shared.IIGOSanctionsScore(data[shared.RuleSanctionPenalty].IntegerData)
		case shared.SanctionAmount:
			c.clientPrint("Got our sanction :( %+v", content)
			c.iigoInfo.sanctions.ourSanction = shared.IIGOSanctionsScore(content.IntegerData)
			c.LocalVariableCache[rules.SanctionExpected] = rules.VariableValuePair{
				VariableName: rules.SanctionExpected,
				Values:       []float64{float64(content.IntegerData)},
			}
		}
	}
}

func (c *client) VoteForRule(matrix rules.RuleMatrix) shared.RuleVoteType {

	newRulesInPlay := make(map[string]rules.RuleMatrix)

	for key, value := range c.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay {
		if key == matrix.RuleName {
			newRulesInPlay[key] = matrix
		} else {
			newRulesInPlay[key] = value
		}
	}

	if _, ok := c.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay[matrix.RuleName]; ok {
		delete(newRulesInPlay, matrix.RuleName)
	} else {
		newRulesInPlay[matrix.RuleName] = matrix
	}

	// TODO: define postion -> list of variables and values associated with the rule (obtained from IIGO communications)
	rulesInPlay := c.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay
	distancetoRulesInPlay := dynamics.CalculateDistanceFromRuleSpace(dynamics.CollapseRuleMap(rulesInPlay), c.locationService.TranslateToInputs(c.LocalVariableCache))
	distancetoNewRulesInPlay := dynamics.CalculateDistanceFromRuleSpace(dynamics.CollapseRuleMap(newRulesInPlay), c.locationService.TranslateToInputs(c.LocalVariableCache))

	c.ruleVotedOn = matrix.RuleName
	if distancetoRulesInPlay < distancetoNewRulesInPlay {
		c.iigoInfo.ruleVotingResults[c.ruleVotedOn] = &ruleVoteInfo{ourVote: shared.Reject}
		return shared.Reject
	} else {
		c.iigoInfo.ruleVotingResults[c.ruleVotedOn] = &ruleVoteInfo{ourVote: shared.Approve}
		return shared.Approve
	}
}

func (c *client) RuleProposal() rules.RuleMatrix {
	if c.params.adv != nil {
		ret, done := c.params.adv.ProposeRule(c.ServerReadHandle.GetGameState().RulesInfo.AvailableRules)
		if done {
			return ret
		}
	}
	return c.generalRuleSelection(c.ServerReadHandle.GetGameState().RulesInfo.AvailableRules)
}

func (c *client) generalRuleSelection(allowedRules map[string]rules.RuleMatrix) rules.RuleMatrix {
	c.locationService.syncGameState(c.ServerReadHandle.GetGameState())
	c.locationService.syncTrustScore(c.trustScore)
	internalMap := copyRulesMap(c.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay)
	inputMap := c.locationService.TranslateToInputs(c.LocalVariableCache)
	c.localInputsCache = inputMap
	shortestSoFar := -2.0
	selectedRule := ""
	if rand.Int()%2 == 0 {
		newMat, success := c.intelligentShift()
		if success {
			return newMat
		}
	}
	for key, rule := range allowedRules {
		if _, ok := c.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay[key]; !ok {
			reqInputs := dynamics.SourceRequiredInputs(rule, inputMap)
			idealLoc, valid := c.locationService.checkIfIdealLocationAvailable(rule, reqInputs)
			if valid {
				ruleDynamics := dynamics.BuildAllDynamics(rule, rule.AuxiliaryVector)
				distance := dynamics.GetDistanceToSubspace(ruleDynamics, idealLoc)
				if distance == -1 {
					return allowedRules[key]
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
		for key := range allowedRules {
			selectedRule = key
			return allowedRules[key]
		}
	}
	return allowedRules[selectedRule]
}

func (c *client) intelligentShift() (rules.RuleMatrix, bool) {
	if len(c.oldBrokenRules) == 0 {
		return rules.RuleMatrix{}, false
	}
	luckyRule := c.oldBrokenRules[0]
	inputMap := c.locationService.TranslateToInputs(c.LocalVariableCache)
	rulesInPlay := c.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay
	return dynamics.Shift(rulesInPlay[luckyRule], inputMap)
}

// RequestAllocation gives how much island is taking from common pool
func (c *client) RequestAllocation() shared.Resources {
	var takenAlloc shared.Resources
	ourAllocation := c.iigoInfo.commonPoolAllocation
	currentState := c.BaseClient.ServerReadHandle.GetGameState()
	distCriticalThreshold := c.criticalThreshold - ourAllocation

	takenAlloc = ourAllocation

	// Escape critical
	if currentState.ClientInfo.LifeStatus == shared.Critical && (ourAllocation < distCriticalThreshold) {
		// Get enough to save ourselves
		takenAlloc = distCriticalThreshold
	} else {
		if c.shouldICheat() {
			// Scale up allocation a bit
			takenAlloc = ourAllocation + shared.Resources(float64(ourAllocation)*c.params.selfishness)
		}
	}

	// Base return - take what we are allocated, but make sure we aren't stolen from!
	if takenAlloc < shared.Resources(0) {
		takenAlloc = shared.Resources(0)
	}
	c.clientPrint("Taking %f from common pool", takenAlloc)

	variablesChanged := map[rules.VariableFieldName]rules.VariableValuePair{
		rules.IslandAllocation: {
			VariableName: rules.IslandAllocation,
			Values:       []float64{float64(takenAlloc)},
		},
		rules.ExpectedAllocation: {
			VariableName: rules.ExpectedAllocation,
			Values:       c.LocalVariableCache[rules.ExpectedAllocation].Values,
		},
	}

	recommendedValues := c.dynamicAssistedResult(variablesChanged)
	resolve := shared.Resources(recommendedValues[rules.IslandAllocation].Values[rules.SingleValueVariableEntry])
	if c.params.complianceLevel > 80 {
		return resolve
	}
	if takenAlloc != resolve {
		rulesInPlay := c.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay

		affectedRules, success := rules.PickUpRulesByVariable(rules.IslandAllocation, rulesInPlay, c.LocalVariableCache)
		if success {
			c.oldBrokenRules = append(c.oldBrokenRules, affectedRules...)
		}
	}
	if c.params.controlLoop {
		return shared.Resources(math.Max(float64(c.account.GetAllocMin()), float64(takenAlloc)))
	}
	return takenAlloc
}

// CommonPoolResourceRequest is called by the President in IIGO to
// request an allocation of resources from the common pool.
func (c *client) CommonPoolResourceRequest() shared.Resources {
	var request shared.Resources

	currentState := c.BaseClient.ServerReadHandle.GetGameState()
	ourResources := currentState.ClientInfo.Resources
	distCriticalThreshold := c.criticalThreshold - ourResources

	//request = c.ServerReadHandle.GetGameConfig().CostOfLiving
	request = shared.Resources(math.Max(float64(c.initialResourcesAtStartOfGame-ourResources), 0))
	// Try to escape critical
	if currentState.ClientInfo.LifeStatus == shared.Critical {
		request += distCriticalThreshold
	}
	if c.shouldICheat() {
		request += shared.Resources(float64(request) * c.params.selfishness)
	}

	if currentState.CommonPool <= request {
		request = shared.Resources(float64(currentState.CommonPool) * c.params.selfishness)
	}

	c.clientPrint("Our Request: %f", request)
	if c.params.controlLoop {
		return shared.Resources(math.Max(float64(c.account.GetAllocMin()), float64(request)))
	}
	return request
}

func (c *client) GetSanctionPayment() shared.Resources {
	if value, ok := c.LocalVariableCache[rules.SanctionExpected]; ok {
		idealVal, available := c.locationService.switchDetermineFunction(rules.SanctionPaid, value.Values)
		if available {
			variablesChanged := map[rules.VariableFieldName]rules.VariableValuePair{
				rules.SanctionPaid: {
					VariableName: rules.SanctionPaid,
					Values:       idealVal,
				},
				rules.SanctionExpected: {
					VariableName: rules.SanctionExpected,
					Values:       c.LocalVariableCache[rules.SanctionExpected].Values,
				},
			}

			recommendedValues := c.dynamicAssistedResult(variablesChanged)
			resolve := shared.Resources(recommendedValues[rules.SanctionPaid].Values[rules.SingleValueVariableEntry])
			if c.params.complianceLevel > 0.05 {
				return resolve
			}
			if shared.Resources(idealVal[rules.SingleValueVariableEntry]) != resolve {
				rulesInPlay := c.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay

				affectedRules, success := rules.PickUpRulesByVariable(rules.SanctionPaid, rulesInPlay, c.LocalVariableCache)
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

func safeDivResources(numerator shared.Resources, denominator shared.Resources) shared.Resources {
	if denominator != 0 {
		return numerator / denominator
	}
	return numerator
}

func safeDivFloat(numerator float64, denominator float64) float64 {
	if denominator != 0 {
		return numerator / denominator
	}
	return numerator
}

package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team3/dynamics"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

type metaStrategy struct {
	conquest    bool
	saviour     bool
	democrat    bool
	independent bool
	generous    bool
	lawful      bool
	legislative bool
	executive   bool
	fury        bool
}

type locator struct {
	islandParamsCache islandParams
	gameState         gamestate.ClientGameState
	locationStrategy  metaStrategy
	trustScore        map[shared.ClientID]float64
}

func (l *locator) syncGameState(newState gamestate.ClientGameState) {
	l.gameState = newState
}

func (l *locator) syncTrustScore(newTrustScore map[shared.ClientID]float64) {
	l.trustScore = newTrustScore
}

func (l *locator) changeStrategy(newParams islandParams) {
	l.islandParamsCache = newParams
	l.calculateMetaStrategy()
}

func (l *locator) calculateMetaStrategy() {
	newMetaStrategy := metaStrategy{}
	currentParams := l.islandParamsCache

	newMetaStrategy.conquest = !currentParams.saveCriticalIsland && currentParams.friendliness < 0.5
	newMetaStrategy.saviour = currentParams.saveCriticalIsland && currentParams.friendliness >= 0.5
	newMetaStrategy.democrat = currentParams.complianceLevel > 0.5 && currentParams.selfishness < 0.5
	newMetaStrategy.generous = currentParams.selfishness < 0.5 && currentParams.saveCriticalIsland && currentParams.friendliness > 0.3
	newMetaStrategy.lawful = currentParams.complianceLevel > 0.50 && currentParams.selfishness < 0.50
	newMetaStrategy.legislative = currentParams.equity > 0.50
	newMetaStrategy.fury = !currentParams.saveCriticalIsland
	l.locationStrategy = newMetaStrategy
}

func (l *locator) checkIfIdealLocationAvailable(ruleMat rules.RuleMatrix, recommendedVars map[rules.VariableFieldName]dynamics.Input) (mat.VecDense, bool) {
	requiredVariables := ruleMat.RequiredVariables
	valid := true
	var idealPosition []float64
	for _, variable := range requiredVariables {
		values, isValid := l.switchDetermineFunction(variable, recommendedVars[variable].Value)
		valid = valid && isValid
		idealPosition = append(idealPosition, values...)
	}
	if !valid {
		return mat.VecDense{}, false
	} else {
		return *mat.NewVecDense(len(idealPosition), idealPosition), true
	}
}

func (l *locator) switchDetermineFunction(name rules.VariableFieldName, recommendedVal []float64) ([]float64, bool) {
	switch name {
	case rules.NumberOfIslandsContributingToCommonPool:
		return l.idealNumberOfIslandsContributingToCommonPool(recommendedVal)
	case rules.NumberOfIslandsAlive:
		return l.idealNumberOfIslandsAlive(recommendedVal)
	case rules.NumberOfBallotsCast:
		return l.idealNumberOfBallotsCast(recommendedVal)
	case rules.AllocationRequestsMade:
		return l.idealAllocationRequestsMade(recommendedVal)
	case rules.AllocationMade:
		return l.idealAllocationsMade(recommendedVal)
	case rules.IslandsAlive:
		return l.idealIslandsAlive(recommendedVal)
	case rules.IslandTaxContribution:
		return l.idealIslandTaxContribution(recommendedVal)
	case rules.IslandAllocation:
		return l.idealIslandAllocation(recommendedVal)
	case rules.SanctionPaid:
		return l.idealIslandSanctionPaid(recommendedVal)
	default:
		return []float64{}, false
	}
}

func (l *locator) idealNumberOfIslandsContributingToCommonPool(recommendedVal []float64) ([]float64, bool) {
	if l.locationStrategy.generous {
		return []float64{l.calcNumberOfNonCriticalIslands()}, true
	} else {
		return []float64{l.calcNumberOfAliveIslands()}, true
	}
}

func (l *locator) idealNumberOfIslandsAlive(recommendedVal []float64) ([]float64, bool) {
	if l.locationStrategy.conquest {
		return []float64{1}, true
	} else if l.locationStrategy.saviour {
		return []float64{l.calcNumberOfAliveIslands()}, true
	} else {
		return []float64{l.calcNumberOfNonCriticalIslands()}, true
	}
}

func (l *locator) idealNumberOfBallotsCast(recommendedVal []float64) ([]float64, bool) {
	if l.locationStrategy.democrat {
		return []float64{l.calcNumberOfAliveIslands()}, true
	} else if l.locationStrategy.conquest {
		return []float64{l.calcNumberOfNonCriticalIslands()}, true
	} else {
		return []float64{}, false
	}
}

func (l *locator) idealAllocationsMade(recommendedVal []float64) ([]float64, bool) {
	return []float64{0}, l.locationStrategy.independent
}

func (l *locator) idealAllocationRequestsMade(recommendedVal []float64) ([]float64, bool) {
	return []float64{0}, l.locationStrategy.independent
}

func (l *locator) idealIslandsAlive(recommendedVal []float64) ([]float64, bool) {
	if l.locationStrategy.conquest {
		return []float64{2}, true
	} else if l.locationStrategy.fury {
		allAliveClients := l.getAllAliveClientIds()
		var allowedToLive []float64
		for _, clientDimID := range allAliveClients {
			if l.trustScore[shared.ClientID(clientDimID)] > 5 {
				allowedToLive = append(allowedToLive, clientDimID)
			}
		}
		return allowedToLive, true
	} else {
		return l.getAllAliveClientIds(), true
	}
}

func (l *locator) idealIslandTaxContribution(recommendedVal []float64) ([]float64, bool) {
	if len(recommendedVal) == 0 {
		return recommendedVal, false
	}
	if l.locationStrategy.lawful {
		return recommendedVal, true
	} else if l.locationStrategy.independent {
		return []float64{0}, true
	} else if l.locationStrategy.fury {
		return []float64{-1 * recommendedVal[0]}, true
	}
	return []float64{0}, false
}

func (l *locator) idealIslandAllocation(recommendedVal []float64) ([]float64, bool) {
	if len(recommendedVal) == 0 {
		return recommendedVal, false
	}
	if l.locationStrategy.lawful {
		return recommendedVal, true
	} else if !l.locationStrategy.generous {
		return []float64{recommendedVal[0] * 2}, true
	} else if l.locationStrategy.fury {
		return []float64{3 * recommendedVal[0]}, true
	}
	return []float64{0}, false
}

func (l *locator) idealIslandSanctionPaid(recommendedVal []float64) ([]float64, bool) {
	if len(recommendedVal) == 0 {
		return recommendedVal, false
	}
	if l.locationStrategy.lawful {
		return recommendedVal, true
	}
	if l.locationStrategy.fury {
		return []float64{0}, true
	}
	if l.locationStrategy.independent {
		return []float64{recommendedVal[rules.SingleValueVariableEntry] * 0.5}, true
	}
	return recommendedVal, false
}

func (l *locator) calcNumberOfAliveIslands() float64 {
	islandStatuses := l.gameState.ClientLifeStatuses
	counter := 0.0
	for _, status := range islandStatuses {
		if status != shared.Dead {
			counter++
		}
	}
	return counter
}

func (l *locator) calcNumberOfNonCriticalIslands() float64 {
	islandStatuses := l.gameState.ClientLifeStatuses
	counter := 0.0
	for _, status := range islandStatuses {
		if status != shared.Critical && status != shared.Dead {
			counter++
		}
	}
	return counter
}

func (l *locator) getAllAliveClientIds() []float64 {
	islandStatuses := l.gameState.ClientLifeStatuses
	var output []float64
	for key := range islandStatuses {
		output = append(output, float64(key))
	}
	return output
}

func (l *locator) TranslateToInputs(variableCache map[rules.VariableFieldName]rules.VariableValuePair) (newVariableCache map[rules.VariableFieldName]dynamics.Input) {
	if variableCache == nil {
		return map[rules.VariableFieldName]dynamics.Input{}
	}
	inputMap := make(map[rules.VariableFieldName]dynamics.Input)
	for key, pair := range variableCache {
		inp, valid := buildInput(pair)
		if valid {
			inputMap[key] = inp
		}
	}
	return inputMap
}

func (l *locator) UpdateCache(oldCache map[rules.VariableFieldName]rules.VariableValuePair, newValues map[rules.VariableFieldName]rules.VariableValuePair) map[rules.VariableFieldName]rules.VariableValuePair {
	for key, val := range newValues {
		oldCache[key] = val
	}
	return oldCache
}

func (l *locator) GetRecommendations(updatedVariables map[rules.VariableFieldName]rules.VariableValuePair, rulesCache map[string]rules.RuleMatrix, allVarsCache map[rules.VariableFieldName]rules.VariableValuePair) map[rules.VariableFieldName]rules.VariableValuePair {
	varList := buildVarList(updatedVariables)
	compliantVals := make(map[rules.VariableFieldName]dynamics.Input)
	compliantVals = l.TranslateToInputs(updatedVariables)
	for _, entry := range varList {
		rulesAffected, success := rules.PickUpRulesByVariable(entry, rulesCache, allVarsCache)
		if success {
			for _, ruleName := range rulesAffected {
				compliantInputs := dynamics.FindClosestApproach(rulesCache[ruleName], compliantVals)
				compliantVals = compliantInputs
			}
		}
	}
	for key, value := range compliantVals {
		updatedVariables[key] = convertInput(value)
	}
	return updatedVariables
}

func buildVarList(updatedVariables map[rules.VariableFieldName]rules.VariableValuePair) []rules.VariableFieldName {
	varList := []rules.VariableFieldName{}
	for key := range updatedVariables {
		varList = append(varList, key)
	}
	return varList
}

func convertInput(inp dynamics.Input) rules.VariableValuePair {
	return rules.VariableValuePair{
		VariableName: inp.Name,
		Values:       inp.Value,
	}
}

func buildInput(pair rules.VariableValuePair) (dynamics.Input, bool) {
	if isAdjustable, ok := IsChangeable()[pair.VariableName]; ok {
		return dynamics.Input{
			Name:             pair.VariableName,
			ClientAdjustable: isAdjustable,
			Value:            pair.Values,
		}, true
	} else {
		return dynamics.Input{
			Name:             pair.VariableName,
			ClientAdjustable: false,
			Value:            pair.Values,
		}, true
	}

}

func IsChangeable() map[rules.VariableFieldName]bool {
	return map[rules.VariableFieldName]bool{
		rules.NumberOfIslandsContributingToCommonPool: false,
		rules.NumberOfFailedForages:                   false,
		rules.NumberOfBrokenAgreements:                false,
		rules.MaxSeverityOfSanctions:                  false,
		rules.NumberOfIslandsAlive:                    false,
		rules.NumberOfBallotsCast:                     false,
		rules.NumberOfAllocationsSent:                 false,
		rules.AllocationRequestsMade:                  false,
		rules.AllocationMade:                          false,
		rules.IslandsAlive:                            false,
		rules.SpeakerSalary:                           true,
		rules.JudgeSalary:                             true,
		rules.PresidentSalary:                         true,
		rules.RuleSelected:                            false,
		rules.VoteCalled:                              false,
		rules.ExpectedTaxContribution:                 false,
		rules.ExpectedAllocation:                      false,
		rules.IslandTaxContribution:                   true,
		rules.IslandAllocation:                        true,
		rules.IslandReportedResources:                 false,
		rules.ConstSanctionAmount:                     false,
		rules.TurnsLeftOnSanction:                     false,
		rules.SanctionPaid:                            true,
		rules.SanctionExpected:                        false,
		rules.TestVariable:                            false,
		rules.JudgeInspectionPerformed:                false,
	}
}

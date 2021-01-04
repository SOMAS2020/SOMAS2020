package team3

import (
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
	newMetaStrategy.conquest = !currentParams.saveCriticalIsland && currentParams.aggression > 0.7 && currentParams.friendliness < 0.5
	newMetaStrategy.saviour = currentParams.saveCriticalIsland && currentParams.friendliness >= 0.5 && currentParams.aggression <= 0.7
	newMetaStrategy.democrat = currentParams.complianceLevel > 0.5 && currentParams.selfishness < 0.5
	newMetaStrategy.generous = currentParams.selfishness < 0.5 && currentParams.saveCriticalIsland && currentParams.recidivism < 0.5 && currentParams.friendliness > 0.5
	newMetaStrategy.lawful = currentParams.complianceLevel > 0.5 && currentParams.selfishness < 0.5
	newMetaStrategy.legislative = currentParams.equity > 0.5
	newMetaStrategy.executive = currentParams.recidivism < 0.5
	newMetaStrategy.fury = !currentParams.saveCriticalIsland && currentParams.aggression > 0.6 && currentParams.anger > 0.8
	newMetaStrategy.independent = currentParams.recidivism > 0.5
	l.locationStrategy = newMetaStrategy
}

func (l *locator) checkIfIdealLocationAvailable(ruleMat rules.RuleMatrix) (mat.VecDense, bool) {
	requiredVariables := ruleMat.RequiredVariables
	valid := true
	var idealPosition []float64
	for _, variable := range requiredVariables {
		values, isValid := l.switchDetermineFunction(variable)
		valid = valid && isValid
		idealPosition = append(idealPosition, values...)
	}
	if !valid {
		return mat.VecDense{}, false
	} else {
		return *mat.NewVecDense(len(idealPosition), idealPosition), true
	}
}

func (l *locator) switchDetermineFunction(name rules.VariableFieldName) ([]float64, bool) {
	switch name {
	case rules.NumberOfIslandsContributingToCommonPool:
		return l.idealNumberOfIslandsContributingToCommonPool()
	case rules.NumberOfIslandsAlive:
		return l.idealNumberOfIslandsAlive()
	case rules.NumberOfBallotsCast:
		return l.idealNumberOfBallotsCast()
	case rules.AllocationRequestsMade:
		return l.idealAllocationRequestsMade()
	case rules.AllocationMade:
		return l.idealAllocationsMade()
	case rules.IslandsAlive:
		return l.idealIslandsAlive()
	default:
		return []float64{}, false
	}
}

func (l *locator) idealNumberOfIslandsContributingToCommonPool() ([]float64, bool) {
	if l.locationStrategy.generous {
		return []float64{l.calcNumberOfNonCriticalIslands()}, true
	} else {
		return []float64{l.calcNumberOfAliveIslands()}, true
	}
}

func (l *locator) idealNumberOfIslandsAlive() ([]float64, bool) {
	if l.locationStrategy.conquest {
		return []float64{1}, true
	} else if l.locationStrategy.saviour {
		return []float64{l.calcNumberOfAliveIslands()}, true
	} else {
		return []float64{l.calcNumberOfNonCriticalIslands()}, true
	}
}

func (l *locator) idealNumberOfBallotsCast() ([]float64, bool) {
	if l.locationStrategy.democrat {
		return []float64{l.calcNumberOfAliveIslands()}, true
	} else if l.locationStrategy.conquest {
		return []float64{l.calcNumberOfNonCriticalIslands()}, true
	} else {
		return []float64{}, false
	}
}

func (l *locator) idealAllocationsMade() ([]float64, bool) {
	return []float64{0}, l.locationStrategy.independent
}

func (l *locator) idealAllocationRequestsMade() ([]float64, bool) {
	return []float64{0}, l.locationStrategy.independent
}

func (l *locator) idealIslandsAlive() ([]float64, bool) {
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

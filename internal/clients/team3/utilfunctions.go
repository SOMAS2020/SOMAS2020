package team3

import (
	"fmt"
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// A mock server handle used for tests
type mockServerReadHandle struct {
	gameState gamestate.ClientGameState
}

func (c *client) clientPrint(format string, a ...interface{}) {
	if printTeam3Logs {
		c.Logf("%v", fmt.Sprintf(format, a...))
	}
}

// getLocalResources retrieves our islands resrouces from server
func (c *client) getLocalResources() shared.Resources {
	currentState := c.BaseClient.ServerReadHandle.GetGameState()
	return currentState.ClientInfo.Resources
}

// getIslandsAlive retrives number of islands still alive
func (c *client) getIslandsAlive() int {
	var lifeStatuses map[shared.ClientID]shared.ClientLifeStatus
	var aliveCount int

	currentState := c.BaseClient.ServerReadHandle.GetGameState()
	lifeStatuses = currentState.ClientLifeStatuses
	for _, statusInfo := range lifeStatuses {
		if statusInfo == shared.Alive {
			aliveCount++
		}
	}
	return aliveCount
}

// areWeCritical returns whether our client is critical or not
func (c *client) areWeCritical() bool {
	currentStatus := c.BaseClient.ServerReadHandle.GetGameState().ClientLifeStatuses
	if currentStatus[c.GetID()] == shared.Critical {
		return true
	}
	return false
}

// getAverage returns the average of the list
func getAverage(lst []float64) float64 {

	if len(lst) == 0 {
		return 0.0
	}

	total := 0.0
	for _, val := range lst {
		total += val
	}

	return (float64(total) / float64(len(lst)))
}

func (m mockServerReadHandle) GetGameState() gamestate.ClientGameState {
	return m.gameState
}

// mostTrusted return the ClientID that corresponds to the highest trust value
func mostTrusted(values map[shared.ClientID][]float64) shared.ClientID {
	var max = -math.MaxFloat64
	var mostTrustedClient shared.ClientID

	for client_id, trust_history := range values {
		var last_trust = trust_history[len(trust_history)-1]
		if last_trust > max {
			max = last_trust
			mostTrustedClient = client_id
		}
	}
	return mostTrustedClient
}

// leastTrusted return the ClientID that corresponds to the smallest trust value
func leastTrusted(values map[shared.ClientID][]float64) shared.ClientID {
	var min = math.MaxFloat64
	var leastTrustedClient shared.ClientID

	for client_id, trust_history := range values {
		var last_trust = trust_history[len(trust_history)-1]
		if last_trust > min {
			min = last_trust
			leastTrustedClient = client_id
		}
	}
	return leastTrustedClient
}

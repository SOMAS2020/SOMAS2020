package team3

import (
	"fmt"
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

//------------------------- Mock Server Read Handle ----------------------------
// A mock server handle used for tests
type mockServerReadHandle struct {
	gameState  gamestate.ClientGameState
	gameConfig config.ClientConfig
}

func (m mockServerReadHandle) GetGameState() gamestate.ClientGameState {
	return m.gameState
}

func (m mockServerReadHandle) GetGameConfig() config.ClientConfig {
	return m.gameConfig
}

// -----------------------------------------------------------------------------

// clientPrint is a wrapper for team3 Logf function, that only prints when
// printTeam3Logs == true
func (c *client) clientPrint(format string, a ...interface{}) {
	turn := c.ServerReadHandle.GetGameState().Turn
	season := c.ServerReadHandle.GetGameState().Season
	if printTeam3Logs {
		c.Logf("[S:%v T%v] %v", season, turn, fmt.Sprintf(format, a...))
	}
}

// getLocalResources retrieves our islands resrouces from server
func (c *client) getLocalResources() shared.Resources {
	currentState := c.BaseClient.ServerReadHandle.GetGameState()
	return currentState.ClientInfo.Resources
}

// getIslandsAliveCount retrives number of islands still alive
func (c *client) getIslandsAliveCount() int {
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
	return c.ServerReadHandle.GetGameState().ClientInfo.Resources == shared.Resources(shared.Critical)
}

// getAliveIslands returns list of alive islands
func (c *client) getAliveIslands() []shared.ClientID {
	var lifeStatuses map[shared.ClientID]shared.ClientLifeStatus
	var aliveIslands = []shared.ClientID{}

	currentState := c.BaseClient.ServerReadHandle.GetGameState()
	lifeStatuses = currentState.ClientLifeStatuses
	for id, statusInfo := range lifeStatuses {
		if statusInfo == shared.Alive {
			aliveIslands = append(aliveIslands, id)
		}
	}
	return aliveIslands
}

// getIslandsCriticalCount retrives number of islands that are critical
func (c *client) getIslandsCriticalCount() int {
	var lifeStatuses map[shared.ClientID]shared.ClientLifeStatus
	var criticalCount int

	currentState := c.BaseClient.ServerReadHandle.GetGameState()
	lifeStatuses = currentState.ClientLifeStatuses
	for _, statusInfo := range lifeStatuses {
		if statusInfo == shared.Critical {
			criticalCount++
		}
	}
	return criticalCount
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

// mostTrusted return the ClientID that corresponds to the highest trust value
func mostTrusted(values map[shared.ClientID]float64) shared.ClientID {
	var max = -math.MaxFloat64
	var mostTrustedClient shared.ClientID

	for clientID, trustScore := range values {
		if trustScore > max {
			max = trustScore
			mostTrustedClient = clientID
		}
	}
	return mostTrustedClient
}

// leastTrusted return the ClientID that corresponds to the smallest trust value
func leastTrusted(values map[shared.ClientID]float64) shared.ClientID {
	var min = math.MaxFloat64
	var leastTrustedClient shared.ClientID

	for clientID, trustScore := range values {
		if trustScore > min {
			min = trustScore
			leastTrustedClient = clientID
		}
	}
	return leastTrustedClient
}

// isClientStatusCritical returns whether or not the ClientID is critical
func (c *client) isClientStatusCritical(ClientID shared.ClientID) bool {
	return c.ServerReadHandle.GetGameState().ClientLifeStatuses[ClientID] == shared.Critical
}

// getIIGOCost returns the miniumum resources in common pool needed to fully run IIGO
func (c *client) getIIGOCost() shared.Resources {
	iigoConfig := c.ServerReadHandle.GetGameConfig().IIGOClientConfig
	sum := shared.Resources(0)

	sum += iigoConfig.GetRuleForSpeakerActionCost
	sum += iigoConfig.BroadcastTaxationActionCost
	sum += iigoConfig.ReplyAllocationRequestsActionCost
	sum += iigoConfig.RequestAllocationRequestActionCost
	sum += iigoConfig.RequestRuleProposalActionCost
	sum += iigoConfig.AppointNextSpeakerActionCost
	sum += iigoConfig.InspectHistoryActionCost
	sum += iigoConfig.InspectBallotActionCost
	sum += iigoConfig.InspectAllocationActionCost
	sum += iigoConfig.AppointNextPresidentActionCost
	sum += iigoConfig.SetVotingResultActionCost
	sum += iigoConfig.SetRuleToVoteActionCost
	sum += iigoConfig.AnnounceVotingResultActionCost
	sum += iigoConfig.UpdateRulesActionCost
	sum += iigoConfig.AppointNextJudgeActionCost
	return sum
}

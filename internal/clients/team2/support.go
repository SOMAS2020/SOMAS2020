package team2

import (
	"math"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) criticalStatus() bool {
	clientInfo := c.gameState().ClientInfo
	if clientInfo.LifeStatus == shared.Critical {
		return true
	}
	return false
}

// If a disaster is reported, append the turn and report of the latest disaster to the disaster history
func (c *client) DisasterNotification(report disasters.DisasterReport, effects disasters.DisasterEffects) {
	disaster := DisasterOccurrence{
		Turn:   c.gameState().Turn,
		Report: report,
	}

	c.disasterHistory = append(c.disasterHistory, disaster)
}

// getIslandsToShareWith returns a slice of the islands we want to share our prediction with.
// We decided to always share our prediction with all islands to improve arhcipelago decisions as a whole.
func (c *client) getIslandsToShareWith() []shared.ClientID {
	islandsToShareWith := make([]shared.ClientID, len(shared.TeamIDs))
	for index, id := range shared.TeamIDs {
		islandsToShareWith[index] = id
	}
	return islandsToShareWith
}

//checkOthersCrit checks if anyone else is critical
func checkOthersCrit(c *client) bool {
	for clientID, status := range c.gameState().ClientLifeStatuses {
		if status == shared.Critical && clientID != c.GetID() {
			return true
		}
	}
	return false
}

func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

func (c *client) gameConfig() config.ClientConfig {
	return c.BaseClient.ServerReadHandle.GetGameConfig()
}

func (c *client) getAliveClients() []shared.ClientID {
	clientStatuses := c.gameState().ClientLifeStatuses
	aliveClients := make([]shared.ClientID, 0)
	for island, status := range clientStatuses {
		if status != shared.Dead {
			aliveClients = append(aliveClients, island)
		}
	}
	return aliveClients
}

func (c *client) getNumAliveClients() int {
	return len(c.getAliveClients())
}

// Stores the information we receive from IIGO
func (c *client) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	c.Communications[sender] = append(c.Communications[sender], data)
	for contentType, content := range data {
		commonPool := CommonPoolInfo{
			tax:             0,
			requestedToPres: 0,
			allocatedByPres: 0,
			takenFromCP:     0,
		}
		switch contentType {
		// How much tax we must pay
		case shared.IIGOTaxDecision:
			if _, ok := c.presCommonPoolHist[c.gameState().PresidentID]; !ok {
				// Create map if it doesn't exist
				c.presCommonPoolHist[c.gameState().PresidentID] = make(map[uint]CommonPoolInfo)

			} else {
				presHist := c.presCommonPoolHist[c.gameState().PresidentID]
				if pastInfo, ok := presHist[c.gameState().Turn]; ok {
					//Take previous values if they exist in the map
					commonPool = pastInfo
				}
			}
			c.taxAmount = shared.Resources(content.IntegerData)
			commonPool.tax = shared.Resources(content.IntegerData)
			// Update the history
			c.presCommonPoolHist[c.gameState().PresidentID][c.gameState().Turn] = commonPool

		// How many resources we've been allocated from the CP by the President
		case shared.IIGOAllocationDecision:
			if _, ok := c.presCommonPoolHist[c.gameState().PresidentID]; !ok {
				// Create map if it doesn't exist
				c.presCommonPoolHist[c.gameState().PresidentID] = make(map[uint]CommonPoolInfo)

			} else {
				presHist := c.presCommonPoolHist[c.gameState().PresidentID]
				if pastInfo, ok := presHist[c.gameState().Turn]; ok {
					//Take previous values if they exist in the map
					commonPool = pastInfo
				}
			}
			c.commonPoolAllocation = shared.Resources(content.IntegerData)
			commonPool.allocatedByPres = shared.Resources(content.IntegerData)
			// Update the history
			c.presCommonPoolHist[c.gameState().PresidentID][c.gameState().Turn] = commonPool

		// What islands have a sanction (and the sanction tier)
		case shared.SanctionClientID:
			islandSanc := IslandSanctionInfo{
				Turn:   c.gameState().Turn,
				Tier:   (content.IntegerData),
				Amount: 0,
			}
			c.islandSanctions[shared.ClientID(content.IntegerData)] = islandSanc

		// The sanction tier "score" for this turn
		case shared.IIGOSanctionTier:
			c.tierLevels[content.IntegerData] = data[shared.IIGOSanctionScore].IntegerData

		// What sanction score we have
		case shared.SanctionAmount:
			if _, ok := c.sanctionHist[c.gameState().JudgeID]; !ok {
				c.sanctionHist[c.gameState().JudgeID] = make([]IslandSanctionInfo, 0)

			}
			sanction := IslandSanctionInfo{
				Turn:   c.gameState().Turn,
				Tier:   c.checkSanctionTier(content.IntegerData),
				Amount: content.IntegerData,
			}
			// Add a new sanction to the sanction hist
			sanctions := c.sanctionHist[c.gameState().JudgeID]
			c.sanctionHist[c.gameState().JudgeID] = append(sanctions, sanction)
		default:
			// will NOT execute logic for other conditions
		}
	}
}

// Gets the current AgentStrategy and returns an AgentStrategy Type
func (c *client) getAgentStrategy() AgentStrategy {
	return c.currStrategy
}

// TODO: it makes more sense to start by free riding to make our selves secure tbh
// Sets the current AgentStrategy and returns an AgentStrategy Type
func (c *client) setAgentStrategy() AgentStrategy {
	currTurn := c.gameState().Turn
	// Factor the common pool must increase by for us to considered free riding
	freeRide := shared.Resources(c.config.SwitchToSelfishFactor)

	// Factor the common pool must drop by for us to consider altruist
	altFactor := c.config.SwitchToSelfishFactor

	// Explore and test limits by playing a selfish strategy for a few turns
	if currTurn <= c.config.SelfishStartTurns {
		return Selfish
	}

	if len(c.commonPoolHistory) != 0 {
		runningMean := shared.Resources(0)

		for _, resources := range c.commonPoolHistory {
			runningMean = runningMean + (resources-runningMean)/shared.Resources(c.gameState().Turn)
		}

		// Percentage change in common pool from previous running mean
		percentageChange := (c.commonPoolHistory[c.gameState().Turn] - runningMean) / runningMean

		// if the pool decreases on average by a factor above altFactor set AgentStrategy to Altruist
		if percentageChange < 0 && math.Abs(float64(percentageChange)) > altFactor {
			return Altruist
		} else if percentageChange > 0 && percentageChange > freeRide {
			return Selfish
		}
	}

	return FairSharer
}

// Takes as input a sanction score and returns what sanction tier is corresponds to from the latest score-sanction threshold we have
// Used to store what sanction tier we're in from the IIGO Communications
func (c *client) checkSanctionTier(score int) int {
	var keys []int

	for k := range c.tierLevels {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for tier := range keys {
		if score >= c.tierLevels[tier] {
			return tier
		}
	}

	// NoSanction
	return 5
}

func Max(i shared.Resources, j shared.Resources) shared.Resources {
	if i >= j {
		return i
	} else {
		return j
	}
}

func Min(i shared.Resources, j shared.Resources) shared.Resources {
	if i < j {
		return i
	} else {
		return j
	}
}

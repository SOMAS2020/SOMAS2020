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

//TODO: how does this work?
func (c *client) DisasterNotification(report disasters.DisasterReport, effects disasters.DisasterEffects) {
	disaster := DisasterOccurrence{
		Turn:   c.gameState().Turn,
		Report: report,
	}
	c.disasterHistory = append(c.disasterHistory, disaster)
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

// MethodOfPlay determines which state we are in: 0=altruist, 1=fair sharer and 2= free rider
// TODO: empathy levels are being used incorrectly here
func (c *client) MethodOfPlay() int {
	currTurn := c.gameState().Turn

	// how many turns at the beginning we cannot free ride for
	minTurnsUntilFreeRide := NoFreeRideAtStart
	// what factor the common pool must increase by for us to considered free riding
	freeRide := shared.Resources(SwitchToFreeRideFactor)
	// what factor the common pool must drop by for us to consider altruist
	altFactor := SwitchToAltruistFactor

	runMeanCommonPool := shared.Resources(0.0)
	div := shared.Resources(0.0)

	for pastTurn, resources := range c.commonPoolHistory {
		if pastTurn == currTurn {
			continue
		}
		diffTurn := shared.Resources(c.gameState().Turn - pastTurn)
		div++

		runMeanCommonPool += (resources/(diffTurn+1) - runMeanCommonPool) / div
	}

	changeCommonPool := (c.commonPoolHistory[currTurn] - runMeanCommonPool) / runMeanCommonPool

	if changeCommonPool < 0 && math.Abs(float64(changeCommonPool)) > altFactor {
		//altruist
		return 0
	} else if changeCommonPool > 0 && changeCommonPool > freeRide && currTurn > minTurnsUntilFreeRide {
		// Free rider
		return 2
	}

	// Default case: Fair Sharer
	return 1
}

func (c *client) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	c.Communications[sender] = append(c.Communications[sender], data)

	for contentType, content := range data {
		switch contentType {
		case shared.IIGOTaxDecision:
			var commonPool CommonPoolInfo
			if _, ok := c.presCommonPoolHist[c.gameState().PresidentID]; !ok {
				c.presCommonPoolHist[c.gameState().PresidentID] = make([]CommonPoolInfo, 0)

			} else {
				presHist := c.presCommonPoolHist[c.gameState().PresidentID]
				if presHist[len(presHist)-1].turn == c.gameState().Turn {
					commonPool = presHist[len(presHist)-1]
				}
			}
			presHist := c.presCommonPoolHist[c.gameState().PresidentID]
			c.taxAmount = shared.Resources(content.IntegerData)

			commonPool = CommonPoolInfo{
				tax:  shared.Resources(content.IntegerData),
				turn: c.gameState().Turn,
			}
			c.presCommonPoolHist[c.gameState().PresidentID] = append(presHist, commonPool)

		case shared.IIGOAllocationDecision:
			var commonPool CommonPoolInfo
			if _, ok := c.presCommonPoolHist[c.gameState().PresidentID]; !ok {
				c.presCommonPoolHist[c.gameState().PresidentID] = make([]CommonPoolInfo, 0)

			} else {
				presHist := c.presCommonPoolHist[c.gameState().PresidentID]
				if presHist[len(presHist)-1].turn == c.gameState().Turn {
					commonPool = presHist[len(presHist)-1]
				}
			}
			presHist := c.presCommonPoolHist[c.gameState().PresidentID]
			c.commonPoolAllocation = shared.Resources(content.IntegerData)
			commonPool = CommonPoolInfo{
				allocatedByPres: shared.Resources(content.IntegerData),
				turn:            c.gameState().Turn,
			}
			c.presCommonPoolHist[c.gameState().PresidentID] = append(presHist, commonPool)

		case shared.SanctionClientID:
			islandSanc := IslandSanctionInfo{
				Turn:   c.gameState().Turn,
				Tier:   (content.IntegerData),
				Amount: 0,
			}
			c.islandSanctions[shared.ClientID(content.IntegerData)] = islandSanc
		case shared.IIGOSanctionTier:
			c.tierLevels[content.IntegerData] = data[shared.IIGOSanctionScore].IntegerData
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

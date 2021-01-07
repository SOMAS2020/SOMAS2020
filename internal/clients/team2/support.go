package team2

import (
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) islandEmpathyLevel() EmpathyLevel {
	clientInfo := c.gameState().ClientInfo

	// switch statement to toggle between three levels
	// change our state based on these cases
	switch {
	case clientInfo.LifeStatus == shared.Critical:
		return Selfish
		// replace with some expression
	case (true):
		return Altruist
	default:
		return FairSharer
	}
}

func criticalStatus(c *client) bool {
	clientInfo := c.gameState().ClientInfo
	if clientInfo.LifeStatus == shared.Critical { //not sure about shared.Critical
		return true
	}
	return false
}

//TODO: how does this work?
func (c *client) DisasterNotification(report disasters.DisasterReport, effects disasters.DisasterEffects) {
	c.disasterHistory[len(c.disasterHistory)] = DisasterOccurence{
		Turn:   float64(c.gameState().Turn),
		Report: report,
	}
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

//MethodOfPlay determine which state we are in 0=altruist, 1=fair sharer and 2= free rider
func (c *client) MethodOfPlay() int {
	ResourceHistory := c.commonPoolHistory
	turn := c.gameState().Turn

	var no_freeride float64 = 3 //how many turns at the beginning we cannot free ride for
	var freeride float64 = 5    //what factor the common pool must increase by for us to considered free riding
	var altfactor float64 = 5   //what factor the common pool must drop by for us to consider altruist

	if turn == 1 { //if there is no historical data then use default strategy
		return 1
	}

	prevTurn := turn - 1
	prevTurn2 := turn - 2
	if ResourceHistory[prevTurn] > (ResourceHistory[turn] * altfactor) { //decreasing common pool means consider altruist
		if ResourceHistory[prevTurn2] > (ResourceHistory[prevTurn] * altfactor) {
			return 0 //altruist
		}
	}

	if float64(turn) > no_freeride { //we will not allow ourselves to use free riding at the start of the game
		if (ResourceHistory[prevTurn] * freeride) < ResourceHistory[turn] {
			if (ResourceHistory[prevTurn2] * freeride) < ResourceHistory[prevTurn] { //two large jumps then we free ride
				return 2 //free rider
			}
		}
	}

	// Else if neither
	return 1
}

func (c *client) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	c.Communications[sender] = append(c.Communications[sender], data)

	for contentType, content := range data {
		switch contentType {
		case shared.IIGOTaxDecision:
			var commonPool CommonPoolInfo
			if val, ok := c.commonPoolHist[c.gameState().Turn]; ok {
				val.tax = shared.Resources(content.IntegerData)
			} else {
				commonPool = CommonPoolInfo{
					tax: shared.Resources(content.IntegerData),
				}
			}
			c.commonPoolHist[c.gameState().Turn] = commonPool
			c.taxAmount = shared.Resources(content.IntegerData)
		case shared.IIGOAllocationDecision:
			var commonPool CommonPoolInfo
			if val, ok := c.commonPoolHist[c.gameState().Turn]; ok {
				val.allocatedByPres = shared.Resources(content.IntegerData)
			} else {
				commonPool = CommonPoolInfo{
					allocatedByPres: shared.Resources(content.IntegerData),
				}
			}
			c.commonPoolHist[c.gameState().Turn] = commonPool
			c.commonPoolAllocation = shared.Resources(content.IntegerData)
		// case shared.RuleName:
		// 	currentRuleID := content.TextData
		// 	// Rule voting
		// 	if _, ok := data[shared.RuleVoteResult]; ok {
		// 		if _, ok := c.iigoInfo.ruleVotingResults[currentRuleID]; ok {
		// 			c.iigoInfo.ruleVotingResults[currentRuleID].resultAnnounced = true
		// 			c.iigoInfo.ruleVotingResults[currentRuleID].result = data[shared.RuleVoteResult].BooleanData
		// 		} else {
		// 			c.iigoInfo.ruleVotingResults[currentRuleID] = &ruleVoteInfo{resultAnnounced: true, result: data[shared.RuleVoteResult].BooleanData}
		// 		}
		// 	}
		// 	// Rule sanctions
		// 	if _, ok := data[shared.IIGOSanctionScore]; ok {
		// 		// c.clientPrint("Received sanction info: %+v", data)
		// 		c.iigoInfo.sanctions.rulePenalties[currentRuleID] = roles.IIGOSanctionScore(data[shared.IIGOSanctionScore].IntegerData)
		// 	}

		// TODO: decide if this is worth it
		// case shared.RoleMonitored:
		// 	c.iigoInfo.monitoringDeclared[content.IIGORoleData] = true
		// 	c.iigoInfo.monitoringOutcomes[content.IIGORoleData] = data[shared.MonitoringResult].BooleanData
		case shared.SanctionClientID:
			sanction := IslandSanctionInfo{
				Turn: c.gameState().Turn,
				Tier: roles.IIGOSanctionTier(data[shared.IIGOSanctionTier].IntegerData),
			}
			c.islandSanctions[shared.ClientID(content.IntegerData)] = append(c.islandSanctions[shared.ClientID(content.IntegerData)], sanction)
		case shared.IIGOSanctionTier:
			c.tierLevels[roles.IIGOSanctionTier(content.IntegerData)] = roles.IIGOSanctionScore(data[shared.IIGOSanctionScore].IntegerData)
		case shared.SanctionAmount:
			sanction := IslandSanctionInfo{
				Turn:   c.gameState().Turn,
				Tier:   c.checkSanctionTier(roles.IIGOSanctionScore(content.IntegerData)),
				Amount: roles.IIGOSanctionScore(content.IntegerData),
			}
			c.sanctionHist = append(c.sanctionHist, sanction)
		}
	}

}

type TierList []roles.IIGOSanctionTier

func (p TierList) Len() int           { return len(p) }
func (p TierList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TierList) Less(i, j int) bool { return p[i] < p[j] }

func (c *client) checkSanctionTier(score roles.IIGOSanctionScore) roles.IIGOSanctionTier {

	var keys TierList
	for k := range c.tierLevels {
		keys = append(keys, k)
	}

	sort.Sort(keys)

	for tier := range keys {
		if score >= c.tierLevels[roles.IIGOSanctionTier(tier)] {
			return roles.IIGOSanctionTier(tier)
		}
	}
	return roles.NoSanction
}

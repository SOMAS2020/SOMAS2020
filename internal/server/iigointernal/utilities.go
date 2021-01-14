package iigointernal

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func broadcastToAllIslands(clients map[shared.ClientID]baseclient.Client, sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent, gameState gamestate.GameState) {
	islandsAlive := gameState.RulesInfo.VariableMap[rules.IslandsAlive]
	for _, v := range islandsAlive.Values {
		communicateWithIslands(clients, shared.TeamIDs[int(v)], sender, data)
	}
}

// func setIIGOClients(clientMap *map[shared.ClientID]baseclient.Client) {
// 	iigoClients = *clientMap
// }

func communicateWithIslands(clients map[shared.ClientID]baseclient.Client, recipientID shared.ClientID, senderID shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {

	if recipientClient, ok := clients[recipientID]; ok {
		recipientClient.ReceiveCommunication(senderID, data)
	}

}

func CheckEnoughInCommonPool(value shared.Resources, gameState *gamestate.GameState) bool {
	return gameState.CommonPool >= value
}

func WithdrawFromCommonPool(value shared.Resources, gameState *gamestate.GameState) (withdrawnAmount shared.Resources, withdrawSuccesful bool) {
	if CheckEnoughInCommonPool(value, gameState) {
		gameState.CommonPool -= value
		return value, true
	} else {
		return shared.Resources(0.0), false
	}
}

func depositIntoClientPrivatePool(value shared.Resources, id shared.ClientID, state *gamestate.GameState) {
	participantInfo := state.ClientInfos[id]
	participantInfo.Resources += value
	state.ClientInfos[id] = participantInfo
}

func depositIntoCommonPool(value shared.Resources, state *gamestate.GameState) {
	state.CommonPool += value
}

func Contains(islandIDSlice []shared.ClientID, islandID shared.ClientID) bool {
	for _, a := range islandIDSlice {
		if a == islandID {
			return true
		}
	}
	return false
}

func boolToFloat(input bool) float64 {
	if input {
		return 1
	}
	return 0
}

func copyClientList(in []shared.ClientID) []shared.ClientID {
	ret := make([]shared.ClientID, len(in))
	copy(ret, in)
	return ret
}

// if an IIGO role is dead, it is replaced with a random living island
func removeDeadBodiesFromOffice(g *gamestate.GameState) {
	aliveClientIds := []shared.ClientID{}
	for clientID, clientGameState := range g.ClientInfos {
		if clientGameState.LifeStatus != shared.Dead {
			aliveClientIds = append(aliveClientIds, clientID)
		}
	}
	if g.ClientInfos[g.PresidentID].LifeStatus == shared.Dead {
		g.PresidentID = aliveClientIds[rand.Intn(len(aliveClientIds))]
	}
	if g.ClientInfos[g.JudgeID].LifeStatus == shared.Dead {
		g.JudgeID = aliveClientIds[rand.Intn(len(aliveClientIds))]
	}
	if g.ClientInfos[g.SpeakerID].LifeStatus == shared.Dead {
		g.SpeakerID = aliveClientIds[rand.Intn(len(aliveClientIds))]
	}
}

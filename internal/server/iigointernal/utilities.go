package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

func broadcastToAllIslands(sender shared.ClientID, data map[int]baseclient.Communication) {
	islandsAlive := rules.VariableMap[rules.IslandsAlive]
	for _, v := range islandsAlive.Values {
		communicateWithIslands(shared.TeamIDs[int(v)], sender, data)
	}
}

func setIIGOClients(clientMap *map[shared.ClientID]baseclient.Client) {
	iigoClients = *clientMap
}

func communicateWithIslands(recipientID shared.ClientID, senderID shared.ClientID, data map[int]baseclient.Communication) {

	clients := iigoClients

	if recipientClient, ok := clients[recipientID]; ok {
		recipientClient.ReceiveCommunication(senderID, data)
	}

}

func CheckEnoughInCommonPool(value shared.Resources, gameState *gamestate.GameState) bool {
	return gameState.CommonPool >= value
}

func WithdrawFromCommonPool(value shared.Resources, gameState *gamestate.GameState) error {
	if CheckEnoughInCommonPool(value, gameState) {
		gameState.CommonPool -= value
		return nil
	} else {
		return errors.Errorf("Not enough resources in the common pool to withdraw the amount '%v'", value)
	}
}

func depositIntoCommonPool(value shared.Resources, state *gamestate.GameState) {
	state.CommonPool += value
}

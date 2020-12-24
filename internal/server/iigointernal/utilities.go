package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

// PickUpRulesByVariable returns a list of rule_id's which are affected by certain variables
func PickUpRulesByVariable(variableName string, ruleStore map[string]rules.RuleMatrix) ([]string, error) {
	var Rules []string
	if _, ok := rules.VariableMap[variableName]; ok {
		for k, v := range ruleStore {
			_, err := searchForStringInArray(variableName, v.RequiredVariables)
			if err != nil {
				Rules = append(Rules, k)
			}
		}
		return Rules, nil
	} else {
		return []string{}, errors.Errorf("Variable name '%v' was not found in the variable cache", variableName)
	}
}

func searchForStringInArray(val string, array []string) (int, error) {
	for i, v := range array {
		if v == val {
			return i, nil
		}
	}
	return 0, errors.Errorf("Not found")
}

func broadcastToAllIslands(sender shared.ClientID, data map[int]baseclient.Communication) {
	islandsAlive := rules.VariableMap["islands_alive"]
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

const (
	BallotID = iota
	PresidentAllocationCheck
	SpeakerID
	RoleConducted
	ResAllocID
	SpeakerBallotCheck
	PresidentID
	RuleName
	RuleVoteResult
	TaxAmount
	AllocationAmount
)

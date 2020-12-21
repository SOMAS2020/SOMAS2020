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

type DataPacket struct {
	integerData int
	textData    string
	booleanData bool
}

type Communication struct {
	recipient int
	sender    int
	data      map[int]DataPacket
}

func broadcastToAllIslands(sender int, data map[int]DataPacket) {
	islandsAlive := rules.VariableMap["islands_alive"]
	for _, v := range islandsAlive.Values {
		communicateWithIslands(int(v), sender, data)
	}
}

func dataPacketToCommunication(d *DataPacket) baseclient.Communication {
	return baseclient.Communication{
		IntegerData: d.integerData,
		TextData:    d.textData,
		BooleanData: d.booleanData,
	}
}

func setIIGOClients(clientMap *map[shared.ClientID]baseclient.Client) {
	iigoClients = *clientMap
}

func communicateWithIslands(recipient int, sender int, data map[int]DataPacket) {
	// for client := range []int{recipient, sender} {
	// 	if client > len(shared.TeamIDs) {
	// 		return errors.Errorf("%v is not a valid TeamID", client)
	// 	}
	// }

	communication := map[int]baseclient.Communication{}
	for k, v := range data {
		communication[k] = dataPacketToCommunication(&v)
	}

	recipientID := shared.TeamIDs[recipient]
	senderID := shared.TeamIDs[sender]
	clients := iigoClients

	if recipientClient, ok := clients[recipientID]; ok {
		recipientClient.ReceiveCommunication(senderID, communication)
	}

}

func CheckEnoughInCommonPool(value int, gameState *gamestate.GameState) bool {
	return gameState.CommonPool >= value
}

func WithdrawFromCommonPool(value int, gameState *gamestate.GameState) error {
	if CheckEnoughInCommonPool(value, gameState) {
		gameState.CommonPool -= value
		return nil
	} else {
		return errors.Errorf("Not enough resources in the common pool to withdraw the amount '%v'", value)
	}
}

func withdrawSalary(value int, gameState *gamestate.GameState) (int, error) {
	return value, WithdrawFromCommonPool(value, gameState)
}

const (
	BallotID                 = iota
	PresidentAllocationCheck = iota
	SpeakerID                = iota
	RoleConducted            = iota
	ResAllocID               = iota
	SpeakerBallotCheck       = iota
	PresidentID              = iota
	RuleName                 = iota
	RuleVoteResult           = iota
	TaxAmount                = iota
	AllocationAmount         = iota
)

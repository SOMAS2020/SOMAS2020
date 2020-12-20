package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
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

type Communication struct {
	recipient int
	sender    int
	data      map[int]int
}

func broadcastToAllIslands(sender int, data map[int]int) {
	islandsAlive := rules.VariableMap["islands_alive"]
	for _, v := range islandsAlive.Values {
		communicateWithIslands(int(v), sender, data)
	}
}

func communicateWithIslands(recipient int, sender int, data map[int]int) {
	communication := Communication{
		recipient: recipient,
		sender:    sender,
		data:      data,
	}
	//Send to islands
	print(communication) //// Get rid of this
}

func collapseBoolean(val bool) int {
	if val {
		return 1
	} else {
		return 0
	}
}

func CheckEnoughInCommonPool(value int, gameState *common.GameState) bool {
	return gameState.CommonPool >= value
}

func WithdrawFromCommonPool(value int, gameState *common.GameState) error {
	if CheckEnoughInCommonPool(value, gameState) {
		gameState.CommonPool -= value
		return nil
	} else {
		return errors.Errorf("Not enough ressources in the common pool to withdraw the amount '%v'", value)
	}
}

func withdrawSalary(value int, gameState *common.GameState) (int, error) {
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
)

package gamestate

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

var TurnHistory []shared.Accountability

func UpdateTurnHistory(clientID shared.ClientID, pairs []rules.VariableValuePair) error {
	TurnHistory = append(TurnHistory, shared.Accountability{
		ClientID: clientID,
		Pairs:    pairs,
	})
	return nil
}

func clearHistory() {
	TurnHistory = []shared.Accountability{}
}

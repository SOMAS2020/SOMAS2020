package gamestate

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Accountability struct {
	ClientID shared.ClientID
	Pairs    []rules.VariableValuePair
}

var TurnHistory []Accountability

func UpdateTurnHistory(clientID shared.ClientID, pairs []rules.VariableValuePair) error {
	TurnHistory = append(TurnHistory, Accountability{
		ClientID: clientID,
		Pairs:    pairs,
	})
	return nil
}

func clearHistory() {
	TurnHistory = []Accountability{}
}

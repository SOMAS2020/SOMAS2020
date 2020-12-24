package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Accountability struct {
	clientID shared.ClientID
	pairs    []rules.VariableValuePair
}

var TurnHistory []Accountability

func UpdateTurnHistory(clientID shared.ClientID, pairs []rules.VariableValuePair) error {
	TurnHistory = append(TurnHistory, Accountability{
		clientID: clientID,
		pairs:    pairs,
	})
	return nil
}

func clearHistory() {
	TurnHistory = []Accountability{}
}

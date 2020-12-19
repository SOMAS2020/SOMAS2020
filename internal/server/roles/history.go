package roles

import "github.com/SOMAS2020/SOMAS2020/internal/common/rules"

type Accountability struct {
	clientID int
	pairs    []rules.VariableValuePair
}

var TurnHistory []Accountability

func updateTurnHistory(clientID int, pairs []rules.VariableValuePair) error {
	TurnHistory = append(TurnHistory, Accountability{
		clientID: clientID,
		pairs:    pairs,
	})
	return nil
}

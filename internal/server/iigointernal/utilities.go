package iigointernal

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

func broadcastToAllIslands(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	islandsAlive := rules.VariableMap[rules.IslandsAlive]
	for _, v := range islandsAlive.Values {
		communicateWithIslands(shared.TeamIDs[int(v)], sender, data)
	}
}

func setIIGOClients(clientMap *map[shared.ClientID]baseclient.Client) {
	iigoClients = *clientMap
}

func communicateWithIslands(recipientID shared.ClientID, senderID shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {

	clients := iigoClients

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

type conversionType int

const (
	tax conversionType = iota
	allocation
)

// convertAmount takes the amount of tax/allocation and converts it into appropriate variable and rule ready to be sent to the client
func convertAmount(amount shared.Resources, amountType conversionType) (rules.VariableValuePair, rules.RuleMatrix) {
	var reqVar rules.VariableFieldName
	var sentVar rules.VariableFieldName
	var ruleVariables []rules.VariableFieldName
	if amountType == tax {
		reqVar, sentVar = rules.IslandTaxContribution, rules.ExpectedTaxContribution
		// Rule in form IslandTaxContribution - ExpectedTaxContribution >= 0
		ruleVariables = []rules.VariableFieldName{reqVar, sentVar}
	} else if amountType == allocation {
		reqVar, sentVar = rules.IslandAllocation, rules.ExpectedAllocation
		// Rule in form ExpectedAllocation - IslandAllocation >= 0
		ruleVariables = []rules.VariableFieldName{sentVar, reqVar}
	}

	v := []float64{1, -1, 0}
	aux := []float64{2}

	rowLength := len(ruleVariables) + 1
	nrows := len(v) / rowLength

	CoreMatrix := mat.NewDense(nrows, rowLength, v)
	AuxiliaryVector := mat.NewVecDense(nrows, aux)

	retRule := rules.RuleMatrix{
		RuleName:          fmt.Sprintf("%s = %.2f", sentVar, amount),
		RequiredVariables: ruleVariables,
		ApplicableMatrix:  *CoreMatrix,
		AuxiliaryVector:   *AuxiliaryVector,
		Mutable:           false,
	}

	retVar := rules.VariableValuePair{
		VariableName: sentVar,
		Values:       []float64{float64(amount)},
	}

	return retVar, retRule
}

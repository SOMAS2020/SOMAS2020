package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// ------ TODO: COMPULSORY ------
func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {
	return c.BaseClient.MakeDisasterPrediction()
}

// ------ TODO: COMPULSORY ------
func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	c.BaseClient.ReceiveDisasterPredictions(receivedPredictions)
}

// OPTIONAL
func (c *client) MakeForageInfo() shared.ForageShareInfo {
	var shareTo []shared.ClientID // containing agents our agent wish to share informationwith

	for id, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		if status != shared.Dead {
			shareTo = append(shareTo, id)
		}
	}

	var lastDecision shared.ForageDecision
	var lastForageOut shared.Resources

	for forageType, results := range c.forageHistory {
		for _, result := range results {
			if uint(result.turn) == c.ServerReadHandle.GetGameState().Turn-1 {
				lastForageOut = result.forageReturn
				lastDecision = shared.ForageDecision{
					Type:         forageType,
					Contribution: result.forageIn,
				}
			}
		}
	}

	forageInfo := shared.ForageShareInfo{
		DecisionMade:     lastDecision,
		ResourceObtained: lastForageOut,
		ShareTo:          shareTo,
	}

	return forageInfo
}

func (c *client) ReceiveForageInfo(forageInfo []shared.ForageShareInfo) {
	for _, val := range forageInfo {
		c.forageHistory[val.DecisionMade.Type] =
			append(
				c.forageHistory[val.DecisionMade.Type],
				ForageResults{
					forageIn:     val.DecisionMade.Contribution,
					forageReturn: val.ResourceObtained,
					turn:         c.ServerReadHandle.GetGameState().Turn,
				},
			)
	}
}

package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// ReceiveDisasterPredictions provides each client with the prediction info, in addition to the source island,
// that they have been granted access to see
func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	for island, prediction := range receivedPredictions {
		updatedHist := append(c.predictionsHist[island], prediction.PredictionMade)
		c.predictionsHist[island] = updatedHist
	}
}

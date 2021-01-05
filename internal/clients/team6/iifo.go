package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {
	return c.BaseClient.MakeDisasterPrediction()
}

func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	c.BaseClient.ReceiveDisasterPredictions(receivedPredictions)
}

func (c *client) MakeForageInfo() shared.ForageShareInfo {
	return c.BaseClient.MakeForageInfo()
}

func (c *client) ReceiveForageInfo(forageInfo []shared.ForageShareInfo) {
	c.BaseClient.ReceiveForageInfo(forageInfo)
}

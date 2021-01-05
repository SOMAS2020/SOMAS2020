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

// ------ TODO: OPTIONAL ------
func (c *client) MakeForageInfo() shared.ForageShareInfo {
	return c.BaseClient.MakeForageInfo()
}

// ------ TODO: OPTIONAL ------
func (c *client) ReceiveForageInfo(forageInfo []shared.ForageShareInfo) {
	c.BaseClient.ReceiveForageInfo(forageInfo)
}

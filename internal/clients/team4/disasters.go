package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
)

// DisasterNotification is an event handler for disasters. Server will notify client
// of the ramifications of a disaster via this method.
// OPTIONAL: Use this method for any tasks you want to happen when a disaster occurs
func (c *client) DisasterNotification(
	dR disasters.DisasterReport,
	effects disasters.DisasterEffects) { // effects contain abs magnitude, prop. mag relative to other islands and CP mitigated mag.
	currTurn := c.ServerReadHandle.GetGameState().Turn
	disasterhappening := baseclient.DisasterInfo{
		CoordinateX: dR.X,
		CoordinateY: dR.Y,
		Magnitude:   dR.Magnitude,
		Turn:        currTurn,
	}

	//for team, prediction := range c.obs.iifoObs.receivedDisasterPredictions {
	// TODO: adjust agentsTrust based on their predictions
	//}

	c.obs.pastDisastersList = append(c.obs.pastDisastersList, disasterhappening)
}

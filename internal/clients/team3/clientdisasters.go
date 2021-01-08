package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
)

func (c *client) DisasterNotification(
	dR disasters.DisasterReport,
	effects disasters.DisasterEffects) {
	c.pastDisastersList = append(c.pastDisastersList, baseclient.DisasterInfo{
		CoordinateX: dR.X,
		CoordinateY: dR.Y,
		Magnitude:   dR.Magnitude,
		Turn:        c.ServerReadHandle.GetGameState().Turn,
	})
}

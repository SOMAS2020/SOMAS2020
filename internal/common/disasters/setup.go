package disasters

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

var envConf = config.GameConfig().DisasterConfig

// InitEnvironment initialises environment according to definitions
func InitEnvironment(islandIDs []shared.ClientID) *Environment {

	ag := ArchipelagoGeography{islands: map[shared.ClientID]Island{}, xBounds: envConf.XBounds, yBounds: envConf.YBounds}
	dp := disasterParameters{globalProb: envConf.GlobalProb, spatialPDF: envConf.SpatialPDFType, magnitudeLambda: envConf.MagnitudeLambda}

	for i, id := range islandIDs {
		island := Island{id, float64(i), float64(0)} // begin with equidistant points on x axis
		ag.islands[id] = island
	}
	//TODO: think about possible security concerns of returning a pointer
	return &Environment{geography: ag, disasterParams: dp, lastDisasterReport: DisasterReport{}} // returning a pointer so that other methods can modify returned Environment instance
}

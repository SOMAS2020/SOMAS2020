package disasters

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// InitEnvironment initialises environment according to definitions
func InitEnvironment(islandIDs []shared.ClientID) Environment {
	envConf := config.GameConfig().DisasterConfig

	ag := ArchipelagoGeography{islands: map[shared.ClientID]IslandLocationInfo{}, xMin: envConf.XMin, xMax: envConf.XMax, yMin: envConf.YMin, yMax: envConf.YMin}
	dp := disasterParameters{globalProb: envConf.GlobalProb, spatialPDF: envConf.SpatialPDFType, magnitudeLambda: envConf.MagnitudeLambda}

	for i, id := range islandIDs {
		island := IslandLocationInfo{id, float64(i), float64(0)} // begin with equidistant points on x axis
		ag.islands[id] = island
	}
	return Environment{Geography: ag, DisasterParams: dp, LastDisasterReport: DisasterReport{}} // returning a pointer so that other methods can modify returned Environment instance
}

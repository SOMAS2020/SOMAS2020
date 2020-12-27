package disasters

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// InitEnvironment initialises environment according to definitions
func InitEnvironment(islandIDs []shared.ClientID, disasterConfig config.DisasterConfig) Environment {
	ag := ArchipelagoGeography{Islands: map[shared.ClientID]IslandLocationInfo{}, XMin: disasterConfig.XMin, XMax: disasterConfig.XMax, YMin: disasterConfig.YMin, YMax: disasterConfig.YMin}
	dp := disasterParameters{globalProb: disasterConfig.GlobalProb, spatialPDF: disasterConfig.SpatialPDFType, magnitudeLambda: disasterConfig.MagnitudeLambda}

	for i, id := range islandIDs {
		island := IslandLocationInfo{id, float64(i), float64(0)} // begin with equidistant points on x axis
		ag.Islands[id] = island
	}
	return Environment{Geography: ag, DisasterParams: dp, LastDisasterReport: DisasterReport{}} // returning a pointer so that other methods can modify returned Environment instance
}

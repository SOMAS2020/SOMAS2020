package disasters

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

var envConf = config.GameConfig().DisasterConfig

// InitEnvironment initialises environment according to definitions
func InitEnvironment(islandIDs []shared.ClientID) (*Environment, error) {

	ag := ArchipelagoGeography{map[shared.ClientID]Island{}, envConf.XBounds, envConf.YBounds}
	dp := disasterParameters{globalProb: envConf.GlobalProb, spatialPDF: envConf.SpatialPDFType, magnitudeLambda: envConf.MagnitudeLambda}

	for i, id := range islandIDs {
		island := Island{id, float64(i), float64(0)} // begin with equidistant points on x axis
		ag.islands[id] = island
	}
	return &Environment{ag, dp, DisasterReport{}}, nil // may want ability to return error in future
}

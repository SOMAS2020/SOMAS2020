package disasters

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// InitEnvironment initialises environment according to definitions
func InitEnvironment(islandIDs []shared.ClientID, envConf config.DisasterConfig) Environment {
	ag := ArchipelagoGeography{
		Islands: map[shared.ClientID]IslandLocationInfo{},
		XMin:    envConf.XMin,
		XMax:    envConf.XMax,
		YMin:    envConf.YMin,
		YMax:    envConf.YMax,
	}
	dp := disasterParameters{
		globalProb:      envConf.GlobalProb,
		spatialPDF:      envConf.SpatialPDFType,
		magnitudeLambda: envConf.MagnitudeLambda,
	}

	xPoints := equidistantPoints(envConf.XMin, envConf.XMax, uint(len(islandIDs)))
	for i, id := range islandIDs {
		island := IslandLocationInfo{id, xPoints[i], float64(0)} // begin with equidistant points on x axis
		ag.Islands[id] = island
	}
	return Environment{Geography: ag, DisasterParams: dp, LastDisasterReport: DisasterReport{}}
}

// get n equally spaced points on a line connecting x0, x1
func equidistantPoints(x0, x1 float64, n uint) (points []float64) {
	delta := (x1 - x0) / math.Max(float64(n-1), 1) // prevent /0 error
	for i := uint(0); i < n; i++ {
		points = append(points, delta*float64(i))
	}
	return points
}

package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Rectangle struct {
	Bottom, Top, Left, Right shared.Coordinate
}

const (
	Min bool = false
	Max bool = true
)

type DisasterVulnerabilityParameters map[shared.ClientID]float64

func GetIslandDVPs() DisasterVulnerabilityParameters {
	archipelagoGeography := c.gamestate().Environment.Geography
	archipelagoCentre := 
	areaOfArchipelago := (archipelagoGeography.XMax - archipelagoGeography.XMin) *
		(archipelagoGeography.YMax - archipelagoGeography.YMin)

	for islandID, locationInfo := range archipelagoGeography.Islands {
		shiftedArchipelagoGeography := rectangle{
			Left:   locationInfo.X + archipelagoGeography.Xmin,
			Right:  locationInfo.X + archipelagoGeography.Xmax,
			Bottom: locationInfo.Y + archipelagoGeography.Ymin,
			Top:    locationInfo.Y + archipelagoGeography.Ymax,
		}
		overlapArchipelagoGeography = rectangle{
			Left: GetMinMax(Max, shiftedArchipelagoGeography.Xmin, archipelagoGeography.Xmin),
		}

	}
}

// GetMinMax returns either the minimum or maximum coordinate out of the two supplied, according to the bool argument
// that is input to the function
func GetMinMax(minOrMax bool, coordinate1 shared.Coordinate, coordinate2 shared.Coordinate) shared.Coordinate {
	if minOrMax == Min {
		if coordinate1 < coordinate2 {
			return coordinate1
		}
		return coordinate2
	}

	if coordinate1 > coordinate2 {
		return coordinate1
	}
	return coordinate2
}

// ReceiveDisasterPredictions provides each client with the prediction info, in addition to the source island,
// that they have been granted access to see
func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	for island, prediction := range receivedPredictions {
		updatedHist := append(c.predictionsHist[island], prediction.PredictionMade)
		c.predictionsHist[island] = updatedHist
	}
}

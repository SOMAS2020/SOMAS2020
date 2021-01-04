package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// DisasterVulnerabilityParametersDict is a map from island ID to an islands DVP
type DisasterVulnerabilityParametersDict map[shared.ClientID]float64

// CartesianCoordinates is a struct that holds the X,Y coordinates of a point
type CartesianCoordinates struct {
	X, Y shared.Coordinate
}

// Outline is a struct that holds the coordinates of the left, right, bottom and top sides
// of a rectangular outline
type Outline struct {
	Left, Right, Bottom, Top shared.Coordinate
}

// Define constant variables for choosing to find maximum or minimum in GetMinMax()
const (
	Min bool = false
	Max bool = true
)

// GetIslandDVPs is used to calculate the disaster vulnerability parameter of each island in the game.
// This only needs to be run at the start of the game because island's positions do not change
func GetIslandDVPs(c *client) DisasterVulnerabilityParametersDict {
	islandDVPs := DisasterVulnerabilityParametersDict{}
	archipelagoGeography := c.gamestate().Environment.Geography
	archipelagoCentre := CartesianCoordinates{
		X: archipelagoGeography.Xmin + (archipelagoGeography.Xmax-archipelagoGeography.Xmin)/2,
		Y: archipelagoGeography.Ymin + (archipelagoGeography.Ymax-archipelagoGeography.Ymin)/2,
	}
	areaOfArchipelago := (archipelagoGeography.XMax - archipelagoGeography.XMin) * (archipelagoGeography.YMax - archipelagoGeography.YMin)

	// For each island, find the overlap between the archipelago and the shifted outline which
	// is centred around the island's position
	for islandID, locationInfo := range archipelagoGeography.Islands {
		relativeOffset := CartesianCoordinates{
			X: locationInfo.X - archipelagoCentre.X,
			Y: locationInfo.Y - archipelagoCentre.Y,
		}
		shiftedArchipelagoOutline := Outline{
			Left:   archipelagoGeography.Xmin + relativeOffset.X,
			Right:  archipelagoGeography.Xmax + relativeOffset.X,
			Bottom: archipelagoGeography.Ymin + relativeOffset.Y,
			Top:    archipelagoGeography.Ymax + relativeOffset.Y,
		}
		overlapArchipelagoOutline := Outline{
			Left:   GetMinMax(Max, shiftedArchipelagoOutline.Left, archipelagoGeography.Xmin),
			Right:  GetMinMax(Min, shiftedArchipelagoOutline.Right, archipelagoGeography.Xmax),
			Bottom: GetMinMax(Max, shiftedArchipelagoOutline.Bottom, archipelagoGeography.Ymin),
			Top:    GetMinMax(Min, shiftedArchipelagoOutline.Top, archipelagoGeography.Ymax),
		}

		areaOfOverlap := (overlapArchipelagoOutline.Right - overlapArchipelagoOutline.Right) * (overlapArchipelagoOutline.Top - overlapArchipelagoOutline.Bottom)
		islandDVPs[islandID] = areaOfOverlap / areaOfArchipelago
	}
	return islandDVPs
}

// GetMinMax returns either the minimum or maximum coordinate of the two supplied, according to the bool argument
// that is input to the function
func GetMinMax(minOrMax bool, coordinate1 shared.Coordinate, coordinate2 shared.Coordinate) shared.Coordinate {
	if (minOrMax == Min && coordinate1 < coordinate2) || (minOrMax == Max && coordinate1 > coordinate2) {
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

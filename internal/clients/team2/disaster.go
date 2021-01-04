package team2

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
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

// MakeDisasterPrediction is used to provide our island's prediction on the next disaster
func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {
	//------------Tuning Parameters---------------//
	tuningParamK := 1.0
	varianceCapTimeRemaining := 10000.0
	tuningParamG := 1.0
	varianceCapMagnitude := 10000.0
	//--------------------------------------------//

	// Get the location prediction
	locationPrediction := GetLocationPrediction(c)

	// Get the time until next disaster prediction
	timeRemainingPrediction, confidenceTimeRemaining := GetTimeRemainingPrediction(c, tuningParamK, varianceCapTimeRemaining)

	// Get the magnitude prediction
	magnitudePrediction, confidenceMagnitude := GetMagnitudePrediction(c, tuningParamG, varianceCapMagnitude)

	// Get the overall confidence in these predictions
	confidencePrediction := GetConfidencePrediction(confidenceTimeRemaining, confidenceMagnitude)

	// Get trusted islands NOTE: CURRENTLY JUST ALL ISLANDS
	trustedislands := GetTrustedIslands()

	// Put everything together and return the whole prediction we have made and teams to share with
	disasterPrediction := shared.DisasterPrediction{
		CoordinateX: locationPrediction.X,
		CoordinateY: locationPrediction.Y,
		Magnitude:   magnitudePrediction,
		TimeLeft:    timeRemainingPrediction,
		Confidence:  confidencePrediction,
	}
	disasterPredictionInfo := shared.DisasterPredictionInfo{
		PredictionMade: disasterPrediction,
		TeamsOfferedTo: trustedislands,
	}
	return disasterPredictionInfo
}

// GetLocationPrediction provides a prediction about the location of the next disaster.
// The prediction is always the the centre of the archipelago
func GetLocationPrediction(c *client) CartesianCoordinates {
	archipelagoGeography := c.gamestate().Environment.Geography
	archipelagoCentre := CartesianCoordinates{
		X: archipelagoGeography.Xmin + (archipelagoGeography.Xmax-archipelagoGeography.Xmin)/2,
		Y: archipelagoGeography.Ymin + (archipelagoGeography.Ymax-archipelagoGeography.Ymin)/2,
	}
	return archipelagoCentre
}

// GetTimeRemainingPrediction provides a prediction about the time remaining until the next disaster.
// The prediction is 1/sample mean of the Bernoulli RV, minus the turns since the last disaster.
func GetTimeRemainingPrediction(c *client, tuningParamK float64, varianceCapTimeRemaining float64) (int, shared.PredictionConfidence) {
	totalTurns := float64(c.gameState().Turn)
	totalDisasters := float64(len(c.disasterHistory))
	sampleMean := math.Round(totalDisasters / totalTurns)
	timeBetweenDisasters := 1 / sampleMean

	// Get the time remaining prediction
	timeRemaining := timeBetweenDisasters - (totalTurns - c.disasterHistory[len(c.disasterHistory)-1].Turn)

	// Get the confidence in this prediction
	estimatedVariance := (1 - sampleMean) / math.Pow(sampleMean, 2)
	confidence := 100.0 - (100.0 * GetMinMaxFloat(Min, estimatedVariance/(tuningParamK*totalTurns), varianceCapTimeRemaining) / varianceCapTimeRemaining)
	return int(timeRemaining), confidence
}

// GetMagnitudePrediction provides a prediction about the magnitude of the next disaster.
// The prediction is the sample mean of the past magnitudes of disasters
func GetMagnitudePrediction(c *client, tuningParamG float64, varianceCapMagnitude float64) (shared.Magnitude, shared.PredictionConfidence) {
	totalTurns := float64(c.gameState().Turn)
	totalMagnitudes := 0.0
	for _, disasterReport := range c.disasterHistory {
		totalMagnitudes += disasterReport.Report.Magnitude
	}
	sampleMean := totalMagnitudes / totalTurns

	// Get the confidence in this prediction
	estimatedVariance := math.Pow(sampleMean, 2)
	confidence := 100.0 - (100.0 * GetMinMaxFloat(Min, estimatedVariance/(tuningParamG*totalTurns), varianceCapMagnitude) / varianceCapMagnitude)
	return sampleMean, confidence
}

// GetConfidencePrediction provides an overall confidence in our prediction.
// The confidence is the average of those from the timeRemaining and Magnitude predictions.
func GetConfidencePrediction(confidenceTimeRemaining shared.PredictionConfidence, confidenceMagnitude shared.PredictionConfidence) shared.PredictionConfidence {
	return (confidenceTimeRemaining + confidenceMagnitude) / 2
}

// GetMinMaxFloat is the same as GetMinMax but works for floats
func GetMinMaxFloat(minOrMax bool, value1 float64, value2 float64) float64 {
	if (minOrMax == Min && value1 < value2) || (minOrMax == Max && value1 > value2) {
		return value1
	}
	return value2
}

// GetTrustedIslands returns a slice of the islands we want to share our prediction with.
// NOTE: CURRENTLY THIS JUST RETURNS ALL ISLANDS.
func GetTrustedIslands() []shared.ClientID {
	trustedIslands := make([]shared.ClientID, len(baseclient.RegisteredClients))
	for index, id := range shared.TeamIDs {
		trustedIslands[index] = id
	}
	return trustedIslands
}

// ReceiveDisasterPredictions provides each client with the prediction info, in addition to the source island,
// that they have been granted access to see
func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	for island, prediction := range receivedPredictions {
		updatedHist := append(c.predictionsHist[island], prediction.PredictionMade)
		c.predictionsHist[island] = updatedHist
	}
}

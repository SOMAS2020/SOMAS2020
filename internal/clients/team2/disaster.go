package team2

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
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

// Define a global variable that holds the last prediction we shared
var LastPredictionMade shared.DisasterPredictionInfo
var CombinedPrediction shared.DisasterPrediction

// GetIslandDVPs is used to calculate the disaster vulnerability parameter of each island in the game.
// This only needs to be run at the start of the game because island's positions do not change
func GetIslandDVPs(archipelagoGeography disasters.ArchipelagoGeography) DisasterVulnerabilityParametersDict {
	islandDVPs := make(DisasterVulnerabilityParametersDict)
	archipelagoCentre := CartesianCoordinates{
		X: archipelagoGeography.XMin + (archipelagoGeography.XMax-archipelagoGeography.XMin)/2,
		Y: archipelagoGeography.YMin + (archipelagoGeography.YMax-archipelagoGeography.YMin)/2,
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
			Left:   archipelagoGeography.XMin + relativeOffset.X,
			Right:  archipelagoGeography.XMax + relativeOffset.X,
			Bottom: archipelagoGeography.YMin + relativeOffset.Y,
			Top:    archipelagoGeography.YMax + relativeOffset.Y,
		}
		overlapArchipelagoOutline := Outline{
			Left:   GetMinMaxCoordinate(Max, shiftedArchipelagoOutline.Left, archipelagoGeography.XMin),
			Right:  GetMinMaxCoordinate(Min, shiftedArchipelagoOutline.Right, archipelagoGeography.XMax),
			Bottom: GetMinMaxCoordinate(Max, shiftedArchipelagoOutline.Bottom, archipelagoGeography.YMin),
			Top:    GetMinMaxCoordinate(Min, shiftedArchipelagoOutline.Top, archipelagoGeography.YMax),
		}

		areaOfOverlap := (overlapArchipelagoOutline.Right - overlapArchipelagoOutline.Left) * (overlapArchipelagoOutline.Top - overlapArchipelagoOutline.Bottom)
		if areaOfArchipelago != 0 {
			islandDVPs[islandID] = areaOfOverlap / areaOfArchipelago
		} else {
			islandDVPs[islandID] = areaOfOverlap
		}
	}
	return islandDVPs
}

// GetMinMaxCoordinate returns either the minimum or maximum coordinate of the two supplied, according to the bool argument
// that is input to the function
func GetMinMaxCoordinate(minOrMax bool, coordinate1 shared.Coordinate, coordinate2 shared.Coordinate) shared.Coordinate {
	if (minOrMax == Min && coordinate1 < coordinate2) || (minOrMax == Max && coordinate1 > coordinate2) {
		return coordinate1
	}
	return coordinate2
}

// GetMinMaxFloat is the same as GetMinMaxCoordinate but works for floats
func GetMinMaxFloat(minOrMax bool, value1 float64, value2 float64) float64 {
	if (minOrMax == Min && value1 < value2) || (minOrMax == Max && value1 > value2) {
		return value1
	}
	return value2
}

// MakeDisasterPrediction is used to provide our island's prediction on the next disaster
func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {
	totalTurns := float64(c.gameState().Turn)

	// If no disasters have occurred then we cannot make a valid prediction
	if len(c.disasterHistory) == 0 {
		return nilPrediction()
	}

	// Get the location prediction
	locationPrediction := GetLocationPrediction(c)

	// Get the time until next disaster prediction and confidence
	sampleMeanX, timeRemainingPrediction := GetTimeRemainingPrediction(c, totalTurns)
	confidenceTimeRemaining := GetTimeRemainingConfidence(totalTurns, sampleMeanX)

	// Get the magnitude prediction and confidence
	sampleMeanM, magnitudePrediction := GetMagnitudePrediction(c, totalTurns)
	confidenceMagnitude := GetMagnitudeConfidence(totalTurns, sampleMeanM)

	// Get the overall confidence in these predictions
	confidencePrediction := GetConfidencePrediction(confidenceTimeRemaining, confidenceMagnitude)

	// Get trusted islands NOTE: CURRENTLY JUST ALL ISLANDS
	islandsToShareWith := GetIslandsToShareWith()

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
		TeamsOfferedTo: islandsToShareWith,
	}
	LastPredictionMade = disasterPredictionInfo
	return disasterPredictionInfo
}

// nilPrediction provides a nil prediction i.e. a prediction containing no information and
// is shared with no teams. We can tell it is a nil prediction because disasterHistory is empty
func nilPrediction() shared.DisasterPredictionInfo {
	nilPrediction := shared.DisasterPredictionInfo{
		PredictionMade: shared.DisasterPrediction{},
		TeamsOfferedTo: []shared.ClientID{},
	}
	return nilPrediction
}

// GetLocationPrediction provides a prediction about the location of the next disaster.
// The prediction is always the the centre of the archipelago
func GetLocationPrediction(c *client) CartesianCoordinates {
	archipelagoGeography := c.gameState().Geography
	archipelagoCentre := CartesianCoordinates{
		X: archipelagoGeography.XMin + (archipelagoGeography.XMax-archipelagoGeography.XMin)/2,
		Y: archipelagoGeography.YMin + (archipelagoGeography.YMax-archipelagoGeography.YMin)/2,
	}
	return archipelagoCentre
}

// GetTimeRemainingPrediction returns a prediction about the time remaining until the next disaster and the sample mean
// of the RV X. The prediction is 1/sample mean of the Bernoulli RV, minus the turns since the last disaster.
func GetTimeRemainingPrediction(c *client, totalTurns float64) (float64, uint) {
	totalDisasters := float64(len(c.disasterHistory))
	sampleMeanX := totalDisasters / totalTurns

	// Get the time remaining prediction
	expectationTd := math.Round(1 / sampleMeanX)
	timeRemaining := expectationTd - (totalTurns - float64(c.disasterHistory[len(c.disasterHistory)-1].Turn))
	if timeRemaining < 0 {
		timeRemaining = 0
	}

	return sampleMeanX, uint(timeRemaining)
}

// GetTimeRemainingConfidence returns the confidence in the time remaining prediction. The formula for this confidence is
// given in the report (can ask Hamish)
func GetTimeRemainingConfidence(totalTurns float64, sampleMeanX float64) shared.PredictionConfidence {
	varianceTd := (1 - sampleMeanX) / math.Pow(sampleMeanX, 2)
	confidence := 100.0 - (100.0 * GetMinMaxFloat(Min, varianceTd/(TuningParamK*totalTurns), VarianceCapTimeRemaining) / VarianceCapTimeRemaining)
	return confidence
}

// GetMagnitudePrediction returns a prediction about the magnitude of the next disaster and the sample mean
// of the RV M. The prediction is the sample mean of the past magnitudes of disasters
func GetMagnitudePrediction(c *client, totalTurns float64) (float64, shared.Magnitude) {
	totalMagnitudes := 0.0
	for _, disasterReport := range c.disasterHistory {
		totalMagnitudes += disasterReport.Report.Magnitude
	}
	sampleMeanM := totalMagnitudes / float64(len(c.disasterHistory))

	// Get the magnitude prediction
	magnitudePrediction := sampleMeanM
	return sampleMeanM, magnitudePrediction
}

// GetMagnitudeConfidence returns the confidence in the magnitude prediction. The formula for this confidence is
// given in the report (can ask Hamish)
func GetMagnitudeConfidence(totalTurns float64, sampleMeanM float64) shared.PredictionConfidence {
	varianceM := math.Pow(sampleMeanM, 2)
	confidence := 100.0 - (100.0 * GetMinMaxFloat(Min, varianceM/(TuningParamG*totalTurns), VarianceCapMagnitude) / VarianceCapMagnitude)
	return confidence
}

// GetConfidencePrediction provides an overall confidence in our prediction.
// The confidence is the average of those from the timeRemaining and Magnitude predictions.
func GetConfidencePrediction(confidenceTimeRemaining shared.PredictionConfidence, confidenceMagnitude shared.PredictionConfidence) shared.PredictionConfidence {
	return (confidenceTimeRemaining + confidenceMagnitude) / 2
}

// islandsToShareWith returns a slice of the islands we want to share our prediction with.
// We decided to always share our prediction with all islands to improve arhcipelago decisions as a whole.
func GetIslandsToShareWith() []shared.ClientID {
	islandsToShareWith := make([]shared.ClientID, len(shared.TeamIDs))
	for index, id := range shared.TeamIDs {
		islandsToShareWith[index] = id
	}
	return islandsToShareWith
}

// ReceiveDisasterPredictions provides each client with the prediction info, in addition to the source island,
// that they have been granted access to see.
// We use this function to combine all predictions into one final prediction (CombinedPrediction) to use for decisions.
func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	UpdatePredictionHistory(c, receivedPredictions)

	// Get the confidence in each island's prediction making ability
	islandConfidences := make(map[shared.ClientID]int)
	for island, _ := range c.opinionHist {
		conf := c.confidence("DisasterPred", island)
		islandConfidences[island] = conf
		c.opinionHist[island].Performances["DisasterPred"] = ExpectationReality{
			exp: conf,
		}
	}

	// Combine each islands prediction
	finalPrediction := CombinePredictions(c, receivedPredictions, islandConfidences)
	CombinedPrediction = finalPrediction
}

// UpdatePredictionHistory updates the history of predictions we have recieved from other islands with
// those recieved this turn.
func UpdatePredictionHistory(c *client, receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	if c.predictionHist == nil {
		c.predictionHist = make(PredictionsHist)
		for _, id := range shared.TeamIDs {
			c.predictionHist[id] = make([]PredictionInfo, 0)
		}
	}

	// Add the prediction to the history
	for _, prediction := range receivedPredictions {
		currPrediction := PredictionInfo{
			Prediction: prediction.PredictionMade,
			Turn:       c.gameState().Turn + uint(prediction.PredictionMade.TimeLeft),
		}
		c.predictionHist[prediction.SharedFrom] = append(c.predictionHist[prediction.SharedFrom], currPrediction)

	}
}

// CombinePredictions combines the predictions recieved from all the islands (including ours) to get
// one final disaster prediction.
// Use our confidence in an island as well as that island's confidence in their prediction to do this
func CombinePredictions(c *client, receivedPredictions shared.ReceivedDisasterPredictionsDict, islandConfidences map[shared.ClientID]int) shared.DisasterPrediction {
	// If confidence in island is zero OR islands confidence in their prediction is zero, for all islands,
	// then take the combined prediction to be our prediction instead
	numZeroTerms := 0
	for islandID, confidence := range islandConfidences {
		if confidence == 0 || receivedPredictions[islandID].PredictionMade.Confidence == 0 {
			numZeroTerms++
		}
	}
	if numZeroTerms == len(islandConfidences) {
		return LastPredictionMade.PredictionMade
	}

	// Add our own prediction to those recieved
	receivedPredictions[c.GetID()] = shared.ReceivedDisasterPredictionInfo{
		PredictionMade: LastPredictionMade.PredictionMade,
		SharedFrom:     c.GetID(),
	}

	// Get the sum of our confidences in other islands
	islandConfidencesSum := 0.0
	for _, confidence := range islandConfidences {
		islandConfidencesSum += float64(confidence)
	}

	// For each recieved prediction, we need the weighted sum (ws) of sub-predictions
	// Confidence must be treated slightly differently however
	wsCoordinateX, wsCoordinateY, wsMagnitude, wsTimeLeft, combinationConfidenceSum := 0.0, 0.0, 0.0, 0.0, 0.0
	for islandID, prediction := range receivedPredictions {

		// Get the combination confidence = (our confidence in island x their confidence in their prediction)/100
		combinationConfidence := (float64(islandConfidences[islandID]) * prediction.PredictionMade.Confidence) / 100

		// Get the weighted sum for each sub-prediction (except confidence)
		wsCoordinateX += combinationConfidence * prediction.PredictionMade.CoordinateX
		wsCoordinateY += combinationConfidence * prediction.PredictionMade.CoordinateY
		wsMagnitude += combinationConfidence * prediction.PredictionMade.Magnitude
		wsTimeLeft += combinationConfidence * float64(prediction.PredictionMade.TimeLeft)

		// Need sum of combination confidence also (sum of weights)
		combinationConfidenceSum += combinationConfidence
	}

	// Finally get the combined prediction by taking the weighted average of each sub-prediction
	finalPrediction := shared.DisasterPrediction{
		CoordinateX: wsCoordinateX / combinationConfidenceSum,
		CoordinateY: wsCoordinateY / combinationConfidenceSum,
		Magnitude:   wsMagnitude / combinationConfidenceSum,
		TimeLeft:    uint((wsTimeLeft / combinationConfidenceSum) + 0.5),
		Confidence:  combinationConfidenceSum / (islandConfidencesSum * 100),
	}

	return finalPrediction
}

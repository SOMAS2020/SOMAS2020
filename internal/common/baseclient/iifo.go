package baseclient

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Disaster defines the disaster location and magnitude.
// These disasters will be stored in PastDisastersDict (maps round number to disaster that occurred)
type DisasterInfo struct {
	CoordinateX float64
	CoordinateY float64
	Magnitude   int
	TurnNumber  int
}

// PastDisastersDict is a helpful construct for
type PastDisastersDict = map[int]DisasterInfo

var ourPredictionInfo shared.PredictionInfo

// MakePrediction is called on each client for them to make a prediction about a disaster
// Prediction includes location, magnitude, confidence etc
// COMPULSORY, you need to implement this method
func (c *BaseClient) MakePrediction() (shared.PredictionInfo, error) {
	// Set up dummy disasters for testing purposes
	pastDisastersDict := make(PastDisastersDict)
	tempDisaster := DisasterInfo{}
	for i := 0; i < 4; i++ {
		tempDisaster.CoordinateX = float64(i)
		tempDisaster.CoordinateY = float64(i)
		tempDisaster.Magnitude = i
		tempDisaster.TurnNumber = i
		pastDisastersDict[i] = tempDisaster
	}

	// Use the sample mean of each field as our prediction
	meanDisaster, err := getMeanDisaster(pastDisastersDict)
	if err != nil {
		return shared.PredictionInfo{}, err
	}

	prediction := shared.Prediction{
		CoordinateX: meanDisaster.CoordinateX,
		CoordinateY: meanDisaster.CoordinateY,
		Magnitude:   meanDisaster.Magnitude,
		TimeLeft:    meanDisaster.TurnNumber,
	}

	// Use (variance limit - mean(sample variance)), where the mean is taken over each field, as confidence
	// Use a variance limit of 100 for now
	varianceLimit := 100.0
	confidence, err := determineConfidence(pastDisastersDict, meanDisaster, varianceLimit)
	if err != nil {
		return shared.PredictionInfo{}, err
	}
	prediction.Confidence = confidence

	/*
		if c.GetID() == 1 {
			c.Logf("CoordinateX Prediction:[%v]\n", prediction.CoordinateX)
			c.Logf("CoordinateY Prediction:[%v]\n", prediction.CoordinateY)
			c.Logf("Magnitude Prediction:[%v]\n", prediction.Magnitude)
			c.Logf("Time Left Prediction:[%v]\n", prediction.TimeLeft)
			c.Logf("Confidence:[%v]\n", prediction.Confidence)
		}*/

	// For MVP, share this prediction with all islands since trust has not yet been implemented
	trustedIslands := make([]shared.ClientID, 6)
	for index, id := range shared.TeamIDs {
		trustedIslands[index] = id
	}

	// Return all prediction info and store our own island's prediction in global variable
	predictionInfo := shared.PredictionInfo{
		PredictionMade: prediction,
		TeamsOfferedTo: trustedIslands,
	}
	ourPredictionInfo = predictionInfo
	return predictionInfo, nil
}

func getMeanDisaster(pastDisastersDict PastDisastersDict) (DisasterInfo, error) {
	totalCoordinateX, totalCoordinateY, totalMagnitude, totalTurnNumber := 0.0, 0.0, 0.0, 0.0
	numberDisastersPassed := float64(len(pastDisastersDict))

	for _, disaster := range pastDisastersDict {
		totalCoordinateX += disaster.CoordinateX
		totalCoordinateY += disaster.CoordinateY
		totalMagnitude += float64(disaster.Magnitude)
		totalTurnNumber += float64(disaster.TurnNumber)
	}

	meanDisaster := DisasterInfo{
		CoordinateX: totalCoordinateX / numberDisastersPassed,
		CoordinateY: totalCoordinateY / numberDisastersPassed,
		Magnitude:   int(totalMagnitude/numberDisastersPassed + 0.5),
		TurnNumber:  int(totalTurnNumber/numberDisastersPassed + 0.5),
	}
	return meanDisaster, nil
}
func determineConfidence(pastDisastersDict PastDisastersDict, meanDisaster DisasterInfo, varianceLimit float64) (int, error) {
	totalCoordinateX, totalCoordinateY, totalMagnitude, totalTurnNumber := 0.0, 0.0, 0.0, 0.0
	numberDisastersPassed := float64(len(pastDisastersDict))

	// Find the sum of the square of the difference between the actual and mean, for each field
	for _, disaster := range pastDisastersDict {
		totalCoordinateX += math.Pow(disaster.CoordinateX-meanDisaster.CoordinateX, 2)
		totalCoordinateY += math.Pow(disaster.CoordinateY-meanDisaster.CoordinateY, 2)
		totalMagnitude += math.Pow(float64(disaster.Magnitude-meanDisaster.Magnitude), 2)
		totalTurnNumber += math.Pow(float64(disaster.TurnNumber-meanDisaster.TurnNumber), 2)
	}

	// Find the sum of the variances and the average variance
	varianceSum := (totalCoordinateX + totalCoordinateY + totalMagnitude + totalTurnNumber) / numberDisastersPassed
	averageVariance := varianceSum / 4

	// Implement the variance cap chosen
	if averageVariance > varianceLimit {
		averageVariance = varianceLimit
	}

	// Return the confidence of the prediction
	return int(varianceLimit - averageVariance + 0.5), nil
}

// RecievePredictions provides each client with the prediction info, in addition to the source island,
// that they have been granted access to see
// COMPULSORY, you need to implement this method
func (c *BaseClient) RecievePredictions(recievedPredictions shared.PredictionInfoDict) error {
	// If we assume that we trust each island equally (including ourselves), then take the final prediction
	// of disaster as being the weighted mean of predictions according to confidence
	numberOfPredictions := float64(len(recievedPredictions) + 1)
	selfConfidence := ourPredictionInfo.PredictionMade.Confidence

	// Initialise running totals using our own island's predictions
	totalCoordinateX := float64(selfConfidence) * ourPredictionInfo.PredictionMade.CoordinateX
	totalCoordinateY := float64(selfConfidence) * ourPredictionInfo.PredictionMade.CoordinateY
	totalMagnitude := selfConfidence * ourPredictionInfo.PredictionMade.Magnitude
	totalTimeLeft := selfConfidence * ourPredictionInfo.PredictionMade.TimeLeft
	totalConfidence := float64(selfConfidence)

	// Add other island's predictions using their confidence values
	for _, prediction := range recievedPredictions {
		totalCoordinateX += float64(prediction.PredictionMade.Confidence) * prediction.PredictionMade.CoordinateX
		totalCoordinateY += float64(prediction.PredictionMade.Confidence) * prediction.PredictionMade.CoordinateY
		totalMagnitude += prediction.PredictionMade.Confidence * prediction.PredictionMade.Magnitude
		totalTimeLeft += prediction.PredictionMade.Confidence * prediction.PredictionMade.TimeLeft
		totalConfidence += float64(prediction.PredictionMade.Confidence)
	}

	// Finally get the final prediction generated by considering predictions from all islands that we have available
	// This result is currently unused but would be used in decision making in full implementation
	finalPrediction := shared.Prediction{
		CoordinateX: totalCoordinateX / totalConfidence,
		CoordinateY: totalCoordinateY / totalConfidence,
		Magnitude:   int((float64(totalMagnitude) / totalConfidence) + 0.5),
		TimeLeft:    int((float64(totalTimeLeft) / totalConfidence) + 0.5),
		Confidence:  int((totalConfidence / numberOfPredictions) + 0.5),
	}

	c.Logf("Final Prediction: [%v]\n", finalPrediction)
	return nil
}

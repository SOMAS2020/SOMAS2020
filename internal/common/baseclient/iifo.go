package baseclient

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Disaster defines the disaster location and magnitude.
// These disasters will be stored in PastDisastersDict (maps round number to disaster that occurred)
// TODO: Agree with environment team on disaster struct representation
type DisasterInfo struct {
	CoordinateX float64
	CoordinateY float64
	Magnitude   float64
	Turn        uint
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
	for i := 0; i < 4; i++ {
		pastDisastersDict[i] = DisasterInfo{
			CoordinateX: float64(i),
			CoordinateY: float64(i),
			Magnitude:   float64(i),
			Turn:        uint(i),
		}
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
		TimeLeft:    int(meanDisaster.Turn),
	}

	// Use (variance limit - mean(sample variance)), where the mean is taken over each field, as confidence
	// Use a variance limit of 100 for now
	varianceLimit := 100.0
	prediction.Confidence, err = determineConfidence(pastDisastersDict, meanDisaster, varianceLimit)
	if err != nil {
		return shared.PredictionInfo{}, err
	}

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
	totalCoordinateX, totalCoordinateY, totalMagnitude, totalTurn := 0.0, 0.0, 0.0, 0.0
	numberDisastersPassed := float64(len(pastDisastersDict))

	for _, disaster := range pastDisastersDict {
		totalCoordinateX += disaster.CoordinateX
		totalCoordinateY += disaster.CoordinateY
		totalMagnitude += float64(disaster.Magnitude)
		totalTurn += float64(disaster.Turn)
	}

	meanDisaster := DisasterInfo{
		CoordinateX: totalCoordinateX / numberDisastersPassed,
		CoordinateY: totalCoordinateY / numberDisastersPassed,
		Magnitude:   totalMagnitude / numberDisastersPassed,
		Turn:        uint(math.Round(totalTurn / numberDisastersPassed)),
	}
	return meanDisaster, nil
}

func determineConfidence(pastDisastersDict PastDisastersDict, meanDisaster DisasterInfo, varianceLimit float64) (float64, error) {
	totalCoordinateX, totalCoordinateY, totalMagnitude, totalTurn := 0.0, 0.0, 0.0, 0.0
	totalDisaster := DisasterInfo{}
	numberDisastersPassed := float64(len(pastDisastersDict))

	// Find the sum of the square of the difference between the actual and mean, for each field
	for _, disaster := range pastDisastersDict {
		totalDisaster.CoordinateX += math.Pow(disaster.CoordinateX-meanDisaster.CoordinateX, 2)
		totalDisaster.CoordinateY += math.Pow(disaster.CoordinateY-meanDisaster.CoordinateY, 2)
		totalDisaster.Magnitude += math.Pow(disaster.Magnitude-meanDisaster.Magnitude, 2)
		totalDisaster.Turn += uint(math.Round(math.Pow(float64(disaster.Turn-meanDisaster.Turn), 2)))
	}

	// Find the sum of the variances and the average variance
	varianceSum := (totalCoordinateX + totalCoordinateY + totalMagnitude + totalTurn) / numberDisastersPassed
	averageVariance := varianceSum / 4

	// Implement the variance cap chosen
	if averageVariance > varianceLimit {
		averageVariance = varianceLimit
	}

	// Return the confidence of the prediction
	return math.Round(varianceLimit - averageVariance), nil
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
	totalCoordinateX := selfConfidence * ourPredictionInfo.PredictionMade.CoordinateX
	totalCoordinateY := selfConfidence * ourPredictionInfo.PredictionMade.CoordinateY
	totalMagnitude := selfConfidence * ourPredictionInfo.PredictionMade.Magnitude
	totalTimeLeft := int(math.Round(selfConfidence)) * ourPredictionInfo.PredictionMade.TimeLeft
	totalConfidence := selfConfidence

	// Add other island's predictions using their confidence values
	for _, prediction := range recievedPredictions {
		totalCoordinateX += prediction.PredictionMade.Confidence * prediction.PredictionMade.CoordinateX
		totalCoordinateY += prediction.PredictionMade.Confidence * prediction.PredictionMade.CoordinateY
		totalMagnitude += prediction.PredictionMade.Confidence * prediction.PredictionMade.Magnitude
		totalTimeLeft += int(math.Round(prediction.PredictionMade.Confidence)) * prediction.PredictionMade.TimeLeft
		totalConfidence += prediction.PredictionMade.Confidence
	}

	// Finally get the final prediction generated by considering predictions from all islands that we have available
	// This result is currently unused but would be used in decision making in full implementation
	finalPrediction := shared.Prediction{
		CoordinateX: totalCoordinateX / totalConfidence,
		CoordinateY: totalCoordinateY / totalConfidence,
		Magnitude:   totalMagnitude / totalConfidence,
		TimeLeft:    int((float64(totalTimeLeft) / totalConfidence) + 0.5),
		Confidence:  totalConfidence / numberOfPredictions,
	}

	c.Logf("Final Prediction: [%v]\n", finalPrediction)
	return nil
}

package team5

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// PastDisastersList is a List of previous disasters.
type PastDisastersList = []disasterInfo

type disasterInfo struct {
	epiX shared.Coordinate // x co-ord of disaster epicentre
	epiY shared.Coordinate // y ""
	mag  shared.Magnitude
	turn uint
}

// MakeDisasterPrediction is called on each client for them to make a prediction about a disaster
// Prediction includes location, magnitude, confidence etc
// COMPULSORY, you need to implement this method
func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {
	// Set up dummy disasters for testing purposes
	pastDisastersList := make(PastDisastersList, 5)
	for i := 0; i < 4; i++ {
		pastDisastersList[i] = (disasterInfo{
			epiX: float64(i),
			epiY: float64(i),
			mag:  float64(i),
			turn: uint(i),
		})
	}

	// Use the sample mean of each field as our prediction
	meanDisaster := getMeanDisaster(pastDisastersList)

	prediction := shared.DisasterPrediction{
		CoordinateX: meanDisaster.epiX,
		CoordinateY: meanDisaster.epiY,
		Magnitude:   meanDisaster.mag,
		TimeLeft:    int(meanDisaster.turn),
	}

	// Use (variance limit - mean(sample variance)), where the mean is taken over each field, as confidence
	// Use a variance limit of 100 for now
	varianceLimit := 100.0
	prediction.Confidence = determineConfidence(pastDisastersList, meanDisaster, varianceLimit)

	// For MVP, share this prediction with all islands since trust has not yet been implemented
	trustedIslands := make([]shared.ClientID, len(RegisteredClients))
	for index, id := range shared.TeamIDs {
		trustedIslands[index] = id
	}

	// Return all prediction info and store our own island's prediction in global variable
	predictionInfo := shared.DisasterPredictionInfo{
		PredictionMade: prediction,
		TeamsOfferedTo: trustedIslands,
	}
	c.predictionInfo = predictionInfo
	return predictionInfo
}

func getMeanDisaster(pastDisastersList PastDisastersList) disasterInfo {
	totalCoordinateX, totalCoordinateY, totalMagnitude, totalTurn := 0.0, 0.0, 0.0, 0.0
	numberDisastersPassed := float64(len(pastDisastersList))

	for _, disaster := range pastDisastersList {
		totalCoordinateX += disaster.epiX
		totalCoordinateY += disaster.epiY
		totalMagnitude += float64(disaster.mag)
		totalTurn += float64(disaster.turn)
	}

	meanDisaster := disasterInfo{
		epiX: totalCoordinateX / numberDisastersPassed,
		epiY: totalCoordinateY / numberDisastersPassed,
		mag:  totalMagnitude / numberDisastersPassed,
		turn: uint(math.Round(totalTurn / numberDisastersPassed)),
	}
	return meanDisaster
}

func determineConfidence(pastDisastersList PastDisastersList, meanDisaster disasterInfo, varianceLimit float64) float64 {
	totalDisaster := disasterInfo{}
	numberDisastersPassed := float64(len(pastDisastersList))

	// Find the sum of the square of the difference between the actual and mean, for each field
	for _, disaster := range pastDisastersList {
		totalDisaster.epiX += math.Pow(disaster.epiX-meanDisaster.epiX, 2)
		totalDisaster.epiY += math.Pow(disaster.epiY-meanDisaster.epiY, 2)
		totalDisaster.mag += math.Pow(disaster.mag-meanDisaster.mag, 2)
		totalDisaster.turn += uint(math.Round(math.Pow(float64(disaster.turn-meanDisaster.turn), 2)))
	}

	// Find the sum of the variances and the average variance
	varianceSum := (totalDisaster.epiX + totalDisaster.epiY + totalDisaster.mag + float64(totalDisaster.turn)) / numberDisastersPassed
	averageVariance := varianceSum / 4

	// Implement the variance cap chosen
	if averageVariance > varianceLimit {
		averageVariance = varianceLimit
	}

	// Return the confidence of the prediction
	return math.Round(varianceLimit - averageVariance)
}

// ReceiveDisasterPredictions provides each client with the prediction info, in addition to the source island,
// that they have been granted access to see
// COMPULSORY, you need to implement this method
func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	// If we assume that we trust each island equally (including ourselves), then take the final prediction
	// of disaster as being the weighted mean of predictions according to confidence
	numberOfPredictions := float64(len(receivedPredictions) + 1)
	selfConfidence := c.predictionInfo.PredictionMade.Confidence

	// Initialise running totals using our own island's predictions
	totalCoordinateX := selfConfidence * c.predictionInfo.PredictionMade.CoordinateX
	totalCoordinateY := selfConfidence * c.predictionInfo.PredictionMade.CoordinateY
	totalMagnitude := selfConfidence * c.predictionInfo.PredictionMade.Magnitude
	totalTimeLeft := int(math.Round(selfConfidence)) * c.predictionInfo.PredictionMade.TimeLeft
	totalConfidence := selfConfidence

	// Add other island's predictions using their confidence values
	for _, prediction := range receivedPredictions {
		totalCoordinateX += prediction.PredictionMade.Confidence * prediction.PredictionMade.CoordinateX
		totalCoordinateY += prediction.PredictionMade.Confidence * prediction.PredictionMade.CoordinateY
		totalMagnitude += prediction.PredictionMade.Confidence * prediction.PredictionMade.Magnitude
		totalTimeLeft += int(math.Round(prediction.PredictionMade.Confidence)) * prediction.PredictionMade.TimeLeft
		totalConfidence += prediction.PredictionMade.Confidence
	}

	// Finally get the final prediction generated by considering predictions from all islands that we have available
	// This result is currently unused but would be used in decision making in full implementation
	finalPrediction := shared.DisasterPrediction{
		CoordinateX: totalCoordinateX / totalConfidence,
		CoordinateY: totalCoordinateY / totalConfidence,
		Magnitude:   totalMagnitude / totalConfidence,
		TimeLeft:    int((float64(totalTimeLeft) / totalConfidence) + 0.5),
		Confidence:  totalConfidence / numberOfPredictions,
	}

	c.Logf("Final Prediction: [%v]", finalPrediction)
}

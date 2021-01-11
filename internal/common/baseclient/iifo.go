package baseclient

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// DisasterInfo defines the disaster location and magnitude.
// These disasters will be stored in PastDisastersDict (maps round number to disaster that occurred)
// TODO: Agree with environment team on disaster struct representation
type DisasterInfo struct {
	CoordinateX shared.Coordinate
	CoordinateY shared.Coordinate
	Magnitude   shared.Magnitude
	Turn        uint
}

// PastDisastersList is a List of previous disasters.
type PastDisastersList = []DisasterInfo

// MakeDisasterPrediction is called on each client for them to make a prediction about a disaster
// Prediction includes location, magnitude, confidence etc
// COMPULSORY, you need to implement this method
func (c *BaseClient) MakeDisasterPrediction() shared.DisasterPredictionInfo {
	// Set up dummy disasters for testing purposes
	pastDisastersList := make(PastDisastersList, 5)
	for i := 0; i < 4; i++ {
		pastDisastersList[i] = (DisasterInfo{
			CoordinateX: float64(i),
			CoordinateY: float64(i),
			Magnitude:   float64(i),
			Turn:        uint(i),
		})
	}

	// Use the sample mean of each field as our prediction
	meanDisaster := getMeanDisaster(pastDisastersList)

	prediction := shared.DisasterPrediction{
		CoordinateX: meanDisaster.CoordinateX,
		CoordinateY: meanDisaster.CoordinateY,
		Magnitude:   meanDisaster.Magnitude,
		TimeLeft:    meanDisaster.Turn,
	}

	// Use (variance limit - mean(sample variance)), where the mean is taken over each field, as confidence
	// Use a variance limit of 100 for now
	varianceLimit := 100.0
	prediction.Confidence = determineConfidence(pastDisastersList, meanDisaster, varianceLimit)

	// For MVP, share this prediction with all islands since trust has not yet been implemented
	trustedIslands := make([]shared.ClientID, len(shared.TeamIDs))
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

func getMeanDisaster(pastDisastersList PastDisastersList) DisasterInfo {
	totalCoordinateX, totalCoordinateY, totalMagnitude, totalTurn := 0.0, 0.0, 0.0, 0.0
	numberDisastersPassed := float64(len(pastDisastersList))
	if numberDisastersPassed == 0 {
		return DisasterInfo{0, 0, 0, 1000}
	}
	for _, disaster := range pastDisastersList {
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
	return meanDisaster
}

func determineConfidence(pastDisastersList PastDisastersList, meanDisaster DisasterInfo, varianceLimit float64) float64 {
	totalDisaster := DisasterInfo{}
	numberDisastersPassed := float64(len(pastDisastersList))

	// Find the sum of the square of the difference between the actual and mean, for each field
	for _, disaster := range pastDisastersList {
		totalDisaster.CoordinateX += math.Pow(disaster.CoordinateX-meanDisaster.CoordinateX, 2)
		totalDisaster.CoordinateY += math.Pow(disaster.CoordinateY-meanDisaster.CoordinateY, 2)
		totalDisaster.Magnitude += math.Pow(disaster.Magnitude-meanDisaster.Magnitude, 2)
		totalDisaster.Turn += uint(math.Round(math.Pow(float64(disaster.Turn-meanDisaster.Turn), 2)))
	}

	// Find the sum of the variances and the average variance
	varianceSum := (totalDisaster.CoordinateX + totalDisaster.CoordinateY + totalDisaster.Magnitude + float64(totalDisaster.Turn)) / numberDisastersPassed
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
func (c *BaseClient) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	// If we assume that we trust each island equally (including ourselves), then take the final prediction
	// of disaster as being the weighted mean of predictions according to confidence
	numberOfPredictions := float64(len(receivedPredictions) + 1)
	selfConfidence := c.predictionInfo.PredictionMade.Confidence

	// Initialise running totals using our own island's predictions
	totalCoordinateX := selfConfidence * c.predictionInfo.PredictionMade.CoordinateX
	totalCoordinateY := selfConfidence * c.predictionInfo.PredictionMade.CoordinateY
	totalMagnitude := selfConfidence * c.predictionInfo.PredictionMade.Magnitude
	totalTimeLeft := uint(math.Round(selfConfidence)) * c.predictionInfo.PredictionMade.TimeLeft
	totalConfidence := selfConfidence

	// Add other island's predictions using their confidence values
	for _, prediction := range receivedPredictions {
		totalCoordinateX += prediction.PredictionMade.Confidence * prediction.PredictionMade.CoordinateX
		totalCoordinateY += prediction.PredictionMade.Confidence * prediction.PredictionMade.CoordinateY
		totalMagnitude += prediction.PredictionMade.Confidence * prediction.PredictionMade.Magnitude
		totalTimeLeft += uint(math.Round(prediction.PredictionMade.Confidence)) * prediction.PredictionMade.TimeLeft
		totalConfidence += prediction.PredictionMade.Confidence
	}

	// Finally get the final prediction generated by considering predictions from all islands that we have available
	// This result is currently unused but would be used in decision making in full implementation
	finalPrediction := shared.DisasterPrediction{
		CoordinateX: totalCoordinateX / totalConfidence,
		CoordinateY: totalCoordinateY / totalConfidence,
		Magnitude:   totalMagnitude / totalConfidence,
		TimeLeft:    uint((float64(totalTimeLeft) / totalConfidence) + 0.5),
		Confidence:  totalConfidence / numberOfPredictions,
	}

	c.Logf("Final Prediction: [%v]", finalPrediction)
}

// MakeForageInfo allows clients to share their most recent foraging DecisionMade, ResourceObtained from it to
// other clients.
// OPTIONAL. If this is not implemented then all values are nil.
func (c *BaseClient) MakeForageInfo() shared.ForageShareInfo {
	contribution := shared.ForageDecision{Type: shared.DeerForageType, Contribution: 0}
	return shared.ForageShareInfo{DecisionMade: contribution, ResourceObtained: 0, ShareTo: []shared.ClientID{}}
}

// ReceiveForageInfo lets clients know what other clients has obtained from their most recent foraging attempt.
// Most recent foraging attempt includes information about: foraging DecisionMade and ResourceObtained as well
// as where this information came from.
// OPTIONAL.
func (c *BaseClient) ReceiveForageInfo(neighbourForaging []shared.ForageShareInfo) {
	// Return on Investment
	roi := map[shared.ClientID]shared.Resources{}
	for _, val := range neighbourForaging {
		if val.DecisionMade.Type == shared.DeerForageType {
			roi[val.SharedFrom] = val.ResourceObtained / shared.Resources(val.DecisionMade.Contribution) * 100
		}
	}
}

package team4

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// MakeDisasterPrediction is called on each client for them to make a prediction about a disaster
// Prediction includes location, magnitude, confidence etc
// COMPULSORY, you need to implement this method
func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {
	// Use the sample mean of each field as our prediction
	meanDisaster := getMeanDisaster(c.obs.pastDisastersList)

	prediction := shared.DisasterPrediction{
		CoordinateX: meanDisaster.CoordinateX,
		CoordinateY: meanDisaster.CoordinateY,
		Magnitude:   meanDisaster.Magnitude,
		TimeLeft:    meanDisaster.Turn,
	}

	// Use (variance limit - mean(sample variance)), where the mean is taken over each field, as confidence
	// Use a variance limit of 100 for now //TODO: tune this
	varianceLimit := 100.0
	prediction.Confidence = determineConfidence(c.obs.pastDisastersList, meanDisaster, varianceLimit)

	// For MVP, share this prediction with all islands since trust has not yet been implemented
	islandsToSend := make([]shared.ClientID, len(c.ServerReadHandle.GetGameState().ClientLifeStatuses))
	for index, id := range shared.TeamIDs {
		islandsToSend[index] = id
	}

	// Return all prediction info and store our own island's prediction in global variable
	predictionInfo := shared.DisasterPredictionInfo{
		PredictionMade: prediction,
		TeamsOfferedTo: islandsToSend,
	}
	c.obs.iifoObs.ourDisasterPrediction = predictionInfo
	return predictionInfo
}

func getMeanDisaster(pastDisastersList baseclient.PastDisastersList) baseclient.DisasterInfo {
	totalCoordinateX, totalCoordinateY, totalMagnitude, totalTurn := 0.0, 0.0, 0.0, 0.0
	numberDisastersPassed := float64(len(pastDisastersList))
	if numberDisastersPassed == 0 {
		return baseclient.DisasterInfo{CoordinateX: 0, CoordinateY: 0, Magnitude: 0, Turn: 1000}
	}
	for _, disaster := range pastDisastersList {
		totalCoordinateX += disaster.CoordinateX
		totalCoordinateY += disaster.CoordinateY
		totalMagnitude += float64(disaster.Magnitude)
		totalTurn += float64(disaster.Turn)
	}

	meanDisaster := baseclient.DisasterInfo{
		CoordinateX: totalCoordinateX / numberDisastersPassed,
		CoordinateY: totalCoordinateY / numberDisastersPassed,
		Magnitude:   totalMagnitude / numberDisastersPassed,
		Turn:        uint(math.Floor(totalTurn/numberDisastersPassed)) - uint(totalTurn)%uint(numberDisastersPassed), // gives the number of turns left until the next disaster
	}
	return meanDisaster
}

func determineConfidence(pastDisastersList baseclient.PastDisastersList, meanDisaster baseclient.DisasterInfo, varianceLimit float64) float64 {
	totalDisaster := baseclient.DisasterInfo{}
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
func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	// If we assume that we trust each island equally (including ourselves), then take the final prediction
	// of disaster as being the weighted mean of predictions according to confidence
	numberOfPredictions := float64(len(receivedPredictions) + 1)
	predictionInfo := c.obs.iifoObs.ourDisasterPrediction.PredictionMade
	selfConfidence := predictionInfo.Confidence

	// Initialise running totals using our own island's predictions
	totalCoordinateX := selfConfidence * predictionInfo.CoordinateX
	totalCoordinateY := selfConfidence * predictionInfo.CoordinateY
	totalMagnitude := selfConfidence * predictionInfo.Magnitude
	totalTimeLeft := uint(math.Round(selfConfidence)) * predictionInfo.TimeLeft
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
	if totalConfidence == 0 {
		totalConfidence = numberOfPredictions
	}
	finalPrediction := shared.DisasterPrediction{
		CoordinateX: totalCoordinateX / totalConfidence,
		CoordinateY: totalCoordinateY / totalConfidence,
		Magnitude:   totalMagnitude / totalConfidence,
		TimeLeft:    uint((float64(totalTimeLeft) / totalConfidence) + 0.5),
		Confidence:  totalConfidence / numberOfPredictions,
	}
	c.obs.iifoObs.finalDisasterPrediction = finalPrediction
}

// MakeForageInfo allows clients to share their most recent foraging DecisionMade, ResourceObtained from it to
// other clients.
// OPTIONAL. If this is not implemented then all values are nil.
func (c *client) MakeForageInfo() shared.ForageShareInfo {

	// Who to share to
	var shareTo []shared.ClientID
	for id, status := range c.getAllLifeStatus() {
		if status != shared.Dead {
			if c.getTurn() < 5 {
				// Send to everyone for the first five rounds
				shareTo = append(shareTo, id)
			} else {
				// Send only to island who sent to us in the previous round
				shareTo = c.returnPreviousForagers()
				// Maybe also add island we trust incase they won't send to us unless we send to them?
				shareTo = append(shareTo, c.trustMatrix.trustedClients(0.70)...)
				shareTo = createClientSet(shareTo)
			}
		}
	}

	// Greediness and selfishness to lie?
	var resourceObtained shared.Resources = 0
	var decisionMade shared.ForageDecision = shared.ForageDecision{}
	if len(c.forage.forageHistory) > 0 {
		lastRound := c.forage.forageHistory[len(c.forage.forageHistory)-1]
		decisionMade = lastRound.decision
		resourceObtained = lastRound.resourceReturn
	}

	forageInfo := shared.ForageShareInfo{
		ShareTo:          shareTo,
		ResourceObtained: resourceObtained,
		DecisionMade:     decisionMade,
		SharedFrom:       c.GetID(),
	}
	return forageInfo
}

//ReceiveForageInfo lets clients know what other clients has obtained from their most recent foraging attempt.
//Most recent foraging attempt includes information about: foraging DecisionMade and ResourceObtained as well
//as where this information came from.

func (c *client) ReceiveForageInfo(neighbourForaging []shared.ForageShareInfo) {
	c.forage.receivedForageData = append(c.forage.receivedForageData, neighbourForaging)
	//Give trust to island that contribute to this?
}

func (c *client) returnPreviousForagers() []shared.ClientID {
	data := c.forage.receivedForageData
	if len(data) < 1 {
		return nil
	}
	lastEntry := data[len(data)-1]

	var shareTo []shared.ClientID
	for _, teamReturns := range lastEntry {

		shareTo = append(shareTo, teamReturns.SharedFrom)
	}
	return shareTo
}

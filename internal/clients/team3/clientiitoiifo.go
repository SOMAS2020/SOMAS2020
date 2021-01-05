package team3

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"

	// 	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	// 	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

/*
	//IIFO: OPTIONAL
	MakeDisasterPrediction() shared.DisasterPredictionInfo
	ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict)
	MakeForageInfo() shared.ForageShareInfo
	ReceiveForageInfo([]shared.ForageShareInfo)

	//IITO: COMPULSORY
	GetGiftRequests() shared.GiftRequestDict
	GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict
	GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict
	UpdateGiftInfo(receivedResponses shared.GiftResponseDict)

	//TODO: THESE ARE NOT DONE yet, how do people think we should implement the actual transfer?
	// The server should handle the below functions maybe?
	SentGift(sent shared.Resources, to shared.ClientID)
	ReceivedGift(received shared.Resources, from shared.ClientID)

*/

func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {
	// Use the sample mean of each field as our prediction
	meanDisaster := getMeanDisaster(c.pastDisastersList)

	prediction := shared.DisasterPrediction{
		CoordinateX: meanDisaster.CoordinateX,
		CoordinateY: meanDisaster.CoordinateY,
		Magnitude:   meanDisaster.Magnitude,
		TimeLeft:    int(meanDisaster.Turn),
	}

	// Use (variance limit - mean(sample variance)), where the mean is taken over each field, as confidence
	// Use a variance limit of 100 for now
	varianceLimit := 100.0
	prediction.Confidence = determineConfidence(c.pastDisastersList, meanDisaster, varianceLimit)

	// For MVP, share this prediction with all islands since trust has not yet been implemented
	trustedIslands := make([]shared.ClientID, len(baseclient.RegisteredClients))
	for index, id := range shared.TeamIDs {
		trustedIslands[index] = id
	}

	// Return all prediction info and store our own island's prediction in global variable
	predictionInfo := shared.DisasterPredictionInfo{
		PredictionMade: prediction,
		TeamsOfferedTo: trustedIslands,
	}

	if len(c.disasterPredictions) < int(c.ServerReadHandle.GetGameState().Turn) {
		c.disasterPredictions = append(c.disasterPredictions, make(map[shared.ClientID]shared.DisasterPrediction))
	}

	c.disasterPredictions[int(c.ServerReadHandle.GetGameState().Turn)-1][c.GetID()] = predictionInfo.PredictionMade
	return predictionInfo
}

func getMeanDisaster(pastDisastersList baseclient.PastDisastersList) baseclient.DisasterInfo {
	totalCoordinateX, totalCoordinateY, totalMagnitude, totalTurn := 0.0, 0.0, 0.0, 0.0
	numberDisastersPassed := float64(len(pastDisastersList))

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
		Turn:        uint(math.Round(totalTurn / numberDisastersPassed)),
	}
	return meanDisaster
}

func determineConfidence(pastDisastersList baseclient.PastDisastersList, meanDisaster baseclient.DisasterInfo, varianceLimit float64) float64 {
	totalCoordinateX, totalCoordinateY, totalMagnitude, totalTurn := 0.0, 0.0, 0.0, 0.0
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
	varianceSum := (totalCoordinateX + totalCoordinateY + totalMagnitude + totalTurn) / numberDisastersPassed
	averageVariance := varianceSum / 4

	// Implement the variance cap chosen
	if averageVariance > varianceLimit {
		averageVariance = varianceLimit
	}

	// Return the confidence of the prediction
	return math.Round(varianceLimit - averageVariance)
}

func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	// Take the final prediction of disaster as being the weighted mean of predictions according to confidence times the opiinion we have of other islands prediction ablity
	numberOfPredictions := float64(len(receivedPredictions) + 1)
	selfConfidence := c.disasterPredictions[int(c.ServerReadHandle.GetGameState().Turn)-1][c.GetID()].Confidence

	// Initialise running totals using our own island's predictions
	totalCoordinateX := c.trustScore[c.GetID()] * selfConfidence * c.disasterPredictions[int(c.ServerReadHandle.GetGameState().Turn)-1][c.GetID()].CoordinateX
	totalCoordinateY := c.trustScore[c.GetID()] * selfConfidence * c.disasterPredictions[int(c.ServerReadHandle.GetGameState().Turn)-1][c.GetID()].CoordinateY
	totalMagnitude := c.trustScore[c.GetID()] * selfConfidence * c.disasterPredictions[int(c.ServerReadHandle.GetGameState().Turn)-1][c.GetID()].Magnitude
	totalTimeLeft := int(math.Round(c.trustScore[c.GetID()]*selfConfidence)) * c.disasterPredictions[int(c.ServerReadHandle.GetGameState().Turn)-1][c.GetID()].TimeLeft
	totalConfidence := c.trustScore[c.GetID()] * selfConfidence

	// Add other island's predictions using their confidence values
	for islandID, prediction := range receivedPredictions {
		totalCoordinateX += c.trustScore[islandID] * prediction.PredictionMade.Confidence * prediction.PredictionMade.CoordinateX
		totalCoordinateY += c.trustScore[islandID] * prediction.PredictionMade.Confidence * prediction.PredictionMade.CoordinateY
		totalMagnitude += c.trustScore[islandID] * prediction.PredictionMade.Confidence * prediction.PredictionMade.Magnitude
		totalTimeLeft += int(math.Round(c.trustScore[islandID]*prediction.PredictionMade.Confidence)) * prediction.PredictionMade.TimeLeft
		totalConfidence += c.trustScore[islandID] * prediction.PredictionMade.Confidence
	}

	// Finally get the final prediction generated by considering predictions from all islands that we have available
	// This result is currently unused but would be used in decision making in full implementation

	c.globalDisasterPredictions = append(c.globalDisasterPredictions, shared.DisasterPrediction{
		CoordinateX: totalCoordinateX / totalConfidence,
		CoordinateY: totalCoordinateY / totalConfidence,
		Magnitude:   totalMagnitude / totalConfidence,
		TimeLeft:    int((float64(totalTimeLeft) / totalConfidence) + 0.5),
		Confidence:  totalConfidence / numberOfPredictions,
	})

	c.Logf("Final Prediction: [%v]", c.globalDisasterPredictions[int(c.ServerReadHandle.GetGameState().Turn)-1])

	// TODO: compare other islands predictions to disaster when info is received and update their trust score
}

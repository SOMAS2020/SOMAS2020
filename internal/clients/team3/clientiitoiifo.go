package team3

import (
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
	meanDisaster, err := getMeanDisaster()
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
	prediction.Confidence, err = determineConfidence(pastDisasters, meanDisaster, varianceLimit)
	if err != nil {
		return shared.PredictionInfo{}, err
	}

	trustedIslands := make([]shared.ClientID, len(RegisteredClients))
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

func getMeanDisaster(pastDisastersList PastDisastersList) DisasterInfo {
	totalCoordinateX, totalCoordinateY, totalMagnitude, totalTurn := 0.0, 0.0, 0.0, 0.0
	numberDisastersPassed := float64(len(pastDisastersList))

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
	totalCoordinateX, totalCoordinateY, totalMagnitude, totalTurn := 0.0, 0.0, 0.0, 0.0
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
	varianceSum := (totalCoordinateX + totalCoordinateY + totalMagnitude + totalTurn) / numberDisastersPassed
	averageVariance := varianceSum / 4

	// Implement the variance cap chosen
	if averageVariance > varianceLimit {
		averageVariance = varianceLimit
	}

	// Return the confidence of the prediction
	return math.Round(varianceLimit - averageVariance)
}
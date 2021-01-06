package team3

import (
	"math"
  "fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
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

// GetGiftOffers allows clients to make offers in response to gift requests by other clients.
// It can offer multiple partial gifts.
// Strategy: We cover the risk that we lose money from the islands that we don’t trust with
// what we get from the islands that we do trust. Also, we don't request any gifts from critical islands.
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	var totalRequestAmt float64

	requests := shared.GiftRequestDict{}

	localPool := c.getLocalResources()

	resourcesNeeded := c.params.localPoolThreshold - float64(localPool)
	//fmt.Println("resources needed: ", resourcesNeeded)
	if resourcesNeeded > 0 {
		resourcesNeeded *= (1 + c.params.giftInflationPercentage)
		totalRequestAmt = resourcesNeeded
	} else {
		totalRequestAmt = c.params.giftInflationPercentage * c.params.localPoolThreshold
	}
	//fmt.Println("total request amount: ", totalRequestAmt)

	avgRequestAmt := totalRequestAmt / float64(c.getIslandsAlive()-c.getIslandsCritical())

	for island, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		if island == shared.Team3 {
			continue
		}
		if status == shared.Critical || status == shared.Dead {
			requests[island] = shared.GiftRequest(0.0)
		} else {
			var requestAmt float64
			requestAmt = avgRequestAmt * math.Pow(c.trustScore[island], c.params.trustParameter) * c.params.trustConstantAdjustor
			requests[island] = shared.GiftRequest(requestAmt)
		}
	}

	c.requestedGiftAmounts = requests
	return requests
}

func findAvgExclMinMax(Requests shared.GiftRequestDict) shared.GiftRequest {
	var sum shared.GiftRequest
	minClient := shared.TeamIDs[0]
	maxClient := shared.TeamIDs[0]

	// Find min and max requests
	for island, request := range Requests {
		if request < Requests[minClient] {
			minClient = island
		}
		if request > Requests[maxClient] {
			maxClient = island
		}
	}

	// Compute average ignoring highest and lowest
	for island, request := range Requests {
		if island != minClient || island != maxClient {
			sum += request
		}
	}

	return shared.GiftRequest(float64(sum) / float64(len(shared.TeamIDs)))
}

// GetGiftResponses allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}

	islandStatusCritical := c.isClientStatusCritical(shared.Team3)

	if islandStatusCritical {
		for island := range receivedRequests {
			offers[island] = 0.0
		}
		return offers
	}

	var amounts = map[shared.ClientID]shared.GiftRequest{}
	var allocations = map[shared.ClientID]float64{}
	var allocWeights = map[shared.ClientID]float64{}
	var allocSum float64
	var totalRequestedAmt float64
	var sumRequest shared.GiftRequest

	for island, request := range receivedRequests {
		sumRequest += request
		//amounts[island] = request * shared.GiftRequest(math.Pow(c.trustScore[island], c.params.trustParameter)+c.params.trustConstantAdjustor)
		amounts[island] = request * shared.GiftRequest(1+(math.Tanh(c.trustScore[island]-50)/100))
	}

	//fmt.Println("original amounts: ", amounts)

	for _, island := range c.getAliveIslands() {
		if island != c.BaseClient.GetID() && amounts[island] == 0.0 {
			amounts[island] = shared.GiftRequest(c.trustScore[island] * c.params.NoRequestGiftParam)
		}
	}

	// delete(amounts, c.GetID())

	fmt.Println("length of amounts map: ", len(amounts))

	avgRequest := findAvgExclMinMax(receivedRequests)
	avgAmount := findAvgExclMinMax(amounts)

	for island, amount := range amounts {
		allocations[island] = float64(avgRequest) + c.params.giftOfferEquity*float64((avgAmount-amount)+(receivedRequests[island]-avgRequest))
		fmt.Printf("allocation for island %v is %v", island, allocations[island])
		allocations[island] = math.Max(float64(receivedRequests[island]), allocations[island])
	}

	//fmt.Println("allocations: ", allocations)

	for _, alloc := range allocations {
		allocSum += alloc
	}
	for island, alloc := range allocations {
		allocWeights[island] = alloc / allocSum
	}

	for _, requests := range c.requestedGiftAmounts {
		totalRequestedAmt += float64(requests)
	}

	localPool := c.getLocalResources()

	giftBudget := (float64(localPool) + totalRequestedAmt) * ((1 - c.params.selfishness) / 2)

	fmt.Println("gift budged amount: ", giftBudget)

	for island := range amounts {
		if float64(sumRequest) < giftBudget {
			offers[island] = shared.GiftOffer(allocWeights[island] * float64(sumRequest))
		} else {
			offers[island] = shared.GiftOffer(allocWeights[island] * giftBudget)
		}
	}

	return offers
}

// GetGiftResponses returns the result of our island accepting/rejecting offered amounts.
// Strategy: we accept all amounts except if the offering island(s) is/are critical, then we
// do not take any of their offered amount(s).
func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
	responses := shared.GiftResponseDict{}
	for client, offer := range receivedOffers {
		responses[client] = shared.GiftResponse{
			AcceptedAmount: shared.Resources(offer),
			Reason:         shared.Accept,
		}
		if c.isClientStatusCritical(client) {
			responses[client] = shared.GiftResponse{
				AcceptedAmount: shared.Resources(0),
				Reason:         shared.DeclineDontNeed,
			}
		}
	}
	return responses
}

// UpdateGiftInfo receives responses from each island from the offers that our island made earlier.
// Strategy: The strategy is to update giftOpinion map on the reasons that gifts are declined or
// ignored for.
func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {
	for clientID, response := range receivedResponses {
		if response.Reason == shared.DeclineDontLikeYou {
			c.giftOpinions[clientID] -= 2
		} else if response.Reason == shared.Ignored {
			c.giftOpinions[clientID]--
		}
	}
}

// SentGift is executed at the end of each turn and notifies clients that
// their gift was successfully sent, along with offer details.
// We store the gift history for the previous turn.
func (c *client) SentGift(sent shared.Resources, to shared.ClientID) {
	c.sentGiftHistory = map[shared.ClientID]shared.Resources{}
	for island, amount := range c.sentGiftHistory {
		c.sentGiftHistory[island] = amount
	}

}

// ReceivedGift will inform us how much amount we have received from the specific islands
// and trust scores are then incremented and decremented based on the received difference.
func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {

	requestedFromIsland := c.requestedGiftAmounts[from]
	trustAdjustor := received - shared.Resources(requestedFromIsland)
	newTrustScore := 10 + (trustAdjustor * 0.2)
	if trustAdjustor >= 0 {
		c.updatetrustMapAgg(from, float64(newTrustScore))
	} else {
		if c.isClientStatusCritical(from) {
			return
		}
		c.updatetrustMapAgg(from, -float64(newTrustScore))
	}
}

// DecideGiftAmount is executed at the end of each turn and asks clients how much
// they wish to fulfill a gift offer they have previously made.
func (c *client) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {
	return giftOffer
}

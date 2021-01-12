package team3

import (
	"math"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {

	var predictionInfo shared.DisasterPredictionInfo
	trustedIslands := make([]shared.ClientID, len(c.BaseClient.ServerReadHandle.GetGameState().ClientLifeStatuses))
	for index, id := range shared.TeamIDs {
		trustedIslands[index] = id
	}

	if len(c.pastDisastersList) == 0 {
		predictionInfo = shared.DisasterPredictionInfo{
			PredictionMade: shared.DisasterPrediction{
				CoordinateX: 0,
				CoordinateY: 0,
				Magnitude:   0,
				TimeLeft:    0,
				Confidence:  0,
			},
			TeamsOfferedTo: trustedIslands,
		}
	} else {
		// Use the sample mean of each field as our prediction
		meanDisaster := getMeanDisaster(c.pastDisastersList)

		prediction := shared.DisasterPrediction{
			CoordinateX: meanDisaster.CoordinateX,
			CoordinateY: meanDisaster.CoordinateY,
			Magnitude:   meanDisaster.Magnitude,
			TimeLeft:    uint(meanDisaster.Turn),
		}

		// Use (variance limit - mean(sample variance)), where the mean is taken over each field, as confidence
		// Use a variance limit of 100 for now
		varianceLimit := 100.0
		prediction.Confidence = determineConfidence(c.pastDisastersList, meanDisaster, varianceLimit)

		// Return all prediction info and store our own island's prediction in global variable
		predictionInfo = shared.DisasterPredictionInfo{
			PredictionMade: prediction,
			TeamsOfferedTo: trustedIslands,
		}
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
	totalCoordinateX := selfConfidence * c.disasterPredictions[int(c.ServerReadHandle.GetGameState().Turn)-1][c.GetID()].CoordinateX
	totalCoordinateY := selfConfidence * c.disasterPredictions[int(c.ServerReadHandle.GetGameState().Turn)-1][c.GetID()].CoordinateY
	totalMagnitude := selfConfidence * c.disasterPredictions[int(c.ServerReadHandle.GetGameState().Turn)-1][c.GetID()].Magnitude
	totalTimeLeft := uint(math.Round(selfConfidence)) * c.disasterPredictions[int(c.ServerReadHandle.GetGameState().Turn)-1][c.GetID()].TimeLeft
	totalConfidence := selfConfidence

	// Add other island's predictions using their confidence values
	for islandID, prediction := range receivedPredictions {
		totalCoordinateX += safeDivFloat(c.trustScore[islandID], 100*prediction.PredictionMade.Confidence*prediction.PredictionMade.CoordinateX)
		totalCoordinateY += safeDivFloat(c.trustScore[islandID], 100*prediction.PredictionMade.Confidence*prediction.PredictionMade.CoordinateY)
		totalMagnitude += safeDivFloat(c.trustScore[islandID], 100*prediction.PredictionMade.Confidence*prediction.PredictionMade.Magnitude)
		totalTimeLeft += uint(math.Round(safeDivFloat(c.trustScore[islandID], 100*prediction.PredictionMade.Confidence))) * prediction.PredictionMade.TimeLeft
		totalConfidence += safeDivFloat(c.trustScore[islandID], 100*prediction.PredictionMade.Confidence)
	}

	// Finally get the final prediction generated by considering predictions from all islands that we have available
	// This result is currently unused but would be used in decision making in full implementation

	c.globalDisasterPredictions = append(c.globalDisasterPredictions, shared.DisasterPrediction{
		CoordinateX: totalCoordinateX / totalConfidence,
		CoordinateY: totalCoordinateY / totalConfidence,
		Magnitude:   totalMagnitude / totalConfidence,
		TimeLeft:    uint((float64(totalTimeLeft) / totalConfidence) + 0.5),
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
	var totalRequestAmt, avgRequestAmt float64

	requests := shared.GiftRequestDict{}

	localPool := c.getLocalResources()

	resourcesNeeded := float64(c.initialResourcesAtStartOfGame - localPool)
	//fmt.Println("resources needed: ", resourcesNeeded)
	if resourcesNeeded > 0 {
		resourcesNeeded *= (1 + c.params.giftInflationPercentage)
		totalRequestAmt = resourcesNeeded
	} else {
		totalRequestAmt = c.params.giftInflationPercentage * float64(c.initialResourcesAtStartOfGame)
	}
	//fmt.Println("total request amount: ", totalRequestAmt)

	// check to avoid division by 0 and only request from alive islands
	if len(c.getAliveIslands()) != 0 {
		avgRequestAmt = totalRequestAmt / float64(len(c.getAliveIslands()))
	} else {
		avgRequestAmt = totalRequestAmt
	}

	for island, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		if island == c.GetID() {
			continue
		}
		if status == shared.Critical || status == shared.Dead {
			requests[island] = shared.GiftRequest(0.0)
		} else {
			var requestAmt float64
			requestAmt = avgRequestAmt * math.Pow(c.trustScore[island], c.params.riskFactor)
			requests[island] = shared.GiftRequest(requestAmt)
		}
	}

	//	c.Logf("[TEAM3]: Actual requests made: %v", requests)
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

// sigmoidAndNormalise returns the normalised number between 0 - 1 based on the
// trust score between 0 - 100.
func (c *client) sigmoidAndNormalise(island shared.ClientID) shared.GiftOffer {

	sigmoid := (1 / (1 + math.Exp(-0.1*(c.trustScore[island]-50))))
	min := (1 / (1 + math.Exp(-0.1*-50)))
	max := (1 / (1 + math.Exp(-0.1*50)))

	normalised := (sigmoid - min) / (max - min)

	return shared.GiftOffer(normalised)
}

// GetGiftResponses allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}

	localPool := c.getLocalResources()
	islandStatusCritical := c.isClientStatusCritical(c.GetID())

	if islandStatusCritical {
		for _, island := range c.getAliveIslands() {
			offers[island] = 0.0
		}
		return offers
	} else if localPool < c.initialResourcesAtStartOfGame*0.1 {
		for _, island := range c.getAliveIslands() {
			offers[island] = 0.01
		}
		return offers
	}

	var amounts = map[shared.ClientID]shared.GiftOffer{}
	var totalRequestedAmt float64
	var sumRequest shared.GiftRequest

	for island, request := range receivedRequests {
		sumRequest += request
		amounts[island] = c.sigmoidAndNormalise(island) * shared.GiftOffer(request)
	}

	//fmt.Println("original amounts: ", amounts)

	for _, island := range c.getAliveIslands() {
		if island != c.GetID() && amounts[island] == 0.0 {
			amounts[island] = shared.GiftOffer(c.trustScore[island] * (c.params.friendliness / 30))
		}
	}

	//fmt.Println("length of amounts map: ", len(amounts))

	for _, requests := range c.requestedGiftAmounts {
		totalRequestedAmt += float64(requests)
	}

	giftBudget := shared.GiftOffer(float64(localPool) * ((1 - c.params.selfishness) / 2))
	rankedIslands := make([]shared.ClientID, 0.0, len(c.trustScore))

	for island := range c.trustScore {
		rankedIslands = append(rankedIslands, island)
	}

	sort.Slice(rankedIslands, func(i, j int) bool {
		return c.trustScore[rankedIslands[i]] > c.trustScore[rankedIslands[j]]
	})

	for _, island := range rankedIslands {
		if giftBudget >= amounts[island] {
			giftBudget -= amounts[island]
			offers[island] = amounts[island]
		} else if giftBudget == 0.0 {
			offers[island] = 0.0
		} else {
			offers[island] = giftBudget
			giftBudget = 0.0
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
			c.updatetrustMapAgg(clientID, -float64(5))
		} else if response.Reason == shared.Ignored {
			c.updatetrustMapAgg(clientID, -float64(2.5))
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

	c.Logf("Requested: %v and received: %v", c.requestedGiftAmounts[from], received)
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

	// Other agent simply look at if we offered them more than suggested, so why not 0.5%
	// more than we originally intended to offer. This would massively upgrade our opinion.
	newOffer := giftOffer * 1.005
	newOffer = shared.Resources(math.Min(float64(newOffer), float64(giftOffer+0.5)))
	return newOffer
}

func (c *client) MakeForageInfo() shared.ForageShareInfo {

	trustedIslands := make([]shared.ClientID, len(c.ServerReadHandle.GetGameState().ClientLifeStatuses))
	for index, id := range shared.TeamIDs {
		trustedIslands[index] = id
	}

	var lastDecision shared.ForageDecision
	var lastForageOutput shared.Resources

	for forageType, data := range c.forageData {
		for _, forage := range data {
			if uint(forage.turn) == c.ServerReadHandle.GetGameState().Turn-1 {
				lastForageOutput = forage.amountReturned
				lastDecision = shared.ForageDecision{
					Type:         forageType,
					Contribution: forage.amountContributed,
				}
			}
		}
	}

	forageInfo := shared.ForageShareInfo{
		DecisionMade:     lastDecision,
		ResourceObtained: lastForageOutput,
		ShareTo:          trustedIslands,
	}

	return forageInfo
}

func (c *client) ReceiveForageInfo(forageInfo []shared.ForageShareInfo) {
	if c.ServerReadHandle.GetGameState().Turn == 1 {
		c.forageData = make(map[shared.ForageType][]ForageData)
	}
	for _, val := range forageInfo {
		c.forageData[val.DecisionMade.Type] =
			append(
				c.forageData[val.DecisionMade.Type],
				ForageData{
					amountContributed: val.DecisionMade.Contribution,
					amountReturned:    val.ResourceObtained,
					turn:              c.ServerReadHandle.GetGameState().Turn,
				},
			)
	}
}

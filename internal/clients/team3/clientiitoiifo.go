package team3

import (
	"fmt"
	"math"

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
	DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources

	//TODO: IITO NON COMPULSORY
	// The server should handle the below functions maybe?
	func (c *BaseClient) SentGift(sent shared.Resources, to shared.ClientID)
	func (c *BaseClient) ReceivedGift(received shared.Resources, from shared.ClientID)

*/

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

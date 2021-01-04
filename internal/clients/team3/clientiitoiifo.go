package team3

import (
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

	// can change these parameters
	c.params.giftInflationPercentage = 0.2
	c.params.localPoolThreshold = 150
	c.params.trustParameter = 0.8
	c.params.giftConstantAdjustor = 10

	requests := shared.GiftRequestDict{}

	resourcesNeeded := c.params.localPoolThreshold - c.localPool
	if resourcesNeeded > 0 {
		resourcesNeeded *= (1 + c.params.giftInflationPercentage)
		totalRequestAmt = resourcesNeeded
	} else {
		totalRequestAmt = c.params.giftInflationPercentage * c.params.localPoolThreshold
	}

	avgRequestAmt := totalRequestAmt / float64(c.getIslandsAlive()-c.getIslandsCritical())

	for island, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		if status == shared.Critical {
			requests[island] = shared.GiftRequest(0.0)
		} else {
			var requestAmt float64
			if c.trustScore[island] >= 50 {
				requestAmt = avgRequestAmt * math.Pow(c.trustScore[island], c.params.trustParameter) * float64(c.params.giftConstantAdjustor)
			} else {
				requestAmt = avgRequestAmt * math.Pow(c.trustScore[island], -c.params.trustParameter) * float64(c.params.giftConstantAdjustor)
			}
			requests[island] = shared.GiftRequest(requestAmt)
		}
	}

	return requests
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
// Strategy: This function is first used to update our localPool with the accepted offer amount,
// and then used to adjust the trust score based on if they accepted or not. However, if an island
// is in critical state then they are exempt from trust score changes if they didn't offer full amount.
func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {
	for clientID := range receivedResponses {
		c.localPool += float64(receivedResponses[clientID].AcceptedAmount)
		trustAdjustor := receivedResponses[clientID].AcceptedAmount - c.requestedGiftAmounts[clientID]
		newTrustScore := 10 + (trustAdjustor * 0.2)
		if trustAdjustor >= 0 {
			c.updatetrustMapAgg(clientID, float64(newTrustScore))
		} else {
			if c.isClientStatusCritical(clientID) {
				continue
			}
			c.updatetrustMapAgg(clientID, -float64(newTrustScore))
		}
	}
}

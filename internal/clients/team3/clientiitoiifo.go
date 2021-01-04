package team3

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

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

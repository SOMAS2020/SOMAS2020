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
	UpdateGiftInfo(receivedResponses shared.GiftResponseDict)
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

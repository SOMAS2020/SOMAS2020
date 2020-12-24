package baseclient

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// FIXME: Wait for implementation of getClientsAlive()
const aliveClients = 6

// GetGiftRequests allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
// COMPULSORY, you need to implement this method
func (c *BaseClient) GetGiftRequests() shared.GiftRequestDict {
	requests := make(shared.GiftRequestDict, aliveClients)
	for key := range requests {
		if c.clientGameState.ClientInfo.LifeStatus == shared.Critical {
			requests[key] = shared.GiftRequest{RequestingTeam: c.GetID(), OfferingTeam: key, RequestAmount: 100.0}
		} else {
			requests[key] = shared.GiftRequest{RequestingTeam: c.GetID(), OfferingTeam: key, RequestAmount: 0.0}
		}
	}
	return requests
}

// GetGiftOffers allows clients to make offers in response to gift requests by other clients.
// It can offer multiple partial gifts.
// COMPULSORY, you need to implement this method. This placeholder implementation offers no gifts.
func (c *BaseClient) GetGiftOffers(receivedRequests shared.GiftRequestDict) (shared.GiftOfferDict, error) {
	offers := make(shared.GiftOfferDict, aliveClients)
	for key := range offers {
		offers[key] = shared.GiftOffer{ReceivingTeam: c.GetID(), OfferingTeam: key, OfferAmount: 0.0}
	}
	return offers, nil
}

// GetGiftResponses allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
// COMPULSORY, you need to implement this method
func (c *BaseClient) GetGiftResponses(receivedOffers shared.GiftOfferDict) (shared.GiftResponseDict, error) {
	responses := shared.GiftResponseDict{}
	for client, offer := range receivedOffers {
		responses[client] = shared.GiftResponse{
			AcceptedAmount: offer.OfferAmount,
			Reason:         shared.Accept,
		}
	}
	return responses, nil
}

// UpdateGiftInfo gives information about the outcome from AcceptGifts.
// This allows for opinion formation.
// COMPULSORY, you need to implement this method
func (c *BaseClient) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) error {
	// PreviousGifts[count] = acceptedGifts
	// count++
	return nil
}

// SendGift is executed at the end of each turn and allows clients to
// send the gifts promised in the IITO
// COMPULSORY, you need to implement this method
func (c *BaseClient) SendGift(offer shared.GiftOffer) error {
	return nil
}

// ReceiveGift is executed at the end of each turn and allows clients to
// receive the gifts promised in the IITO
// COMPULSORY, you need to implement this method
func (c *BaseClient) ReceiveGift(offer shared.GiftOffer) error {
	return nil
}

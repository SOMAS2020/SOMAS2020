package baseclient

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// GetGiftRequests allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
// COMPULSORY, you need to implement this method
func (c *BaseClient) GetGiftRequests() []shared.GiftRequest {
	// FIXME: Make this the appropriate length, the loop below doesn't run otherwise.
	requests := []shared.GiftRequest{}
	for fromTeam := range requests {
		if c.clientGameState.ClientInfo.LifeStatus == shared.Critical {
			requests[fromTeam] = shared.GiftRequest{RequestFrom: shared.ClientID(fromTeam), RequestAmount: 100.0}
		} else {
			requests[fromTeam] = shared.GiftRequest{RequestFrom: shared.ClientID(fromTeam), RequestAmount: 0.0}
		}
	}
	return requests
}

// GetGiftOffers allows clients to make offers in response to gift requests by other clients.
// It can offer multiple partial gifts.
// COMPULSORY, you need to implement this method. This placeholder implementation offers no gifts.
func (c *BaseClient) GetGiftOffers(receivedRequests []shared.GiftRequest) ([]shared.GiftOffer, error) {
	// FIXME: Make this the appropriate length, the loop below doesn't run otherwise.
	offers := []shared.GiftOffer{}
	for toTeam := range offers {
		offers[toTeam] = shared.GiftOffer{ReceivingTeam: shared.ClientID(toTeam), OfferAmount: 0.0}
	}
	return offers, nil
}

// GetGiftResponses allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
// COMPULSORY, you need to implement this method
func (c *BaseClient) GetGiftResponses(receivedOffers []shared.GiftOffer) ([]shared.GiftResponse, error) {
	// FIXME: Make this the appropriate length, the loop below doesn't run otherwise.
	responses := []shared.GiftResponse{}
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
func (c *BaseClient) UpdateGiftInfo(receivedResponses []shared.GiftResponse) error {
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

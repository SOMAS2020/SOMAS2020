package baseclient

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// GetGiftRequests allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
// COMPULSORY, you need to implement this method
func (c *BaseClient) GetGiftRequests() shared.GiftRequestDict {
	requests := shared.GiftRequestDict{}

	// You can fetch the clients which are alive like this:
	for team, status := range c.serverReadHandle.GetGameState().ClientLifeStatuses {
		if status == shared.Critical {
			requests[team] = shared.GiftRequest(100.0)
		} else {
			requests[team] = shared.GiftRequest(0.0)
		}
	}
	return requests
}

// GetGiftOffers allows clients to make offers in response to gift requests by other clients.
// It can offer multiple partial gifts.
// COMPULSORY, you need to implement this method. This placeholder implementation offers no gifts,
// unless another team is critical.
func (c *BaseClient) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}

	// You can fetch the clients which are alive like this:
	for team, status := range c.serverReadHandle.GetGameState().ClientLifeStatuses {
		if status == shared.Critical {
			offers[team] = shared.GiftOffer(100.0)
		} else {
			offers[team] = shared.GiftOffer(0.0)
		}
	}
	return offers
}

// GetGiftResponses allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
// COMPULSORY, you need to implement this method
func (c *BaseClient) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
	responses := shared.GiftResponseDict{}
	for client, offer := range receivedOffers {
		responses[client] = shared.GiftResponse{
			AcceptedAmount: shared.Resources(offer),
			Reason:         shared.Accept,
		}
	}
	return responses
}

// UpdateGiftInfo will be called by the server with all the responses you received
// from the gift session. It is up to you to complete the transactions with the methods
// that you will implement yourself below. This allows for opinion formation.
// COMPULSORY, you need to implement this method
func (c *BaseClient) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) error {
	return nil
}

// SentGift is executed at the end of each turn and notifies clients that
// their gift was successfully sent, along with the offer details.
// COMPULSORY, you need to implement this method
func (c *BaseClient) SentGift(sent shared.Resources, to shared.ClientID) error {
	// You can check your updated resources like this:
	// myResources := c.serverReadHandle.GetGameState().ClientInfo.Resources
	return nil
}

// ReceivedGift is executed at the end of each turn and notifies clients that
// their gift was successfully received, along with the offer details.
// COMPULSORY, you need to implement this method
func (c *BaseClient) ReceivedGift(received shared.Resources, from shared.ClientID) error {
	// You can check your updated resources like this:
	// myResources := c.serverReadHandle.GetGameState().ClientInfo.Resources
	return nil
}

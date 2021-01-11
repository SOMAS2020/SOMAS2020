package baseclient

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// GetGiftRequests allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
// COMPULSORY, you need to implement this method
func (c *BaseClient) GetGiftRequests() shared.GiftRequestDict {
	requests := shared.GiftRequestDict{}

	// You can fetch the clients which are alive like this:
	for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
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
	for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
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
func (c *BaseClient) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {
}

// SentGift is executed at the end of each turn and notifies clients that
// their gift was successfully sent, along with the offer details.
// COMPULSORY, you need to implement this method
func (c *BaseClient) SentGift(sent shared.Resources, to shared.ClientID) {
	// You can check your updated resources like this:
	// myResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
}

// ReceivedGift is executed at the end of each turn and notifies clients that
// their gift was successfully received, along with the offer details.
// COMPULSORY, you need to implement this method
func (c *BaseClient) ReceivedGift(received shared.Resources, from shared.ClientID) {
	// You can check your updated resources like this:
	// myResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
}

// ShareIntendedContribution is called on each client to give them the
// option of sharing how much they intend to contribute to the common pool. They must also
// choose what islands they want to share the information with
// OPTIONAL, you can implement this if you want to partake in the iito session
func (c *BaseClient) ShareIntendedContribution() shared.IntendedContribution {

	// For MVP, share this prediction with all islands since trust has not yet been implemented
	trustedIslands := make([]shared.ClientID, len(shared.TeamIDs))
	for index, id := range shared.TeamIDs {
		trustedIslands[index] = id
	}

	contribution := shared.IntendedContribution{
		Contribution:   shared.Resources(rand.Float64()),
		TeamsOfferedTo: trustedIslands,
	}

	c.intendedContribution = contribution
	return contribution

}

// ReceiveIntendedContribution provides each client with the intended common pool contribution from the islands
// that have chosen to share it with them
// OPTIONAL, you can implement this if you want to partake in the iito session
func (c *BaseClient) ReceiveIntendedContribution(receivedIntendedContribution shared.ReceivedIntendedContributionDict) {
	// You can check the other's common pool contributions like this
	// intededContributions := c.intendedContribution
}

// DecideGiftAmount is executed at the end of each turn and asks clients how much
// they want to fulfill a gift offer they have made.
// COMPULSORY, you need to implement this method
func (c *BaseClient) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {
	return giftOffer
}

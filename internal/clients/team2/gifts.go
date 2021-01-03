package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// GetGiftRequests allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	turn := c.gameState().Turn
	requests := shared.GiftRequestDict{}

	// check our critical and threshold - if either is off - request
	ourAgentCritical := status == shared.Critical
	requestAmount := internalThreshold(c) - c.gameState().ClientInfo.Resources

	// who we request to: using trust + whether or not they are critical
	
	
	if ourAgentCritical || requestAmount > 0 {
		// how much we request depends on: how much is in common pool + how much we have
		// ask for difference between threshold and what we currently have? 
		// can split that amount amongst the number of people we want to ask for help from
		// but not linearly - e.g. might not be 50/50
		// assume that smaller requests are more likely to be approved
		// notice if an agent is accepting large requests and abuse it - could use confidence

		
		// Split up the goal request amount between islands to maximise likelihood of receiving gift
		for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
			
			

			if isCritical  {
				// the amount that we request should be related to how much we need/want
				// low requests are more likely to be approved so it may be better to make more *small* requests
	
				requests[team] = shared.GiftRequest(100.0)
			} else {
				requests[team] = shared.GiftRequest(0.0)
			}
		}
	}

	for island, requestAmount := range requests {
		newGiftRequest := GiftInfo{
			requested: requestAmount,
		}
		c.giftHist[island].OurRequest[turn] = newGiftRequest
	}
	return requests
}

// GetGiftOffers allows clients to make offers in response to gift requests by other clients.
// It can offer multiple partial gifts.
// COMPULSORY, you need to implement this method. This placeholder implementation offers no gifts,
// unless another team is critical.
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}
	turn := c.gameState().Turn

	// if we are critical do not offer gifts-> there should be a way to see which other islands are critical
	// if we are not critical and another island is critical offer gift
	// do not offer more than proportion of total resources we have
// unless another team is critical.
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}
	turn := c.gameState().Turn

	// if we are critical do not offer gifts-> there should be a way to see which other islands are critical
	// if we are not critical and another island is critical offer gift
	// do not offer more than proportion of total resources we have

	// You can fetch the clients which are alive like this:
	for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		if status == shared.Critical {
			offers[team] = shared.GiftOffer(100.0)
		} else {
			offers[team] = shared.GiftOffer(0.0)
		}
	}

	for island, offeredAmount := range offers {
		newGiftRequest := GiftInfo{
			requested: receivedRequests[island],
			gifted:    offeredAmount,
		}
		c.giftHist[island].IslandRequest[turn] = newGiftRequest
	}

	return offers
}

// GetGiftResponses allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
// COMPULSORY, you need to implement this method

func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
	responses := shared.GiftResponseDict{}
	turn := c.gameState().Turn

	for client, offer := range receivedOffers {
		responses[client] = shared.GiftResponse{
			AcceptedAmount: shared.Resources(offer),
			Reason:         shared.Accept,
		}
	}

	for island, response := range responses {
		newGiftRequest := GiftInfo{
			requested: c.giftHist[island].OurRequest[turn].requested,
			gifted:    shared.GiftOffer(response.AcceptedAmount),
			reason:    shared.AcceptReason(response.AcceptedAmount),
		}
		c.giftHist[island].OurRequest[turn] = newGiftRequest
	}
	return responses
}

// UpdateGiftInfo will be called by the server with all the responses you received
// from the gift session. It is up to you to complete the transactions with the methods
// that you will implement yourself below. This allows for opinion formation.
// COMPULSORY, you need to implement this method
func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {
	turn := c.gameState().Turn

	for island, response := range receivedResponses {
		newGiftRequest := GiftInfo{
			requested: c.giftHist[island].IslandRequest[turn].requested,
			gifted:    shared.GiftOffer(response.AcceptedAmount),
			reason:    shared.AcceptReason(response.AcceptedAmount),
		}
		c.giftHist[island].IslandRequest[turn] = newGiftRequest
	}
}

// SentGift is executed at the end of each turn and notifies clients that
// their gift was successfully sent, along with the offer details.
// COMPULSORY, you need to implement this method
func (c *client) SentGift(sent shared.Resources, to shared.ClientID) {
	// You can check your updated resources like this:
	// myResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources

}

// ReceivedGift is executed at the end of each turn and notifies clients that
// their gift was successfully received, along with the offer details.
// COMPULSORY, you need to implement this method
func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {
	// You can check your updated resources like this:
	// myResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources

}

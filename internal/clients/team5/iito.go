package team5

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

//=============================================================
//=======GIFTS=================================================
//=============================================================

/*
	COMPULSORY:
	GetGiftRequests() shared.GiftRequestDict
		- signal that we want a gift. This info is shared to all other clients.
	GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict
		- allows us to make an offer to other teams if they request.
	GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict
		- allows us to consider gift offers from other teams and accept/decline
		- need to provide reasoning for not accepting full amount if that is the case
	UpdateGiftInfo(receivedResponses shared.GiftResponseDict)
		- allows client to notify server of all gift responses we received. We would need to store these
		transactions to be able to inform server when it calls this method
	SentGift(sent shared.Resources, to shared.ClientID)
		- notifies us that gift was successfully sent. Can use this as a time to check resource deduction.
	ReceivedGift(sent shared.Resources, from shared.ClientID)
		- notifies us that gift was successfully received. Can use this as a time to check resource addition.
	DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources
		- called by server at end of round. Allows us to choose how much of our gift offers we actually want to fulfill.
*/

// GetGiftRequests allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
// COMPULSORY, you need to implement this method
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	requests := shared.GiftRequestDict{}
	switch {
	case c.gameState().ClientInfo.LifeStatus == shared.Critical:
		for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
			if status == shared.Critical {
				requests[team] = shared.GiftRequest(0.0)
			} else {
				requests[team] = shared.GiftRequest(c.config.DyingGiftRequest)
			}
		}
	case c.wealth() == ImperialStudent: // Poor
		for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
			if status == shared.Critical {
				requests[team] = shared.GiftRequest(0.0)
			} else {
				requests[team] = shared.GiftRequest(c.config.ImperialGiftRequest)
			}
		}
	default:
		for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
			if status == shared.Critical {
				requests[team] = shared.GiftRequest(0.0)
			} else {
				requests[team] = shared.GiftRequest(c.config.MiddleGiftRequest)
			}
		}
	}
	c.Logf("Team 5 is requesting %v", requests)
	return requests
}

// GetGiftOffers allows clients to make offers in response to gift requests by other clients.
// It can offer multiple partial gifts.
// COMPULSORY, you need to implement this method. This placeholder implementation offers no gifts,
// unless another team is critical.
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}
	if c.wealth() >= 2 { // Middle class and JB
		for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
			if status == shared.Critical {
				offers[team] = shared.GiftOffer(3.0)
			} else {
				offers[team] = shared.GiftOffer(1.0) // can we abuse the fact that they look at the amount of gifts?
			}
		}
	} else {
		for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
			if status == shared.Critical {
				offers[team] = shared.GiftOffer(1.5)
			} else {
				offers[team] = shared.GiftOffer(0.5) // can we abuse the fact that they look at the amount of gifts?
			}
		}
	}
	return offers
}

// GetGiftResponses allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
// COMPULSORY, you need to implement this method
func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
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
// func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {
// 	turn := c.gameState().Turn

// 	for island, response := range receivedResponses {
// 		newGiftRequest := GiftInfo{
// 			requested: c.giftHistory[island].IslandRequest[turn].requested,
// 			gifted:    shared.GiftOffer(response.AcceptedAmount),
// 			reason:    shared.AcceptReason(response.AcceptedAmount),
// 		}
// 		c.giftHistory[island].IslandRequest[turn] = newGiftRequest
// 	}
// }

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

// DecideGiftAmount is executed at the end of each turn and asks clients how much
// they want to fulfill a gift offer they have made.
// COMPULSORY, you need to implement this method
func (c *client) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {
	if c.resourceHistory[c.gameState().Turn-1] < c.ServerReadHandle.GetGameState().ClientInfo.Resources {
		if c.wealth() >= 2 { //Middle class and JB, fulfill all offers
			return giftOffer
		} else if c.wealth() == 1 { //When Imperial Student fulfill all offers but divide Team 3's by 3
			if toTeam == 2 {
				return (giftOffer / 3)
			} else {
				return giftOffer
			}
		} else { //Reject all offers
			return 0
		}
	} else { //Reject all offers
		return 0
	}
}

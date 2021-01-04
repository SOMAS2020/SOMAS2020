package team5

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

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

// Requesting gifts from other islands
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	requests := shared.GiftRequestDict{}
	switch {
	case c.gameState().ClientInfo.LifeStatus == shared.Critical: // Case we are critical
		for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
			if status == shared.Critical { // Other island are critical
				requests[team] = shared.GiftRequest(0.0) // Dont ask for money
			} else {
				requests[team] = shared.GiftRequest(c.config.DyingGiftRequestAmount) // Ask for money cus we dying
			}
		}
	case c.wealth() == ImperialStudent: // We are poor
		for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
			if status == shared.Critical {
				requests[team] = shared.GiftRequest(0.0) // Dont ask from people if they are dying
			} else {
				requests[team] = shared.GiftRequest(c.config.ImperialGiftRequestAmount) // Ask for money
			}
		}
	default:
		for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
			if status == shared.Critical {
				requests[team] = shared.GiftRequest(0.0)
			} else {
				requests[team] = shared.GiftRequest(c.config.MiddleGiftRequestAmount) // Ask for money
			}
		}
	}
	c.Logf("[Debug] Team 5 Gift request: %v", requests)
	return requests
}

// GetGiftOffers allows clients to make offers in response to gift requests by other clients.
// It can offer multiple partial gifts.
// COMPULSORY, you need to implement this method. This placeholder implementation offers no gifts,
// unless another team is critical.
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}
	if c.wealth() >= 2 { // We are >= Middle class / JB
		for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
			if status == shared.Critical {
				offers[team] = shared.GiftOffer(3.0) // Give 3 to dying islands
			} else {
				offers[team] = shared.GiftOffer(1.0) // gift a dollar
			} // can we abuse the fact that they look at the amount of gifts?
		}
	} else { // Other cases we gift half the amount
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
			Reason:         shared.Accept, // Accept all gifts duh
		}
	}
	return responses
}

// =====================================Gifting history to be made==========================================

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
// ==================================== Gifting history to be made =========================================

// ===================================== Has sending / recv gifts been implemented? ===============================
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

// ===================================== Has sending / recv gifts been implemented? ===============================

// DecideGiftAmount is executed at the end of each turn and asks clients how much
// they want to fulfill a gift offer they have made.
// COMPULSORY, you need to implement this method
func (c *client) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {
	if c.resourceHistory[c.gameState().Turn-1] < c.ServerReadHandle.GetGameState().ClientInfo.Resources { // Only if we have more resources than the previous round
		if c.wealth() >= 2 { //Middle class and JB, fulfill all offers
			return giftOffer // Give what we said
		} else if c.wealth() == 1 { //When Imperial Student fulfill all offers but divide Team 3's by 3
			if toTeam == shared.Team3 {
				return (giftOffer / 3)
			}
			return giftOffer // Return all offers if we are in imperial class but Team 3 / 3
		} else {
			return 0 // Reject all offers if we are in critical state
		}
	} else { //Reject all offers if we have less than we did last round
		return 0
	}
}

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
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
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
func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {
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

// DecideGiftAmount is executed at the end of each turn and asks clients how much
// they want to fulfill a gift offer they have made.
// COMPULSORY, you need to implement this method
func (c *client) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {
	return giftOffer
}

// GiftAccceptance
func (c *client) giftAcceptance() shared.GiftResponse {

	if c.wealth() == JeffBezos {
		c.Logf("I don't need any gifts you peasants")
	} else if c.wealth() == MiddleClass {
		c.Logf("We only accept Gifts from Team 1, 2, 4, 6")
	} else if c.wealth() == ImperialStudent {
		c.Logf("We accept Gifts from everyone")
	} else {
		c.Logf("We accept Gifts from everyone")
	}

	GiftResponse := shared.GiftResponse{
		AcceptedAmount: 0,
		Reason:         shared.DeclineDontLikeYou,
	}
	return GiftResponse
}

/*Gift Requests*/
func (c *client) giftRequests() shared.GiftRequest {
	var giftrequest shared.GiftRequest

	if c.wealth() == JeffBezos {
		c.Logf("I don't need any gifts you peasants")
	} else if c.wealth() == MiddleClass {
		c.Logf("We only accept Gifts from Team 1, 2, 4, 6")
	} else if c.wealth() == ImperialStudent {
		c.Logf("We accept Gifts from everyone")
	} else {
		c.Logf("We accept Gifts from everyone")
	}

	return giftrequest
}

package team5

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

/*
	COMPULSORY:
	GetGiftRequests() shared.GiftRequestDict
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

/*
	For refrence
	OurRequest giftInfo
	struct {
		requested      shared.GiftRequest    - GetGiftRequest - Amount WE request 1
		offered        shared.GiftOffer 	- GetGiftResponses - Amount offered TO US 3
		response       shared.GiftResponse	- GetGiftResponses - Amount WE accepeted and reason 3
		actualRecieved shared.Resources		- ReceivedGift - Amount WE actually received 6
	}

	TheirRequest giftInfo
	struct {
		requested      shared.GiftRequest    - GetGiftOffers - Amount THEY requested  2
		offered        shared.GiftOffer		- GetGiftOffers - Amount WE offered 2
		response       shared.GiftResponse 	- UpdateGiftInfo - Amount THEY accepeted and reason 4
		actualRecieved shared.Resources		- DecideGiftAmount - Amount THEY actually get 5
}
*/

// creates initial opinions of clients and sets the values to 0 (prevents nil mapping)
func (c *client) initGiftHist() {
	c.giftHistory = giftHistory{}

	for _, team := range c.getAliveTeams(true) {
		ourGiftInfo := giftInfo{
			requested:      0,
			offered:        0,
			response:       shared.GiftResponse{AcceptedAmount: 0, Reason: 0},
			actualRecieved: 0,
		}
		theirGiftInfo := giftInfo{
			requested:      0,
			offered:        0,
			response:       shared.GiftResponse{AcceptedAmount: 0, Reason: 0},
			actualRecieved: 0,
		}
		ourReq := map[uint]giftInfo{c.gameState().Turn: ourGiftInfo}
		theirReq := map[uint]giftInfo{c.gameState().Turn: theirGiftInfo}
		c.giftHistory[team] = giftExchange{
			OurRequest:   ourReq,
			TheirRequest: theirReq,
		}
	}
}

// GetGiftRequests we want gifts!
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	requests := shared.GiftRequestDict{}
	switch {
	case c.gameState().ClientInfo.LifeStatus == shared.Critical: // Case we are critical
		for team, status := range c.gameState().ClientLifeStatuses {
			if status == shared.Critical { // Other island are critical
				requests[team] = shared.GiftRequest(0.0) // Dont ask for money
			} else {
				requests[team] = shared.GiftRequest(c.config.dyingGiftRequestAmount) // Ask for money cus we dying
			}
		}
	case c.wealth() == imperialStudent: // We are poor
		for team, status := range c.gameState().ClientLifeStatuses {
			if status == shared.Critical {
				requests[team] = shared.GiftRequest(0.0) // Dont ask from people if they are dying
			} else {
				requests[team] = shared.GiftRequest(c.config.imperialGiftRequestAmount) // Ask for money
			}
		}
	default:
		for team, status := range c.gameState().ClientLifeStatuses {
			if status == shared.Critical {
				requests[team] = shared.GiftRequest(0.0)
			} else {
				requests[team] = shared.GiftRequest(c.config.middleGiftRequestAmount) // Ask for money
			}
		}
	}

	// History
	for team := range c.gameState().ClientLifeStatuses {
		newGiftRequest := giftInfo{
			requested: requests[team], // Amount WE request
		}
		c.giftHistory[team].OurRequest[c.gameState().Turn] = newGiftRequest
		c.Logf("GetGiftRequestsTeam [%v] we REQUEST this much %v", team, c.giftHistory[team].OurRequest[c.gameState().Turn])
	}

	return requests
}

// GetGiftOffers allows clients to make offers in response to gift requests by other clients.
// It can offer multiple partial gifts.
// COMPULSORY, you need to implement this method. This placeholder implementation offers no gifts,
// unless another team is critical.
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}
	if c.wealth() >= 2 { // We are >= Middle class / JB
		for team, status := range c.gameState().ClientLifeStatuses {
			if status == shared.Critical {
				offers[team] = shared.GiftOffer(3.0) // Give 3 to dying islands
			} else {
				offers[team] = shared.GiftOffer(1.0) // gift a dollar
			} // can we abuse the fact that they look at the amount of gifts?
		}
	} else { // Other cases we gift half the amount
		for team, status := range c.gameState().ClientLifeStatuses {
			if status == shared.Critical {
				offers[team] = shared.GiftOffer(1.5)
			} else {
				offers[team] = shared.GiftOffer(0.5) // can we abuse the fact that they look at the amount of gifts?
			}
		}
	}

	// Store the offers we gave into giftHistory
	for team := range c.gameState().ClientLifeStatuses {
		newGiftRequest := giftInfo{
			requested: receivedRequests[team], // Amount THEY requested
			offered:   offers[team],           // Amount WE offered
		}
		c.giftHistory[team].TheirRequest[c.gameState().Turn] = newGiftRequest
		c.Logf("GetGiftOffers [%v] we OFFER this much %v", team, c.giftHistory[team].TheirRequest[c.gameState().Turn])
	}

	return offers
}

// GetGiftResponses allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
// COMPULSORY, you need to implement this method
func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
	responses := shared.GiftResponseDict{}
	for team, offer := range receivedOffers { // For all the clients we look at the offers
		responses[team] = shared.GiftResponse{
			AcceptedAmount: shared.Resources(offer), // Accept all they gave us
			Reason:         shared.Accept,           // Accept all gifts duh
		}
	}

	for team := range c.gameState().ClientLifeStatuses {
		newGiftRequest := giftInfo{ // 	For each client create a new gift info
			requested: c.giftHistory[team].OurRequest[c.gameState().Turn].requested, // Amount We requested
			offered:   receivedOffers[team],                                         // Amount offered TO US
			response:  responses[team],                                              // Amount and reason WE accepted
		}
		c.giftHistory[team].OurRequest[c.gameState().Turn] = newGiftRequest
		c.Logf("GetGiftResponses [%v] %v", team, c.giftHistory[team].OurRequest[c.gameState().Turn])

	}

	return responses
}

// =====================================Gifting history to be made==========================================

// UpdateGiftInfo will be called by the server with all the responses you received
// from the gift session. It is up to you to complete the transactions with the methods
// that you will implement yourself below. This allows for opinion formation.
// COMPULSORY, you need to implement this method

func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {

	for team := range c.gameState().ClientLifeStatuses { // for each ID
		newGiftRequest := giftInfo{
			requested: c.giftHistory[team].TheirRequest[c.gameState().Turn].requested, // Amount THEY requested
			offered:   c.giftHistory[team].TheirRequest[c.gameState().Turn].offered,   // Amount WE offered them
			response:  receivedResponses[team],                                        // Amount THEY accepted and REASON
		}
		c.giftHistory[team].TheirRequest[c.gameState().Turn] = newGiftRequest
		c.Logf("UpdateGiftInfo [%v] %v", team, c.giftHistory[team].TheirRequest[c.gameState().Turn])
	}
}

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

	if c.wealth() < 2 { // Imperial or Dying
		giftOffer = 0
	}

	// Regardless if we have less than last turn we dont give gifts away
	if c.resourceHistory[c.gameState().Turn-1] > c.ServerReadHandle.GetGameState().ClientInfo.Resources {
		giftOffer = 0
	}

	// History
	newGiftRequest := giftInfo{ // 	For each client create a new gift info
		requested:      c.giftHistory[toTeam].TheirRequest[c.gameState().Turn].requested, // Amount We requested
		offered:        c.giftHistory[toTeam].TheirRequest[c.gameState().Turn].offered,   // Amount offered TO US
		response:       c.giftHistory[toTeam].TheirRequest[c.gameState().Turn].response,  // Amount and reason WE accepted
		actualRecieved: giftOffer,
	}
	c.giftHistory[toTeam].TheirRequest[c.gameState().Turn] = newGiftRequest

	c.Logf("DecideGiftAmount [%v]", c.giftHistory[toTeam].TheirRequest[c.gameState().Turn])

	return giftOffer
}

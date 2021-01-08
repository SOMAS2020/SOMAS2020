package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

/*
	COMPULSORY:
	GetGiftRequests() shared.GiftRequestDict
	GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict
	GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict
	UpdateGiftInfo(receivedResponses shared.GiftResponseDict)
	DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources
	SentGift(sent shared.Resources, to shared.ClientID)
	ReceivedGift(sent shared.Resources, from shared.ClientID)

	For reference
	ourRequest giftInfo
		requested      shared.GiftRequest    - GetGiftRequest - Amount WE request 1
		offered        shared.GiftOffer 	- GetGiftResponses - Amount offered TO US 3
		response       shared.GiftResponse	- GetGiftResponses - Amount WE accepeted and reason 3
		actualReceived shared.Resources		- ReceivedGift - Amount WE actually received 6

	theirRequest giftInfo
		requested      shared.GiftRequest    - GetGiftOffers - Amount THEY requested  2
		offered        shared.GiftOffer		- GetGiftOffers - Amount WE offered 2
		response       shared.GiftResponse 	- UpdateGiftInfo - Amount THEY accepeted and reason 4
		actualReceived shared.Resources		- DecideGiftAmount - Amount THEY actually get 5
*/

// creates initial opinions of clients and sets the values to 0 (prevents nil mapping)
func (c *client) initGiftHist() {
	c.giftHistory = giftHistory{}

	for _, team := range c.getAliveTeams(true) {
		ourGiftInfo := giftInfo{
			requested:      0,
			offered:        0,
			response:       shared.GiftResponse{AcceptedAmount: 0, Reason: 0},
			actualReceived: 0,
		}
		theirGiftInfo := giftInfo{
			requested:      0,
			offered:        0,
			response:       shared.GiftResponse{AcceptedAmount: 0, Reason: 0},
			actualReceived: 0,
		}
		ourReq := map[uint]giftInfo{c.getTurn(): ourGiftInfo}
		theirReq := map[uint]giftInfo{c.getTurn(): theirGiftInfo}
		c.giftHistory[team] = giftExchange{
			ourRequest:   ourReq,
			theirRequest: theirReq,
		}
	}
}

// GetGiftRequests we want gifts!
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	requests := shared.GiftRequestDict{}
	for team, status := range c.gameState().ClientLifeStatuses {
		if status != shared.Critical { // We are not dying
			if c.opinions[team].getScore() > c.config.opinionThresholdRequest { // They are nice people
				switch c.wealth() { // Look at our wealth
				case dying: // Case we are Critical (dying)
					requests[team] = shared.GiftRequest(c.config.dyingGiftRequestAmount * 0.8) // Ask  for dying amount
				case imperialStudent:
					requests[team] = shared.GiftRequest( // Scale our request
						c.config.imperialGiftRequestAmount /
							float64(c.opinions[team].getScore()) *
							c.config.opinionRequestMultiplier) // Scale down the request if we like them
				case middleClass:
					requests[team] = shared.GiftRequest(
						c.config.middleGiftRequestAmount /
							float64(c.opinions[team].getScore()) *
							c.config.opinionRequestMultiplier) // Scale down the request if we like them
				default:
					requests[team] = shared.GiftRequest(0)
				}
			} else if c.opinions[team].getScore() < (-c.config.opinionThresholdRequest) { // They are trashy people
				switch c.wealth() {
				case dying: // Case we are Critical
					requests[team] = shared.GiftRequest(c.config.dyingGiftRequestAmount * 1.2)
				case imperialStudent: // Case we are Poor
					requests[team] = shared.GiftRequest(
						c.config.imperialGiftRequestAmount /
							float64(c.opinions[team].getScore()) *
							c.config.opinionRequestMultiplier) // Scale down the request if we like them
				case middleClass:
					requests[team] = shared.GiftRequest(
						c.config.middleGiftRequestAmount /
							float64(c.opinions[team].getScore()) *
							c.config.opinionRequestMultiplier) // Scale down the request if we like them
				default:
					requests[team] = shared.GiftRequest(0)
				}
			} else { // Normal people
				switch c.wealth() {
				case dying: // Case we are Critical
					requests[team] = shared.GiftRequest(c.config.dyingGiftRequestAmount)
				case imperialStudent:
					requests[team] = shared.GiftRequest(c.config.imperialGiftRequestAmount) // Scale down the request if we like them
				case middleClass:
					requests[team] = shared.GiftRequest(c.config.middleGiftRequestAmount) // Scale down the request if we like them
				default:
					requests[team] = shared.GiftRequest(0)
				}
			}
		}
		// History
		newGiftRequest := giftInfo{
			requested: requests[team], // Amount WE request
		}
		c.giftHistory[team].ourRequest[c.getTurn()] = newGiftRequest
	}
	return requests
}

// GetGiftOffers allows clients to make offers in response to gift requests by other clients.
// It can offer multiple partial gifts.
// COMPULSORY, you need to implement this method. This placeholder implementation offers no gifts,
// unless another team is critical.

func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}
	for team, status := range c.gameState().ClientLifeStatuses {
		status := shared.ClientLifeStatus(status)
		switch {
		case c.wealth() >= 2: // case we are rich
			if status == shared.Critical {
				offers[team] = shared.GiftOffer(c.config.offertoDyingIslands) // Give 3 to dying islands
			} else {
				offers[team] = shared.GiftOffer(1.0) // gift a dollar
			}
		default:
			if status == shared.Critical { // Case we are poor
				offers[team] = shared.GiftOffer(c.config.offertoDyingIslands / 2) // Half the amount if we are poor
			} else {
				offers[team] = shared.GiftOffer(0.5) // can we abuse the fact that they look at the amount of gifts?
			}
		}
		// History
		newGiftRequest := giftInfo{
			requested: receivedRequests[team], // Amount THEY requested
			offered:   offers[team],           // Amount WE offered
		}
		c.giftHistory[team].theirRequest[c.getTurn()] = newGiftRequest
	}
	return offers
}

// GetGiftResponses allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
// COMPULSORY, you need to implement this method
func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
	// receivedOffers := shared.GiftOfferDict{}  // For future use when actually considering peoples offers
	responses := shared.GiftResponseDict{}
	for team, offer := range receivedOffers { // For all the clients we look at the offers
		if offer >= 0 {
			responses[team] = shared.GiftResponse{
				AcceptedAmount: shared.Resources(offer), // Accept all they gave us
				Reason:         shared.Accept,           // Accept all gifts duh
			}
		} else {
			responses[team] = shared.GiftResponse{
				AcceptedAmount: 0,                      // Accept all they gave us
				Reason:         shared.DeclineDontNeed, // Accept all gifts duh
			}
		}
	}
	// History
	for _, team := range c.getAliveTeams(true) { // Could be in the loop above but
		newGiftRequest := giftInfo{ // for consistency with other functions its here
			requested: c.giftHistory[team].ourRequest[c.getTurn()].requested, // Amount We requested
			offered:   receivedOffers[team],                                  // Amount offered TO US
			response:  responses[team],                                       // Amount and reason WE accepted
		}
		c.giftHistory[team].ourRequest[c.getTurn()] = newGiftRequest
	}
	return responses
}

// =====================================Gifting history to be made==========================================

// UpdateGiftInfo will be called by the server with all the responses you received
// from the gift session. It is up to you to complete the transactions with the methods
// that you will implement yourself below. This allows for opinion formation.
// COMPULSORY, you need to implement this method

func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {
	for _, team := range c.getAliveTeams(true) {
		newGiftRequest := giftInfo{
			requested: c.giftHistory[team].theirRequest[c.getTurn()].requested, // Amount THEY requested
			offered:   c.giftHistory[team].theirRequest[c.getTurn()].offered,   // Amount WE offered them
			response:  receivedResponses[team],                                 // Amount THEY accepted and REASON
		}
		c.giftHistory[team].theirRequest[c.getTurn()] = newGiftRequest
	}
}

// ===================================== Has sending / recv gifts been implemented? ===============================

// DecideGiftAmount is executed at the end of each turn and asks clients how much
// they want to fulfill a gift offer they have made.
// COMPULSORY, you need to implement this method
func (c *client) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {
	var giftOff shared.Resources
	if c.resourceHistory[c.gameState().Turn-1] < c.gameState().ClientInfo.Resources { // if resources are higher that previous' rounds resources
		if c.wealth() >= wealthTier(c.config.middleThreshold) { //this is only fulfilled if we are wealthy enough Mid and JB
			if c.opinions[toTeam].getScore() > 0 && c.opinions[toTeam].getScore() <= 0.5 { //if twe are walthy (>=2) and our opinion on the island is between 0 and 0.5 then fulfill full offer
				giftOff = giftOffer
			} else if c.opinions[toTeam].getScore() > 0.5 && c.opinions[toTeam].getScore() <= 1 { //if we are wealthy (>=2) and we have a high opinion on the island, then boost the gift a little by 1.4
				giftOff = giftOffer * c.config.giftBoosting
			} else {
				giftOff = 0
			}
		} else if c.wealth() == wealthTier(c.config.imperialThreshold) { //this is only fulfilled if we are ICL students rich
			if c.opinions[toTeam].getScore() > 0 && c.opinions[toTeam].getScore() <= 0.5 { //if wealth is one but opinion is between 0 and 0.5 then give half the offerr
				giftOff = giftOffer * c.config.giftReduct
			} else if c.opinions[toTeam].getScore() > 0.5 && c.opinions[toTeam].getScore() <= 1 { //if wealth is 1 and opinion is 0.5 to 1 then give fulfill whole offer
				giftOff = giftOffer
			} else {
				giftOff = 0
			}
		} else { //Reject all offers if opinions are below zero and if wealth is below 1
			giftOff = 0
		}
	} else { //Reject all offers if we have less than we did last round
		giftOff = 0
	}
	// History
	newGiftRequest := giftInfo{ // 	For each client create a new gift info
		requested:      c.giftHistory[toTeam].theirRequest[c.getTurn()].requested, // Amount We requested
		offered:        c.giftHistory[toTeam].theirRequest[c.getTurn()].offered,   // Amount offered TO US
		response:       c.giftHistory[toTeam].theirRequest[c.getTurn()].response,  // Amount and reason WE accepted
		actualReceived: giftOffer,                                                 // Amount they ACTUALLY receive
	}
	c.giftHistory[toTeam].theirRequest[c.getTurn()] = newGiftRequest

	return giftOff
	// Debugging for gift
	// c.Logf("[Debug] ourRequest [%v]", c.giftHistory[toTeam].ourRequest[c.getTurn()])
	// c.Logf("[Debug] theirRequest [%v]", c.giftHistory[toTeam].theirRequest[c.getTurn()])
}

// ===================================== Has sending / recv gifts been implemented? ===============================
// SentGift is executed at the end of each turn and notifies clients that
// their gift was successfully sent, along with the offer details.
// COMPULSORY, you need to implement this method
func (c *client) SentGift(sent shared.Resources, to shared.ClientID) {
	// You can check your updated resources like this:
	// myResources := c.gameState().ClientInfo.Resources

	newGiftRequest := giftInfo{ // 	For each client create a new gift info
		requested:      c.giftHistory[to].theirRequest[c.getTurn()].requested, // Amount We requested
		offered:        c.giftHistory[to].theirRequest[c.getTurn()].offered,   // Amount offered TO US
		response:       c.giftHistory[to].theirRequest[c.getTurn()].response,  // Amount and reason WE accepted
		actualReceived: sent,                                                  // Amount they actually receive according to server
	}
	c.giftHistory[to].theirRequest[c.getTurn()] = newGiftRequest
}

// ReceivedGift is executed at the end of each turn and notifies clients that
// their gift was successfully received, along with the offer details.
// COMPULSORY, you need to implement this method
func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {
	// You can check your updated resources like this:
	// myResources := c.gameState().ClientInfo.Resources
	newGiftRequest := giftInfo{ // 	For each client create a new gift info
		requested:      c.giftHistory[from].ourRequest[c.getTurn()].requested, // Amount We requested
		offered:        c.giftHistory[from].ourRequest[c.getTurn()].offered,   // Amount offered TO US
		response:       c.giftHistory[from].ourRequest[c.getTurn()].response,  // Amount and reason WE accepted
		actualReceived: received,                                              // Amount they actually GAVE us
	}
	c.giftHistory[from].ourRequest[c.getTurn()] = newGiftRequest

	c.giftOpinions()
}

func (c *client) giftOpinions() {
	for team := range c.gameState().ClientLifeStatuses { // for each ID
		// ======================= Bad =======================
		// If we get OFFERED LESS than we Requested
		if shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].offered) <
			shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].requested) {
			c.opinions[team].updateOpinion(generalBasis, -0.05)
		}

		// If we ACTUALLY get LESS than they OFFERED us
		if shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].actualReceived) <
			shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].offered) {
			c.opinions[team].updateOpinion(generalBasis, -0.1)
		}

		// If they REQUEST the MOST compared to other islands
		highestRequest := shared.Team1
		if c.giftHistory[highestRequest].theirRequest[c.getTurn()].requested <
			c.giftHistory[team].theirRequest[c.getTurn()].requested {
			highestRequest = team
		}
		c.opinions[highestRequest].updateOpinion(generalBasis, -0.05)

		// ======================= Good =======================
		// If they GIVE MORE than OFFERED then increase it a bit (can be abused)
		if shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].actualReceived) >
			shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].offered) {
			c.opinions[team].updateOpinion(generalBasis, 0.025)
		}

		// If we RECEIVE MORE than WE REQUESTED and they OFFERED
		if shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].actualReceived) >
			shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].offered) &&
			shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].actualReceived) >
				shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].requested) {
			c.opinions[team].updateOpinion(generalBasis, 0.02)
		}

		// If they REQUEST the LEAST compared to other islands
		lowestRequest := shared.Team1
		if c.giftHistory[lowestRequest].theirRequest[c.getTurn()].requested >
			c.giftHistory[team].theirRequest[c.getTurn()].requested {
			lowestRequest = team
		}
		c.opinions[lowestRequest].updateOpinion(generalBasis, 0.05)
	}
}

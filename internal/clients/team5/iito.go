package team5

import (
	"math"

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
	if c.getTurn() > 1 {
		c.updateGiftOpinions()
	}

	requests := shared.GiftRequestDict{}
	for team, status := range c.gameState().ClientLifeStatuses {
		if status != shared.Critical { // THEY are not dying

			// They are nice people
			if c.opinions[team].getScore() > c.config.opinionThresholdRequest {
				switch c.wealth() { // Look at our wealth
				case dying: // Case we are Critical (dying)
					requests[team] = shared.GiftRequest(c.config.dyingGiftRequestAmount * 0.8) // Ask  for dying amount
				case imperialStudent:
					requests[team] = shared.GiftRequest( // Scale our request
						(c.config.imperialGiftRequestAmount /
							float64(c.opinions[team].getScore())) *
							c.config.opinionRequestMultiplier) // Scale down the request if we like them
				case middleClass:
					requests[team] = shared.GiftRequest(
						(c.config.middleGiftRequestAmount /
							float64(c.opinions[team].getScore())) *
							c.config.opinionRequestMultiplier) // Scale down the request if we like them
				case jeffBezos:
					requests[team] = shared.GiftRequest(0)
				}

				// They are trashy people
			} else if c.opinions[team].getScore() < (-c.config.opinionThresholdRequest) {
				switch c.wealth() {
				case dying: // Case we are Dying
					requests[team] = shared.GiftRequest(c.config.dyingGiftRequestAmount * 1.2)
				case imperialStudent: // Case we are Poor
					requests[team] = shared.GiftRequest(
						(c.config.imperialGiftRequestAmount *
							float64(-c.opinions[team].getScore())) /
							c.config.opinionRequestMultiplier) // Scale up the request if dont like them
				case middleClass:
					requests[team] = shared.GiftRequest(
						(c.config.middleGiftRequestAmount *
							float64(-c.opinions[team].getScore())) /
							c.config.opinionRequestMultiplier) // Scale up the request if dont like them
				case jeffBezos:
					requests[team] = shared.GiftRequest(0)
				}
			} else { // Normal Opinion
				switch c.wealth() {
				case dying: // We are DYING
					requests[team] = shared.GiftRequest(c.config.dyingGiftRequestAmount)
				case imperialStudent:
					requests[team] = shared.GiftRequest(c.config.imperialGiftRequestAmount) // No scale
				case middleClass:
					requests[team] = shared.GiftRequest(c.config.middleGiftRequestAmount) // No scale
				case jeffBezos:
					requests[team] = shared.GiftRequest(0)
				}
			}
		} else {
			requests[team] = shared.GiftRequest(0)
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
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}
	for team, status := range c.gameState().ClientLifeStatuses {
		status := shared.ClientLifeStatus(status)
		switch {
		// case we are RICH or Middle class
		case c.wealth() >= 2:
			opinionMulti := mapToRange(float64(c.opinions[team].getScore()), -1, 1, 0.2, 1) // Opinion = 0 then you get what we say, opinion = 1 get what they ask for
			amount := ((opinionMulti * float64(receivedRequests[team])) +                   // opinion = 1 then they get what they asked for
				((1 - opinionMulti) * c.config.offertoDyingIslands)) // opinion = -1 they get 0 of what they want and all of what we pay them
			if status == shared.Critical {
				if c.opinions[team].getScore() >= 0 {
					offers[team] = shared.GiftOffer(math.Min(
						0.10*float64(c.gameState().ClientInfo.Resources), // max is 10% of our worth
						amount))
				} else { // opinions less than 0
					offers[team] = shared.GiftOffer(math.Min(
						0.075*float64(c.gameState().ClientInfo.Resources), // Max is 7.5% of our worth
						amount))
				}
			} else { // THEY are NOT CRITICAL but we have money
				if c.opinions[team].getScore() >= 0 {
					offers[team] = shared.GiftOffer(math.Min(
						0.05*float64(c.gameState().ClientInfo.Resources), // max is 5% of our worth
						amount))
				} else { // opinions less than 0
					offers[team] = shared.GiftOffer(math.Min(
						0.025*float64(c.gameState().ClientInfo.Resources), // Max is 2.5% of our worth
						amount))
				}
			}
		// we are POOR af people
		default:
			opinionMulti := mapToRange(float64(c.opinions[team].getScore()), -1, 1, 0.1, 0.5) // Opinion = 0 then you get what we say, opinion = 1 get what they ask for
			amount := ((opinionMulti * float64(receivedRequests[team])) +                     // opinion = 1 then they get half of what they wanted
				((1 - opinionMulti) * c.config.offertoDyingIslands)) // opinion = -1 they get 0 of what they want and all of what we pay them
			if status == shared.Critical {
				if c.opinions[team].getScore() >= 0 {
					offers[team] = shared.GiftOffer(math.Min(
						0.05*float64(c.gameState().ClientInfo.Resources), // max is 5% of our worth
						amount))
				} else { // opinions less than 0
					offers[team] = shared.GiftOffer(math.Min(
						0.025*float64(c.gameState().ClientInfo.Resources), // Max is 2.5% of our worth
						amount))
				}
			} else { // THEY are NOT CRITICAL but we have no money
				if c.opinions[team].getScore() >= 0 {
					offers[team] = shared.GiftOffer(math.Min(
						0.025*float64(c.gameState().ClientInfo.Resources), // max is 5% of our worth
						amount))
				} else { // opinions less than 0
					offers[team] = 0
				}
			} // End We Poor They Rich
		} // End of Switch
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
	responses := shared.GiftResponseDict{}
	for team, offer := range receivedOffers { // For all the clients we look at the offers
		if offer > 0 {
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
		if receivedResponses[team].Reason >= 2 {
			c.opinions[team].updateOpinion(generalBasis, c.changeOpinion(-0.05))
		} // why did they decline our offer?

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
	giftOff := giftOffer
	// History
	newGiftRequest := giftInfo{ // 	For each client create a new gift info
		requested:      c.giftHistory[toTeam].theirRequest[c.getTurn()].requested, // Amount We requested
		offered:        c.giftHistory[toTeam].theirRequest[c.getTurn()].offered,   // Amount offered TO US
		response:       c.giftHistory[toTeam].theirRequest[c.getTurn()].response,  // Amount and reason WE accepted
		actualReceived: giftOff,                                                   // Amount they ACTUALLY receive
	}
	c.giftHistory[toTeam].theirRequest[c.getTurn()] = newGiftRequest

	return giftOffer

}

// ===================================== Has sending / recv gifts been implemented? ===============================
// SentGift is executed at the end of each turn and notifies clients that
func (c *client) SentGift(sent shared.Resources, to shared.ClientID) {
	newGiftRequest := giftInfo{ // 	For each client create a new gift info
		requested:      c.giftHistory[to].theirRequest[c.getTurn()].requested, // Amount We requested
		offered:        c.giftHistory[to].theirRequest[c.getTurn()].offered,   // Amount offered TO US
		response:       c.giftHistory[to].theirRequest[c.getTurn()].response,  // Amount and reason WE accepted
		actualReceived: sent,                                                  // Amount they actually receive according to server
	}
	c.giftHistory[to].theirRequest[c.getTurn()] = newGiftRequest

	c.Logf("[SentGift]: Sent to %v amount = %v", to, sent)
}

// ReceivedGift is executed at the end of each turn and notifies clients that
func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {
	newGiftRequest := giftInfo{ // 	For each client create a new gift info
		requested:      c.giftHistory[from].ourRequest[c.getTurn()].requested, // Amount We requested
		offered:        c.giftHistory[from].ourRequest[c.getTurn()].offered,   // Amount offered TO US
		response:       c.giftHistory[from].ourRequest[c.getTurn()].response,  // Amount and reason WE accepted
		actualReceived: received,                                              // Amount they actually GAVE us
	}
	c.giftHistory[from].ourRequest[c.getTurn()] = newGiftRequest
	c.Logf("[ReceivedGift]: Received from %v amount = %v", from, received)
}

func (c *client) updateGiftOpinions() {
	lastTurn := c.getTurn() - 1
	var highestRequest shared.ClientID
	var lowestRequest shared.ClientID
	for _, team := range c.getAliveTeams(false) { // for each ID
		// c.Logf("Opinion of %v BEFORE gifts = %v", team, c.opinions[team].getScore())
		// ======================= Bad =======================
		// If we get OFFERED LESS than we Requested
		if shared.Resources(c.giftHistory[team].ourRequest[lastTurn].offered) <
			shared.Resources(c.giftHistory[team].ourRequest[lastTurn].requested) {
			c.opinions[team].updateOpinion(generalBasis, c.changeOpinion(-0.025))
		}
		// If we ACTUALLY get LESS than they OFFERED us
		if shared.Resources(c.giftHistory[team].ourRequest[lastTurn].actualReceived) <
			shared.Resources(c.giftHistory[team].ourRequest[lastTurn].offered) {
			c.opinions[team].updateOpinion(generalBasis, c.changeOpinion(-0.1))
		}

		// If they REQUEST the MOST compared to other islands
		if c.giftHistory[highestRequest].theirRequest[lastTurn].requested <
			c.giftHistory[team].theirRequest[lastTurn].requested {
			highestRequest = team
		}

		// ======================= Good =======================
		// If they GIVE MORE than OFFERED then increase it a bit (can be abused)
		if shared.Resources(c.giftHistory[team].ourRequest[lastTurn].actualReceived) >=
			shared.Resources(c.giftHistory[team].ourRequest[lastTurn].offered) {
			c.opinions[team].updateOpinion(generalBasis, c.changeOpinion(0.01))
		}

		// If we RECEIVE MORE than WE REQUESTED
		if shared.Resources(c.giftHistory[team].ourRequest[lastTurn].offered) >
			shared.Resources(c.giftHistory[team].ourRequest[lastTurn].requested) &&
			shared.Resources(c.giftHistory[team].ourRequest[lastTurn].actualReceived) >
				shared.Resources(c.giftHistory[team].ourRequest[lastTurn].requested) {
			c.opinions[team].updateOpinion(generalBasis, c.changeOpinion(0.1))
		}

		// If they REQUEST the LEAST compared to other islands
		if c.giftHistory[lowestRequest].theirRequest[lastTurn].requested >
			c.giftHistory[team].theirRequest[lastTurn].requested {
			lowestRequest = team
		}
		// c.Logf("Opinion of %v AFTER gifts = %v", team, c.opinions[team].getScore())
	}
	c.opinions[highestRequest].updateOpinion(generalBasis, c.changeOpinion(-0.025))
	c.opinions[lowestRequest].updateOpinion(generalBasis, c.changeOpinion(0.025))
}

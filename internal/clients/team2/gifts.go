package team2

import (
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type IslandTrust struct {
	island shared.ClientID
	trust  int
}

//Game Method Implementation: altruist, fair sharer and free rider
//LINE 230: this function determines how generous we will be based off of current method
//LINE 33: methodConfig is called
//LINE 99: methodConfig is called

type IslandTrustList []IslandTrust

// Overwrite default sort implementation
func (p IslandTrustList) Len() int           { return len(p) }
func (p IslandTrustList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p IslandTrustList) Less(i, j int) bool { return p[i].trust < p[j].trust }

// GetGiftRequests allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	turn := c.gameState().Turn
	requests := shared.GiftRequestDict{}

	// check our critical and threshold - if either is off - request
	ourAgentCritical := shared.Critical == shared.ClientLifeStatus(1)
	requestAmount := determineAllocation(c) * methodConfGift(c)

	// confidence[island] * requestAmount until -> target
	if ourAgentCritical || requestAmount > 0 {
		target := 1.5 * requestAmount //TODO config: cushion
		var trustRank IslandTrustList

		for team, status := range c.gameState().ClientLifeStatuses {

			// get our confidence in the team
			// TODO: this should be whether or not that team is critical not us
			if status != shared.ClientLifeStatus(2) {
				islandConf := IslandTrust{
					island: team,
					trust:  c.confidence("Gifts", team),
				}
				trustRank = append(trustRank, islandConf)
			}

		}

		// keep a ranked list of the teams
		sort.Sort(trustRank)

		// request confidence (0-1)*amountNeeded to consecutive islands in rank
		// until amountRequested = (factor eg 1.5) * amountNeeded (to accommodate for some islands not giving us a gift)
		targetFulfilled := target
		for i := 0; i < len(trustRank); i++ {
			requestAmount := shared.Resources(trustRank[i].trust) / 100 * target
			requestedTo := trustRank[i].island
			requests[requestedTo] = shared.GiftRequest(requestAmount)

			// to keep track in our history
			newGiftRequest := GiftInfo{
				requested: shared.GiftRequest(requestAmount),
			}
			c.giftHist[requestedTo].OurRequest[turn] = newGiftRequest

			targetFulfilled -= requestAmount

			if targetFulfilled <= 0 {
				return requests
			}
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
	turn := c.gameState().Turn

	// if we are critical do not offer gifts-> there should be a way to see which other islands are critical
	// if we are not critical and another island is critical offer gift
	// do not offer more than proportion of total resources we have
	ourAgentCritical := shared.Critical == shared.ClientLifeStatus(1)

	// prioritize giving gifts to islands we trust (for now confidence)

	// Give no more than half of amount before we reach threshold
	maxToGive := (c.gameState().ClientInfo.Resources - c.agentThreshold()) / (1 / (methodConfGift(c) / 2))

	var trustRank IslandTrustList
	if !ourAgentCritical || maxToGive <= 0 {
		for team := range receivedRequests {
			// max would be 200
			// c.confidence("ReceivedRequests", team) should reflect the status of an island and int,float64 requests hist
			// c.confidence("GiftWeRequest", team) should reflect how they respond to our requests
			confidenceMetric := c.confidence("Gifts", team)
			otherTeamCritical := c.gameState().ClientLifeStatuses[team] == shared.Critical
			if confidenceMetric > 50 || otherTeamCritical {
				islandConf := IslandTrust{
					island: team,
					trust:  confidenceMetric,
				}
				trustRank = append(trustRank, islandConf)
			}
		}

		sort.Sort(trustRank)

		// TODO: need to factor in the size of a request in decisions - above we were just sorting them by confidence
		for i := 0; i < len(trustRank); i++ {
			offeredTo := trustRank[i].island
			offeredAmount := receivedRequests[offeredTo]

			if offeredAmount >= shared.GiftRequest(maxToGive) {
				offeredAmount = shared.GiftRequest(maxToGive)
				maxToGive = 0
			}

			offers[offeredTo] = shared.GiftOffer(offeredAmount)

			// to keep track in our history
			newGiftRequest := GiftInfo{
				requested: receivedRequests[offeredTo],
				gifted:    shared.GiftOffer(offeredAmount),
			}
			c.giftHist[offeredTo].IslandRequest[turn] = newGiftRequest

			maxToGive -= shared.Resources(offeredAmount)

			if maxToGive <= 0 {
				return offers
			}
		}
	}

	return offers
}

// GetGiftResponses allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
// COMPULSORY, you need to implement this method
// sanctions could be a reason to deny a gift offer
func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
	responses := shared.GiftResponseDict{}
	turn := c.gameState().Turn

	for client, offer := range receivedOffers {
		responses[client] = shared.GiftResponse{
			AcceptedAmount: shared.Resources(offer),
			Reason:         shared.Accept,
		}
		newGiftRequest := GiftInfo{
			// it could potentially crash if we receive a gift we didn't ask for... this entry would be a null pointer
			requested: c.giftHist[client].OurRequest[turn].requested,
			gifted:    shared.GiftOffer(responses[client].AcceptedAmount),
			reason:    shared.AcceptReason(responses[client].AcceptedAmount),
		}
		c.giftHist[client].OurRequest[turn] = newGiftRequest
	}

	return responses
}

// UpdateGiftInfo will be called by the server with all the responses you received
// from the gift session. It is up to you to complete the transactions with the methods
// that you will implement yourself below. This allows for opinion formation.
// COMPULSORY, you need to implement this method
func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {
	turn := c.gameState().Turn

	// we should update our opinion of something if they reject a gift
	// instead for now, we base our decisions/opinions on actions not words
	// so we disregard what people say they will do and only store what they actually do
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
	turn := c.gameState().Turn
	newGiftRequest := GiftInfo{
		requested: c.giftHist[to].IslandRequest[turn].requested,
		gifted:    shared.GiftOffer(sent),
	}
	c.giftHist[to].IslandRequest[turn] = newGiftRequest
	// because received gift is called first we call this here
	c.updateGiftConfidence(to)
}

// // ReceivedGift is executed at the end of each turn and notifies clients that
// // their gift was successfully received, along with the offer details.
// // COMPULSORY, you need to implement this method
func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {
	turn := c.gameState().Turn
	newGiftRequest := GiftInfo{
		// it could potentially crash if we receive a gift we didn't ask for... this entry would be a null pointer
		requested: c.giftHist[from].OurRequest[turn].requested,
		gifted:    shared.GiftOffer(received),
	}
	c.giftHist[from].OurRequest[turn] = newGiftRequest

}

func (c *client) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {
	// Give no more than half of amount before we reach threshold
	maxToGive := (c.gameState().ClientInfo.Resources - c.agentThreshold()) / 2
	if giftOffer <= maxToGive {
		return giftOffer
	}
	return shared.Resources(0.0)
}

func methodConfGift(c *client) float64 {
	var modeMult float64
	switch c.MethodOfPlay() {
	case 0:
		modeMult = 0.8 //we do not want to take from the pool as it is struggling, request the majority from gifts
	case 1:
		modeMult = 0.6 //pool is doing average, request 60% from gifts
	case 2:
		modeMult = 0.2 //there is plenty in the pool so we take from the pool rather than request from people
	}
	return modeMult
}

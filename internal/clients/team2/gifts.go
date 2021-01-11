package team2

import (
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type IslandTrust struct {
	island shared.ClientID
	trust  int
}

type IslandTrustList []IslandTrust

// Overwrite default sort implementation
func (p IslandTrustList) Len() int           { return len(p) }
func (p IslandTrustList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p IslandTrustList) Less(i, j int) bool { return p[i].trust < p[j].trust }

// TODO: we still seem to be giving negative gifts whatever that means lol

// GetGiftRequests allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	turn := c.gameState().Turn
	requests := shared.GiftRequestDict{}

	// check our critical and threshold - if either is off - request
	ourAgentCritical := c.criticalStatus()
	requestAmount := c.determineBaseCommonPoolRequest() * c.giftReliance()
	c.Logf("stupid amount: %v", requestAmount)
	c.Logf("island chcik status: %v", ourAgentCritical)

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
				c.Logf("doink requests : %v", requests)
				return requests
			}
		}
	}
	c.Logf("doink requests 2 : %v", requests)
	return requests
}

// GetGiftOffers allows clients to make offers in response to gift requests by other clients.
// It can offer multiple partial gifts.
// COMPULSORY, you need to implement this method. This placeholder implementation offers no gifts,
// unless another team is critical.
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}
	turn := c.gameState().Turn

	for island := range offers {
		offers[island] = 0
	}

	// if we are critical do not offer gifts-> there should be a way to see which other islands are critical
	// if we are not critical and another island is critical offer gift
	// do not offer more than proportion of total resources we have
	ourAgentCritical := c.criticalStatus()
	// ourAgentCritical := true

	// prioritize giving gifts to islands we trust (for now confidence)

	// Give no more than half of amount before we reach threshold
	maxToGive := (c.gameState().ClientInfo.Resources - c.agentThreshold()) * (c.giftReliance() / shared.Resources(2))
	c.Logf("MAXTOGIVE : %v", maxToGive)
	c.Logf("our agent critical: %v", ourAgentCritical)
	var trustRank IslandTrustList
	//checks if we are not critical or maxtogive is -ve
	if !ourAgentCritical && maxToGive >= 0 {
		c.Logf("WAKA WAKA ENtered LOOP")
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

		//TODO: Our trust rank is always empty
		if len(trustRank) != 0 {
			sort.Sort(trustRank)
		}

		c.Logf("TRUSTRANK : %v", trustRank)

		// TODO: need to factor in the size of a request in decisions - above we were just sorting them by confidence
		for i := 0; i < len(trustRank); i++ {
			offeredTo := trustRank[i].island
			offeredAmount := receivedRequests[offeredTo] * shared.GiftRequest(trustRank[i].trust/100)

			if offeredAmount >= shared.GiftRequest(maxToGive) {
				offeredAmount = shared.GiftRequest(maxToGive)
				maxToGive = 0
			}

			offers[offeredTo] = shared.GiftOffer(offeredAmount)
			c.Logf("Predatory Lending 1: %v", offers)

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
	c.Logf("Predatory Lending 2: %v", receivedRequests)
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
		weRequested := shared.GiftRequest(0)
		if val, ok := c.giftHist[client].OurRequest[turn]; ok {
			weRequested = val.requested
		}
		newGiftRequest := GiftInfo{
			// it could potentially crash if we receive a gift we didn't ask for... this entry would be a null pointer
			requested: weRequested,
			gifted:    shared.GiftOffer(responses[client].AcceptedAmount),
			reason:    shared.AcceptReason(responses[client].Reason),
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
			reason:    shared.AcceptReason(response.Reason),
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

// TODO: DOUBLE CHECK THE LOGIC ON GIFT AMOUNTS
func (c *client) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {
	// Give no more than half of amount before we reach threshold
	maxToGive := (c.gameState().ClientInfo.Resources - c.agentThreshold()) / (shared.Resources(1) / (c.giftReliance() / shared.Resources(2)))
	if maxToGive > 0 || giftOffer <= maxToGive {
		return giftOffer
	}
	return shared.Resources(0)
}

func (c *client) giftReliance() shared.Resources {
	var multiplier shared.Resources

	switch c.setAgentStrategy() {
	case Altruist:
		// Pool Level is LOW, rely 80% on gifts
		multiplier = 0.8
	case FairSharer:
		// Pool Level is OK, rely 60% on gifts
		multiplier = 0.6
	case Selfish:
		// Pool Level is HIGH, rely 20% on gifts
		multiplier = 0.2
	}

	return multiplier
}

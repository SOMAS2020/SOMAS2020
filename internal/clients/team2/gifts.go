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
	ourAgentCritical := c.criticalStatus()
	giftTarget := c.giftReliance() * (c.taxAmount + c.gameConfig().CostOfLiving + c.gameConfig().MinimumResourceThreshold)

	// c.Logf("TARGETGIFT : %v", giftTarget)

	if ourAgentCritical || giftTarget != 0 {
		var trustRank IslandTrustList

		for _, team := range c.getAliveClients() {
			if team != c.GetID() {
				islandConf := IslandTrust{
					island: team,
					trust:  c.confidence("Gifts", team),
				}

				trustRank = append(trustRank, islandConf)
			}
		}

		// keep a ranked list of the teams
		if len(trustRank) != 0 {
			sort.Sort(trustRank)
		}

		targetToFulfill := giftTarget

		for i := 0; i < len(trustRank); i++ {
			requestAmount := shared.Resources(trustRank[i].trust) / shared.Resources(100) * giftTarget
			requestedTo := trustRank[i].island
			requests[requestedTo] = shared.GiftRequest(requestAmount)

			// To keep track in our history (for each island)
			newGiftRequest := GiftInfo{
				requested: requests[requestedTo],
				gifted:    0,             // to be changed later (when we receive responses)
				reason:    shared.Accept, // to be changed later (when we receive responses)
			}

			c.giftHist[requestedTo].OurRequest[turn] = newGiftRequest

			targetToFulfill -= requestAmount

			if targetToFulfill <= 0 {
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

	for island := range offers {
		offers[island] = 0
	}
	ourAgentCritical := c.criticalStatus()
	excess := c.gameState().ClientInfo.Resources - c.taxAmount + c.gameConfig().CostOfLiving + c.gameConfig().MinimumResourceThreshold
	maxToGive := excess * c.config.MaxGiftOffersMultiplier

	if maxToGive > 50 {
		maxToGive = 50
	}

	// prioritize giving gifts to islands we trust (for now confidence)
	var trustRank IslandTrustList
	if !ourAgentCritical && maxToGive > 0 {
		for team := range receivedRequests {
			teamLifeStatus := c.gameState().ClientLifeStatuses[team]
			islandConf := c.confidence("Gifts", team)
			if islandConf > 50 || teamLifeStatus == shared.Critical {
				islandConf := IslandTrust{
					island: team,
					trust:  islandConf,
				}
				trustRank = append(trustRank, islandConf)
			}
		}

		if len(trustRank) != 0 {
			sort.Sort(trustRank)
		}

		for i := 0; i < len(trustRank); i++ {
			offeredTo := trustRank[i].island
			offeredAmount := receivedRequests[offeredTo] * (shared.GiftRequest(trustRank[i].trust) / 100)

			//returning the minimum of costpfliving, offeredamount and maxtogive
			//costofliving is just 10, maybe increase a bit
			if offeredAmount >= shared.GiftRequest(c.gameConfig().CostOfLiving) && maxToGive >= c.gameConfig().CostOfLiving {
				offeredAmount = shared.GiftRequest(c.gameConfig().CostOfLiving)
				maxToGive -= c.gameConfig().CostOfLiving
			} else if offeredAmount <= shared.GiftRequest(c.gameConfig().CostOfLiving) && shared.GiftRequest(maxToGive) >= offeredAmount {
				maxToGive -= shared.Resources(offeredAmount)
			} else if offeredAmount >= shared.GiftRequest(maxToGive) {
				offeredAmount = shared.GiftRequest(maxToGive)
				maxToGive = 0
			}

			offers[offeredTo] = shared.GiftOffer(offeredAmount)

			// to keep track in our history
			newGiftRequest := GiftInfo{
				requested: receivedRequests[offeredTo],
				gifted:    shared.GiftOffer(offeredAmount),
				reason:    shared.Accept,
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
	for island, response := range receivedResponses {
		c.Logf("Island[%v] accepted [%v] Resources ", island, response.AcceptedAmount)
	}
}

// SentGift is executed at the end of each turn and notifies clients that
// their gift was successfully sent, along with the offer details.
// COMPULSORY, you need to implement this method
func (c *client) SentGift(sent shared.Resources, to shared.ClientID) {
	turn := c.gameState().Turn
	theyRequested := shared.GiftRequest(0)

	if val, ok := c.giftHist[to].IslandRequest[turn]; ok {
		theyRequested = val.requested
	}

	newGiftRequest := GiftInfo{
		requested: theyRequested,
		gifted:    shared.GiftOffer(sent),
		reason:    shared.Accept,
	}

	c.giftHist[to].IslandRequest[turn] = newGiftRequest
	// because received gift is called first we call this here
	c.updateGiftConfidence(to)
}

// ReceivedGift is executed at the end of each turn and notifies clients that
// their gift was successfully received, along with the offer details.
// COMPULSORY, you need to implement this method
func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {
	turn := c.gameState().Turn
	weRequested := shared.GiftRequest(0)

	if val, ok := c.giftHist[from].OurRequest[turn]; ok {
		weRequested = val.requested
	}

	newGiftRequest := GiftInfo{
		requested: weRequested,
		gifted:    shared.GiftOffer(received),
		reason:    shared.Accept,
	}

	c.giftHist[from].OurRequest[turn] = newGiftRequest
}

// If we have less resources to give than expected (after actually using our resources this turn)
// we can decide not to give the previously agreed upon amount
func (c *client) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {
	updatedExcess := c.gameState().ClientInfo.Resources - c.taxAmount + c.gameConfig().CostOfLiving + c.gameConfig().MinimumResourceThreshold
	updatedMaxToGive := updatedExcess * c.config.MaxGiftOffersMultiplier

	if updatedMaxToGive > 0 || giftOffer <= updatedMaxToGive {
		return giftOffer
	}

	return shared.Resources(0)
}

func (c *client) giftReliance() shared.Resources {
	var multiplier shared.Resources

	switch c.getAgentStrategy() {
	case Altruist:
		// Pool Level is LOW, rely 80% on gifts
		multiplier = c.config.AltruistMultiplier
	case FairSharer:
		// Pool Level is OK, rely 60% on gifts
		multiplier = c.config.FairSharerMultipler
	case Selfish:
		// Pool Level is HIGH, rely 20% on gifts
		multiplier = c.config.FreeRiderMultipler
	}

	return multiplier
}

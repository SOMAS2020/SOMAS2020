package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) GetGiftRequests() shared.GiftRequestDict {
	requests := shared.GiftRequestDict{}

	ourStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	ourPersonality := c.getPersonality()

	minThreshold := c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold
	costOfLiving := c.ServerReadHandle.GetGameConfig().CostOfLiving
	criticalCounter := c.ServerReadHandle.GetGameState().ClientInfo.CriticalConsecutiveTurnsCounter
	maxCriticalCounter := c.ServerReadHandle.GetGameConfig().MaxCriticalConsecutiveTurns

	friendshipCoffesOnRequest := make(map[shared.ClientID]shared.GiftRequest)

	for team, fsc := range c.getFriendshipCoeffs() {
		friendshipCoffesOnRequest[team] = shared.GiftRequest(fsc + float64(1))
	}

	// requests gifts from all islands if we are in critical status
	for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		if team == c.GetID() {
			continue
		}

		if ourStatus == shared.Critical {
			if criticalCounter == maxCriticalCounter-uint(1) {
				// EMERGENCY!! will try to get a minimum number so that we can survive
				requests[team] = shared.GiftRequest(minThreshold)
			} else {
				// not emergency, will be more greedy, but we are still in critical status
				requests[team] = shared.GiftRequest(minThreshold + costOfLiving)
			}
		} else if status != shared.Alive {
			// won't ask a gift from a critical or dead island
			continue
		} else {
			// TODO: do we need to request gifts when we are not critical anyway?
			if ourPersonality == Selfish {
				// asks more when we are selfish
				requests[team] = shared.GiftRequest(costOfLiving) * friendshipCoffesOnRequest[team]
			} else if ourPersonality == Normal {
				// asks a regular amount
				requests[team] = shared.GiftRequest(costOfLiving)
			} else if ourPersonality == Generous {
				// asks for a cost of living
				requests[team] = shared.GiftRequest(costOfLiving) / friendshipCoffesOnRequest[team]
			}
		}

		// records the requests into our history book
		if _, found := c.giftsRequestedHistory[team]; found {
			c.giftsRequestedHistory[team] += shared.Resources(requests[team])
		} else {
			c.giftsRequestedHistory[team] = shared.Resources(requests[team])
		}
	}

	return requests
}

func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}

	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	ourStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	ourPersonality := c.getPersonality()

	minThreshold := c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold

	friendshipCoffesOnOffer := make(map[shared.ClientID]shared.GiftOffer)

	for team, fsc := range c.getFriendshipCoeffs() {
		friendshipCoffesOnOffer[team] = shared.GiftOffer(fsc)
	}

	if ourStatus == shared.Critical {
		// rejects all offers if we are in critical status or will be
		return offers
	}

	for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		amountOffer := shared.GiftOffer(0.0)

		if team == c.GetID() {
			continue
		}

		if status == shared.Critical {
			// offers a minimum resources to all islands which are in critical status
			offers[team] = shared.GiftOffer(minThreshold)
		} else if amountRequested, found := receivedRequests[team]; found {
			// otherwise gifts will be offered based on our friendship and the amount
			// TODO: offers gifts to the best friends first
			if ourResources-shared.Resources(amountRequested) <= minThreshold {
				// cannot offer the gift if it let us fall into the minimum resources threshold
				break
			} else if shared.Resources(amountRequested) > ourResources/shared.Resources(c.getNumOfAliveIslands()) {
				amountOffer = shared.GiftOffer(ourResources / shared.Resources(c.getNumOfAliveIslands()))
			} else {
				amountOffer = shared.GiftOffer(amountRequested)
			}

			if ourPersonality == Selfish {
				// introduces high penalty on the gift offers if we are selfish
				// yes, we are very stingy in this case ;)
				offers[team] = amountOffer * friendshipCoffesOnOffer[team]
			} else if ourPersonality == Normal {
				// introduces normal penalty
				offers[team] = amountOffer
			}
		}
	}

	if c.ServerReadHandle.GetGameState().Turn == 1 {
		for _, team := range shared.TeamIDs[:] {
			offers[team] = shared.GiftOffer(1)
		}
	}

	if ourPersonality == Generous {
		// intorduces no penalty - we are rich!
		for _, team := range shared.TeamIDs[:] {
			offers[team] = shared.GiftOffer(c.ServerReadHandle.GetGameConfig().CostOfLiving)
		}
	}

	return offers
}

func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
	responses := shared.GiftResponseDict{}

	costOfLiving := c.ServerReadHandle.GetGameConfig().CostOfLiving
	ourStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus

	for client, offer := range receivedOffers {
		if ourStatus == shared.Critical {
			responses[client] = shared.GiftResponse{
				AcceptedAmount: shared.Resources(offer),
				Reason:         shared.Accept,
			}
		} else if c.ServerReadHandle.GetGameState().ClientLifeStatuses[client] == shared.Critical {
			responses[client] = shared.GiftResponse{
				AcceptedAmount: 0,
				Reason:         shared.DeclineDontNeed,
			}
		} else if c.friendship[client] == c.clientConfig.minFriendship {
			// TODO: is this stupid?
			responses[client] = shared.GiftResponse{
				AcceptedAmount: shared.Resources(offer),
				Reason:         shared.Accept,
			}
		} else if c.friendship[client] == c.clientConfig.maxFriendship {
			// prevents our friend island from running out of resources
			responses[client] = shared.GiftResponse{
				AcceptedAmount: shared.Resources(offer) - costOfLiving,
				Reason:         shared.Accept,
			}
		} else {
			responses[client] = shared.GiftResponse{
				AcceptedAmount: shared.Resources(offer),
				Reason:         shared.Accept,
			}
		}
	}

	return responses
}

// TODO: anything else for this function?
func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {
	for client, response := range receivedResponses {
		c.Logf("The island[%v] accepts the amount[%v] because they [%v]", client, response.AcceptedAmount, response.Reason)
	}
}

func (c *client) SentGift(sent shared.Resources, to shared.ClientID) {
	if _, found := c.giftsSentHistory[to]; found {
		c.giftsSentHistory[to] += sent
	} else {
		c.giftsSentHistory[to] = sent
	}

	c.updateFriendship(-sent, to)
}

func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {
	if _, found := c.giftsReceivedHistory[from]; found {
		c.giftsReceivedHistory[from] += received
	} else {
		c.giftsReceivedHistory[from] = received
	}

	c.updateFriendship(received, from)
}

func (c *client) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {
	giftAmount := giftOffer

	if c.ServerReadHandle.GetGameState().ClientLifeStatuses[toTeam] == shared.Critical {
		if giftOffer < c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold {
			giftAmount = c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold
		}
	}

	return giftAmount
}

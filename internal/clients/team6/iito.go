package team6

import (

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared")



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

func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	ourStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	ourPersonality := c.getPersonality()
	friendshipCoffesOnOffer := make(map[shared.ClientID]shared.GiftOffer)

	for team, fsc := range c.getFriendshipCoeffs() {
		friendshipCoffesOnOffer[team] = shared.GiftOffer(fsc)
	}

	if ourStatus == shared.Critical {
		// rejects all offers if we are in critical status
		return offers
	}

	for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		amountOffer := shared.GiftOffer(0.0)

		if status == shared.Critical {
			// offers a minimum resources to all islands which are in critical status
			// TODO: what is the best amount to offer???
			offers[team] = shared.GiftOffer(10.0)
		} else if amountRequested, found := receivedRequests[team]; found {
			// otherwise gifts will be offered based on our friendship and the amount
			if shared.Resources(amountRequested) > ourResources/6.0 {
				amountOffer = shared.GiftOffer(ourResources / 6.0)
			} else {
				amountOffer = shared.GiftOffer(amountRequested)
			}

			if ourPersonality == Selfish {
				// introduces high penalty on the gift offers if we are selfish
				// yes, we are very stingy in this case ;)
				offers[team] = shared.GiftOffer(friendshipCoffesOnOffer[team] * friendshipCoffesOnOffer[team] * amountOffer)
			} else if ourPersonality == Normal {
				// introduces normal penalty
				offers[team] = shared.GiftOffer(friendshipCoffesOnOffer[team] * amountOffer)
			} else {
				// intorduces no penalty - we are rich!
				offers[team] = shared.GiftOffer(amountOffer)
			}

		}
	}

	// TODO: implement different wealth state for our island so we can decided
	// whether to be generous or not --finished
	return offers

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

func (c *client) UpdateGiftInfo(ReceivedResponses shared.GiftResponseDict) {}

func (c *client) SentGift(sent shared.Resources, to shared.ClientID) {
	//myResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
}

func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {
	//myResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
}

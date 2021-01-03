package team1

import (
	"math"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Helpers

// optionOnTeamsDict contains the opinion of Team1 about other teams.
// 0 is neutral, Positive -> Positive Opinion, Negative -> Negative Opinion
type opinionOnTeams struct {
	clientID shared.ClientID
	opinion  int
}

type sortByOpinion []opinionOnTeams

func (a sortByOpinion) Len() int           { return len(a) }
func (a sortByOpinion) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByOpinion) Less(i, j int) bool { return a[i].opinion > a[j].opinion }

// GetGiftRequests allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	requests := shared.GiftRequestDict{}
	// TODO: Add a flag so that agent can just request gifts all the time.
	if c.emotionalState() == Desperate {
		for clientID, status := range c.gameState().ClientLifeStatuses {
			if status != shared.Dead && clientID != c.GetID() {
				// TODO: Probably best to request a portion of Living Cost + Tax?
				requests[id] = shared.GiftRequest(25.0)
			}
		}
	}
	return requests
}

// GetGiftOffers allows clients to make offers in response to gift requests by other clients.
// It can offer multiple partial gifts.
// This will first iterate through a sorted array of opinionTeams starting from most favoured -> least favoured.
// + If the opinion of that team is large enough, we just give resources and don't question it.
// + If the opinion of that team is too low, we don't give resources :)
// + If the opinion of that team is inbetween, we give resources if they are critical.
// 					''						 and our resources is large enough (2000), we give resources freely
// + Else, ignore.
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}
	resourcesAvailable := c.gameState().ClientInfo.Resources
	teamStatus := c.gameState().ClientLifeStatuses
	sort.Sort(sortByOpinion(c.opinionTeams))
	for _, teams := range c.opinionTeams {
		for id, request := range receivedRequests {
			if teams.clientID == id && resourcesAvailable > c.config.anxietyThreshold && resourcesAvailable > shared.Resources(request) {
				// TODO: Fix arbitary value to a flag.
				if teams.opinion > 10 {
					offers[id] = shared.GiftOffer(request)
					resourcesAvailable -= shared.Resources(request)
				} else if teams.opinion < -10 {
					// Skip the giftOffer. We don't like them >:)
					continue
				} else if teamStatus[id] == shared.Critical {
					offers[id] = shared.GiftOffer(request)
					resourcesAvailable -= shared.Resources(request)
				} else {
					offerResource := giveLeftoverResources(resourcesAvailable, c.config.anxietyThreshold, shared.Resources(request))
					if offerResource != -1 {
						offers[id] = shared.GiftOffer(offerResource)
						resourcesAvailable -= offerResource
					}
				}
			}
		}
	}
	return offers
}

func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
	responses := shared.GiftResponseDict{}
	for id, amount := range receivedOffers {
		responses[id] = shared.GiftResponse{
			Reason:         shared.Accept,
			AcceptedAmount: shared.Resources(amount),
		}
	}
	return responses
}

// giveLeftoverResources finds the ratio between available resources and anxiety threshold. Using that ratio, the agent decides the max
// amount that it is willing to give away.
// Return -1 if you don't want to give any resources away.
func giveLeftoverResources(resourcesAvailable shared.Resources, anxietyThreshold shared.Resources, request shared.Resources) shared.Resources {
	if resourcesAvailable-request < anxietyThreshold {
		return -1
	}

	anxietyRatio := float64(resourcesAvailable / anxietyThreshold)
	giveResource := math.Min(float64(request), (float64(resourcesAvailable) * (anxietyRatio / 100)))
	return shared.Resources(giveResource)
}

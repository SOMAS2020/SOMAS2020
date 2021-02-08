package team1

import (
	"math"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

/**************************/
/*** 		Types	 	***/
/**************************/

// optionOnTeam contains the opinion of Team1 about another team.
// 0 is neutral, Positive -> Positive Opinion, Negative -> Negative Opinion
type opinionOnTeam struct {
	clientID shared.ClientID
	opinion  Opinion
}

type sortByOpinion []opinionOnTeam

// implemenent sort.Interface
func (a sortByOpinion) Len() int           { return len(a) }
func (a sortByOpinion) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByOpinion) Less(i, j int) bool { return a[i].opinion > a[j].opinion }

/**************************/
/*** 		Helpers	 	***/
/**************************/

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

/**************************/
/*** 		IITO	 	***/
/**************************/

// GetGiftRequests allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	requests := shared.GiftRequestDict{}
	// TODO: Malicious !! Add a flag so that agent can just request gifts all the time.
	if c.emotionalState() == Desperate {
		for clientID, status := range c.gameState().ClientLifeStatuses {
			if status != shared.Dead && clientID != c.GetID() {
				// TODO: Probably best to request a portion of Living Cost + Tax?
				requests[clientID] = shared.GiftRequest(2 * c.gameConfig().CostOfLiving)
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
// + If the opinion of that team is in between, we give resources if they are critical.
// 					''						 and our resources is large enough, we give a the maximum
//   amount of resources freely depending on our anxiety.
// + Else, ignore.
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}
	resourcesAvailable := c.gameState().ClientInfo.Resources
	teamStatus := c.gameState().ClientLifeStatuses

	// Sort so that we go through those we like first.
	teamOpinionList := []opinionOnTeam{}
	for id, opinion := range c.teamOpinions {
		teamOpinionList = append(teamOpinionList, opinionOnTeam{
			clientID: id,
			opinion: opinion,
		})
	}
	sort.Sort(sortByOpinion(teamOpinionList))

	for _, teamOpinion := range teamOpinionList {
		teamID := teamOpinion.clientID
		opinion := teamOpinion.opinion
		request, teamMadeRequest := receivedRequests[teamID]

		switch {
		case !teamMadeRequest:
		case resourcesAvailable <= c.config.anxietyThreshold:
		case resourcesAvailable <= shared.Resources(request):
		case opinion < -c.config.maxOpinion:
			// Do not make an offer
		case opinion > c.config.maxOpinion:
			offers[teamID] = shared.GiftOffer(request)
			resourcesAvailable -= shared.Resources(request)
		case teamStatus[teamID] == shared.Critical:
			offers[teamID] = shared.GiftOffer(request)
			resourcesAvailable -= shared.Resources(request)
		default:
			offerResource := giveLeftoverResources(resourcesAvailable, c.config.anxietyThreshold, shared.Resources(request))
			if offerResource != -1 {
				offers[teamID] = shared.GiftOffer(offerResource)
				resourcesAvailable -= offerResource
			}
		}
	}
	return offers
}

// GetGiftResponses allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
// Accept all GiftResponses since why not?
func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
	responses := shared.GiftResponseDict{}
	for id, amount := range receivedOffers {
		responses[id] = shared.GiftResponse{
			Reason:         shared.Accept,
			AcceptedAmount: shared.Resources(amount),
		}
		c.receivedOffer[id] = shared.Resources(amount)
	}
	return responses
}

// UpdateGiftInfo will be called by the server with all the responses you received
// from the gift session. It is up to you to complete the transactions with the methods
// that you will implement yourself below. This allows for opinion formation.
func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {
	for id, response := range receivedResponses {
		if response.Reason == shared.DeclineDontLikeYou {
			c.teamOpinions[id]--
		}
	}
}

// SentGift is executed at the end of each turn and notifies clients that
// their gift was successfully sent, along with the offer details.
func (c *client) SentGift(sent shared.Resources, to shared.ClientID) {
	// TODO: Do we actually need this?
}

// ReceivedGift is executed at the end of each turn and notifies clients that
// their gift was successfully received, along with the offer details.
func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {
	if received > shared.Resources(c.receivedOffer[from]) {
		// We love them cause they gave more than they promised.
		c.teamOpinions[from] += c.config.maxOpinion / 5
	} else if received > 0 {
		c.teamOpinions[from]++
	} else if received <= 0 {
		c.teamOpinions[from]--
	}
}

// DecideGiftAmount is executed at the end of each turn and asks clients how much
// they want to fulfill a gift offer they have made.
// Very similar to GetGiftOffers(). Only difference is the in between opinionated teams.
// Resources given will be different to the promised offers since this depends on how
// much resources we currently have.
func (c *client) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {
	resourcesAvailable := c.gameState().ClientInfo.Resources
	teamStatus := c.gameState().ClientLifeStatuses

	switch {
	case resourcesAvailable <= c.config.anxietyThreshold:
		return 0
	case resourcesAvailable <= giftOffer:
		return 0
	case c.teamOpinions[toTeam] < -c.config.maxOpinion:
		// Skip the giftOffer. We don't like them >:)
		return 0
	case c.teamOpinions[toTeam] > c.config.maxOpinion:
		return giftOffer
	case teamStatus[toTeam] == shared.Critical:
		// We are trying to be nice.
		return giftOffer
	default:
		offerResource := giveLeftoverResources(resourcesAvailable, c.config.anxietyThreshold, giftOffer)
		if offerResource != -1 {
			return offerResource
		}
		return shared.Resources(0)
	}
}

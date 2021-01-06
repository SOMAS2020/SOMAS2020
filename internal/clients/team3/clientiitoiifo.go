package team3

import (
	"math"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// GetGiftOffers allows clients to make offers in response to gift requests by other clients.
// It can offer multiple partial gifts.
// Strategy: We cover the risk that we lose money from the islands that we don’t trust with
// what we get from the islands that we do trust. Also, we don't request any gifts from critical islands.
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	var totalRequestAmt float64

	requests := shared.GiftRequestDict{}

	localPool := c.getLocalResources()

	resourcesNeeded := c.params.localPoolThreshold - float64(localPool)
	//fmt.Println("resources needed: ", resourcesNeeded)
	if resourcesNeeded > 0 {
		resourcesNeeded *= (1 + c.params.giftInflationPercentage)
		totalRequestAmt = resourcesNeeded
	} else {
		totalRequestAmt = c.params.giftInflationPercentage * c.params.localPoolThreshold
	}
	//fmt.Println("total request amount: ", totalRequestAmt)

	avgRequestAmt := totalRequestAmt / float64(c.getIslandsAlive()-c.getIslandsCritical())

	for island, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		if island == shared.Team3 {
			continue
		}
		if status == shared.Critical || status == shared.Dead {
			requests[island] = shared.GiftRequest(0.0)
		} else {
			var requestAmt float64
			requestAmt = avgRequestAmt * math.Pow(c.trustScore[island], c.params.trustParameter) * c.params.trustConstantAdjustor
			requests[island] = shared.GiftRequest(requestAmt)
		}
	}

	c.requestedGiftAmounts = requests
	return requests
}

func findAvgExclMinMax(Requests shared.GiftRequestDict) shared.GiftRequest {
	var sum shared.GiftRequest
	minClient := shared.TeamIDs[0]
	maxClient := shared.TeamIDs[0]

	// Find min and max requests
	for island, request := range Requests {
		if request < Requests[minClient] {
			minClient = island
		}
		if request > Requests[maxClient] {
			maxClient = island
		}
	}

	// Compute average ignoring highest and lowest
	for island, request := range Requests {
		if island != minClient || island != maxClient {
			sum += request
		}
	}

	return shared.GiftRequest(float64(sum) / float64(len(shared.TeamIDs)))
}

// sigmoidAndNormalise returns the normalised number between 0 - 1 based on the
// trust score between 0 - 100.
func (c *client) sigmoidAndNormalise(island shared.ClientID) shared.GiftOffer {

	sigmoid := (1 / (1 + math.Exp(-0.1*(c.trustScore[island]-50))))
	min := (1 / (1 + math.Exp(-0.1*-50)))
	max := (1 / (1 + math.Exp(-0.1*50)))

	normalised := (sigmoid - min) / (max - min)

	return shared.GiftOffer(normalised)
}

// GetGiftResponses allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}

	islandStatusCritical := c.isClientStatusCritical(id)

	if islandStatusCritical {
		for _, island := range c.getAliveIslands() {
			offers[island] = 0.0
		}
		return offers
	}

	var amounts = map[shared.ClientID]shared.GiftOffer{}
	var totalRequestedAmt float64
	var sumRequest shared.GiftRequest

	for island, request := range receivedRequests {
		sumRequest += request
		amounts[island] = c.sigmoidAndNormalise(island) * shared.GiftOffer(request)
	}

	//fmt.Println("original amounts: ", amounts)

	for _, island := range c.getAliveIslands() {
		if island != id && amounts[island] == 0.0 {
			amounts[island] = shared.GiftOffer(c.trustScore[island] * c.params.NoRequestGiftParam)
		}
	}

	//fmt.Println("length of amounts map: ", len(amounts))

	for _, requests := range c.requestedGiftAmounts {
		totalRequestedAmt += float64(requests)
	}

	localPool := c.getLocalResources()
	giftBudget := shared.GiftOffer((float64(localPool) + totalRequestedAmt) * ((1 - c.params.selfishness) / 2))

	rankedIslands := make([]shared.ClientID, 0.0, len(c.trustScore))

	for island := range c.trustScore {
		rankedIslands = append(rankedIslands, island)
	}

	sort.Slice(rankedIslands, func(i, j int) bool {
		return c.trustScore[rankedIslands[i]] > c.trustScore[rankedIslands[j]]
	})

	for _, island := range rankedIslands {
		if giftBudget >= amounts[island] {
			giftBudget -= amounts[island]
			offers[island] = amounts[island]
		} else if giftBudget == 0.0 {
			offers[island] = 0.0
		} else {
			offers[island] = giftBudget
			giftBudget = 0.0
		}
	}

	return offers
}

// GetGiftResponses returns the result of our island accepting/rejecting offered amounts.
// Strategy: we accept all amounts except if the offering island(s) is/are critical, then we
// do not take any of their offered amount(s).
func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
	responses := shared.GiftResponseDict{}
	for client, offer := range receivedOffers {
		responses[client] = shared.GiftResponse{
			AcceptedAmount: shared.Resources(offer),
			Reason:         shared.Accept,
		}
		if c.isClientStatusCritical(client) {
			responses[client] = shared.GiftResponse{
				AcceptedAmount: shared.Resources(0),
				Reason:         shared.DeclineDontNeed,
			}
		}
	}
	return responses
}

// UpdateGiftInfo receives responses from each island from the offers that our island made earlier.
// Strategy: The strategy is to update giftOpinion map on the reasons that gifts are declined or
// ignored for.
func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {
	for clientID, response := range receivedResponses {
		if response.Reason == shared.DeclineDontLikeYou {
			c.giftOpinions[clientID] -= 2
		} else if response.Reason == shared.Ignored {
			c.giftOpinions[clientID]--
		}
	}
}

// SentGift is executed at the end of each turn and notifies clients that
// their gift was successfully sent, along with offer details.
// We store the gift history for the previous turn.
func (c *client) SentGift(sent shared.Resources, to shared.ClientID) {
	c.sentGiftHistory = map[shared.ClientID]shared.Resources{}
	for island, amount := range c.sentGiftHistory {
		c.sentGiftHistory[island] = amount
	}

}

// ReceivedGift will inform us how much amount we have received from the specific islands
// and trust scores are then incremented and decremented based on the received difference.
func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {

	requestedFromIsland := c.requestedGiftAmounts[from]
	trustAdjustor := received - shared.Resources(requestedFromIsland)
	newTrustScore := 10 + (trustAdjustor * 0.2)
	if trustAdjustor >= 0 {
		c.updatetrustMapAgg(from, float64(newTrustScore))
	} else {
		if c.isClientStatusCritical(from) {
			return
		}
		c.updatetrustMapAgg(from, -float64(newTrustScore))
	}
}

// DecideGiftAmount is executed at the end of each turn and asks clients how much
// they wish to fulfill a gift offer they have previously made.
func (c *client) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {
	return giftOffer
}

package team4

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

// GetGiftRequests allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
// COMPULSORY, you need to implement this method
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	requests := shared.GiftRequestDict{}
	ourResources := c.getOurResources()
	importance := c.importances.getGiftRequestsImportance

	parameters := mat.NewVecDense(4, []float64{
		c.internalParam.greediness,
		c.internalParam.selfishness,
		c.internalParam.fairness,
		c.internalParam.collaboration,
	})
	// greedinessLevel designed to be greater than 1 for high selfishness/greediness
	//or low collaboration/fairness.
	greedinessLevel := mat.Dot(importance, parameters)

	wealthGoal := shared.Resources(greedinessLevel) * ourResources
	giftSize := wealthGoal - ourResources

	// You can fetch the clients which are alive like this:
	for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {

		if status == shared.Alive && wealthGoal > ourResources {
			requests[team] = shared.GiftRequest(giftSize)
		} else if status == shared.Critical {
			requests[team] = shared.GiftRequest(c.getSafeResourceLevel() * 2)
		} else {
			requests[team] = shared.GiftRequest(0.0)
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

	ourResources := c.getOurResources()
	// same as for requests. Now we will consider when you are more collaborative than greedy.
	importance := c.importances.getGiftRequestsImportance

	parameters := mat.NewVecDense(4, []float64{
		c.internalParam.greediness,
		c.internalParam.selfishness,
		c.internalParam.fairness,
		c.internalParam.collaboration,
	})
	// selflessLevel is the opposite of greedinessLevel in GetGiftRequests function.
	selflessLevel := mat.Dot(importance, parameters)

	// Initialising thresholds used.
	trustThreshold := 0.5
	resourcesThreshold := 100.0

	// have a cap at 100 in case we are too selfless
	wealthGoal := math.Max(selflessLevel*float64(ourResources), resourcesThreshold)

	// giftSize is the maximum we can afford giving away.
	giftSize := float64(ourResources) - wealthGoal

	// You can fetch the clients which are alive like this:
	for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		teamTrust := c.getTrust(team)
		offerSize := 0.0
		if status == shared.Critical &&
			giftSize > 0 {
			// If an island is critical takes precedence!
			// We donate everything we can according to our giftSize until we satisfy their
			// request.
			offerSize = math.Min(float64(receivedRequests[team]), giftSize)
			offers[team] = shared.GiftOffer(offerSize)
		} else if giftSize > 0 &&
			teamTrust > trustThreshold &&
			status == shared.Alive {
			// If an island is not critical and we can still donate something.
			// We offer them a proportional amount based on how much we trust them.
			// teamTrust is less or equal to 1.
			offerSize = giftSize * teamTrust
			if c.internalParam.giftExtra {
				offers[team] = shared.GiftOffer(offerSize + 10)
			} else {
				offers[team] = shared.GiftOffer(offerSize)
			}
			c.Logf("Team4 gifted %v %v resources", team, offerSize)
		}
		// else if giftSize < 0 && c.internalParam.giftExtra && (status == shared.Alive || status == shared.Critical) {
		// 	offers[team] = shared.GiftOffer(10)
		// }
		// Subtract the offer just made from the giftSize so that we don't offer too much.
		giftSize -= offerSize
	}
	c.Logf("Team4 gift offers: %v", offers)
	return offers
}

// DecideGiftAmount is executed at the end of each turn and asks clients how much
// they want to fulfill a gift offer they have made.
// COMPULSORY, you need to implement this method
func (c *client) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {

	// THIS JUST REPLICATES THE CODE IN "GetGiftOffers" to find our wealthGoal
	ourResources := c.getOurResources()
	// same as for requests. Now we will consider when you are more collaborative than greedy.
	importance := c.importances.getGiftRequestsImportance

	parameters := mat.NewVecDense(4, []float64{
		c.internalParam.greediness,
		c.internalParam.selfishness,
		c.internalParam.fairness,
		c.internalParam.collaboration,
	})
	// selflessLevel is the opposite of greedinessLevel in GetGiftRequests function.
	selflessLevel := mat.Dot(importance, parameters)

	// Initialising thresholds used.
	resourcesThreshold := 100.0

	// have a cap at 100 in case we are too selfless
	wealthGoal := math.Max(selflessLevel*float64(ourResources), resourcesThreshold)

	// giftSize is the maximum we can afford giving away.
	giftSize := float64(ourResources) - wealthGoal

	// Lets just return the min of "giftSize" and giftOffer.
	// The idea is that we are very careful in the function that
	// makes offers, so here we just keep our word as long as it does not
	// put us in a stressful state / our selfishness and greediness have
	// increased since we made the offer.
	actualGift := math.Min(float64(giftOffer), giftSize)
	return shared.Resources(actualGift)
}

// SentGift is executed at the end of each turn and notifies clients that
// their gift was successfully sent, along with the offer details.
// COMPULSORY, you need to implement this method

// func (c *client) SentGift(sent shared.Resources, to shared.ClientID) {
// 	// You can check your updated resources like this:
// 	myResources := c.getOurResources()
// 	percentageDonated := 0.0
// 	// avoid div 0
// 	if float64(sent+myResources) > 0 {
// 		percentageDonated = float64(sent / (sent + myResources))
// 	}
// 	bumpThreshold := 0.2
// 	currentGreediness := c.internalParam.greediness
// 	maxGreedBump := math.Min(currentGreediness+bumpThreshold, 1.0)

// 	currentSelfishness := c.internalParam.selfishness
// 	maxSelfishBump := math.Min(currentSelfishness+bumpThreshold, 1.0)

// 	// Bump up greediness and selfishness.
// 	// We use a maximum bump to cap at 1 and reduce the volatility.
// 	c.internalParam.greediness = math.Min(currentGreediness+percentageDonated, maxGreedBump)
// 	c.internalParam.selfishness = math.Min(currentSelfishness+percentageDonated, maxSelfishBump)

// }

// ReceivedGift is executed at the end of each turn and notifies clients that
// their gift was successfully received, along with the offer details.
// COMPULSORY, you need to implement this method
func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {
	// You can check your updated resources like this:
	myResources := c.getOurResources()
	percentageDonated := 0.0
	// avoid div 0
	if float64(myResources) > 0 {
		percentageDonated = float64(received / myResources)
	}
	bumpThreshold := 0.2

	currentCollab := c.internalParam.collaboration
	maxCollabBump := math.Min(currentCollab+bumpThreshold, 1.0)

	currentAgentTrust := c.trustMatrix.GetClientTrust(from)
	maxTrustBump := math.Min(currentAgentTrust+bumpThreshold, 1.0)
	// Bump up collaboration and agent trust.
	// We use a maximum bump to cap at 1 and reduce the volatility.
	c.internalParam.collaboration = math.Min(currentCollab+percentageDonated, maxCollabBump)
	c.trustMatrix.SetClientTrust(from, math.Min(currentAgentTrust+percentageDonated, maxTrustBump))
}

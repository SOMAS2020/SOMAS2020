package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// TODO: initialise confidence to 50
// TODO: Future work - initialise confidence to something more (based on strategy)

// Calculates the confidence we have in an island based on our past experience with them
// Depending on the situation we need to judge, we look at a different history
// The values in the histories should be updated in retrospect
func (c *client) confidence(situation Situation, otherIsland shared.ClientID) int {
	islandHist := c.opinionHist[otherIsland].Histories
	situationHist := islandHist[situation]
	sum := 0
	div := 0
	// TODO: change list iteration to just look at the turns we have info abt
	for i := len(situationHist); i > 0; i-- {
		sum += (situationHist[i-1]) * i
		div += i
	}

	average := sum / div

	islandSituationPerf := c.opinionHist[otherIsland].Performances[situation]
	islandSituationPerf.exp = average
	c.opinionHist[otherIsland].Performances[situation] = islandSituationPerf

	return average

}

func setLimits(confidence int) int {
	if confidence < 0 {
		return 0
	} else if confidence > 100 {
		return 100
	}
	return confidence
}

// Updates the HISTORY of an island in the required situation by comparing the expected
// performance with the reality
// Should be called after an action (with an island) has occurred
func (c *client) confidenceRestrospect(situation Situation, otherIsland shared.ClientID) {
	islandHist := c.opinionHist[otherIsland].Histories
	situationHist := islandHist[situation]

	islandSituationPerf := c.opinionHist[otherIsland].Performances[situation]
	situationExp := islandSituationPerf.exp
	situationReal := islandSituationPerf.real
	confidenceFactor := 0.5 // Factor by which the confidence increases/decreases, can be changed

	var updatedHist []int
	//TODO: make sure to check the range of values coz that might affect percentagediff when situationExp is 0
	percentageDiff := situationReal
	if situationExp != 0 {
		// between -100 and 100
		percentageDiff = 100 * (situationReal - situationExp) / situationExp
	}
	newConf := int(float64(percentageDiff)*confidenceFactor + float64(situationExp))
	updatedHist = append(situationHist, setLimits(newConf))

	c.opinionHist[otherIsland].Histories[situation] = updatedHist
}

// The implementation of this function (if needed) depends on where (and how) the confidence
// function is called in the first place
// func (c *client) confidenceReality(situation string, otherIsland shared.ClientID) {

// }

func max(numbers map[int]GiftInfo) int {
	var maxNumber int
	for maxNumber = range numbers {
		break
	}
	for n := range numbers {
		if n > maxNumber {
			maxNumber = n
		}
	}
	return maxNumber
}

// [[(1, toUS: 30, fromUS: 40), (2, toUS: 10, fromUS: 10), (3, toUS: 40, fromUS: 50)],
// [(1, toUS: 50, fromUS: 10), (2, toUS: 10, fromUS: 30), (3, toUS: 40, fromUS: 50)],
// [(1, toUS: 40, fromUS: 40), (2, toUS: 30, fromUS: 10), (3, toUS: 40, fromUS: 50)]]
// interactions with team toUs ... fromUs
// 1: 120 ...  90
// 2: 50 ... 50
// 3: 120 ... 150

// this just means the confidence we have in others while requesting gifts, not the trust we have on them
// Updates the confidence of an island regarding gifts
func (c *client) updateGiftConfidence(island shared.ClientID) int {
	turn := c.gameState().Turn
	pastConfidence := c.confidence("Gifts", island)

	ourLastTurnRequested := max(c.giftHist[island].OurRequest)
	theirLastTurnRequested := max(c.giftHist[island].IslandRequest)

	ourDonation := c.giftHist[island].IslandRequest[theirLastTurnRequested].gifted
	ourRequest := c.giftHist[island].OurRequest[ourLastTurnRequested].request
	theirDonation := c.giftHist[island].OurRequest[ourLastTurnRequested].gifted
	theirRequest := c.giftHist[island].IslandRequest[theirLastTurnRequested].request

	ourDonationRatio := ourDonation / theirRequest   // min: 0, max: 1
	theirDonationRatio := theirDonation / ourRequest // min: 0, max: 1

	//should take into account the dimensions of the contributions ratio of 0.9
	// 9 out of 10 != 90 out of 100
	if mod(ourDonation-theirDonation) > 10 {

	}

	switch {
	case ourLastTurnRequested == theirLastTurnRequested: //we both requested a gift from eachother:
		switch {
		case theirDonationRatio > ourDonationRatio:
			// return itwentamazing
		case ourDonationRatio > theirDonationRatio:
			// return itwentokay
		case numberoftheirrequests << ournumberofrequests:
			// return itwentreasonable // < itwentokay
		case theiramountrequested << ouramountrequested:
			// return itwentreasonable
		case theyarecritical:
			// return theyshouldknowbetterbutfine // < itwentreasonable
		default:
			// return bigdisappointment // < theyshouldknowbetterbutfine
		}
	case ourRequest > 0 && theirRequest == 0:

	case ourRequest == 0 && theirRequest > 0:
	}

	// how much they give us vs how much we requested, how much we give them,
	// how many times they request, how much they've previously requested, state they're in

	return total
}

//func (c *client) credibility(situation Situation, otherIsland shared.ClientID) int {
//Situation
func (c *client) credibility(situation Situation, otherIsland shared.ClientID) int {
	// Situation
	// Long term vs short term importance
	// how much they have gifted in general
	// their transparency, ethical behaviour as an island (have they shared their foraging predictions, their cp intended contributions, etc)
	// their empathy level
	// how they acted during a role
	// performance (how well they are doing)
	return 0
}

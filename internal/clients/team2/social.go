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

	islandHist, ok := c.opinionHist[otherIsland].Histories
	if(!ok){
		
	}
	situationHist := islandHist[situation]
	sum := 0
	div := 0
	// TODO: change list iteration to just look at the turns we have info abt
	for i := len(situationHist); i > 0; i-- {
		sum += (situationHist[i-1])*i
		div += i
	}

	average := sum / div

	islandSituationPerf := c.opinionHist[otherIsland].Performances[situation]
	islandSituationPerf.exp = average
	c.opinionHist[otherIsland].Performances[situation] = islandSituationPerf

	return average

}

func setLimits(confidence int) int {
	if(confidence < 0){
		return 0
	} else if (confidence > 100) {
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
	percetageDiff := situationReal
	if(situationExp != 0){
		// between -100 and 100
		percentageDiff = 100*(situationReal - situationExp)/situationExp
	}
	newConf := int(percentageDiff*confidenceFactor + situationExp)
	updatedHist = append(situationHist, setLimits(newConf))

	c.opinionHist[otherIsland].Histories[situation] = updatedHist
}

// The implementation of this function (if needed) depends on where (and how) the confidence
// function is called in the first place
// func (c *client) confidenceReality(situation string, otherIsland shared.ClientID) {

// }

func (c *client) giftRequestedConfidence(island shared.ClientID) int  {
	pastConfidence := c.confidence("ReceivedRequests", island)
	// how much they gave us in past, how many times they request, how much they've previously requested, state they're in
	[[(1, toUS: 30, fromUS: 40), (2, toUS: 10, fromUS: 10), (3, toUS: 40, fromUS: 50)],
	[(1, toUS: 50, fromUS: 10), (2, toUS: 10, fromUS: 30), (3, toUS: 40, fromUS: 50)],
	[(1, toUS: 40, fromUS: 40), (2, toUS: 30, fromUS: 10), (3, toUS: 40, fromUS: 50)]]
	interactions with team toUs ... fromUs
	1: 120 ...  90
	2: 50 ... 50
	3: 120 ... 150
	
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

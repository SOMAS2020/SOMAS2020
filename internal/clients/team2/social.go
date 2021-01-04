package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Calculates the confidence we have in an island based on our past experience with them
// Depending on the situation we need to judge, we look at a different history
// The values in the histories should be updated in retrospect
func (c *client) confidence(situation Situation, otherIsland shared.ClientID) int {
	islandHist := c.opinionHist[otherIsland].Histories
	situationHist := islandHist[situation]
	sum := 0
	for i := 0; i < len(situationHist); i++ {
		sum += (situationHist[i])
	}

	average := sum / (len(situationHist))

	islandSituationPerf := c.opinionHist[otherIsland].Performances[situation]
	islandSituationPerf.exp = average
	c.opinionHist[otherIsland].Performances[situation] = islandSituationPerf

	return average

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
	confidenceFactor := 5 // Factor by which the confidence increases/decreases, can be changed

	var updatedHist []int
	if situationExp > situationReal { // We expected more
		diff := situationExp - situationReal
		updatedHist = append(situationHist, situationExp-diff*confidenceFactor)
	} else if situationExp < situationReal {
		diff := situationReal - situationExp
		updatedHist = append(situationHist, situationExp+diff*confidenceFactor)
	} else {
		updatedHist = append(situationHist, situationExp)
	}

	c.opinionHist[otherIsland].Histories[situation] = updatedHist
}

// The implementation of this function (if needed) depends on where (and how) the confidence
// function is called in the first place
// func (c *client) confidenceReality(situation string, otherIsland shared.ClientID) {

// }

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

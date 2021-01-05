package team2

import (
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// TODO: Future work - initialise confidence to something more (based on strategy)

type IslandTrustMap map[int]GiftInfo

// Overwrite default sort implementation
func (p IslandTrustMap) Len() int           { return len(p) }
func (p IslandTrustMap) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p IslandTrustMap) Less(i, j int) bool { return p[i].trust < p[j].trust }

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

func max(numbers map[uint]GiftInfo) uint {
	var maxNumber uint
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
	const bufferLen = 10

	// runningMean := previousMean + (nthValue - previousMean)/numberofTurnsBefore
	runMeanTheyReq := 0
	runMeanTheyDon := 0
	runMeanWeReq := 0
	runMeanWeDon := 0

	ourReqMap := c.giftHist[island].OurRequest
	theirReqMap := c.giftHist[island].IslandRequest

	// Sort the keys in decreasing order
	sort.Ints(ourReqMap)
	sort.Ints(theirReqMap)

	// Idea is to look at the past 10 transactions, not turns
	ourLastTurns := [bufferLen]int{}
	theirLastTurns := [bufferLen]int{}
	numKeys := 0
	// Take keys (turns) of our last 10 transactions
	for key := range ourReqMap {
		ourLastTurns = append(ourLastTurns, key)
		if numKeys > bufferLen {
			break
		}
	}

	numKeys = 0
	// Take keys (turns) of their last 10 transactions
	for key := range theirReqMap {
		theirLastTurns = append(theirLastTurns, key)
		if numKeys > bufferLen {
			break
		}
	}

	// Take running average of the interactions
	// The individual turn values will be scaled wrt to the "distance" from the current turn
	// ie transactions further in the past are valued less
	for i := 0; i < bufferLen; i++ {
		// Get the transaction distance to the previous transaction
		theirTransDist := turn - theirLastTurns[i]
		ourTransDist := turn - ourLastTurns[i]
		// Update the respective running mean factoring in the transactionDistance (inv proportioanl to transactionDistance so farther transactions are weighted less)
		runMeanTheyReq = runMeanTheyReq + (theirReqMap[theirLastTurns[i]].requested/theirTransDist-runMeanTheyReq)/(i+1)
		runMeanTheyDon = runMeanTheyDon + (ourReqMap[ourLastTurns[i]].gifted/ourTransDist-runMeanTheyDon)/(i+1)
		runMeanWeReq = runMeanWeReq + (ourReqMap[ourLastTurns[i]].requested/ourTransDist-runMeanWeReq)/(i+1)
		runMeanWeDon = runMeanWeDon + (theirReqMap[theirLastTurns[i]].gifted/theirTransDist-runMeanWeDon)/(i+1)
	}

	usRatio := runMeanTheyDon / runMeanWeReq
	themRatio := runMeanWeDon / runMeanTheyReq

	if usRatio >= themRatio {
		// confidence increases
	} else {
		// confidence decreases
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

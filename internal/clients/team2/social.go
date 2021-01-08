package team2

import (
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// TODO: Future work - initialise confidence to something more (based on strategy)

type IslandTrustMap map[int]GiftInfo

// Overwrite default sort implementation
func (p IslandTrustMap) Len() int      { return len(p) }
func (p IslandTrustMap) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// func (p IslandTrustMap) Less(i, j int) bool { return p[i].trust < p[j].trust }

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
	c.Logf("%v", c.opinionHist)
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
func MinInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

// this just means the confidence we have in others while requesting gifts, not the trust we have on them
// Updates the confidence of an island regarding gifts
func (c *client) updateGiftConfidence(island shared.ClientID) int {
	turn := c.gameState().Turn
	pastConfidence := c.confidence("Gifts", island)

	var bufferLen = 0
	if turn < 10 {
		bufferLen = int(turn)
	} else {
		bufferLen = 10
	}

	runMeanTheyReq := 0.0
	runMeanTheyDon := 0.0
	runMeanWeReq := 0.0
	runMeanWeDon := 0.0

	ourReqMap := c.giftHist[island].OurRequest
	theirReqMap := c.giftHist[island].IslandRequest

	ourKeys := make([]int, 0)
	for k, _ := range ourReqMap {
		ourKeys = append(ourKeys, int(k))
	}

	theirKeys := make([]int, 0)
	for k, _ := range theirReqMap {
		theirKeys = append(theirKeys, int(k))
	}

	// Sort the keys in decreasing order
	sort.Ints(ourKeys)
	sort.Ints(theirKeys)

	// Take running average of the interactions
	// The individual turn values will be scaled wrt to the "distance" from the current turn
	// ie transactions further in the past are valued less
	if MinInt(len(ourKeys), len(theirKeys)) == 0 {
		return pastConfidence
	}
	c.Logf("Bufferlen %v", bufferLen)
	for i := 0; i < MinInt(bufferLen, len(ourKeys)); i++ {
		// Get the transaction distance to the previous transaction
		ourTransDist := turn - uint(ourKeys[i])
		// Update the respective running mean factoring in the transactionDistance (inv proportioanl to transactionDistance so farther transactions are weighted less)
		runMeanTheyDon = runMeanTheyDon + (float64(ourReqMap[uint(ourKeys[i])].gifted)/float64(ourTransDist)-float64(runMeanTheyDon))/float64(i+1)
		runMeanWeReq = runMeanWeReq + (float64(ourReqMap[uint(ourKeys[i])].requested)/float64(ourTransDist)-float64(runMeanWeReq))/float64(i+1)
	}
	for i := 0; i < MinInt(bufferLen, len(theirKeys)); i++ {
		// Get the transaction distance to the previous transaction
		theirTransDist := turn - uint(theirKeys[i])
		// Update the respective running mean factoring in the transactionDistance (inv proportioanl to transactionDistance so farther transactions are weighted less)
		runMeanTheyReq = runMeanTheyReq + (float64(theirReqMap[uint(theirKeys[i])].requested)/float64(theirTransDist)-float64(runMeanTheyReq))/float64(i+1)
		runMeanWeDon = runMeanWeDon + (float64(theirReqMap[uint(theirKeys[i])].gifted))/float64(theirTransDist) - float64(runMeanWeDon)/float64(i+1)
	}

	// TODO: is there a potential divide by 0 here?
	usRatio := runMeanTheyDon / runMeanWeReq   // between 0 and 1
	themRatio := runMeanWeDon / runMeanTheyReq // between 0 and 1

	diff := usRatio - themRatio // between -1 and 1
	// confidence increases if usRatio >= themRatio
	// confidence decreases if not

	// e.g. 1 pastConfidnece = 50%
	// diff = 100% in our favour 1.0
	// inc pastConfidence = (50 + 100)/2 = 75

	// e.g. 2 pastConfidence = 90%
	// diff = 70% in our favour
	// inc pastConfidence = (90 + 70)/2 = 80

	// e.g. 3 pastConfidence = 80%
	// diff = 30% against us
	// inc pastConfidence = (80 - 30)/2 = 25

	// e.g. 4 pastConfidence = 100%
	// diff = 100% against us
	// inc pastConfidence = (100 - 100)/2 = 0

	// e.g. 5 pastConfidence = 0%
	// diff = 100% in our favour
	// inc pastConfidence = (0 + 100)/2 = 50

	// TODO: improve how ratios are used to improve pastConfidence
	// pastConfidence = (pastConfidence + sensitivity*diff*100) / 2
	pastConfidence = int((pastConfidence + int(diff*100)) / 2)

	return pastConfidence
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

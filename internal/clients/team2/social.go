package team2

import (
	"math"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// TODO: Future work - initialise confidence to something more (based on strategy)

type IslandTrustMap map[int]GiftInfo

// Overwrite default sort implementation
func (p IslandTrustMap) Len() int      { return len(p) }
func (p IslandTrustMap) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// func (p IslandTrustMap) Less(i, j int) bool { return p[i].trust < p[j].trust }

func (c *client) initialiseOpinionForIsland(otherIsland shared.ClientID) {
	Histories := make(map[Situation][]int)
	Histories["President"] = []int{50}
	Histories["RoleOpinion"] = []int{50}
	Histories["Judge"] = []int{50}
	Histories["Foraging"] = []int{50}
	Histories["Gifts"] = []int{50}
	c.opinionHist[otherIsland] = Opinion{
		Histories:    Histories,
		Performances: map[Situation]ExpectationReality{},
	}
}

// Calculates the confidence we have in an island based on our past experience with them
// Depending on the situation we need to judge, we look at a different history
// The values in the histories should be updated in retrospect
func (c *client) confidence(situation Situation, otherIsland shared.ClientID) int {
	trust := 50
	// The default for no data is to trust other islands 50%

	if _, ok := c.opinionHist[otherIsland]; !ok {
		c.initialiseOpinionForIsland(otherIsland)
		c.Logf("Initialised opinions: ", c.opinionHist)
	}

	// We have initialised the histories
	islandHist := c.opinionHist[otherIsland].Histories
	c.Logf("Situation ", situation)
	c.Logf("Island ", otherIsland)
	c.Logf("History with island ", islandHist)
	situationHist := islandHist[situation]

	// Check if there is a history to take a weighted average from
	if len(situationHist) != 0 {
		sum := 0
		div := 0

		for i := len(situationHist); i > 0; i-- {
			sum += (situationHist[i-1]) * i
			div += i
		}
		trust = sum / div
	}

	// Set the expectation to the weighted past average or the default 50
	islandSituationPerf := ExpectationReality{
		exp:  trust,
		real: 0,
	}
	c.Logf("[Situation Performance (exp)]:", islandSituationPerf)

	//c.Logf("island opinions", c.opinionHist[shared.ClientID(4)])
	c.opinionHist[otherIsland].Performances[situation] = islandSituationPerf
	c.Logf("[returning full object]:", c.opinionHist[otherIsland])
	return trust

}

func (c *client) setLimits(confidence int) int {
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
	c.Logf("[our opinion of others]:", c.opinionHist[otherIsland])
	if opinion, ok := c.opinionHist[otherIsland]; ok {
		islandHist := opinion.Histories
		situationHist := islandHist[situation]
		c.Logf("situation we're dealing with ", situation)
		c.Logf("[situation before]:", situationHist)
		islandSituationPerf := opinion.Performances[situation]
		situationExp := islandSituationPerf.exp
		situationReal := islandSituationPerf.real
		c.Logf("[current performance]:", islandSituationPerf)

		percentageDiff := situationReal
		if situationExp != 0 { // Forgiveness principle: if we had 0 expectation, give them a chance to improve
			// between -100 and 100
			percentageDiff = 100 * (situationReal - situationExp) / situationExp
		}
		newConf := int(float64(percentageDiff)*c.config.ConfidenceRetrospectFactor + float64(situationExp))
		updatedHist := append(situationHist, c.setLimits(newConf))
		c.Logf("[situation after]:", updatedHist)
		c.Logf("[appended yuhuuu]:", len(updatedHist))

		c.opinionHist[otherIsland].Histories[situation] = updatedHist
		c.Logf("[did it really append?]:", c.opinionHist[otherIsland].Histories[situation])
	}
}

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

	var bufferLen int

	if turn < 10 {
		bufferLen = int(turn)
	} else {
		bufferLen = 10
	}

	runMeanTheyReq := 0.0
	runMeanTheyDon := 0.0
	runMeanWeReq := 0.0
	runMeanWeDon := 0.0

	var ourReqMap, theirReqMap map[uint]GiftInfo

	if hist, ok := c.giftHist[island]; ok {
		ourReqMap = hist.OurRequest
		theirReqMap = hist.IslandRequest
		c.Logf("ourReqMap: ", ourReqMap)
		c.Logf("theirReqMap: ", theirReqMap)
	} else {
		// TODO: this could be something
		return pastConfidence
	}

	ourKeys := make([]int, 0)
	for k := range ourReqMap {
		ourKeys = append(ourKeys, int(k))
	}

	theirKeys := make([]int, 0)
	for k := range theirReqMap {
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
		ourTransDist := turn - uint(ourKeys[i]) + 1
		// Update the respective running mean factoring in the transactionDistance (inv proportioanl to transactionDistance so farther transactions are weighted less)
		runMeanTheyDon = runMeanTheyDon + (float64(ourReqMap[uint(ourKeys[i])].gifted)/float64(ourTransDist)-float64(runMeanTheyDon))/float64(i+1)
		runMeanWeReq = runMeanWeReq + (float64(ourReqMap[uint(ourKeys[i])].requested)/float64(ourTransDist)-float64(runMeanWeReq))/float64(i+1)
	}

	for i := 0; i < MinInt(bufferLen, len(theirKeys)); i++ {
		// Get the transaction distance to the previous transaction
		theirTransDist := turn - uint(theirKeys[i]) + 1
		// Update the respective running mean factoring in the transactionDistance (inv proportioanl to transactionDistance so farther transactions are weighted less)
		runMeanTheyReq = runMeanTheyReq + (float64(theirReqMap[uint(theirKeys[i])].requested)/float64(theirTransDist)-float64(runMeanTheyReq))/float64(i+1)
		runMeanWeDon = runMeanWeDon + (float64(theirReqMap[uint(theirKeys[i])].gifted))/float64(theirTransDist) - float64(runMeanWeDon)/float64(i+1)
	}

	// TODO: if this checkDivZero returns 1 this is not between 0 and 1
	usRatio := runMeanTheyDon / checkDivZero(runMeanWeReq)   // between 0 and 1
	themRatio := runMeanWeDon / checkDivZero(runMeanTheyReq) // between 0 and 1
	diff := usRatio - themRatio                              // between -1 and 1
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
	c.Logf("[Past confidence (GIFTS)]:", pastConfidence)
	c.Logf("[Diff added (GIFTS)]: ", diff)
	pastConfidence = int((pastConfidence + int(diff*100)) / 2)
	updatedHist := append(c.opinionHist[island].Histories["Gifts"], c.setLimits(pastConfidence))
	c.Logf("[Gift situation updated]:", updatedHist)
	c.Logf("[Length of gift situation after update]:", len(updatedHist))

	c.opinionHist[island].Histories["Gifts"] = updatedHist

	return pastConfidence
}

func (c *client) updatePresidentTrust() {
	currPres := c.gameState().PresidentID

	// Default value for Opinion if we have no history
	reality := 50

	if presHist, ok := c.presCommonPoolHist[currPres]; ok {
		c.Logf("c.presCommonPoolHist", presHist)
		runMeanTax := shared.Resources(0)
		runMeanWeRequest := shared.Resources(0)
		runMeanWeAllocated := shared.Resources(0)
		runMeanWeTake := shared.Resources(0)
		counter := shared.Resources(1)

		// Running average m(n) = m(n-1) + (a(n) - m(n-1))/n
		for _, commonPool := range presHist {
			runMeanTax = runMeanTax + (commonPool.tax-runMeanTax)/shared.Resources(counter)
			runMeanWeRequest = runMeanWeRequest + (commonPool.requestedToPres-runMeanWeRequest)/shared.Resources(counter)
			runMeanWeAllocated = runMeanWeAllocated + (commonPool.allocatedByPres-runMeanWeAllocated)/shared.Resources(counter)
			runMeanWeTake = runMeanWeTake + (commonPool.takenFromCP-runMeanWeTake)/shared.Resources(counter)
			counter++
		}

		percChangeTax := 100 * (c.taxAmount - runMeanTax) / shared.Resources(checkDivZero(float64(runMeanTax)))

		// How much less we're giveen
		percWeGet := 100 * (runMeanWeRequest - runMeanWeAllocated) / shared.Resources(checkDivZero(float64(runMeanWeAllocated)))

		// How much more we've taken
		percWeTake := 100 * (runMeanWeAllocated - runMeanWeTake) / shared.Resources(checkDivZero(float64(runMeanWeTake)))

		reality = c.setLimits(int(100 - percWeGet - percChangeTax + percWeTake))
	}

	islandSituationPerf := ExpectationReality{
		exp:  50, // Would not get overwritten if we have no current expectation, so default should be 50
		real: reality,
	}

	if perf, ok := c.opinionHist[currPres].Performances["President"]; ok {
		islandSituationPerf.exp = perf.exp
	}

	c.Logf("Previous performance for President", c.opinionHist[currPres].Performances["President"])
	c.Logf("Updated performance for President", islandSituationPerf)
	c.opinionHist[currPres].Performances["President"] = islandSituationPerf
}

func (c *client) updateJudgeTrust() {
	currJudge := c.gameState().JudgeID

	numConsecTier := 0
	numDiffTiers := 0
	avgTurnsPerTier := 0
	runMeanScore := 0
	reality := 50

	if _, ok := c.sanctionHist[currJudge]; ok {
		c.Logf("sanction hist for judge ", c.sanctionHist[currJudge])
		prevTier := c.sanctionHist[currJudge][0].Tier
		for i, sanction := range c.sanctionHist[currJudge] {
			turn := int(c.gameState().Turn - sanction.Turn)
			div := i + 1

			runMeanScore += (sanction.Amount/turn - runMeanScore) / div
			if prevTier == sanction.Tier {
				numConsecTier++
			} else {
				numDiffTiers++
				avgTurnsPerTier += (numConsecTier - avgTurnsPerTier) / numDiffTiers
				prevTier = sanction.Tier
				numConsecTier = 0
			}
		}

		// We don't want to be sanctioned
		// We don't want a sanction to last too long
		// We want the judge to be "fair"

		lastScore := c.sanctionHist[currJudge][len(c.sanctionHist[currJudge])-1]
		percChangeScore := int(100 * int((lastScore.Amount-runMeanScore)/runMeanScore))
		reality = c.setLimits(100 - (avgTurnsPerTier * percChangeScore))
	}
	islandSituationPerf := ExpectationReality{
		exp:  50, // Would not get overwritten if we have no current expectation, so default should be 50
		real: reality,
	}
	if perf, ok := c.opinionHist[currJudge].Performances["Judge"]; ok {
		islandSituationPerf.exp = perf.exp
	}

	islandSituationPerf.real = reality

	c.Logf("Updated performance for Judge", islandSituationPerf)

	c.opinionHist[currJudge].Performances["Judge"] = islandSituationPerf
	c.confidenceRestrospect("Judge", currJudge)

}

type AccountabilityInfo struct {
	AllocationRequestsMade         []float64
	AllocationMade                 []float64
	ExpectedTaxContribution        []float64
	ExpectedAllocation             []float64
	IslandTaxContribution          []float64
	IslandAllocation               []float64
	SanctionPaid                   []float64
	SanctionExpected               []float64
	IslandActualPrivateResources   []float64
	IslandReportedPrivateResources []float64
}

func (c *client) getWeightedAverage(list []float64) int {
	if len(list) == 0 {
		return 0
	}

	div := 1
	total := 0
	for i, item := range list {
		total += i * int(item)
		div += i
	}
	return total / div
}

func (c *client) updateRoleTrust(iigoHistory []shared.Accountability) {

	//Interested in how much they took vs how much they were allowed to
	// Interested in how much they said they have vs how much they actually have
	// How much they've been sanctioned vs How much theey're paying
	islandInfo := make(map[shared.ClientID]*AccountabilityInfo)
	emptyInt := []float64{}
	// Initialise islandInfo for all Alive islands
	for _, island := range c.getAliveClients() {
		islandInfo[island] = &AccountabilityInfo{
			AllocationRequestsMade:         emptyInt,
			AllocationMade:                 emptyInt,
			ExpectedTaxContribution:        emptyInt,
			ExpectedAllocation:             emptyInt,
			IslandTaxContribution:          emptyInt,
			IslandAllocation:               emptyInt,
			SanctionPaid:                   emptyInt,
			SanctionExpected:               emptyInt,
			IslandActualPrivateResources:   emptyInt,
			IslandReportedPrivateResources: emptyInt,
		}
	}

	c.Logf("IslandInfo ", islandInfo)

	for _, accountability := range iigoHistory {
		if c.isAlive(accountability.ClientID) {
			for _, pair := range accountability.Pairs {
				switch pair.VariableName {
				case 24: // ExpectedTaxContribution
					islandInfo[accountability.ClientID].ExpectedTaxContribution = pair.Values
				case 25: // ExpectedAllocation
					islandInfo[accountability.ClientID].ExpectedAllocation = pair.Values
				case 26: // IslandTaxContribution
					islandInfo[accountability.ClientID].IslandTaxContribution = pair.Values
				case 27: // IslandAllocation
					islandInfo[accountability.ClientID].IslandAllocation = pair.Values
				case 31: // SanctionPaid
					islandInfo[accountability.ClientID].SanctionPaid = pair.Values
				case 32: // SanctionExpected
					islandInfo[accountability.ClientID].SanctionExpected = pair.Values
				case 52: // IslandActualPrivateResources
					islandInfo[accountability.ClientID].IslandActualPrivateResources = pair.Values
				case 53: // IslandReportedPrivateResources
					islandInfo[accountability.ClientID].IslandReportedPrivateResources = pair.Values
				}
			}
		}

	}
	for island, accountability := range islandInfo {
		allocationDiff := 0
		taxContribDiff := 0
		sanctionDiff := 0
		islandResourceDiff := 0
		if len(accountability.ExpectedTaxContribution) != 0 && len(accountability.IslandTaxContribution) != 0 {
			avgExpected := c.getWeightedAverage(accountability.ExpectedTaxContribution)
			avgActual := c.getWeightedAverage(accountability.IslandTaxContribution)
			if avgActual != 0 {
				taxContribDiff = 100 * (avgExpected - avgActual) / avgActual
			}
		}
		if len(accountability.ExpectedAllocation) != 0 && len(accountability.IslandAllocation) != 0 {
			avgExpected := c.getWeightedAverage(accountability.ExpectedAllocation)
			avgActual := c.getWeightedAverage(accountability.IslandAllocation)
			if avgActual != 0 {
				allocationDiff = 100 * (avgExpected - avgActual) / avgActual
			}
		}
		if len(accountability.SanctionPaid) != 0 && len(accountability.SanctionExpected) != 0 {
			avgExpected := c.getWeightedAverage(accountability.SanctionExpected)
			avgActual := c.getWeightedAverage(accountability.SanctionPaid)
			if avgActual != 0 {
				sanctionDiff = 100 * (avgExpected - avgActual) / avgActual
			}
		}
		if len(accountability.IslandActualPrivateResources) != 0 && len(accountability.IslandReportedPrivateResources) != 0 {
			avgExpected := c.getWeightedAverage(accountability.IslandReportedPrivateResources)
			avgActual := c.getWeightedAverage(accountability.IslandActualPrivateResources)
			if avgActual != 0 {
				islandResourceDiff = 100 * (avgExpected - avgActual) / avgActual
			}
		}

		c.Logf("Role trust for turn ", c.gameState().Turn)
		c.Logf("For island ", island)
		c.Logf("Accountability map", accountability)
		c.Logf("IIGO History", iigoHistory)
		c.Logf("taxContribDiff", taxContribDiff)
		c.Logf("allocationDiff", allocationDiff)
		c.Logf("sanctionDiff", sanctionDiff)
		c.Logf("islandResourceDiff", islandResourceDiff)
		reality := c.setLimits(100 - taxContribDiff - allocationDiff - sanctionDiff - islandResourceDiff)
		islandSituationPerf := ExpectationReality{
			exp:  50,
			real: reality,
		}
		if _, ok := c.opinionHist[island]; ok {
			islandSituationPerf.exp = c.opinionHist[island].Performances["RoleOpinion"].exp
		} else {
			c.initialiseOpinionForIsland(island)
		}
		c.Logf("Updated performance ", islandSituationPerf)
		c.opinionHist[island].Performances["RoleOpinion"] = islandSituationPerf

	}

}

// disasters:

// how accurate they were -- once a disaster happens
// 	magnitude and time remaining until disaster

// check avg of their predictions over that season vs magnitude that actually occurred
// number of turns off
// confidence

//This function is called when a disaster occurs to update our confidence on others' predictions
func (c *client) updateDisasterConf() {

	// If a disaster has occurred
	disasterMag := 0.0
	disasterTurn := uint(0)
	if len(c.disasterHistory) > 1 {
		disasterMag = c.disasterHistory[len(c.disasterHistory)-1].Report.Magnitude
		disasterTurn = c.disasterHistory[len(c.disasterHistory)-1].Turn

	}

	for island, predictions := range c.predictionHist {
		// If we have received predictions from others
		avgMag := 0.0
		avgConf := 0.0
		avgTurn := 0
		if len(predictions) > 0 {
			for _, prediction := range predictions {
				avgMag += prediction.Prediction.Magnitude
				avgConf += prediction.Prediction.Confidence
				avgTurn += int(prediction.Turn)
			}

			// The three metrics we will assess an island by
			avgTurn = avgTurn / len(predictions)
			avgMag = avgMag / float64(len(predictions))
			avgConf = avgConf / float64(len(predictions))
		}

		magError := int(100 * math.Abs(avgMag-disasterMag) / checkDivZero(disasterMag))                           // percentage error
		turnError := int(100 * math.Abs(float64(uint(avgTurn)-disasterTurn)/checkDivZero(float64(disasterTurn)))) // percentage error

		predError := int(avgConf) * (magError + turnError)

		predConf := 100 - c.setLimits(predError)

		c.opinionHist[island].Performances["DisasterPred"] = ExpectationReality{
			real: predConf,
		}

		c.confidenceRestrospect("DisasterPred", island)

	}

}

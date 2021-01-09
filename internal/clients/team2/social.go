package team2

import (
	"math"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"

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
	if len(situationHist) == 0 {
		return 50
	}
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
	updatedHist = append(situationHist, c.setLimits(newConf))

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
	for i := 0; i < bufferLen; i++ {
		// Get the transaction distance to the previous transaction
		theirTransDist := turn - uint(theirKeys[i])
		ourTransDist := turn - uint(ourKeys[i])
		// Update the respective running mean factoring in the transactionDistance (inv proportioanl to transactionDistance so farther transactions are weighted less)
		runMeanTheyReq = runMeanTheyReq + (float64(theirReqMap[uint(theirKeys[i])].requested)/float64(theirTransDist)-float64(runMeanTheyReq))/float64(i+1)
		runMeanTheyDon = runMeanTheyDon + (float64(ourReqMap[uint(ourKeys[i])].gifted)/float64(ourTransDist)-float64(runMeanTheyDon))/float64(i+1)
		runMeanWeReq = runMeanWeReq + (float64(ourReqMap[uint(ourKeys[i])].requested)/float64(ourTransDist)-float64(runMeanWeReq))/float64(i+1)
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

func (c *client) updatePresidentTrust() {
	currPres := c.gameState().PresidentID
	// Take weighted average of past turns

	runMeanTax := shared.Resources(0.0)
	runMeanWeRequest := shared.Resources(0.0)
	runMeanWeAllocated := shared.Resources(0.0)
	runMeanWeTake := shared.Resources(0.0)

	for i, commonPool := range c.commonPoolHist[currPres] {
		turn := shared.Resources(c.gameState().Turn - commonPool.turn)
		div := shared.Resources(i + 1)

		runMeanTax += (commonPool.tax/turn - runMeanTax) / div
		runMeanWeRequest += (commonPool.requestedToPres/turn - runMeanWeRequest) / div
		runMeanWeAllocated += (commonPool.allocatedByPres/turn - runMeanWeAllocated) / div
		runMeanWeTake += (commonPool.takenFromCP/turn - runMeanWeTake) / div

	}

	percChangeTax := 100 * (c.taxAmount - runMeanTax) / runMeanTax
	percWeGet := 100 * (runMeanWeRequest - runMeanWeAllocated) / runMeanWeAllocated // How much less we're giveen
	percWeTake := 100 * (runMeanWeAllocated - runMeanWeTake) / runMeanWeTake        // How much more we've taken

	reality := c.setLimits(int(100 - percWeGet - percChangeTax + percWeTake))

	islandSituationPerf := c.opinionHist[currPres].Performances["President"]
	islandSituationPerf.real = reality
	c.opinionHist[currPres].Performances["President"] = islandSituationPerf

}

func (c *client) updateJudgeTrust() {
	currJudge := c.gameState().JudgeID

	prevTier := c.sanctionHist[currJudge][0].Tier
	numConsecTier := 0
	numDiffTiers := 0
	avgTurnsPerTier := 0
	runMeanScore := roles.IIGOSanctionScore(0.0)

	for i, sanction := range c.sanctionHist[currJudge] {
		turn := roles.IIGOSanctionScore(c.gameState().Turn - sanction.Turn)
		div := roles.IIGOSanctionScore(i + 1)

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
	percChangeScore := int(100 * (lastScore.Amount - runMeanScore) / runMeanScore)
	reality := c.setLimits(100 - (avgTurnsPerTier * percChangeScore))

	islandSituationPerf := c.opinionHist[currJudge].Performances["Judge"]
	islandSituationPerf.real = reality
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
		total *= i * int(item)
		div += i
	}
	return total / div
}

func (c *client) updateRoleTrust(iigoHistory []shared.Accountability) {

	//Interested in how much they took vs how much they were allowed to
	// Interested in how much they said they have vs how much they actually have
	// How much they've been sanctioned vs How much theey're paying

	for _, info := range iigoHistory {
		island := info.ClientID
		islandInfo := info.Pairs
		var islAccInfo AccountabilityInfo
		for _, pair := range islandInfo {
			switch pair.VariableName {
			case 24: // ExpectedTaxContribution
				islAccInfo.ExpectedTaxContribution = pair.Values
			case 25: // ExpectedAllocation
				islAccInfo.ExpectedAllocation = pair.Values
			case 26: // IslandTaxContribution
				islAccInfo.IslandTaxContribution = pair.Values
			case 27: // IslandAllocation
				islAccInfo.IslandAllocation = pair.Values
			case 31: // SanctionPaid
				islAccInfo.SanctionPaid = pair.Values
			case 32: // SanctionExpected
				islAccInfo.SanctionExpected = pair.Values
			case 52: // IslandActualPrivateResources
				islAccInfo.IslandActualPrivateResources = pair.Values
			case 53: // IslandReportedPrivateResources
				islAccInfo.IslandReportedPrivateResources = pair.Values
			}
		}
		allocationDiff := 0
		taxContribDiff := 0
		sanctionDiff := 0
		islandResourceDiff := 0
		if islAccInfo.ExpectedTaxContribution != nil && islAccInfo.IslandTaxContribution != nil {
			avgExpected := c.getWeightedAverage(islAccInfo.ExpectedTaxContribution)
			avgActual := c.getWeightedAverage(islAccInfo.IslandTaxContribution)
			if avgActual != 0 {
				taxContribDiff = 100 * (avgExpected - avgActual) / avgActual
			}
		}
		if islAccInfo.ExpectedAllocation != nil && islAccInfo.IslandAllocation != nil {
			avgExpected := c.getWeightedAverage(islAccInfo.ExpectedAllocation)
			avgActual := c.getWeightedAverage(islAccInfo.IslandAllocation)
			if avgActual != 0 {
				allocationDiff = 100 * (avgExpected - avgActual) / avgActual
			}
		}
		if islAccInfo.SanctionPaid != nil && islAccInfo.SanctionExpected != nil {
			avgExpected := c.getWeightedAverage(islAccInfo.SanctionExpected)
			avgActual := c.getWeightedAverage(islAccInfo.SanctionPaid)
			if avgActual != 0 {
				sanctionDiff = 100 * (avgExpected - avgActual) / avgActual
			}
		}
		if islAccInfo.IslandActualPrivateResources != nil && islAccInfo.IslandReportedPrivateResources != nil {
			avgExpected := c.getWeightedAverage(islAccInfo.IslandReportedPrivateResources)
			avgActual := c.getWeightedAverage(islAccInfo.IslandActualPrivateResources)
			if avgActual != 0 {
				islandResourceDiff = 100 * (avgExpected - avgActual) / avgActual
			}
		}
		reality := c.setLimits(100 - taxContribDiff - allocationDiff - sanctionDiff - islandResourceDiff)
		islandSituationPerf := c.opinionHist[island].Performances["RoleOpinion"]
		islandSituationPerf.real = reality
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
	disasterMag := c.disasterHistory[len(c.disasterHistory)-1].Report.Magnitude
	disasterTurn := c.disasterHistory[len(c.disasterHistory)-1].Turn
	for island, predictions := range c.predictionHist {
		avgMag := 0.0
		avgConf := 0.0
		avgTurn := 0
		for _, prediction := range predictions {
			avgMag += prediction.Prediction.Magnitude
			avgConf += prediction.Prediction.Confidence
			avgTurn += int(prediction.Turn)
		}

		// The three metrics we will assess an island by
		avgTurn = avgTurn / len(predictions)
		avgMag = avgMag / float64(len(predictions))
		avgConf = avgConf / float64(len(predictions))

		magError := int(100 * math.Abs(avgMag-disasterMag) / disasterMag)                    // percentage error
		turnError := int(100 * math.Abs(float64((uint(avgTurn)-disasterTurn)/disasterTurn))) // percentage error

		predError := int(avgConf) * (magError + turnError)

		predConf := 100 - setLimits(predError)

		c.opinionHist[island].Performances["DisasterPred"] = ExpectationReality{
			real: predConf,
		}

		c.confidenceRestrospect("DisasterPred", island)

	}

}

package team5

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// ---------------- Request resource from the President -----------------------------//
// This function asking permission from the President to take resource from the commonpool legally
// The President will reply with an allocation amount
func (c *client) CommonPoolResourceRequest() shared.Resources {
	// c.evaluateRoles()
	// Initially, request the minimum
	turn := c.getTurn()
	season := c.getSeason()
	currentCP := c.getCP()
	currentTier := c.wealth()
	reqAmount := c.calculateRequestToPresident(turn, season, currentTier, currentCP)

	c.Logf("Submitting CP resource request of %v resources", reqAmount)
	//Update request history
	c.cpRequestHistory[turn] = reqAmount
	return reqAmount
}

// calculate the amount of resource to request from the president
func (c *client) calculateRequestToPresident(turn uint, season uint, currentTier wealthTier, currentCP shared.Resources) (reqAmount shared.Resources) {
	if turn == 1 && season == 1 {
		reqAmount = c.config.imperialThreshold
	} else if currentTier == imperialStudent || currentTier == dying {
		// If we are as poor as imperial student, request more resource from cp (whichever number is higher)
		reqAmount = shared.Resources(math.Max(float64(c.config.imperialThreshold), float64(currentCP/6)))
	} else if turn > 1 {
		// For other scenarios, look at the history and make decisions based on that
		lastAllocation := c.cpAllocationHistory[turn-1]
		if lastAllocation == 0 {
			reqAmount = c.config.imperialThreshold
		} else {
			reqAmount = lastAllocation + 10
		}
	}
	return
}

//--------------------------------------------------------------------------------------//

//----------------------Take away resource from the CP --------------------------------//
// Take resource from the common pool ideally from the allocation given by President
// If we are in imperial state, we may take resources regardless what the Presdient say (steal mode)
func (c *client) RequestAllocation() shared.Resources {
	// calculate the amount to take away from common pool
	currentCP := c.getCP()
	currentTier := c.wealth()
	turn := c.getTurn()
	allocationMade := c.BaseClient.LocalVariableCache[rules.AllocationMade].Values[0] != 0
	allocationAmount := shared.Resources(c.BaseClient.LocalVariableCache[rules.ExpectedAllocation].Values[0])
	allocation := c.calculateAllocationFromCP(currentTier, currentCP, allocationMade, allocationAmount)
	if float64(allocation) < 0 {
		allocation = 0
	}
	// update Allocation hisotry
	c.cpAllocationHistory[c.getTurn()] = allocationAmount

	// update opinion on President based on allocation
	presidentID := c.gameState().PresidentID
	if allocationAmount <= 0 { // if no allocation, then minus
		c.opinions[presidentID].updateOpinion(generalBasis, c.changeOpinion(-0.1))
	} else if allocationAmount < c.cpRequestHistory[turn] && allocationAmount >= currentCP/6 { // if some allocation, then just a bit of score
		c.opinions[presidentID].updateOpinion(generalBasis, c.changeOpinion(0.05))
	} else if allocationAmount >= c.cpRequestHistory[turn] {
		c.opinions[presidentID].updateOpinion(generalBasis, c.changeOpinion(0.1))
	}

	//Debug
	c.Logf("Current CP: %v", c.getCP())
	c.Logf("Taking %v from common pool out of %v permitted", allocation, allocationAmount)
	return allocation
}

// calculate the amount to take away from the common pool
func (c *client) calculateAllocationFromCP(currentTier wealthTier, currentCP shared.Resources, allocationMade bool, allocationAmount shared.Resources) (allocation shared.Resources) {
	if currentTier == imperialStudent || currentTier == dying {
		if c.config.imperialThreshold < (currentCP / 6) {
			if (currentCP / 6) > allocationAmount {
				allocation = currentCP / 6
			} else {
				allocation = allocationAmount
			}
		} else {
			if c.config.imperialThreshold > allocationAmount {
				allocation = c.config.imperialThreshold
			} else {
				allocation = allocationAmount
			}
		}
	} else {
		if allocationMade {
			allocation = allocationAmount
		}

	}
	return allocation
}

//------------------------------------------------------------------------------------------//

//---------------------------Pay tax and other contribution to CP --------------------------//
// Pay tax and contribution to the common pool

func (c *client) GetTaxContribution() shared.Resources {
	turn := c.getTurn()
	season := c.getSeason()
	currentTier := c.wealth()
	expectedTax, taxDecisionMade := c.expectedTaxContribution()
	contribution := c.calculateCPContribution(turn, season, currentTier) //CP contribution
	currentResource := c.gameState().ClientInfo.Resources
	geography := c.gameState().Geography
	lastDisaster := c.disasterHistory.getLastDisasterTurn()
	disasterMitigation := c.calculateDisasterContributionCP(turn, currentResource, geography, lastDisaster)

	c.Logf("Disaster Prediction: %v", c.forecastHistory[turn-1])
	// only contribute to the CP if tax decision isn't made
	if !taxDecisionMade {
		return contribution + disasterMitigation
	}
	actualTax := calculateTaxContribution(expectedTax, turn, season, currentTier)
	// c.Logf("%v+v", c.getCP())
	c.Logf("Team 5 paying tax %v out of %v, contribute %v for community, %v for disaster mitigation", actualTax, expectedTax, contribution, disasterMitigation)
	return actualTax + contribution + disasterMitigation
}

// get tax info
func (c *client) expectedTaxContribution() (shared.Resources, bool) {
	expectedTax := shared.Resources(c.BaseClient.LocalVariableCache[rules.ExpectedTaxContribution].Values[0])
	taxDecisionMade := c.BaseClient.LocalVariableCache[rules.TaxDecisionMade].Values[0] != 0
	return expectedTax, taxDecisionMade //TODO: not sure if this is correct tax amount to use
}

// Calculate how much tax to pay according to our wealth state
func calculateTaxContribution(expectedTax shared.Resources, turn uint, season uint, currentTier wealthTier) (contribution shared.Resources) {
	if currentTier == imperialStudent || currentTier == dying {
		return 0
	}
	return expectedTax
}

// Calculate our contribution to common pool based on CP's cashflow
func (c *client) calculateCPContribution(turn uint, season uint, currentTier wealthTier) (contribution shared.Resources) {
	// Day 1 we don't contribute anything
	if turn == 1 && season == 1 {
		contribution = 0
	} else {
		if currentTier != imperialStudent && currentTier != dying {
			// other days we contribute based on cashflow of commonpool
			difference := c.cpResourceHistory[turn] - c.cpResourceHistory[turn-1]
			c.Logf("CP Cashflow: %v - %v = %v", c.cpResourceHistory[turn], c.cpResourceHistory[turn-1], difference)
			if difference < 0 {
				contribution = 0
			} else {
				contribution = difference / 6
			}
		} else {
			contribution = 0
		}

	}
	return contribution
}

// Contribution to CP to prepare for disaster
// Calculate the effect of disaster on our islands using the predicted epicenters
// The contribution amount is based on our confidence intervals, financial status, and predicted period
// Since we don't want our contribution in CP to get taken away by the mob, we will contribute on the last day & the day of the incident
// The day of the incident we can contribute because IIGO happens before disaster.
func (c *client) calculateDisasterContributionCP(currentTurn uint, currentResource shared.Resources, geography disasters.ArchipelagoGeography, lastDisaster uint) shared.Resources {
	if len(c.forecastHistory) == 0 {
		return 0
	}
	predictionInfo := c.forecastHistory[currentTurn-1] // turn - 1 because IIFO happens before IIGO
	idealContribution := shared.Resources(0)
	contribution := shared.Resources(0)
	//how far is the epicenter from our island
	minDistance := math.Sqrt(math.Pow(geography.Islands[shared.Team5].X-predictionInfo.epiX, 2) + math.Pow(geography.Islands[shared.Team5].Y-predictionInfo.epiY, 2))
	if minDistance < 1 {
		minDistance = 1
	}

	//ideal contribution based on disaster magnitude
	if predictionInfo.confidence <= 0.5 && predictionInfo.confidence > 0.2 {
		idealContribution = shared.Resources((float64(predictionInfo.mag) * 50) / minDistance)
	} else if predictionInfo.confidence > 0.5 && predictionInfo.confidence < 0.8 {
		idealContribution = shared.
			Resources((float64(predictionInfo.mag) * 100) / minDistance)
		c.Logf("prediction mag: %v, minDisatance: %v idealContribution %v", predictionInfo.mag, minDistance, idealContribution)
	} else if predictionInfo.confidence >= 0.8 {
		idealContribution = shared.Resources((float64(predictionInfo.mag) * 200) / minDistance)
	}

	//update idealContribution based on our financial status
	// if currentResource < idealContribution {
	// 	idealContribution = shared.Resources(0.5) * currentResource
	// }

	//Only contribute if disaster is coming in 1 days
	// Contribute mostly on the day of the disaster to prevent people taking CP resource
	if currentTurn-lastDisaster == predictionInfo.period-1 {
		contribution = shared.Resources(0.2) * idealContribution
	} else if currentTurn-lastDisaster == predictionInfo.period {
		contribution = shared.Resources(0.8) * idealContribution
	} else {
		contribution = 0
		return contribution
	}
	if currentResource < contribution {
		contribution = shared.Resources(0.5) * currentResource
	}
	return contribution
}

//------------------------------------------------------------------------------------------//

//-----------------------------------Monitor IIGO------------------------------------------//
// Only monitor the foes
func (c *client) MonitorIIGORole(roleName shared.Role) bool {
	var roleID shared.ClientID
	if roleName == shared.Speaker {
		roleID = c.gameState().SpeakerID
	} else if roleName == shared.Judge {
		roleID = c.gameState().JudgeID
	} else {
		roleID = c.gameState().PresidentID
	}
	if roleID == shared.Team5 {
		return false
	} else if c.opinions[roleID].getScore() > 0.5 {
		return false
	}
	return true
}

//DecideIIGOMonitoringAnnouncement decides whether to share the result of monitoring a role and what result to share
//COMPULSORY: must be implemented
// always broadcast monitoring result
func (c *client) DecideIIGOMonitoringAnnouncement(monitoringResult bool) (resultToShare bool, announce bool) {
	resultToShare = monitoringResult
	announce = true
	return
}

func (c *client) GetSanctionPayment() shared.Resources {
	return c.BaseClient.GetSanctionPayment()
}

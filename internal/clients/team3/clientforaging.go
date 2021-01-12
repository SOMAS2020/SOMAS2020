package team3

import (
	"math"
	// "github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	// "github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	// "github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

/*
	DecideForage() (shared.ForageDecision, error)
	ForageUpdate(shared.ForageDecision, shared.Resources)
*/

func (c *client) DecideForage() (shared.ForageDecision, error) {

	// No risk -> minimum is 2 times the critical threshold, Full risk -> minimum is the critical threshold
	safetyFactor := 2.0 - c.params.riskFactor

	//we want to have more than the critical threshold leftover after foraging
	var minimumLeftoverResources = float64(c.criticalThreshold) * safetyFactor

	var foragingInvestment = 0.0
	//for now we invest everything we can, because foraging is iffy.
	if c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus == shared.Alive {
		foragingInvestment = math.Max(float64(c.ServerReadHandle.GetGameState().ClientInfo.Resources)-minimumLeftoverResources, 0)
	}

	c.clientPrint("Foraging investment is %v", foragingInvestment*c.params.riskFactor)

	var forageType shared.ForageType

	fishingROI := c.computeRecentExpectedROI(shared.FishForageType)
	deerHuntingROI := c.computeRecentExpectedROI(shared.DeerForageType)
	if deerHuntingROI != 0 && fishingROI != 0 || (deerHuntingROI > 100 || fishingROI > 100) {
		if deerHuntingROI > fishingROI {
			forageType = shared.DeerForageType
		} else {
			forageType = shared.FishForageType
		}
	} else {
		if fishingROI == 0 {
			forageType = shared.FishForageType
		}
		if deerHuntingROI == 0 {
			forageType = shared.DeerForageType
		}
	}

	if len(c.getAliveIslands()) == 1 {
		forageType = shared.FishForageType
	}

	coef := c.params.riskFactor
	var decay float64
	var sumOfCaught uint
	var numberOfHunters uint
	for _, forage := range c.forageData[forageType] {
		if uint(forage.turn) == c.ServerReadHandle.GetGameState().Turn-1 || uint(forage.turn) == c.ServerReadHandle.GetGameState().Turn-2 {
			sumOfCaught += forage.caught
			numberOfHunters++
		}
	}

	if numberOfHunters == 0 {
		decay = 0
	} else {
		decay = (float64(sumOfCaught) + 1) * float64(numberOfHunters) * 0.5
	}

	// A bit arbitrary
	if decay == 0 {
		coef = 0.1
	} else {
		coef = math.Max(c.params.riskFactor-0.01*(1-c.params.riskFactor)/4*decay, 0.1)
	}

	finalForagingInvestment := foragingInvestment * coef

	if c.getLocalResources() < c.minimumResourcesWeWant || c.computeRecentExpectedROI(forageType) < 100 {
		finalForagingInvestment = 0.01
	}

	c.Logf("coef: %v, sumOfCaught: %v, numberOfHunter: %v, decay: %v", coef, sumOfCaught, numberOfHunters, decay)

	return shared.ForageDecision{
		Type:         forageType,
		Contribution: shared.Resources(finalForagingInvestment),
	}, nil
}

func (c *client) computeRecentExpectedROI(forageType shared.ForageType) float64 {
	var sumOfROI float64
	var numberOfROI uint

	for _, forage := range c.forageData[forageType] {
		if uint(forage.turn) == c.ServerReadHandle.GetGameState().Turn-1 || uint(forage.turn) == c.ServerReadHandle.GetGameState().Turn-2 {
			if forage.amountContributed != 0 {
				sumOfROI += float64((forage.amountReturned / forage.amountContributed) * 100)
				numberOfROI++
			}
		}
	}

	if numberOfROI == 0 {
		return 0
	}

	c.Logf("Expected return of %v: %v per cent", forageType, (sumOfROI / float64(numberOfROI)))
	return sumOfROI / float64(numberOfROI)
}

func (c *client) ForageUpdate(forageDecision shared.ForageDecision, outcome shared.Resources, numberCaught uint) {
	c.forageData[forageDecision.Type] =
		append(
			c.forageData[forageDecision.Type],
			ForageData{
				amountContributed: forageDecision.Contribution,
				amountReturned:    outcome,
				turn:              c.ServerReadHandle.GetGameState().Turn,
				caught:            numberCaught,
			},
		)
}

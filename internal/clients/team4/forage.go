package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type forageStorage struct {
	forageHistory      []forageHistory
	receivedForageData [][]shared.ForageShareInfo
}

type forageHistory struct {
	decision       shared.ForageDecision
	resourceReturn shared.Resources
	numberCaught   uint
}

func (c *client) analyseHistory() shared.ForageType {
	constLookBack := 5
	if c.getTurn() < 5 {
		return shared.DeerForageType
	}

	totalResources := make(map[shared.ForageType]float64)
	totalResources[shared.DeerForageType] = 1.1
	totalResources[shared.FishForageType] = 1

	if len(c.forage.forageHistory) < constLookBack {
		constLookBack = len(c.forage.forageHistory)
	}
	for _, e := range c.forage.forageHistory[len(c.forage.forageHistory)-constLookBack:] {
		var ratio float64
		if e.decision.Contribution <= 0 {
			ratio = 1
		} else {
			ratio = float64(e.resourceReturn) / float64(e.decision.Contribution)
		}
		totalResources[e.decision.Type] = totalResources[e.decision.Type] * ratio
	}

	if len(c.forage.receivedForageData) < constLookBack {
		constLookBack = len(c.forage.receivedForageData)
	}
	for _, e := range c.forage.receivedForageData[len(c.forage.receivedForageData)-constLookBack:] {
		for _, teamEntry := range e {
			var ratio float64
			if teamEntry.DecisionMade.Contribution <= 0 {
				ratio = 1
			} else {
				ratio = float64(teamEntry.ResourceObtained) / float64(teamEntry.DecisionMade.Contribution)
			}
			if ratio > 2 {
				ratio = 2
			}
			totalResources[teamEntry.DecisionMade.Type] = totalResources[teamEntry.DecisionMade.Type] * ratio
		}
	}
	c.Logf("Foraging: FishingRatio %v, DeerRatio %v", totalResources[shared.DeerForageType], totalResources[shared.FishForageType])
	if totalResources[shared.DeerForageType] >= totalResources[shared.FishForageType] {
		return shared.DeerForageType
	}
	return shared.FishForageType
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	ft := c.analyseHistory()
	scale := 5 * c.getSafeResourceLevel()
	forageContribution := c.getResources() - (2-shared.Resources(c.internalParam.riskTaking))*scale

	if c.getSafeResourceLevel()*2 > c.getResources()-forageContribution {
		forageContribution = c.getResources() - c.getSafeResourceLevel()*2
	}
	if forageContribution < 0 {
		forageContribution = 0
	}
	c.Logf("Foraging: Decision %v, Resources %v", ft, forageContribution)
	return shared.ForageDecision{
		Type:         shared.ForageType(ft),
		Contribution: shared.Resources(forageContribution),
	}, nil
}

// ForageUpdate is called by the server upon completion of a foraging session. This handler can be used by clients to
// analyse their returns - resources returned to them, as well as number of fish/deer caught.
func (c *client) ForageUpdate(initialDecision shared.ForageDecision, resourceReturn shared.Resources, numberCaught uint) {
	c.forage.forageHistory = append(c.forage.forageHistory, forageHistory{
		decision:       initialDecision,
		resourceReturn: resourceReturn,
		numberCaught:   numberCaught,
	})
}

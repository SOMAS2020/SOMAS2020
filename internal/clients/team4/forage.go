package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type forageStorage struct {
	preferedForageMethod shared.ForageType
	forageHistory        []forageHistory
	receivedForageData   [][]shared.ForageShareInfo
	turnsSinceChange     int
}

type forageHistory struct {
	decision       shared.ForageDecision
	resourceReturn shared.Resources
	numberCaught   uint
}

func (c *client) analyseHistory() {
	constLookBack := 5

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
			totalResources[teamEntry.DecisionMade.Type] = totalResources[teamEntry.DecisionMade.Type] * ratio
		}
		if totalResources[shared.DeerForageType] >= totalResources[shared.FishForageType] {
			c.forage.preferedForageMethod = shared.DeerForageType
		} else {
			c.forage.preferedForageMethod = shared.DeerForageType
		}
	}
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	c.analyseHistory()
	ft := c.forage.preferedForageMethod
	scale := 5 * c.getSafeResourceLevel()
	resources := c.getResources() - c.getSafeResourceLevel()*(2-shared.Resources(c.internalParam.riskTaking)*scale)
	return shared.ForageDecision{
		Type:         shared.ForageType(ft),
		Contribution: shared.Resources(resources),
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

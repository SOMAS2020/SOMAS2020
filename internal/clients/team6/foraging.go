package team6

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// ForageHistory stores our forage history
// type ForageHistory map[shared.ForageType][]ForageResults

// type ForageResults struct {
// 	forageIn     shared.Resources
// 	forageReturn shared.Resources
// }
func (result ForageResults) calcROI() float64 {
	if result.forageIn == 0 {
		return 0
	}
	return float64(result.forageReturn/result.forageIn) - 1
}

func (c *client) changeForageType() shared.ForageType {
	var deerRoiTotal float64
	var fishRoiTotal float64
	var deerParticipant uint = 0
	var fishParticipant uint = 0

	var deerRoiTotal2 float64
	var fishRoiTotal2 float64
	var deerParticipant2 uint = 0
	var fishParticipant2 uint = 0

	var deerAverageRoi, deerAverageRoi2, fishAverageRoi, fishAverageRoi2 float64
	for _, deerResults := range c.forageHistory[shared.DeerForageType] {
		if deerResults.turn == c.ServerReadHandle.GetGameState().Turn-1 {
			deerRoiTotal += deerResults.calcROI()
			deerParticipant++
		}
		if deerResults.turn == c.ServerReadHandle.GetGameState().Turn-2 {
			deerRoiTotal2 += deerResults.calcROI()
			deerParticipant2++
		}
	}
	deerAverageRoi = deerRoiTotal / float64(deerParticipant)
	deerAverageRoi2 = deerRoiTotal2 / float64(deerParticipant2)

	for _, fishResults := range c.forageHistory[shared.FishForageType] {
		if fishResults.turn == c.ServerReadHandle.GetGameState().Turn-1 {
			fishRoiTotal += fishResults.calcROI()
			fishParticipant++
		}
		if fishResults.turn == c.ServerReadHandle.GetGameState().Turn-2 {
			fishRoiTotal2 += fishResults.calcROI()
			fishParticipant2++
		}
	}

	fishAverageRoi = fishRoiTotal / float64(fishParticipant)
	fishAverageRoi2 = fishRoiTotal2 / float64(fishParticipant2)

	if fishAverageRoi < deerAverageRoi {
		if deerAverageRoi < deerAverageRoi2 {
			c.clientConfig.multiplier -= 0.03
		}
		if deerAverageRoi > deerAverageRoi2 {
			c.clientConfig.multiplier += 0.03
		}
		return shared.DeerForageType
	}

	if fishAverageRoi < fishAverageRoi2 {
		c.clientConfig.multiplier -= 0.05
	}
	if fishAverageRoi > fishAverageRoi2 {
		c.clientConfig.multiplier += 0.05
	}
	return shared.DeerForageType
}

func (c *client) decideContribution() shared.Resources {

	var safetyBuffer shared.Resources = 10
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	return shared.Resources(c.clientConfig.multiplier) * (ourResources - safetyBuffer)
}

func (c *client) randomForage() shared.ForageDecision {
	var resources shared.Resources
	var forageType shared.ForageType

	if c.ServerReadHandle.GetGameState().Turn == 2 {
		forageType = shared.FishForageType
	} else {
		forageType = shared.DeerForageType
	}
	tmp := rand.Float64()
	if tmp > 0.2 { //up to 20% resources
		resources = 0.2 * c.ServerReadHandle.GetGameState().ClientInfo.Resources
	} else {
		resources = shared.Resources(tmp) * c.ServerReadHandle.GetGameState().ClientInfo.Resources
	}

	return shared.ForageDecision{
		Type:         shared.ForageType(forageType),
		Contribution: shared.Resources(resources),
	}
}

func (c *client) noramlForage() shared.ForageDecision {
	ft := c.changeForageType()
	amt := c.decideContribution()
	return shared.ForageDecision{
		Type:         shared.ForageType(ft),
		Contribution: shared.Resources(amt),
	}
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	if c.ServerReadHandle.GetGameState().Turn < 3 { //the agent will randomly forage at the start
		return c.randomForage(), nil
	}
	return c.noramlForage(), nil

}

func (c *client) ForageUpdate(forageDecision shared.ForageDecision, outcome shared.Resources, numberCaught uint) {
	currTurn := c.ServerReadHandle.GetGameState().Turn

	c.forageHistory[forageDecision.Type] =
		append(
			c.forageHistory[forageDecision.Type],
			ForageResults{
				forageIn:     forageDecision.Contribution,
				forageReturn: outcome,
				turn:         currTurn,
			},
		)

	c.Logf(
		"Forage History Updated: Type %v ,Conribution: %v ,Return: %v",
		forageDecision.Type,
		forageDecision.Contribution,
		outcome,
	)
}

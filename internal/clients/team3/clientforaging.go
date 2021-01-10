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

	safetyFactor := 1.0 + (0.5/100)*c.params.riskFactor

	//we want to have more than the critical threshold leftover after foraging
	var minimumLeftoverResources = float64(c.criticalThreshold) * safetyFactor

	var foragingInvestment = 0.0
	//for now we invest everything we can, because foraging is iffy.
	if c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus == shared.Alive {
		foragingInvestment = math.Max(float64(c.ServerReadHandle.GetGameState().ClientInfo.Resources)-minimumLeftoverResources, 0)
	}

	// ## When fishing foraging is implemnented ##

	// deerForagingInvestment := p.riskFactor * foragingInvestment
	// deerForagingInvestment := (1 - p.riskFactor) * foragingInvestment

	// math.Max(deerForagingInvestment, p.minimunInvestment)
	// math.Max(fishForagingInvestment, p.minimunInvestment)

	return shared.ForageDecision{
		Type:         shared.DeerForageType,
		Contribution: shared.Resources(foragingInvestment),
	}, nil
}

// Increment p.minimumInvestment when receiving forage update and no deer/fish was caughgt.

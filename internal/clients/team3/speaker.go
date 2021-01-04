package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type speaker struct {
	// Base implementation
	*baseclient.BaseSpeaker
	// Our client
	c *client
}

// Override functions here, see president.go for examples

func (s *speaker) PayJudge(salary shared.Resources) (shared.Resources, bool) {
	// Use the base implementation
	return s.BaseSpeaker.PayJudge(salary)
}

func (s *speaker) DecideForage() {
	
	safetyFactor := 1.0 + (0.5/100)*c.params.riskFactor

	//we want to have more than the critical threshold leftover after foraging
	minimumLeftoverResources := criticalThreshold * safetyFactor

	//for now we invest everything we can, because foraging is iffy.
	if lifeStatus != critical :
		foragingInvestment = ourResources - minimumLeftoverResources
}

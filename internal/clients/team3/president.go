package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type president struct {
	// Base implementation
	*baseclient.BasePresident
	// Our client
	c *client
}

func (p *president) PaySpeaker(salary shared.Resources) (shared.Resources, bool) {
	// Use the base implementation
	return p.BasePresident.PaySpeaker(salary)
}

func (p *president) DecideNextSpeaker(winner shared.ClientID) shared.ClientID {
	p.c.clientPrint("choosing speaker")
	// Naively choose group 0
	return mostTrusted(p.c.trustMapAgg)
}

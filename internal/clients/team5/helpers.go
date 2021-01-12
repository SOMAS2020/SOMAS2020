package team5

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// shorthand to get current turn as it's needed often
func (c client) getTurn() uint {
	return c.gameState().Turn
}

// shorthand to get current turn as it's needed often
func (c client) getSeason() uint {
	return c.gameState().Season
}

// shorthand to get current turn as it's needed often
func (c client) getCP() shared.Resources {
	return c.gameState().CommonPool
}

// shorthand to get our current life status
func (c client) getLifeStatus() shared.ClientLifeStatus {
	return c.gameState().ClientInfo.LifeStatus
}

func (c client) getGameConfig() config.ClientConfig {
	return c.ServerReadHandle.GetGameConfig()
}

func (c client) getAliveTeams(includeUs bool) (aliveTeams []shared.ClientID) {
	for team, status := range c.gameState().ClientLifeStatuses {
		if status == shared.Alive {
			if includeUs || team != c.GetID() {
				aliveTeams = append(aliveTeams, team)
			}
		}
	}
	return aliveTeams
}

// checks if a given client is alive
func (c client) isClientAlive(id shared.ClientID) bool {
	for _, cl := range c.getAliveTeams(false) {
		if cl == id {
			return true
		}
	}
	return false
}

func roundTo(x float64, decPlaces uint) float64 {
	x *= math.Pow(10, float64(decPlaces))
	y := math.Round(x)
	return y / math.Pow(10, float64(decPlaces))
}

// caps magnitude of val to absThresh
func absoluteCap(val, absThresh float64) float64 {
	if val > 0 {
		return math.Min(val, absThresh)
	}
	return math.Max(val, absThresh*-1)
}

func uintsAsFloats(x []uint) []float64 {
	out := make([]float64, len(x))
	for i, el := range x {
		out[i] = float64(el)
	}
	return out
}

func floatsAsUints(x []float64) []uint {
	out := make([]uint, len(x))
	for i, el := range x {
		out[i] = uint(el)
	}
	return out
}

func (c client) getMood() float64 {
	if c.gameState().ClientInfo.Resources >= getClientConfig().jbThreshold {
		return mapToRange(float64(c.gameState().ClientInfo.Resources),
			float64(c.gameState().ClientInfo.Resources), 0, 0.5, 1.5)
	}
	return mapToRange(float64(c.gameState().ClientInfo.Resources),
		float64(c.config.jbThreshold), 0, 0.5, 1.5)
}

func mapToRange(x, inMin, inMax, outMin, outMax float64) float64 {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

// changeOpinion true = positive , false = negative
func (c *client) changeOpinion(opinionChange float64) float64 {
	switch c.config.agentMentality {
	case okBoomer: // Strict opinion (greedy)
		if opinionChange >= 0 { // positive case
			opinionChange = opinionChange * 0.1 * c.getMood() // less emphasis on positive
		} else {
			opinionChange = opinionChange * 2 * c.getMood() // more emphasis on negative
		}
	case millennial: // You get a positive opinion, You get a positive opinion, everyone gets a positive opinion
		if opinionChange >= 0 { // positive case
			opinionChange = opinionChange * 1.5 * c.getMood() // more emphasis on positive
		} else {
			opinionChange = opinionChange * 0.5 * c.getMood() // less emphasis on negative
		}
	}

	return opinionChange
}

package team5

import (
	"math"
	"sort"

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

func (c client) getAliveTeams(includeUs bool) (aliveTeams []shared.ClientID) {
	for team, status := range c.gameState().ClientLifeStatuses {
		if status == shared.Alive {
			if includeUs || team != ourClientID {
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

func minMaxCap(val, absThresh float64) float64 {
	if val > 0 {
		return math.Min(val, absThresh)
	}
	return math.Max(val, absThresh*-1)
}

// returns a slice of sorted keys for maps with uint keys
func sortUintSlice(m []uint) []uint {
	keys := make([]int, len(m))
	for k := range m {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	finalKeys := make([]uint, len(m))
	for k := range keys {
		finalKeys = append(finalKeys, uint(k))
	}
	return finalKeys
}

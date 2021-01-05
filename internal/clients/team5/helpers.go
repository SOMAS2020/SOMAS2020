package team5

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

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

func (c client) getAliveTeams(includeUs bool) (aliveTeams []shared.ClientID) {
	for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
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

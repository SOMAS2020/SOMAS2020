package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// increases the friendship level with some other islands
func (c *client) raiseFriendshipLevel(clientID shared.ClientID, increment FriendshipLevel) {
	currFriendship := c.friendship[clientID]
	logIncrement := (increment / (c.friendship[clientID] + increment)) * (c.config.maxFriendship / 5)
	raisedFriendship := currFriendship + logIncrement

	if raisedFriendship > c.config.maxFriendship {
		// c.Logf("Friendship with island[%v] is at maximum!", clientID)
		c.friendship[clientID] = c.config.maxFriendship
	} else {
		c.friendship[clientID] = raisedFriendship
	}
}

// decreases the friendship level with some other islands
func (c *client) lowerFriendshipLevel(clientID shared.ClientID, deduction FriendshipLevel) {
	currFriendship := c.friendship[clientID]
	logDeduction := (deduction / (c.friendship[clientID] + deduction)) * (c.config.maxFriendship / 5)
	loweredFriendship := currFriendship - logDeduction

	if loweredFriendship < c.config.minFriendship {
		// c.Logf("Friendship with island[%v] is at minimum!", clientID)
		c.friendship[clientID] = c.config.minFriendship
	} else {
		c.friendship[clientID] = loweredFriendship
	}
}

// gets friendship cofficients
func (c client) getFriendshipCoeffs() map[shared.ClientID]float64 {
	friendshipCoeffs := make(map[shared.ClientID]float64)

	for team, fs := range c.friendship {
		friendshipCoeffs[team] = float64(fs / c.config.maxFriendship)
	}

	return friendshipCoeffs
}

// gets our personality
func (c client) getPersonality() Personality {
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources

	if ourResources <= c.config.selfishThreshold {
		return Selfish
	} else if ourResources <= c.config.normalThreshold {
		return Normal
	}

	// TODO: more cases to have?
	return Generous
}

// gets the number of alive islands
func (c client) getNumOfAliveIslands() uint {
	num := uint(0)

	for _, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		if status == shared.Alive {
			num++
		}
	}

	return num
}

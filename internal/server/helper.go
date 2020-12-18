package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

func getClientInfosAndMapFromRegisteredClients(
	registeredClients map[shared.ClientID]baseclient.Client,
) (map[shared.ClientID]gamestate.ClientInfo, map[shared.ClientID]baseclient.Client) {
	clientInfos := map[shared.ClientID]gamestate.ClientInfo{}
	clientMap := map[shared.ClientID]baseclient.Client{}

	for id, c := range registeredClients {
		clientInfos[id] = gamestate.ClientInfo{
			Resources:  config.InitialResources,
			LifeStatus: shared.Alive,
		}
		clientMap[id] = c
	}

	return clientInfos, clientMap
}

// anyClientsAlive returns true if any one client is Alive (including critical).
func anyClientsAlive(clientInfos map[shared.ClientID]gamestate.ClientInfo) bool {
	for _, ci := range clientInfos {
		if ci.LifeStatus != shared.Dead {
			return true
		}
	}
	return false
}

// updateIslandLivingStatusForClient returns an updated copy of the clientInfo after updating
// the Alive, Critical, and CriticalConsecutiveTurnsLeft attribs according to the resource levels and
// the game's configuration.
func updateIslandLivingStatusForClient(
	ci gamestate.ClientInfo,
	minimumResourceThreshold int,
	maxCriticalConsecutiveTurns uint,
) (gamestate.ClientInfo, error) {
	switch ci.LifeStatus {
	case shared.Alive:
		if ci.Resources < minimumResourceThreshold {
			ci.LifeStatus = shared.Critical
			ci.CriticalConsecutiveTurnsCounter = 0
		}
		return ci, nil

	case shared.Critical:
		if ci.Resources < minimumResourceThreshold {
			if ci.CriticalConsecutiveTurnsCounter == maxCriticalConsecutiveTurns {
				ci.LifeStatus = shared.Dead
			} else {
				ci.CriticalConsecutiveTurnsCounter++
			}
			return ci, nil
		}
		ci.LifeStatus = shared.Alive
		ci.CriticalConsecutiveTurnsCounter = 0
		return ci, nil

	case shared.Dead:
		// dead clients are not resurrected
		return ci, nil

	default:
		return ci,
			errors.Errorf("updateIslandLivingStatusForClient not implemented for LifeStatus %v",
				ci.LifeStatus)
	}
}

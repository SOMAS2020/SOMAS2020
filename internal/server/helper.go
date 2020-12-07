package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func getClientInfoFromRegisteredClients(registeredClients map[shared.ClientID]common.Client) map[shared.ClientID]common.ClientInfo {
	clientInfos := map[shared.ClientID]common.ClientInfo{}

	for id, c := range registeredClients {
		clientInfos[id] = common.ClientInfo{
			Client:                       c,
			Resources:                    config.InitialResources,
			Alive:                        true,
			Critical:                     false,
			CriticalConsecutiveTurnsLeft: config.MaxCriticalConsecutiveTurns,
		}
	}

	return clientInfos
}

// anyClientsAlive returns true if any one client is Alive (including critical).
func anyClientsAlive(clientInfos map[shared.ClientID]common.ClientInfo) bool {
	for _, ci := range clientInfos {
		if ci.Alive {
			return true
		}
	}
	return false
}

// updateIslandLivingStatusForClient returns an updated copy of the clientInfo after updating
// the Alive, Critical, and CriticalConsecutiveTurnsLeft attribs according to the resource levels and
// the game's configuration.
func updateIslandLivingStatusForClient(ci common.ClientInfo, minimumResourceThreshold int, MaxCriticalConsecutiveTurns uint) common.ClientInfo {
	// we don't resurrect dead clients!
	if !ci.Alive {
		return ci
	}

	if ci.Resources < minimumResourceThreshold {
		if !ci.Critical {
			ci.Critical = true
			ci.CriticalConsecutiveTurnsLeft = MaxCriticalConsecutiveTurns
		} else {
			// Critical!
			if ci.CriticalConsecutiveTurnsLeft == 0 {
				// RIP!
				ci.Alive = false
			} else {
				ci.CriticalConsecutiveTurnsLeft--
			}
		}
	} else {
		// have got above threshold!
		if ci.Critical {
			ci.Critical = false
			ci.CriticalConsecutiveTurnsLeft = MaxCriticalConsecutiveTurns
		}
	}

	return ci
}

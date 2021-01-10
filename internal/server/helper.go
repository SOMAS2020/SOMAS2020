package server

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

type clientInfoUpdateResult struct {
	ID  shared.ClientID
	Ci  gamestate.ClientInfo
	Err error
}

func getClientInfosAndMapFromRegisteredClients(
	registeredClients map[shared.ClientID]baseclient.Client,
	initialResources shared.Resources,
) (map[shared.ClientID]gamestate.ClientInfo, map[shared.ClientID]baseclient.Client) {
	clientInfos := map[shared.ClientID]gamestate.ClientInfo{}
	clientMap := map[shared.ClientID]baseclient.Client{}

	for id, c := range registeredClients {
		clientInfos[id] = gamestate.ClientInfo{
			Resources:  initialResources,
			LifeStatus: shared.Alive,
		}
		clientMap[id] = c
	}

	return clientInfos, clientMap
}

// anyClientsAlive returns true if any one client is Alive (including critical).
func anyClientsAlive(clientInfos map[shared.ClientID]gamestate.ClientInfo) bool {
	return len(getNonDeadClientIDs(clientInfos)) != 0
}

// getNonDeadClients returns a map of all clients with a non-dead status
func getNonDeadClients(clientInfos map[shared.ClientID]gamestate.ClientInfo,
	clientMap map[shared.ClientID]baseclient.Client) map[shared.ClientID]baseclient.Client {

	nonDeadClientMap := map[shared.ClientID]baseclient.Client{}
	clientIDs := getNonDeadClientIDs(clientInfos)
	for _, id := range clientIDs {
		nonDeadClientMap[id] = clientMap[id]
	}
	return nonDeadClientMap
}

// updateIslandLivingStatusForClient returns an updated copy of the clientInfo after updating
// the Alive, Critical, and CriticalConsecutiveTurnsLeft attribs according to the resource levels and
// the game's configuration.
func updateIslandLivingStatusForClient(
	ci gamestate.ClientInfo,
	minimumResourceThreshold shared.Resources,
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

// getNonDeadClients return ClientIDs of clients that are not dead (alive + critical).
// The result is NOT ordered.
func getNonDeadClientIDs(clientInfos map[shared.ClientID]gamestate.ClientInfo) []shared.ClientID {
	nonDeadClients := []shared.ClientID{}

	for id, ci := range clientInfos {
		if ci.LifeStatus != shared.Dead {
			nonDeadClients = append(nonDeadClients, id)
		}
	}

	return nonDeadClients
}

// giveResources takes resources to client, logging it and mentioning reason
func (s *SOMASServer) takeResources(clientID shared.ClientID, resources shared.Resources, reason string) error {
	s.logf("Trying to take %v from %v (reason: %s)", resources, clientID, reason)
	if math.IsNaN(float64(resources)) || resources < 0 {
		return errors.Errorf("Cannot take invalid number of resources %v from client %v", resources, clientID)
	}

	participantInfo := s.gameState.ClientInfos[clientID]
	if participantInfo.Resources < resources {
		return errors.Errorf(
			"Client %v did not have enough resources. Requested %v, only had %v",
			clientID, resources, participantInfo.Resources,
		)
	}
	participantInfo.Resources -= resources
	s.gameState.ClientInfos[clientID] = participantInfo
	s.logf("Took %v from %v (reason: %s)", resources, clientID, reason)
	return nil
}

// giveResources gives resources to client, logging it and mentioning reason
func (s *SOMASServer) giveResources(clientID shared.ClientID, resources shared.Resources, reason string) error {
	s.logf("Trying to give %v to %v (reason: %s)", resources, clientID, reason)
	if math.IsNaN(float64(resources)) || resources < 0 {
		return errors.Errorf("Cannot give invalid number of resources %v to client %v", resources, clientID)
	}

	participantInfo := s.gameState.ClientInfos[clientID]
	participantInfo.Resources += resources
	s.gameState.ClientInfos[clientID] = participantInfo

	s.logf("Gave %v to %v (reason: %s)", resources, clientID, reason)
	return nil
}

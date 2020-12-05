package server

import (
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func getClientInfoFromRegisteredClients(registeredClients map[shared.ClientID]common.Client) map[shared.ClientID]common.ClientInfo {
	clientInfos := map[shared.ClientID]common.ClientInfo{}

	for id, c := range registeredClients {
		clientInfos[id] = common.ClientInfo{
			Client:    c,
			Resources: common.DefaultResources,
			Alive:     true,
		}
	}

	return clientInfos
}

// anyClientsAlive returns true if any client is alive in the provided ClientInfos map.
func anyClientsAlive(clientInfos map[shared.ClientID]common.ClientInfo) bool {
	for id, ci := range clientInfos {
		if ci.Alive {
			log.Printf("%v", id)
			return true
		}
	}
	return false
}

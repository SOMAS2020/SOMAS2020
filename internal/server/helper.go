package server

import (
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
)

func getClientInfoFromRegisteredClients(registeredClients map[common.ClientID]common.Client) map[common.ClientID]common.ClientInfo {
	clientInfos := map[common.ClientID]common.ClientInfo{}

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
func anyClientsAlive(clientInfos map[common.ClientID]common.ClientInfo) bool {
	for id, ci := range clientInfos {
		if ci.Alive {
			log.Printf("%v", id)
			return true
		}
	}
	return false
}

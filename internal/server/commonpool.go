package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

//islandDeplete depletes island's resource based on the severity of the storm
func (s *SOMASServer) islandDeplete(porportionalEffect map[shared.ClientID]float64) {
	clientMap := getNonDeadClients(s.gameState.ClientInfos, s.clientMap)
	for clientID := range clientMap {
		deduction := shared.Resources(porportionalEffect[clientID]) // min resources = 0
		s.takeResources(clientID, deduction, "disaster damage")
	}
}

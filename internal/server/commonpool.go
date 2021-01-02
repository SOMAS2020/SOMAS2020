package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// islandDeplete depletes island's resource based on the severity of the storm (after CP mitigation)
func (s *SOMASServer) islandDeplete(cpMitigatedEffect map[shared.ClientID]float64) {
	clientMap := getNonDeadClients(s.gameState.ClientInfos, s.clientMap)
	for clientID := range clientMap {
		deduction := shared.Resources(cpMitigatedEffect[clientID]) // min resources = 0
		if deduction > 0 {                                         // don't create pointless call if no deduction applicable
			err := s.takeResources(clientID, deduction, "disaster damage")
			if err != nil {
				s.logf("Error taking resources from %v for disaster damage: %v", clientID, err)
			}
		}
	}
}

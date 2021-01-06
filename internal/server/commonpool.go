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
			ci := s.gameState.ClientInfos[clientID]
			if ci.Resources < deduction {
				ci.Resources = 0
			} else {
				ci.Resources -= deduction
			}
			s.gameState.ClientInfos[id] = ci
			s.logf("[DISASTER]: %v reduced to %v resources due to disaster damage of %v", clientID, ci.Resources, deduction)
		}
	}
}

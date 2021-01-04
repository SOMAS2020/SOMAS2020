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
				s.logf("[DISASTER]: %v reduced to 0 resources due to disaster damage", clientID)
				ci.Resources = 0
			} else {
				s.logf("[DISASTER]: %v reduced by %v resources due to disaster damage", clientID, deduction)
				ci.Resources -= deduction
			}
		}
	}
}

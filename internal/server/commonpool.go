package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

//islandDistribute distributes resources to island as requested
func (s *SOMASServer) islandDistribute(resource shared.Resources, islandID shared.ClientID) {
	ci := s.gameState.ClientInfos[islandID]
	ci.Resources = ci.Resources + resource
	s.gameState.ClientInfos[islandID] = ci
	s.gameState.Environment.CommonPool.Resources -= resource
}

//islandDeplete depletes island's resource based on the severity of the storm
func (s *SOMASServer) islandDeplete(porportionalEffect map[shared.ClientID]float64) {
	for id, ci := range s.gameState.ClientInfos {
		ci.Resources = shared.Resources(float64(ci.Resources) - porportionalEffect[id])
		s.gameState.ClientInfos[id] = ci
	}
}

// islandContribute takes resource donation from individual island into the common-pool
func (s *SOMASServer) islandContribute(resource shared.Resources, islandID shared.ClientID) {
	ci := s.gameState.ClientInfos[islandID]
	ci.Resources = ci.Resources - resource
	s.gameState.ClientInfos[islandID] = ci
	s.gameState.Environment.CommonPool.Resources += shared.Resources(resource)
}

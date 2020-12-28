package server
import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

//islandDistribute distributes resources to island as requested 
func (s *SOMASServer) islandDistribute(resource int, islandID shared.ClientID) { 
	ci := s.gameState.ClientInfos[islandID]
	ci.Resources = ci.Resources + resource 				
	s.gameState.ClientInfos[islandID] = ci
	s.gameState.Environment.CommonPool.Resource -= uint(resource)
}

//islandDeplete depletes island's resource based on the severity of the storm
//Can resource go to negative?
func (s *SOMASServer) islandDeplete(porportionalEffect map[shared.ClientID]float64) {
	temp := 0.0
	for id, ci := range s.gameState.ClientInfos {
		temp = float64(ci.Resources)-porportionalEffect[id] 
		ci.Resources = int(temp)
		s.gameState.ClientInfos[id] = ci
	}
}

// islandContribute takes resource donation from individual island into the common-pool
func (s *SOMASServer) islandContribute(resource int, islandID shared.ClientID) {
	ci := s.gameState.ClientInfos[islandID]
	ci.Resources = ci.Resources - resource 					
	s.gameState.ClientInfos[islandID] = ci
	s.gameState.Environment.CommonPool.Resource += uint(resource)
}
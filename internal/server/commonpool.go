package server
import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)


//islandDistribute distributes resources to island as requested 
func (s *SOMASServer) islandDistribute(resource int, islandID shared.ClientID) { 

	resChan := make(chan clientInfoUpdateResult, 1)

	ci := s.gameState.ClientInfos[islandID]
	ci.Resources += resource 	
	resChan <- clientInfoUpdateResult{
		ID:  islandID,
		Ci:  ci,
		Err: nil,
	}			
	s.gameState.ClientInfos[islandID] = ci
	s.gameState.Environment.CommonPool.Resource -= uint(resource)
}

//islandDeplete depletes island's resource based on the severity of the storm
//Can resource go to negative?
func (s *SOMASServer) islandDeplete(porpotionEffect map[shared.ClientID]float64) {
	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)
	N := len(nonDeadClients)
	resChan := make(chan clientInfoUpdateResult, N)
	temp := 0.0

	for id, ci := range s.gameState.ClientInfos {
		temp = float64(ci.Resources)-porpotionEffect[id] 
		ci.Resources = int(temp)
			resChan <- clientInfoUpdateResult{
			ID:  id,
			Ci:  ci,
			Err: nil,
		}
		s.gameState.ClientInfos[id] = ci
	}

	close(resChan)
	for res := range resChan {
		id, ci := res.ID, res.Ci
		s.gameState.ClientInfos[id] = ci
		
	}
}

// islandContribute takes resource donation from individual island into the common-pool
func (s *SOMASServer) islandContribute(resource int, islandID shared.ClientID) {
	resChan := make(chan clientInfoUpdateResult, 1)

	ci := s.gameState.ClientInfos[islandID]
	ci.Resources -= resource 	
	resChan <- clientInfoUpdateResult{
		ID:  islandID,
		Ci:  ci,
		Err: nil,
	}			
	s.gameState.ClientInfos[islandID] = ci
	s.gameState.Environment.CommonPool.Resource += uint(resource)
}
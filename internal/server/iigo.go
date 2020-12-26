package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/server/iigointernal"
)

// runIIGO : IIGO decides rule changes, elections, sanctions
func (s *SOMASServer) runIIGO() error {
	s.logf("start runIIGO")
	nonDead := getNonDeadClientIDs(s.gameState.ClientInfos)
	updateAliveIslands(nonDead)
	_ = iigointernal.RunIIGO(&s.gameState, &s.clientMap)
	defer s.logf("finish runIIGO")
	return nil
}

func (s *SOMASServer) runIIGOEndOfTurn() error {
	s.logf("start runIIGOEndOfTurn")
	defer s.logf("finish runIIGOEndOfTurn")
	clientMap := getNonDeadClients(s.gameState.ClientInfos, s.clientMap)
	for clientID, v := range clientMap {
		tax := v.GetTaxContribution()
		s.gameState.CommonPool += tax
		newGameState := s.gameState.GetClientGameStateCopy(clientID)
		newGameState.ClientInfo.Resources -= tax
		v.GameStateUpdate(newGameState)
		gamestate.UpdateTurnHistory(clientID, []rules.VariableValuePair{
			{
				VariableName: rules.IslandTaxContribution,
				Values:       []float64{float64(tax)},
			},
			{
				VariableName: rules.ExpectedTaxContribution,
				Values:       []float64{float64(iigointernal.TaxAmountMapExport[clientID])},
			},
		})
	}
	return nil
}

func (s *SOMASServer) runIIGOAllocations() error {
	s.logf("start runIIGOAllocations")
	defer s.logf("finish runIIGOAllocations")
	clientMap := getNonDeadClients(s.gameState.ClientInfos, s.clientMap)
	for clientID, v := range clientMap {
		allocation := v.RequestAllocation()
		s.gameState.CommonPool -= allocation
		newGameState := s.gameState.GetClientGameStateCopy(clientID)
		newGameState.ClientInfo.Resources += allocation
		v.GameStateUpdate(newGameState)
		if allocation <= s.gameState.CommonPool {
			s.gameState.CommonPool -= allocation
			gamestate.UpdateTurnHistory(clientID, []rules.VariableValuePair{
				{
					VariableName: rules.IslandAllocation,
					Values:       []float64{float64(allocation)},
				},
				{
					VariableName: rules.ExpectedAllocation,
					Values:       []float64{float64(iigointernal.AllocationAmountMapExport[clientID])},
				},
			})
		}
	}
	return nil
}

func updateAliveIslands(aliveIslands []shared.ClientID) {
	var forVariables []float64
	for _, v := range aliveIslands {
		forVariables = append(forVariables, float64(v))
	}
	_ = rules.UpdateVariable(rules.IslandsAlive, rules.VariableValuePair{
		VariableName: rules.IslandsAlive,
		Values:       forVariables,
	})
}

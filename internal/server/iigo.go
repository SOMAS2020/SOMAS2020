package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/server/iigointernal"
)

// runIIGO : IIGO decides rule changes, elections, sanctions
func (s *SOMASServer) runIIGO() error {
	s.logf("start runIIGO")
	_ = iigointernal.RunIIGO(&s.gameState, &s.clientMap)
	defer s.logf("finish runIIGO")
	return nil
}

func (s *SOMASServer) runIIGOEndOfTurn() error {
	s.logf("start runIIGOEndOfTurn")
	defer s.logf("finish runIIGOEndOfTurn")
	clientMap := getNonDeadClients(s.gameState.ClientInfos, s.clientMap)
	for i, v := range clientMap {
		tax := v.GetTaxContribution()
		s.gameState.CommonPool += tax
		newGameState := s.gameState.GetClientGameStateCopy(i)
		newGameState.ClientInfo.Resources -= tax
		v.GameStateUpdate(newGameState)
		_ = iigointernal.UpdateTurnHistory(int(i), []rules.VariableValuePair{
			{
				VariableName: "island_tax_contribution",
				Values:       []float64{float64(tax)},
			},
			{
				VariableName: "expected_tax_contribution",
				Values:       []float64{float64(iigointernal.TaxAmountMapExport[int(i)])},
			},
		})
	}
	return nil
}

func (s *SOMASServer) runIIGOAllocations() error {
	s.logf("start runIIGOAllocations")
	defer s.logf("finish runIIGOAllocations")
	clientMap := getNonDeadClients(s.gameState.ClientInfos, s.clientMap)
	for i, v := range clientMap {
		allocation := v.RequestAllocation()
		s.gameState.CommonPool -= allocation
		newGameState := s.gameState.GetClientGameStateCopy(i)
		newGameState.ClientInfo.Resources += allocation
		v.GameStateUpdate(newGameState)
		if allocation <= s.gameState.CommonPool {
			s.gameState.CommonPool -= allocation
			_ = iigointernal.UpdateTurnHistory(int(i), []rules.VariableValuePair{
				{
					VariableName: "island_allocation",
					Values:       []float64{float64(allocation)},
				},
				{
					VariableName: "expected_allocation",
					Values:       []float64{float64(iigointernal.AllocationAmountMapExport[int(i)])},
				},
			})
		}
	}
	return nil
}

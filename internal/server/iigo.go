package server

import (
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

		err := s.takeResources(clientID, tax, "tax")
		if err != nil {
			s.gameState.CommonPool += tax
		} else {
			s.logf("Error getting tax from %v: %v", clientID, err)
		}
		_ = iigointernal.UpdateTurnHistory(clientID, []rules.VariableValuePair{
			{
				VariableName: "island_tax_contribution",
				Values:       []float64{float64(tax)},
			},
			{
				VariableName: "expected_tax_contribution",
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
		s.giveResources(clientID, allocation, "allocation")
		if allocation <= s.gameState.CommonPool {
			s.gameState.CommonPool -= allocation
			_ = iigointernal.UpdateTurnHistory(clientID, []rules.VariableValuePair{
				{
					VariableName: "island_allocation",
					Values:       []float64{float64(allocation)},
				},
				{
					VariableName: "expected_allocation",
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
	_ = rules.UpdateVariable("islands_alive", rules.VariableValuePair{
		VariableName: "islands_alive",
		Values:       forVariables,
	})
}

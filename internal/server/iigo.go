package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/server/iigointernal"
)

// runIIGO : IIGO decides rule changes, elections, sanctions
func (s *SOMASServer) runIIGO() error {
	s.logf("start runIIGO")
	defer s.logf("finish runIIGO")

	nonDead := getNonDeadClientIDs(s.gameState.ClientInfos)
	updateAliveIslands(nonDead)
	iigoSuccessful, iigoStatus := iigointernal.RunIIGO(&s.gameState, &s.clientMap)
	if !iigoSuccessful {
		s.logf(iigoStatus)
	}
	return nil
}

func (s *SOMASServer) updateIIGOTurnHistory(clientID shared.ClientID, pairs []rules.VariableValuePair) {
	s.gameState.IIGOHistory = append(s.gameState.IIGOHistory,
		shared.Accountability{
			ClientID: clientID,
			Pairs:    pairs,
		},
	)
}

func (s *SOMASServer) runIIGOTax() error {
	s.logf("start runIIGOTax")
	defer s.logf("finish runIIGOTax")
	clientMap := getNonDeadClients(s.gameState.ClientInfos, s.clientMap)
	for clientID, v := range clientMap {
		tax := v.GetTaxContribution()
		sanction := v.GetSanctionPayment()
		err := s.takeResources(clientID, tax, "tax")
		err2 := s.takeResources(clientID, sanction, "sanction")
		if err == nil && err2 == nil {
			s.gameState.CommonPool += tax + sanction
			s.clientMap[clientID].TaxTaken(tax)
		} else {
			s.logf("Error getting tax or sanctions from %v: %v sanction error: %v", clientID, err, err2)
		}
		s.updateIIGOTurnHistory(clientID, []rules.VariableValuePair{
			{
				VariableName: rules.IslandTaxContribution,
				Values:       []float64{float64(tax)},
			},
			{
				VariableName: rules.ExpectedTaxContribution,
				Values:       []float64{float64(iigointernal.TaxAmountMapExport[clientID])},
			},
		})
		s.updateIIGOTurnHistory(clientID, []rules.VariableValuePair{
			{
				VariableName: rules.SanctionPaid,
				Values:       []float64{float64(sanction)},
			},
			{
				VariableName: rules.SanctionExpected,
				Values:       []float64{float64(iigointernal.SanctionAmountMapExport[clientID])},
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

		if allocation <= s.gameState.CommonPool {
			s.giveResources(clientID, allocation, "allocation")
			s.gameState.CommonPool -= allocation

			s.updateIIGOTurnHistory(clientID, []rules.VariableValuePair{
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

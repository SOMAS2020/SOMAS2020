package server

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/server/iigointernal"
	"github.com/pkg/errors"
)

// runIIGO : IIGO decides rule changes, elections, sanctions
func (s *SOMASServer) runIIGO() error {
	s.logf("start runIIGO")
	defer s.logf("finish runIIGO")

	nonDead := getNonDeadClientIDs(s.gameState.ClientInfos)
	updateAliveIslands(nonDead, s.gameState)
	iigoSuccessful, iigoStatus := iigointernal.RunIIGO(s.logf, &s.gameState, &s.clientMap, &s.gameConfig)
	if !iigoSuccessful {
		s.logf(iigoStatus)
	}
	s.gameState.IIGORunStatus = iigoStatus
	return nil
}

func (s *SOMASServer) updateIIGOHistoryAndRules(clientID shared.ClientID, pairs []rules.VariableValuePair) {
	s.gameState.IIGOHistory[s.gameState.Turn] = append(s.gameState.IIGOHistory[s.gameState.Turn],
		shared.Accountability{
			ClientID: clientID,
			Pairs:    pairs,
		},
	)
}

func (s *SOMASServer) runIIGOTax() error {
	s.logf("start runIIGOTaxCommonPool")
	defer s.logf("finish runIIGOTaxCommonPool")
	clientMap := getNonDeadClients(s.gameState.ClientInfos, s.clientMap)
	for clientID, v := range clientMap {
		var taxPaid shared.Resources
		var sanctionPaid shared.Resources
		tax := v.GetTaxContribution()
		sanction := v.GetSanctionPayment()
		clientTaxErr := s.takeResources(clientID, tax, "tax")
		if clientTaxErr != nil {
			s.logf("Error getting tax from %v: %v", clientID, clientTaxErr)
			taxPaid = 0
		} else {
			s.gameState.CommonPool += tax
			taxPaid = tax
		}
		clientSanctionErr := s.takeResources(clientID, sanction, "sanction")
		if clientSanctionErr != nil {
			s.logf("Error getting sanctions from %v: %v ", clientID, clientSanctionErr)
			sanctionPaid = 0
		} else {
			s.gameState.CommonPool += sanction
			sanctionPaid = sanction
		}

		s.updateIIGOHistoryAndRules(clientID, []rules.VariableValuePair{
			{
				VariableName: rules.IslandTaxContribution,
				Values:       []float64{float64(taxPaid)},
			},
			{
				VariableName: rules.ExpectedTaxContribution,
				Values:       []float64{float64(s.gameState.IIGOTaxAmount[clientID])},
			},
		})
		s.updateIIGOHistoryAndRules(clientID, []rules.VariableValuePair{
			{
				VariableName: rules.SanctionPaid,
				Values:       []float64{float64(sanctionPaid)},
			},
			{
				VariableName: rules.SanctionExpected,
				Values:       []float64{float64(s.gameState.IIGOSanctionMap[clientID])},
			},
		})

	}
	return nil
}

func (s *SOMASServer) runIIGOAllocations() error {
	s.logf("start runIIGOAllocations")
	defer s.logf("finish runIIGOAllocations")
	clientMap := getNonDeadClients(s.gameState.ClientInfos, s.clientMap)
	allocationMap := make(map[shared.ClientID]shared.Resources)
	for clientID, v := range clientMap {
		allocation := v.RequestAllocation()
		if allocation < 0 || math.IsNaN(float64(allocation)) {
			s.logf("Invalid allocation of %v by %v. Changing allocation to 0", allocation, clientID)
			allocation = 0
		}
		if allocation <= s.gameState.CommonPool {
			err := s.giveResources(clientID, allocation, "allocation")
			if err != nil {
				return errors.Errorf("Failed to give resources: %v", err)
			}
			s.gameState.CommonPool -= allocation

			if s.gameState.IIGOAllocationMade {
				s.updateIIGOHistoryAndRules(clientID, []rules.VariableValuePair{
					{
						VariableName: rules.IslandAllocation,
						Values:       []float64{float64(allocation)},
					},
					{
						VariableName: rules.ExpectedAllocation,
						Values:       []float64{float64(s.gameState.IIGOAllocationMap[clientID])},
					},
				})
			} else {
				allocationMap[clientID] = allocation
				//Allow islands to take what they want if no allocations made
				s.updateIIGOHistoryAndRules(clientID, []rules.VariableValuePair{
					{
						VariableName: rules.IslandAllocation,
						Values:       []float64{float64(allocation)},
					},
					{
						VariableName: rules.ExpectedAllocation,
						Values:       []float64{float64(allocation)},
					},
				})
			}
		}
	}

	//Update gamestate for visualisation
	if !s.gameState.IIGOAllocationMade {
		s.gameState.IIGOAllocationMap = allocationMap
	}
	return nil
}

func updateAliveIslands(aliveIslands []shared.ClientID, gamestate gamestate.GameState) {
	var forVariables []float64
	for _, v := range aliveIslands {
		forVariables = append(forVariables, float64(v))
	}
	_ = gamestate.UpdateVariable(rules.IslandsAlive, rules.VariableValuePair{
		VariableName: rules.IslandsAlive,
		Values:       forVariables,
	})
}

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

	// .----------------.  .----------------.  .----------------.   .----------------.  .----------------.
	// | .--------------. || .--------------. || .--------------. | | .--------------. || .--------------. |
	// | |  _________   | || |     _____    | || |  ____  ____  | | | | ____    ____ | || |  _________   | |
	// | | |_   ___  |  | || |    |_   _|   | || | |_  _||_  _| | | | ||_   \  /   _|| || | |_   ___  |  | |
	// | |   | |_  \_|  | || |      | |     | || |   \ \  / /   | | | |  |   \/   |  | || |   | |_  \_|  | |
	// | |   |  _|      | || |      | |     | || |    > `' <    | | | |  | |\  /| |  | || |   |  _|  _   | |
	// | |  _| |_       | || |     _| |_    | || |  _/ /'`\ \_  | | | | _| |_\/_| |_ | || |  _| |___/ |  | |
	// | | |_____|      | || |    |_____|   | || | |____||____| | | | ||_____||_____|| || | |_________|  | |
	// | |              | || |              | || |              | | | |              | || |              | |
	// | '--------------' || '--------------' || '--------------' | | '--------------' || '--------------' |
	//  '----------------'  '----------------'  '----------------'   '----------------'  '----------------'

	// ISSUES
	// ..#..#.....##.....####.....##.............#..#.....##.....####...######.
	// .######...###........##...##.............######...###........##.....##..
	// ..#..#.....##.....####...#####............#..#.....##.....####.....##...
	// .######....##....##......##..##..........######....##....##.......##....
	// ..#..#...######..######...####............#..#...######..######..##.....

	// panic: Could not run IIGO since President has no resoruces to spend
	// goroutine 1 [running]:
	// github.com/SOMAS2020/SOMAS2020/internal/server.(*SOMASServer).runIIGO(0xc000048b60, 0x0, 0x0)
	// github.com/SOMAS2020/SOMAS2020/internal/server.(*SOMASServer).runOrgs(0xc000048b60, 0x0, 0x0)
	// github.com/SOMAS2020/SOMAS2020/internal/server.(*SOMASServer).runTurn(0xc000048b60, 0x0, 0x0)
	// github.com/SOMAS2020/SOMAS2020/internal/server.(*SOMASServer).EntryPoint(0xc000048b60, 0xc000048b60, 0x10, 0x20, 0x2d4a9250108, 0xc0000045e0)
	// main.main()
	// exit status 2
	_ = iigointernal.RunIIGO(&s.gameState, &s.clientMap)
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

func (s *SOMASServer) runIIGOEndOfTurn() error {
	s.logf("start runIIGOEndOfTurn")
	defer s.logf("finish runIIGOEndOfTurn")
	clientMap := getNonDeadClients(s.gameState.ClientInfos, s.clientMap)
	for clientID, v := range clientMap {
		tax := v.GetTaxContribution()
		err := s.takeResources(clientID, tax, "tax")
		if err == nil {
			s.gameState.CommonPool += tax
			s.clientMap[clientID].TaxTaken(tax)
		} else {
			s.logf("Error getting tax from %v: %v", clientID, err)
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

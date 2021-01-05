package gamestate

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestGetClientGameStateCopy(t *testing.T) {
	islandLocation := map[shared.ClientID]disasters.IslandLocationInfo{
		shared.Team1: {ID: shared.Team1, X: shared.Coordinate(5), Y: shared.Coordinate(0)},
		shared.Team2: {ID: shared.Team2, X: shared.Coordinate(6), Y: shared.Coordinate(1)},
		shared.Team3: {ID: shared.Team3, X: shared.Coordinate(7), Y: shared.Coordinate(2)},
	}

	geography := disasters.ArchipelagoGeography{Islands: islandLocation}
	env := disasters.Environment{Geography: geography}

	gameState := GameState{
		Season:      1,
		Turn:        4,
		Environment: env,
		ClientInfos: map[shared.ClientID]ClientInfo{
			shared.Team1: {
				Resources:                       10,
				LifeStatus:                      shared.Alive,
				CriticalConsecutiveTurnsCounter: 0,
			},
			shared.Team2: {
				Resources:                       20,
				LifeStatus:                      shared.Critical,
				CriticalConsecutiveTurnsCounter: 1,
			},
			shared.Team3: {
				Resources:                       30,
				LifeStatus:                      shared.Dead,
				CriticalConsecutiveTurnsCounter: 2,
			},
		},
		IIGORolesBudget: map[shared.Role]shared.Resources{
			shared.Judge:     shared.Resources(10),
			shared.President: shared.Resources(30),
			shared.Speaker:   shared.Resources(40),
		},
		CommonPool: 20,
	}

	lifeStatuses := map[shared.ClientID]shared.ClientLifeStatus{
		shared.Team1: gameState.ClientInfos[shared.Team1].LifeStatus,
		shared.Team2: gameState.ClientInfos[shared.Team2].LifeStatus,
		shared.Team3: gameState.ClientInfos[shared.Team3].LifeStatus,
	}

	cases := []shared.ClientID{shared.Team1, shared.Team2, shared.Team3}

	for _, tc := range cases {
		t.Run(tc.String(), func(t *testing.T) {
			expectClientGS := ClientGameState{
				Season:             gameState.Season,
				Turn:               gameState.Turn,
				ClientInfo:         gameState.ClientInfos[tc],
				ClientLifeStatuses: lifeStatuses,
				CommonPool:         gameState.CommonPool,
				IslandLocations:    gameState.Environment.Geography.Islands,
				IIGORolesBudget:    gameState.IIGORolesBudget,
			}

			gotClientGS := gameState.GetClientGameStateCopy(tc)

			if !reflect.DeepEqual(gotClientGS, expectClientGS) {
				t.Errorf(
					`Got unexpected ClientGameState.
					Got: %v
					Expected: %v`,
					gotClientGS, expectClientGS)
			}
		})
	}
}

package server

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
)

func TestDeductCostOfLiving(t *testing.T) {
	const costOfLiving = 42
	clientInfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Team1: {
			Resources:  43,
			LifeStatus: shared.Alive,
		},
		shared.Team2: {
			Resources:  44,
			LifeStatus: shared.Critical,
		},
		shared.Team3: {
			Resources:  45,
			LifeStatus: shared.Dead,
		},
		shared.Team4: {
			Resources:  20,
			LifeStatus: shared.Alive,
		},
	}
	wantClientInfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Team1: {
			Resources:  1,
			LifeStatus: shared.Alive,
		},
		shared.Team2: {
			Resources:  2,
			LifeStatus: shared.Critical,
		},
		shared.Team3: {
			Resources:  45,
			LifeStatus: shared.Dead,
		},
		shared.Team4: {
			Resources:  0,
			LifeStatus: shared.Alive,
		},
	}

	s := SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
	}

	s.deductCostOfLiving(costOfLiving)

	if !reflect.DeepEqual(wantClientInfos, s.gameState.ClientInfos) {
		t.Errorf("want '%v' got '%v'", wantClientInfos, s.gameState.ClientInfos)
	}
}

func TestUpdateIslandLivingStatus(t *testing.T) {
	// this does not test for updateIslandLivingStatusForClient
	// those are covered in TestUpdateIslandLivingStatusForClient
	clientInfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Team1: {
			LifeStatus: shared.Alive,
			Resources:  config.GameConfig().MinimumResourceThreshold - 1,
		},
		shared.Team2: {
			LifeStatus: shared.Critical,
			Resources:  config.GameConfig().MinimumResourceThreshold,
		},
	}
	wantClientInfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Team1: {
			LifeStatus: shared.Critical,
			Resources:  config.GameConfig().MinimumResourceThreshold - 1,
		},
		shared.Team2: {
			LifeStatus: shared.Alive,
			Resources:  config.GameConfig().MinimumResourceThreshold,
		},
	}

	s := SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
	}

	err := s.updateIslandLivingStatus()
	testutils.CompareTestErrors(nil, err, t)
	if !reflect.DeepEqual(s.gameState.ClientInfos, wantClientInfos) {
		t.Errorf("want '%v' got '%v'", wantClientInfos, s.gameState.ClientInfos)
	}
}

func TestGameOver(t *testing.T) {
	const maxTurns = 10
	const maxSeasons = 10
	cases := []struct {
		name        string
		clientInfos map[shared.ClientID]gamestate.ClientInfo
		turn        uint
		season      uint
		want        bool
	}{
		{
			name: "game not over",
			clientInfos: map[shared.ClientID]gamestate.ClientInfo{
				shared.Team2: {
					LifeStatus: shared.Alive,
				},
			},
			turn:   10,
			season: 10,
			want:   false,
		},
		{
			name: "all clients dead",
			clientInfos: map[shared.ClientID]gamestate.ClientInfo{
				shared.Team2: {
					LifeStatus: shared.Dead,
				},
			},
			turn:   10,
			season: 10,
			want:   true,
		},
		{
			name: "maxTurns reached",
			clientInfos: map[shared.ClientID]gamestate.ClientInfo{
				shared.Team2: {
					LifeStatus: shared.Alive,
				},
			},
			turn:   11,
			season: 10,
			want:   true,
		},
		{
			name: "maxSeasons reached",
			clientInfos: map[shared.ClientID]gamestate.ClientInfo{
				shared.Team2: {
					LifeStatus: shared.Alive,
				},
			},
			turn:   10,
			season: 11,
			want:   true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := SOMASServer{
				gameState: gamestate.GameState{
					ClientInfos: tc.clientInfos,
					Turn:        tc.turn,
					Season:      tc.season,
				},
			}

			got := server.gameOver(maxTurns, maxSeasons)
			if tc.want != got {
				t.Errorf("want '%v' got '%v'", tc.want, got)
			}
		})
	}
}

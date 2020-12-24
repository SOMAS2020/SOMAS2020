package server

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
)

type mockClientUpdate struct {
	baseclient.Client
	StartOfTurnUpdateCalled bool
	GameStateUpdateCalled   bool
}

func (c *mockClientUpdate) StartOfTurnUpdate(g gamestate.ClientGameState) {
	c.StartOfTurnUpdateCalled = true
}

func (c *mockClientUpdate) GameStateUpdate(g gamestate.ClientGameState) {
	c.GameStateUpdateCalled = true
}

func TestStartOfTurnUpdate(t *testing.T) {
	clientInfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Team1: {
			LifeStatus: shared.Alive,
		},
		shared.Team2: {
			LifeStatus: shared.Critical,
		},
		shared.Team3: {
			LifeStatus: shared.Dead,
		},
	}
	clientMap := map[shared.ClientID]baseclient.Client{
		shared.Team1: &mockClientUpdate{},
		shared.Team2: &mockClientUpdate{},
		shared.Team3: &mockClientUpdate{},
	}
	wantCalled := map[shared.ClientID]bool{
		shared.Team1: true,
		shared.Team2: true,
		shared.Team3: false,
	}

	s := &SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
		clientMap: clientMap,
	}

	s.startOfTurnUpdate()

	for id, want := range wantCalled {
		v, ok := clientMap[id].(*mockClientUpdate)
		if !ok {
			t.Errorf("Can't coerce type!")
		}
		got := v.StartOfTurnUpdateCalled
		if want != got {
			t.Errorf("For id %v, want '%v' got '%v'", id, want, got)
		}
	}
}

func TestGameStateUpdate(t *testing.T) {
	clientInfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Team1: {
			LifeStatus: shared.Alive,
		},
		shared.Team2: {
			LifeStatus: shared.Critical,
		},
		shared.Team3: {
			LifeStatus: shared.Dead,
		},
	}
	clientMap := map[shared.ClientID]baseclient.Client{
		shared.Team1: &mockClientUpdate{},
		shared.Team2: &mockClientUpdate{},
		shared.Team3: &mockClientUpdate{},
	}
	wantCalled := map[shared.ClientID]bool{
		shared.Team1: true,
		shared.Team2: true,
		shared.Team3: false,
	}

	s := &SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
		clientMap: clientMap,
	}

	s.gameStateUpdate()

	for id, want := range wantCalled {
		v, ok := clientMap[id].(*mockClientUpdate)
		if !ok {
			t.Errorf("Can't coerce type!")
		}
		got := v.GameStateUpdateCalled
		if want != got {
			t.Errorf("For id %v, want '%v' got '%v'", id, want, got)
		}
	}
}
func TestIncrementTurnAndSeason(t *testing.T) {
	cases := []struct {
		name             string
		initialGameState gamestate.GameState
		disasterHappened bool
		wantTurn         uint
		wantSeason       uint
	}{
		{
			name: "no disaster",
			initialGameState: gamestate.GameState{
				Turn:   42,
				Season: 69,
			},
			disasterHappened: false,
			wantTurn:         43,
			wantSeason:       69,
		},
		{
			name: "disaster happened",
			initialGameState: gamestate.GameState{
				Turn:   42,
				Season: 69,
			},
			disasterHappened: false,
			wantTurn:         43,
			wantSeason:       69,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := SOMASServer{
				gameState: tc.initialGameState,
			}

			s.incrementTurnAndSeason(tc.disasterHappened)

			gotTurn := s.gameState.Turn
			gotSeason := s.gameState.Season

			if gotTurn != tc.wantTurn {
				t.Errorf("Turn: want '%v' got '%v'", tc.wantTurn, gotTurn)
			}
			if gotSeason != tc.wantSeason {
				t.Errorf("Season: want '%v' got '%v'", tc.wantSeason, gotSeason)
			}
		})
	}
}

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

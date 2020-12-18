package server

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type mockClientUpdate struct {
	common.Client
	StartOfTurnUpdateCalled bool
	GameStateUpdateCalled   bool
}

func (c *mockClientUpdate) StartOfTurnUpdate(g common.GameState) {
	c.StartOfTurnUpdateCalled = true
}

func (c *mockClientUpdate) GameStateUpdate(g common.GameState) {
	c.GameStateUpdateCalled = true
}

func TestStartOfTurnUpdate(t *testing.T) {
	clientInfos := map[shared.ClientID]common.ClientInfo{
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
	clientMap := map[shared.ClientID]common.Client{
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
		gameState: common.GameState{
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
	clientInfos := map[shared.ClientID]common.ClientInfo{
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
	clientMap := map[shared.ClientID]common.Client{
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
		gameState: common.GameState{
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

func TestGameOver(t *testing.T) {
	const maxTurns = 10
	const maxSeasons = 10
	cases := []struct {
		name        string
		clientInfos map[shared.ClientID]common.ClientInfo
		turn        uint
		season      uint
		want        bool
	}{
		{
			name: "game not over",
			clientInfos: map[shared.ClientID]common.ClientInfo{
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
			clientInfos: map[shared.ClientID]common.ClientInfo{
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
			clientInfos: map[shared.ClientID]common.ClientInfo{
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
			clientInfos: map[shared.ClientID]common.ClientInfo{
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
				gameState: common.GameState{
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

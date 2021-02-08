package server

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
	"github.com/pkg/errors"
)

func TestAnyClientsAlive(t *testing.T) {
	cases := []struct {
		name string
		cis  map[shared.ClientID]gamestate.ClientInfo
		want bool
	}{
		{
			name: "all alive",
			cis: map[shared.ClientID]gamestate.ClientInfo{
				shared.Team1: {
					LifeStatus: shared.Alive,
				},
				shared.Team2: {
					LifeStatus: shared.Critical, // still alive
				},
			},
			want: true,
		},
		{
			name: "one alive",
			cis: map[shared.ClientID]gamestate.ClientInfo{
				shared.Team1: {
					LifeStatus: shared.Alive,
				},
				shared.Team2: {
					LifeStatus: shared.Dead,
				},
			},
			want: true,
		}, {
			name: "none alive",
			cis: map[shared.ClientID]gamestate.ClientInfo{
				shared.Team1: {
					LifeStatus: shared.Dead,
				},
				shared.Team2: {
					LifeStatus: shared.Dead,
				},
			},
			want: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := anyClientsAlive(tc.cis)
			if tc.want != got {
				t.Errorf("want '%v' got '%v'", tc.want, got)
			}
		})
	}
}

func TestUpdateIslandLivingStatusForClient(t *testing.T) {
	const minResThres = 10
	const maxConsCritTurns = 3
	cases := []struct {
		name    string
		ci      gamestate.ClientInfo
		want    gamestate.ClientInfo
		wantErr error
	}{
		{
			name: "alive and well",
			ci: gamestate.ClientInfo{
				LifeStatus:                      shared.Alive,
				Resources:                       minResThres,
				CriticalConsecutiveTurnsCounter: 0,
			},
			want: gamestate.ClientInfo{
				LifeStatus:                      shared.Alive,
				Resources:                       minResThres,
				CriticalConsecutiveTurnsCounter: 0,
			},
			wantErr: nil,
		},
		{
			name: "already dead",
			ci: gamestate.ClientInfo{
				LifeStatus:                      shared.Dead,
				Resources:                       minResThres - 1,
				CriticalConsecutiveTurnsCounter: 3,
			},
			want: gamestate.ClientInfo{
				LifeStatus:                      shared.Dead,
				Resources:                       minResThres - 1,
				CriticalConsecutiveTurnsCounter: 3,
			},
			wantErr: nil,
		},
		{
			name: "turn critical",
			ci: gamestate.ClientInfo{
				LifeStatus:                      shared.Alive,
				Resources:                       minResThres - 1,
				CriticalConsecutiveTurnsCounter: 100,
			},
			want: gamestate.ClientInfo{
				LifeStatus:                      shared.Critical,
				Resources:                       minResThres - 1,
				CriticalConsecutiveTurnsCounter: 0,
			},
			wantErr: nil,
		},
		{
			name: "add critical",
			ci: gamestate.ClientInfo{
				LifeStatus:                      shared.Critical,
				Resources:                       minResThres - 1,
				CriticalConsecutiveTurnsCounter: 1,
			},
			want: gamestate.ClientInfo{
				LifeStatus:                      shared.Critical,
				Resources:                       minResThres - 1,
				CriticalConsecutiveTurnsCounter: 2,
			},
			wantErr: nil,
		},
		{
			name: "ran out of critical turns",
			ci: gamestate.ClientInfo{
				LifeStatus:                      shared.Critical,
				Resources:                       minResThres - 1,
				CriticalConsecutiveTurnsCounter: 3,
			},
			want: gamestate.ClientInfo{
				LifeStatus:                      shared.Dead,
				Resources:                       minResThres - 1,
				CriticalConsecutiveTurnsCounter: 3,
			},
			wantErr: nil,
		},
		{
			name: "turn non-critical",
			ci: gamestate.ClientInfo{
				LifeStatus:                      shared.Critical,
				Resources:                       minResThres,
				CriticalConsecutiveTurnsCounter: 2,
			},
			want: gamestate.ClientInfo{
				LifeStatus:                      shared.Alive,
				Resources:                       minResThres,
				CriticalConsecutiveTurnsCounter: 0,
			},
			wantErr: nil,
		},
		{
			name: "bogus LifeStatus",
			ci: gamestate.ClientInfo{
				LifeStatus:                      99999,
				Resources:                       minResThres,
				CriticalConsecutiveTurnsCounter: 2,
			},
			want: gamestate.ClientInfo{
				LifeStatus:                      99999,
				Resources:                       minResThres,
				CriticalConsecutiveTurnsCounter: 2,
			},
			wantErr: errors.Errorf("updateIslandLivingStatusForClient not implemented " +
				"for LifeStatus UNKNOWN ClientLifeStatus '99999'"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := updateIslandLivingStatusForClient(tc.ci, minResThres, maxConsCritTurns)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want '%#v' got '%#v'", tc.want, got)
			}
			testutils.CompareTestErrors(tc.wantErr, err, t)
		})
	}
}

func TestGetNonDeadClientIDs(t *testing.T) {
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
	want := []shared.ClientID{
		shared.Team1, shared.Team2,
	}
	got := getNonDeadClientIDs(clientInfos)

	sort.Sort(shared.SortClientByID(want))
	sort.Sort(shared.SortClientByID(got))

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want '%v' got '%v'", want, got)
	}
}

func TestTakeResources(t *testing.T) {
	cases := []struct {
		name      string
		resources shared.Resources
		takeAmt   shared.Resources
		want      shared.Resources
		wantErr   error
	}{
		{
			name:      "normal",
			resources: 42,
			takeAmt:   3,
			want:      39,
		},
		{
			name:      "take 0",
			resources: 42,
			takeAmt:   0,
			want:      42,
		},
		{
			name:      "Go to 0",
			resources: 42,
			takeAmt:   42,
			want:      0,
		},
		{
			name:      "Go below 0",
			resources: 10,
			takeAmt:   15,
			want:      10,
			wantErr:   errors.Errorf("Client %v did not have enough resources. Requested %v, only had %v", shared.Team1, 15, 10),
		},
		{
			name:      "NaN",
			resources: 42,
			takeAmt:   shared.Resources(math.NaN()),
			want:      42,
			wantErr: errors.Errorf("Cannot take invalid number of resources %v from client %v",
				math.NaN(), shared.Team1),
		},
		{
			name:      "try take negative",
			resources: 42,
			takeAmt:   shared.Resources(-42),
			want:      42,
			wantErr: errors.Errorf("Cannot take invalid number of resources %v from client %v",
				-42, shared.Team1),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := SOMASServer{
				gameState: gamestate.GameState{
					ClientInfos: map[shared.ClientID]gamestate.ClientInfo{
						shared.Team1: {Resources: tc.resources},
					},
				},
			}

			err := s.takeResources(shared.Team1, tc.takeAmt, tc.name)
			got := s.gameState.ClientInfos[shared.Team1].Resources

			testutils.CompareTestErrors(tc.wantErr, err, t)

			if tc.want != got {
				t.Errorf("want '%v' got '%v'", tc.want, got)
			}
		})
	}
}

func TestGiveResources(t *testing.T) {
	cases := []struct {
		name      string
		resources shared.Resources
		giveAmt   shared.Resources
		want      shared.Resources
		wantErr   error
	}{
		{
			name:      "normal",
			resources: 42,
			giveAmt:   3,
			want:      45,
		},
		{
			name:      "Give 0",
			resources: 42,
			giveAmt:   0,
			want:      42,
		},
		{
			name:      "Give NaN",
			resources: 42,
			giveAmt:   shared.Resources(math.NaN()),
			want:      42,
			wantErr: errors.Errorf("Cannot give invalid number of resources %v to client %v",
				math.NaN(), shared.Team1),
		},
		{
			name:      "Give -42",
			resources: 42,
			giveAmt:   -42,
			want:      42,
			wantErr: errors.Errorf("Cannot give invalid number of resources %v to client %v",
				-42, shared.Team1),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := SOMASServer{
				gameState: gamestate.GameState{
					ClientInfos: map[shared.ClientID]gamestate.ClientInfo{
						shared.Team1: {Resources: tc.resources},
					},
				},
			}

			err := s.giveResources(shared.Team1, tc.giveAmt, tc.name)
			got := s.gameState.ClientInfos[shared.Team1].Resources

			testutils.CompareTestErrors(tc.wantErr, err, t)

			if tc.want != got {
				t.Errorf("want '%v' got '%v'", tc.want, got)
			}
		})
	}
}

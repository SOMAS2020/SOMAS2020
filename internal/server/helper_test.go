package server

import (
	"reflect"
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

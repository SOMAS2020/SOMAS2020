package server

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestAnyClientsAlive(t *testing.T) {
	cases := []struct {
		name string
		cis  map[shared.ClientID]common.ClientInfo
		want bool
	}{
		{
			name: "all alive",
			cis: map[shared.ClientID]common.ClientInfo{
				shared.Team1: {
					Alive: true,
				},
				shared.Team2: {
					Alive: true,
				},
			},
			want: true,
		},
		{
			name: "one alive",
			cis: map[shared.ClientID]common.ClientInfo{
				shared.Team1: {
					Alive: false,
				},
				shared.Team2: {
					Alive: true,
				},
			},
			want: true,
		}, {
			name: "none alive",
			cis: map[shared.ClientID]common.ClientInfo{
				shared.Team1: {
					Alive: false,
				},
				shared.Team2: {
					Alive: false,
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
		name string
		ci   common.ClientInfo
		want common.ClientInfo
	}{
		{
			name: "alive and well",
			ci: common.ClientInfo{
				Alive:                        true,
				Critical:                     false,
				Resources:                    minResThres,
				CriticalConsecutiveTurnsLeft: 42,
			},
			want: common.ClientInfo{
				Alive:                        true,
				Critical:                     false,
				Resources:                    minResThres,
				CriticalConsecutiveTurnsLeft: 42,
			},
		},
		{
			name: "already dead",
			ci: common.ClientInfo{
				Alive:                        false,
				Critical:                     true,
				Resources:                    minResThres - 1,
				CriticalConsecutiveTurnsLeft: 42,
			},
			want: common.ClientInfo{
				Alive:                        false,
				Critical:                     true,
				Resources:                    minResThres - 1,
				CriticalConsecutiveTurnsLeft: 42,
			},
		},
		{
			name: "turn critical",
			ci: common.ClientInfo{
				Alive:                        true,
				Critical:                     false,
				Resources:                    minResThres - 1,
				CriticalConsecutiveTurnsLeft: 100,
			},
			want: common.ClientInfo{
				Alive:                        true,
				Critical:                     true,
				Resources:                    minResThres - 1,
				CriticalConsecutiveTurnsLeft: maxConsCritTurns,
			},
		},
		{
			name: "deduct critical",
			ci: common.ClientInfo{
				Alive:                        true,
				Critical:                     true,
				Resources:                    minResThres - 1,
				CriticalConsecutiveTurnsLeft: 100,
			},
			want: common.ClientInfo{
				Alive:                        true,
				Critical:                     true,
				Resources:                    minResThres - 1,
				CriticalConsecutiveTurnsLeft: 99,
			},
		},
		{
			name: "ran out of critical turns",
			ci: common.ClientInfo{
				Alive:                        true,
				Critical:                     true,
				Resources:                    minResThres - 1,
				CriticalConsecutiveTurnsLeft: 0,
			},
			want: common.ClientInfo{
				Alive:                        false,
				Critical:                     true,
				Resources:                    minResThres - 1,
				CriticalConsecutiveTurnsLeft: 0,
			},
		},
		{
			name: "turn non-critical",
			ci: common.ClientInfo{
				Alive:                        true,
				Critical:                     true,
				Resources:                    minResThres,
				CriticalConsecutiveTurnsLeft: 0,
			},
			want: common.ClientInfo{
				Alive:                        true,
				Critical:                     false,
				Resources:                    minResThres,
				CriticalConsecutiveTurnsLeft: maxConsCritTurns,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := updateIslandLivingStatusForClient(tc.ci, minResThres, maxConsCritTurns)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want '%#v' got '%#v'", tc.want, got)
			}
		})
	}
}

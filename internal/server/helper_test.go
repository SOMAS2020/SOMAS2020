package server

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
)

func TestAnyClientsAlive(t *testing.T) {
	cases := []struct {
		name string
		cis  map[common.ClientID]common.ClientInfo
		want bool
	}{
		{
			name: "all alive",
			cis: map[common.ClientID]common.ClientInfo{
				common.Team1: {
					Alive: true,
				},
				common.Team2: {
					Alive: true,
				},
			},
			want: true,
		},
		{
			name: "one alive",
			cis: map[common.ClientID]common.ClientInfo{
				common.Team1: {
					Alive: false,
				},
				common.Team2: {
					Alive: true,
				},
			},
			want: true,
		}, {
			name: "none alive",
			cis: map[common.ClientID]common.ClientInfo{
				common.Team1: {
					Alive: false,
				},
				common.Team2: {
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

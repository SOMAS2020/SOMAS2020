package server

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
)

type mockClientForage struct {
	common.Client
	forageInvestment uint
}

func (c *mockClientForage) GetForageInvestment(g common.GameState) uint {
	return c.forageInvestment
}

func TestGetTeamForageInvestments(t *testing.T) {
	gs := common.GameState{
		ClientInfos: map[common.ClientID]common.ClientInfo{
			common.Team1: { // OK
				Client: &mockClientForage{
					forageInvestment: 9,
				},
				Resources: 10,
				Alive:     true,
			},
			common.Team2: { // DEAD
				Client: &mockClientForage{
					forageInvestment: 9,
				},
				Resources: 10,
				Alive:     false,
			},
			common.Team3: { // OVERINVEST
				Client: &mockClientForage{
					forageInvestment: 1000,
				},
				Resources: 10,
				Alive:     true,
			},
		},
	}

	want := map[common.ClientID]uint{
		common.Team1: 9,
		common.Team3: 10,
	}

	var fakeLogger = func(format string, a ...interface{}) {}

	got, err := getTeamForageInvestments(fakeLogger, gs)
	testutils.CompareTestErrors(nil, err, t)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want '%v' got '%v'", want, got)
	}
}

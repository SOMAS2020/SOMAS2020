package team4

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestGetTrustMatrix(t *testing.T) {
	cases := []struct {
		name                string
		clientMap           map[shared.ClientID]float64
		expectedTrustMatrix []float64
	}{
		{
			name: "simple",
			clientMap: map[shared.ClientID]float64{
				shared.Team1: 0.5,
				shared.Team2: 0.5,
				shared.Team3: 0.5,
			},
			expectedTrustMatrix: []float64{0.5, 0.5, 0.5},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tr := trust{tc.clientMap}
			gotTrustMatrix := tr.GetTrustMatrix()
			if !reflect.DeepEqual(gotTrustMatrix, tc.expectedTrustMatrix) {
				t.Errorf("Expected trust matrix: %v, got: %v", tc.expectedTrustMatrix, gotTrustMatrix)
			}
		})
	}
}

func TestNormalise(t *testing.T) {
	cases := []struct {
		name      string
		clientMap map[shared.ClientID]float64
	}{
		{
			name: "simple normalise",
			clientMap: map[shared.ClientID]float64{
				shared.Team1: 0.5,
				shared.Team2: 0.7,
				shared.Team3: 0.5,
				shared.Team4: 0.5,
			},
		},
		{
			name: "simple normalise 2",
			clientMap: map[shared.ClientID]float64{
				shared.Team1: 0.2,
				shared.Team2: 0.7,
				shared.Team3: 0.6,
				shared.Team4: 0.55,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tr := trust{tc.clientMap}
			tr.normalise()

			totalTrust := 0.0
			for _, trust := range tr.trustMap {
				totalTrust += trust
			}

			expectedTrust := tr.expectedTrustSum()

			if !floatEqual(totalTrust, expectedTrust) {
				t.Errorf("Normalization failed. Expected trust sum: %v. Got trust sum: %v", expectedTrust, totalTrust)
			}

		})
	}
}

func TestIncreaseClientTrust(t *testing.T) {
	cases := []struct {
		name      string
		clientMap map[shared.ClientID]float64
		change    map[shared.ClientID]float64
	}{
		{
			name: "",
			clientMap: map[shared.ClientID]float64{
				shared.Team1: 0.5,
				shared.Team2: 0.5,
				shared.Team3: 0.5,
				shared.Team4: 0.5,
			},
			change: map[shared.ClientID]float64{
				shared.Team1: -0.4,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tr := trust{tc.clientMap}

			for client, diff := range tc.change {
				tr.ChangeClientTrust(client, diff)
			}

			t.Logf("%v", tr.trustMap)

			if !floatEqual(tr.totalTrustSum(), tr.expectedTrustSum()) {
				t.Errorf("Trust increase failed. Expected trust sum: %v. Got trust sum: %v", tr.expectedTrustSum(), tr.totalTrustSum())
			}

		})
	}
}

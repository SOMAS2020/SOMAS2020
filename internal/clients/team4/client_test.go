package team4

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestGetJudgePointer(t *testing.T) {
	testClient := newClientInternal(shared.Team4, honest)
	testServer := fakeServerHandle{}
	testClient.Initialise(testServer)
	j := testClient.GetClientJudgePointer()

	winner := j.DecideNextPresident(shared.Team1)

	if winner != shared.Team1 {
		t.Errorf("Got wrong judge pointer. Winner is %v", winner)
	}
}

func TestUpdateTrustFromHistory(t *testing.T) {
	cases := []struct {
		name         string
		savedHistory accountabilityHistory
	}{
		{
			name: "simple test",
			savedHistory: accountabilityHistory{
				history: map[uint]map[shared.ClientID]judgeHistoryInfo{
					1: {
						shared.Team1: {
							LawfulRatio: 0.2,
						},
						shared.Team2: {
							LawfulRatio: 0.4,
						},
						shared.Team3: {
							LawfulRatio: 0.87,
						},
						shared.Team4: {
							LawfulRatio: 1,
						},
						shared.Team5: {
							LawfulRatio: 0.675,
						},
						shared.Team6: {
							LawfulRatio: 0.55,
						},
					},
				},
				updated:     true,
				updatedTurn: 1,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clients := []shared.ClientID{}
			for clientID := range tc.savedHistory.getNewInfo() {
				clients = append(clients, clientID)
			}

			testServer := fakeServerHandle{clients: clients}
			testClient := newClientInternal(shared.Team4, honest)
			testClient.Initialise(testServer)

			testClient.savedHistory = &tc.savedHistory
			testClient.updateTrustFromSavedHistory()

			for clientID, trust := range testClient.trustMatrix.trustMap {
				// t.Logf("Updated trust for %v: %v", clientID, trust)
				if floatEqual(trust, 0.5) {
					t.Errorf("Update from history failed. trust didn't change for %v", clientID)
				}
			}

			if trustSum := testClient.trustMatrix.totalTrustSum(); trustSum > 0 {
				averageTrust := trustSum / float64(len(testClient.trustMatrix.trustMap))
				if !floatEqual(averageTrust, 0.5) {
					t.Errorf("Update from history failed. Average trust: %v", averageTrust)
				}
			}

		})
	}
}

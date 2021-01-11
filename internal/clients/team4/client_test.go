package team4

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestGetJudgePointer(t *testing.T) {
	testClient := newClientInternal(id)
	testServer := fakeServerHandle{}
	testClient.Initialise(testServer)
	j := testClient.GetClientJudgePointer()

	winner := j.DecideNextPresident(shared.Team1)

	if winner != shared.Team1 {
		t.Errorf("Got wrong judge pointer. Winner is %v", winner)
	}
}

// func TestUpdateTrustFromHistory(t *testing.T) {
// 	cases := []struct {
// 		name         string
// 		savedHistory accountabilityHistory
// 	}{
// 		{
// 			name: "simple test",
// 			savedHistory: accountabilityHistory{
// 				history: map[uint]map[shared.ClientID]judgeHistoryInfo{
// 					1: {
// 						shared.Team1: {
// 							TruthfulRatio: 0.2,
// 						},
// 						shared.Team2: {
// 							TruthfulRatio: 0.4,
// 						},
// 						shared.Team3: {
// 							TruthfulRatio: 0.87,
// 						},
// 						shared.Team4: {
// 							TruthfulRatio: 1,
// 						},
// 						shared.Team5: {
// 							TruthfulRatio: 0.6319,
// 						},
// 						shared.Team6: {
// 							TruthfulRatio: 0.55,
// 						},
// 					},
// 				},
// 				updated:     true,
// 				updatedTurn: 1,
// 			},
// 		},
// 	}

// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			clients := []shared.ClientID{}
// 			for clientID := range tc.savedHistory.getNewInfo() {
// 				clients = append(clients, clientID)
// 			}

// 			testServer := fakeServerHandle{clients: clients}
// 			testClient := newClientInternal(id, t)
// 			testClient.Initialise(testServer)

// 			testClient.savedHistory = &tc.savedHistory
// 			testClient.updateTrustFromSavedHistory()

// 		})
// 	}
// }

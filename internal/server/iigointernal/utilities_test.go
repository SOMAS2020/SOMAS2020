package iigointernal

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestWithdrawFromCommonPoolThrowsError(t *testing.T) {
	fakeGameState := gamestate.GameState{CommonPool: 100}
	// Withdraw more than we have in it
	valueToWithdraw := shared.Resources(120)
	_, withdrawSuccessful := WithdrawFromCommonPool(valueToWithdraw, &fakeGameState)
	if withdrawSuccessful {
		t.Errorf("We can withdraw more from the common pool than it actually has.")
	}
}

func TestWithdrawFromCommonPoolDeductsValue(t *testing.T) {
	cases := []struct {
		name              string
		fakeGameState     gamestate.GameState
		valueToWithdraw   shared.Resources
		expectedRemainder shared.Resources
		expectedResult    bool
	}{
		{
			name:              "Basic Common Pool deduction test",
			fakeGameState:     gamestate.GameState{CommonPool: 100},
			valueToWithdraw:   60,
			expectedRemainder: 40,
			expectedResult:    true,
		},
		{
			name:              "Checking error created for withdrawing more than available",
			fakeGameState:     gamestate.GameState{CommonPool: 10},
			valueToWithdraw:   50,
			expectedRemainder: 10,
			expectedResult:    false,
		},
		{
			name:              "Checking that you can extract entire common pool",
			fakeGameState:     gamestate.GameState{CommonPool: 100},
			valueToWithdraw:   100,
			expectedRemainder: 0,
			expectedResult:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, withdrawSuccessful := WithdrawFromCommonPool(tc.valueToWithdraw, &tc.fakeGameState)
			got := tc.fakeGameState.CommonPool
			if tc.expectedRemainder != got {
				t.Errorf("remainder: want '%v' got '%v'", tc.expectedRemainder, got)
			}
			if tc.expectedResult != withdrawSuccessful {
				t.Errorf("Expected withdrawalResult to be '%v' got '%v'", tc.expectedResult, withdrawSuccessful)
			}
		})
	}
}

func TestCommunicateWithIslands(t *testing.T) {

	dataA := map[shared.CommunicationFieldName]shared.CommunicationContent{
		3:  {IntegerData: 123, TextData: "Hello World - dataA", BooleanData: true},
		1:  {TextData: "SOMAS", BooleanData: false},
		22: {BooleanData: false},
		14: {IntegerData: 420, BooleanData: false},
	}

	dataB := map[shared.CommunicationFieldName]shared.CommunicationContent{
		0: {IntegerData: 11, TextData: "SOMAS", BooleanData: true},
	}

	dataC := map[shared.CommunicationFieldName]shared.CommunicationContent{
		5:  {BooleanData: true},
		4:  {TextData: "communication test"},
		16: {IntegerData: 7832},
		73: {IntegerData: 234511, TextData: "dataC", BooleanData: false},
	}

	dataEmpty := map[shared.CommunicationFieldName]shared.CommunicationContent{}

	cases := []struct {
		name           string
		sendersPayload map[int][]map[shared.CommunicationFieldName]shared.CommunicationContent
		receiver       int
	}{
		{
			name: "single transmission",
			sendersPayload: map[int][]map[shared.CommunicationFieldName]shared.CommunicationContent{
				1: {dataA},
			},
			receiver: 4,
		},
		{
			name: "2 senders, 1 transmission each",
			sendersPayload: map[int][]map[shared.CommunicationFieldName]shared.CommunicationContent{
				1: {dataA},
				2: {dataB},
			},
			receiver: 5,
		},
		{
			name: "1 sender, 2 transmissions",
			sendersPayload: map[int][]map[shared.CommunicationFieldName]shared.CommunicationContent{
				4: {dataA, dataC},
			},
			receiver: 0,
		},
		{
			name: "multiple transmissions",
			sendersPayload: map[int][]map[shared.CommunicationFieldName]shared.CommunicationContent{
				1: {dataA, dataC, dataA, dataC},
				2: {dataB, dataB, dataC},
				3: {dataA, dataB, dataC, dataC},
				4: {dataB, dataC, dataC, dataC},
			},
			receiver: 5,
		},
		{
			name: "multiple transmissions v2",
			sendersPayload: map[int][]map[shared.CommunicationFieldName]shared.CommunicationContent{
				1: {dataA, dataC, dataA, dataC},
				2: {dataB, dataB, dataC},
				3: {dataA, dataB, dataC, dataC},
				4: {dataB, dataC, dataC, dataC},
				5: {dataB, dataC, dataC, dataC},
			},
			receiver: 0,
		},
		{
			name: "1 sender, many transmissions",
			sendersPayload: map[int][]map[shared.CommunicationFieldName]shared.CommunicationContent{
				1: {dataA, dataC, dataA, dataC, dataB, dataB, dataC, dataA, dataB, dataC,
					dataC, dataB, dataC, dataC, dataC, dataA, dataC, dataA, dataC, dataB,
					dataB, dataC, dataA, dataB, dataC, dataC, dataB, dataC, dataC, dataC},
			},
			receiver: 0,
		},
		{
			name: "Empty transmission",
			sendersPayload: map[int][]map[shared.CommunicationFieldName]shared.CommunicationContent{
				1: {dataEmpty},
			},
			receiver: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fakeClientMap := map[shared.ClientID]baseclient.Client{}

			// Register clients
			receiverID := shared.TeamIDs[tc.receiver]
			fakeClientMap[receiverID] = baseclient.NewClient(receiverID)

			for sender := range tc.sendersPayload {
				senderID := shared.TeamIDs[sender]
				fakeClientMap[senderID] = baseclient.NewClient(senderID)
			}

			setIIGOClients(&fakeClientMap)

			// Perform communications + build expected output
			expectedResult := map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent{}

			for sender, dataList := range tc.sendersPayload {
				senderID := shared.TeamIDs[sender]
				for _, data := range dataList {
					communicateWithIslands(receiverID, senderID, data)

					expectedResult[senderID] = append(expectedResult[senderID], data)
				}
			}

			// Check internals of clients
			recieverGot := *(fakeClientMap[receiverID]).GetCommunications()

			if !reflect.DeepEqual(expectedResult, recieverGot) {
				t.Errorf("CommunicationContent failed. Sent: %v\nGot: %v", expectedResult, recieverGot)
			}

		})
	}
}

func TestWithdrawFromCommonPool(t *testing.T) {
	cases := []struct {
		name             string
		gamestate        *gamestate.GameState
		inputValue       shared.Resources
		expectedResource shared.Resources
		expectedAmount   shared.Resources
		expectedState    bool
	}{
		{
			name: "Withdraw amount",
			gamestate: &gamestate.GameState{
				CommonPool: 300,
				ClientInfos: map[shared.ClientID]gamestate.ClientInfo{
					shared.Team1: gamestate.ClientInfo{Resources: 100},
					shared.Team2: gamestate.ClientInfo{Resources: 100},
					shared.Team3: gamestate.ClientInfo{Resources: 100},
				},
			},
			inputValue:       100,
			expectedAmount:   100,
			expectedResource: 200,
			expectedState:    true,
		},
		{
			name: "Withdraw negative amount",
			gamestate: &gamestate.GameState{
				CommonPool: 300,
				ClientInfos: map[shared.ClientID]gamestate.ClientInfo{
					shared.Team1: gamestate.ClientInfo{Resources: 100},
					shared.Team2: gamestate.ClientInfo{Resources: 100},
					shared.Team3: gamestate.ClientInfo{Resources: 100},
				},
			},
			inputValue:       -80,
			expectedAmount:   -80,
			expectedResource: 380,
			expectedState:    true,
		},
		{
			name: "Withdraw more than common pool",
			gamestate: &gamestate.GameState{
				CommonPool: 100,
				ClientInfos: map[shared.ClientID]gamestate.ClientInfo{
					shared.Team1: gamestate.ClientInfo{Resources: 100},
					shared.Team2: gamestate.ClientInfo{Resources: 70},
					shared.Team3: gamestate.ClientInfo{Resources: 30},
				},
			},
			inputValue:       250,
			expectedAmount:   0,
			expectedResource: 100,
			expectedState:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			amount, state := WithdrawFromCommonPool(tc.inputValue, tc.gamestate)
			commonPool := tc.gamestate.CommonPool
			if amount != tc.expectedAmount || commonPool != tc.expectedResource || state != tc.expectedState {
				t.Errorf("%v - Failed. Got '%v %v %v', but expected '%v %v %v'",
					tc.name, amount, commonPool, state,
					tc.expectedAmount, tc.expectedResource, tc.expectedState)
			}
		})
	}
}

func TestDepositIntoClientPrivatePool(t *testing.T) {
	cases := []struct {
		name             string
		gamestate        *gamestate.GameState
		inputValue       shared.Resources
		inputID          shared.ClientID
		expectedResource shared.Resources
	}{
		{
			name: "Deposit amount",
			gamestate: &gamestate.GameState{
				ClientInfos: map[shared.ClientID]gamestate.ClientInfo{
					shared.Team1: gamestate.ClientInfo{Resources: 100},
					shared.Team2: gamestate.ClientInfo{Resources: 100},
					shared.Team3: gamestate.ClientInfo{Resources: 100},
				},
			},
			inputValue:       80,
			inputID:          shared.Team2,
			expectedResource: 180,
		},
		{
			name: "Negative amount",
			gamestate: &gamestate.GameState{
				ClientInfos: map[shared.ClientID]gamestate.ClientInfo{
					shared.Team1: gamestate.ClientInfo{Resources: 100},
					shared.Team2: gamestate.ClientInfo{Resources: 100},
					shared.Team3: gamestate.ClientInfo{Resources: 100},
				},
			},
			inputValue:       -80,
			inputID:          shared.Team1,
			expectedResource: 20,
		},
		{
			name: "Negative resources",
			gamestate: &gamestate.GameState{
				ClientInfos: map[shared.ClientID]gamestate.ClientInfo{
					shared.Team1: gamestate.ClientInfo{Resources: 100},
					shared.Team2: gamestate.ClientInfo{Resources: 70},
					shared.Team3: gamestate.ClientInfo{Resources: 30},
				},
			},
			inputValue:       -80,
			inputID:          shared.Team3,
			expectedResource: -50,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			depositIntoClientPrivatePool(tc.inputValue, tc.inputID, tc.gamestate)
			participantResources := tc.gamestate.ClientInfos[tc.inputID].Resources
			if participantResources != tc.expectedResource {
				t.Errorf("%v - Failed. Got '%v', but expected '%v'", tc.name, participantResources, tc.expectedResource)
			}
		})
	}
}

func TestRemoveDead(t *testing.T) {
	clientinfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Team1: {
			LifeStatus: shared.Alive,
		},
		shared.Team2: {
			LifeStatus: shared.Alive,
		},
		shared.Team3: {
			LifeStatus: shared.Alive,
		},
		shared.Team4: {
			LifeStatus: shared.Dead,
		},
		shared.Team5: {
			LifeStatus: shared.Dead,
		},
	}
	cases := []struct {
		name          string
		testgamestate gamestate.GameState
	}{
		{
			name: "no_dead_roles",
			testgamestate: gamestate.GameState{
				PresidentID: shared.ClientID(1),
				SpeakerID:   shared.ClientID(2),
				JudgeID:     shared.ClientID(3),
				ClientInfos: clientinfos,
			},
		},
		{
			name: "dead_president",
			testgamestate: gamestate.GameState{
				PresidentID: shared.ClientID(4),
				SpeakerID:   shared.ClientID(2),
				JudgeID:     shared.ClientID(3),
				ClientInfos: clientinfos,
			},
		},
		{
			name: "all_dead",
			testgamestate: gamestate.GameState{
				PresidentID: shared.ClientID(4),
				SpeakerID:   shared.ClientID(5),
				JudgeID:     shared.ClientID(5),
				ClientInfos: clientinfos,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			removeDeadBodiesFromOffice(&tc.testgamestate)
			if !reflect.DeepEqual(tc.testgamestate.ClientInfos[tc.testgamestate.PresidentID].LifeStatus, shared.Alive) {
				t.Errorf("Expected President to be %v got %v", shared.Alive, tc.testgamestate.ClientInfos[tc.testgamestate.PresidentID].LifeStatus)
			}
			if !reflect.DeepEqual(tc.testgamestate.ClientInfos[tc.testgamestate.SpeakerID].LifeStatus, shared.Alive) {
				t.Errorf("Expected Speaker to be %v got %v", shared.Alive, tc.testgamestate.ClientInfos[tc.testgamestate.SpeakerID].LifeStatus)
			}
			if !reflect.DeepEqual(tc.testgamestate.ClientInfos[tc.testgamestate.JudgeID].LifeStatus, shared.Alive) {
				t.Errorf("Expected Judge to be %v got %v", shared.Alive, tc.testgamestate.ClientInfos[tc.testgamestate.JudgeID].LifeStatus)
			}
		})
	}
}

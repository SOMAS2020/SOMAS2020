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
	valueToWithdraw := 120
	err := WithdrawFromCommonPool(valueToWithdraw, &fakeGameState)
	if err == nil {
		t.Errorf("We can withdraw more from the common pool than it actually has.")
	}
}

func TestWithdrawFromCommonPoolDeductsValue(t *testing.T) {
	fakeGameState := gamestate.GameState{CommonPool: 100}
	valueToWithdraw := 60
	_ = WithdrawFromCommonPool(valueToWithdraw, &fakeGameState)
	unexpectedAmountRemaining := fakeGameState.CommonPool != 40
	if unexpectedAmountRemaining == true {
		t.Errorf("Not withdrawing resources from CommonPool correctly.")
	}
}

// func communicateHelper(sender int, receiver int, packetCount int) map[int]map[int][]map[int]DataPacket {
// 	charSet := "abcdedfghijklmnopqrst"
// 	strLength := 10
// 	connections := map[int]map[int][]map[int]DataPacket{}

// 	for i := 0; i < packetCount; i++ {
// 		var output strings.Builder
// 		for j := 0; j < strLength; j++ {
// 			random := rand.Intn(len(charSet))
// 			randomChar := charSet[random]
// 			output.WriteString(string(randomChar))
// 		}

// 		connections[sender][receiver] = append(connections[sender][receiver], DataPacket{})

// 	}
// }

func TestCommunicateWithIslands(t *testing.T) {

	dataA := map[int]DataPacket{
		3:  {integerData: 123, textData: "Hello World - dataA", booleanData: true},
		1:  {textData: "SOMAS", booleanData: false},
		22: {booleanData: false},
		14: {integerData: 420, booleanData: false},
	}

	dataB := map[int]DataPacket{
		0: {integerData: 11, textData: "SOMAS", booleanData: true},
	}

	dataC := map[int]DataPacket{
		5:  {booleanData: true},
		4:  {textData: "communication test"},
		16: {integerData: 7832},
		73: {integerData: 234511, textData: "dataC", booleanData: false},
	}

	cases := []struct {
		name           string
		senders        []int
		sendersPayload map[int][]map[int]DataPacket
		receiver       int
	}{
		{
			name:    "single transmission",
			senders: []int{1},
			sendersPayload: map[int][]map[int]DataPacket{
				1: {dataA},
			},
			receiver: 4,
		},
		{
			name:    "2 transmissions",
			senders: []int{1, 2},
			sendersPayload: map[int][]map[int]DataPacket{
				1: {dataA},
				2: {dataB},
			},
			receiver: 5,
		},
		{
			name:    "1 sender, 2 transmissions",
			senders: []int{4},
			sendersPayload: map[int][]map[int]DataPacket{
				4: {dataA, dataC},
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
			expectedResult := map[shared.ClientID][]map[int]baseclient.Communication{}

			for sender, dataList := range tc.sendersPayload {
				for _, data := range dataList {
					communicateWithIslands(tc.receiver, sender, data)

					dataComm := map[int]baseclient.Communication{}
					for k, dp := range data {
						dataComm[k] = dataPacketToCommunication(&dp)
					}

					expectedResult[shared.TeamIDs[sender]] = append(expectedResult[shared.TeamIDs[sender]], dataComm)
				}
			}

			// Check internals of clients
			recieverGot := *(fakeClientMap[receiverID]).GetCommunications()

			if !reflect.DeepEqual(expectedResult, recieverGot) {
				t.Errorf("Communication failed. Sent: %v\nGot: %v", expectedResult, recieverGot)
			}

		})
	}
}

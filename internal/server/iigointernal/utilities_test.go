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

func TestCommunicateWithIslands(t *testing.T) {
	fakeClientMap := map[shared.ClientID]baseclient.Client{}
	sender := 1
	senderID := shared.TeamIDs[sender]
	senderClient := baseclient.NewClient(senderID)
	receiver := 3
	receiverID := shared.TeamIDs[receiver]
	receiverClient := baseclient.NewClient(receiverID)

	fakeClientMap[senderID] = senderClient
	fakeClientMap[receiverID] = receiverClient

	data := map[int]DataPacket{
		0: {integerData: 5, textData: "Hello World", booleanData: true},
		1: {integerData: 22, textData: "SOMAS", booleanData: false},
	}

	dataComm := map[int]baseclient.Communication{}
	for k, dp := range data {
		dataComm[k] = dataPacketToCommunication(&dp)
	}

	setIIGOClients(&fakeClientMap)
	communicateWithIslands(receiver, sender, data)

	recieverGot := receiverClient.GetCommunications()
	// t.Log(fakeClientMap[receiverID])
	// t.Log(fakeClientMap[senderID])
	recievedFromSender := (*recieverGot)[senderID][0]

	if !reflect.DeepEqual(dataComm, recievedFromSender) {
		t.Errorf("Communication failed. Sent: %v\nGot: %v", data, recievedFromSender)
	}
}

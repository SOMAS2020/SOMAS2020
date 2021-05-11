package server

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type mockClientIIFO struct {
	baseclient.Client
	foragingValues                shared.ForageShareInfo
	otherIslandInfo               []shared.ForageShareInfo
	otherIslandDisasterPrediction shared.ReceivedDisasterPredictionsDict
}

func (c *mockClientIIFO) MakeForageInfo() shared.ForageShareInfo {
	return c.foragingValues
}

func (c *mockClientIIFO) ReceiveForageInfo(otherIslandInfo []shared.ForageShareInfo) {
	c.otherIslandInfo = otherIslandInfo
}

func (c *mockClientIIFO) getOtherIslandInfo() []shared.ForageShareInfo {
	return c.otherIslandInfo
}

func (c *mockClientIIFO) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	c.otherIslandDisasterPrediction = receivedPredictions
}

func (c *mockClientIIFO) getOtherIslandDisasterPrediction() shared.ReceivedDisasterPredictionsDict {
	return c.otherIslandDisasterPrediction
}

func makeForagingInfo(contribution shared.Resources, resources shared.Resources, shareTo []shared.ClientID) shared.ForageShareInfo {
	if len(shareTo) > 0 {
		return shared.ForageShareInfo{
			DecisionMade:     shared.ForageDecision{Type: shared.DeerForageType, Contribution: contribution},
			ResourceObtained: resources,
			ShareTo:          shareTo,
		}
	}
	// People can be selfish and choose not to share their foraging information
	return shared.ForageShareInfo{
		DecisionMade:     shared.ForageDecision{Type: shared.DeerForageType, Contribution: contribution},
		ResourceObtained: resources,
	}
}

func receiveForagingInfo(contribution shared.Resources, resources shared.Resources, sharedFrom shared.ClientID) shared.ForageShareInfo {
	return shared.ForageShareInfo{
		DecisionMade:     shared.ForageDecision{Type: shared.DeerForageType, Contribution: contribution},
		ResourceObtained: resources,
		SharedFrom:       sharedFrom,
	}
}

func makeDisasterPrediction(prediction shared.DisasterPrediction, shareTo []shared.ClientID) shared.DisasterPredictionInfo {
	if len(shareTo) > 0 {
		return shared.DisasterPredictionInfo{
			PredictionMade: prediction,
			TeamsOfferedTo: shareTo,
		}
	}
	// People can be selfish and choose not to share their foraging information
	return shared.DisasterPredictionInfo{
		PredictionMade: prediction,
	}
}

func receiveDisasterPrediction(prediction shared.DisasterPrediction, sharedFrom shared.ClientID) shared.ReceivedDisasterPredictionInfo {
	return shared.ReceivedDisasterPredictionInfo{
		PredictionMade: prediction,
		SharedFrom:     sharedFrom,
	}
}

func TestGetForageSharingWorks(t *testing.T) {
	clientInfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Teams["Team1"]: {
			LifeStatus: shared.Alive,
		},
		shared.Teams["Team2"]: {
			LifeStatus: shared.Critical,
		},
		shared.Teams["Team3"]: {
			LifeStatus: shared.Dead,
		},
	}

	clientMap := map[shared.ClientID]baseclient.Client{
		shared.Teams["Team1"]: &mockClientIIFO{
			foragingValues: makeForagingInfo(52.7, 64, []shared.ClientID{shared.Teams["Team2"], shared.Teams["Team3"]}),
		},
		shared.Teams["Team2"]: &mockClientIIFO{
			foragingValues: makeForagingInfo(22.2, 22.3, []shared.ClientID{}),
		},
		shared.Teams["Team3"]: &mockClientIIFO{
			foragingValues: makeForagingInfo(33.2, 233.3, []shared.ClientID{shared.Teams["Team2"]}),
		},
	}

	want := shared.ForagingOfferDict{
		shared.Teams["Team1"]: makeForagingInfo(52.7, 64, []shared.ClientID{shared.Teams["Team2"], shared.Teams["Team3"]}),
		shared.Teams["Team2"]: makeForagingInfo(22.2, 22.3, []shared.ClientID{}),
	}

	server := &SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
		clientMap: clientMap,
	}

	got := server.getForageSharing()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want '%#v' got '%#v'", want, got)
	}
}

func TestDistributeForageSharing(t *testing.T) {
	clientInfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Teams["Team1"]: {
			LifeStatus: shared.Alive,
		},
		shared.Teams["Team2"]: {
			LifeStatus: shared.Critical,
		},
		shared.Teams["Team3"]: {
			LifeStatus: shared.Dead,
		},
	}

	mockClient := map[shared.ClientID]*mockClientIIFO{
		shared.Teams["Team1"]: {},
		shared.Teams["Team2"]: {},
		shared.Teams["Team3"]: {},
	}

	clientMap := map[shared.ClientID]baseclient.Client{
		shared.Teams["Team1"]: mockClient[shared.Teams["Team1"]],
		shared.Teams["Team2"]: mockClient[shared.Teams["Team2"]],
		shared.Teams["Team3"]: mockClient[shared.Teams["Team3"]],
	}

	input := shared.ForagingOfferDict{
		shared.Teams["Team1"]: makeForagingInfo(52.7, 64, []shared.ClientID{shared.Teams["Team2"], shared.Teams["Team3"]}),
		shared.Teams["Team2"]: makeForagingInfo(22.2, 22.3, []shared.ClientID{}),
	}

	want := shared.ForagingReceiptDict{
		shared.Teams["Team1"]: []shared.ForageShareInfo(nil),
		shared.Teams["Team2"]: []shared.ForageShareInfo{receiveForagingInfo(52.7, 64, shared.Teams["Team1"])},
	}

	server := &SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
		clientMap: clientMap,
	}

	//	getOtherIslandInfo
	server.distributeForageSharing(input)
	for id := range mockClient {
		got := mockClient[id].getOtherIslandInfo()
		w := want[id]
		if !reflect.DeepEqual(w, got) {
			t.Errorf("want '%#v' got '%#v' for %#v", w, got, id)
		}
	}
}

func TestDistributePredictions(t *testing.T) {
	clientInfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Teams["Team1"]: {
			LifeStatus: shared.Alive,
		},
		shared.Teams["Team2"]: {
			LifeStatus: shared.Critical,
		},
		shared.Teams["Team3"]: {
			LifeStatus: shared.Dead,
		},
	}

	mockClient := map[shared.ClientID]*mockClientIIFO{
		shared.Teams["Team1"]: {},
		shared.Teams["Team2"]: {},
		shared.Teams["Team3"]: {},
	}

	clientMap := map[shared.ClientID]baseclient.Client{
		shared.Teams["Team1"]: mockClient[shared.Teams["Team1"]],
		shared.Teams["Team2"]: mockClient[shared.Teams["Team2"]],
		shared.Teams["Team3"]: mockClient[shared.Teams["Team3"]],
	}
	team1Prediction := shared.DisasterPrediction{
		CoordinateX: 0,
		CoordinateY: 1,
		Magnitude:   1,
		TimeLeft:    4,
		Confidence:  100,
	}
	input := shared.DisasterPredictionInfoDict{
		shared.Teams["Team1"]: makeDisasterPrediction(team1Prediction, []shared.ClientID{shared.Teams["Team2"], shared.Teams["Team3"]}),
		shared.Teams["Team2"]: makeDisasterPrediction(shared.DisasterPrediction{}, []shared.ClientID{}),
		shared.Teams["Team3"]: makeDisasterPrediction(shared.DisasterPrediction{}, []shared.ClientID{shared.Teams["Team1"], shared.Teams["Team2"]}),
	}
	want := map[shared.ClientID]shared.ReceivedDisasterPredictionsDict{
		shared.Teams["Team1"]: shared.ReceivedDisasterPredictionsDict(nil),
		shared.Teams["Team2"]: {shared.Teams["Team1"]: receiveDisasterPrediction(team1Prediction, shared.Teams["Team1"])},
	}

	server := &SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
		clientMap: clientMap,
	}

	//	distributePredictions and verify results
	server.distributePredictions(input)
	for id := range mockClient {
		got := mockClient[id].getOtherIslandDisasterPrediction()
		w := want[id]
		if !reflect.DeepEqual(w, got) {
			t.Errorf("want '%#v' got '%#v' for %#v", w, got, id)
		}
	}

}

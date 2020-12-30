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
	ownPrediction                 shared.PredictionInfo
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

func (c *mockClientIIFO) ReceivePredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
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

func makeDisasterPrediction(prediction shared.Prediction, shareTo []shared.ClientID) shared.PredictionInfo {
	if len(shareTo) > 0 {
		return shared.PredictionInfo{
			PredictionMade: prediction,
			TeamsOfferedTo: shareTo,
		}
	}
	// People can be selfish and choose not to share their foraging information
	return shared.PredictionInfo{
		PredictionMade: prediction,
	}
}

func receiveDisasterPrediction(prediction shared.Prediction, sharedFrom shared.ClientID) shared.ReceivedDisasterPredictionInfo {
	return shared.ReceivedDisasterPredictionInfo{
		PredictionMade: prediction,
		SharedFrom:     sharedFrom,
	}
}

func TestGetForageSharingWorks(t *testing.T) {
	clientInfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Team1: {
			LifeStatus: shared.Alive,
		},
		shared.Team2: {
			LifeStatus: shared.Critical,
		},
		shared.Team3: {
			LifeStatus: shared.Dead,
		},
	}

	clientMap := map[shared.ClientID]baseclient.Client{
		shared.Team1: &mockClientIIFO{
			foragingValues: makeForagingInfo(52.7, 64, []shared.ClientID{shared.Team2, shared.Team3}),
		},
		shared.Team2: &mockClientIIFO{
			foragingValues: makeForagingInfo(22.2, 22.3, []shared.ClientID{}),
		},
		shared.Team3: &mockClientIIFO{
			foragingValues: makeForagingInfo(33.2, 233.3, []shared.ClientID{shared.Team2}),
		},
	}

	want := shared.ForagingOfferDict{
		shared.Team1: makeForagingInfo(52.7, 64, []shared.ClientID{shared.Team2, shared.Team3}),
		shared.Team2: makeForagingInfo(22.2, 22.3, []shared.ClientID{}),
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
		shared.Team1: {
			LifeStatus: shared.Alive,
		},
		shared.Team2: {
			LifeStatus: shared.Critical,
		},
		shared.Team3: {
			LifeStatus: shared.Dead,
		},
	}

	mockClient := map[shared.ClientID]*mockClientIIFO{
		shared.Team1: {},
		shared.Team2: {},
		shared.Team3: {},
	}

	clientMap := map[shared.ClientID]baseclient.Client{
		shared.Team1: mockClient[shared.Team1],
		shared.Team2: mockClient[shared.Team2],
		shared.Team3: mockClient[shared.Team3],
	}

	input := shared.ForagingOfferDict{
		shared.Team1: makeForagingInfo(52.7, 64, []shared.ClientID{shared.Team2, shared.Team3}),
		shared.Team2: makeForagingInfo(22.2, 22.3, []shared.ClientID{}),
	}

	want := shared.ForagingReceiptDict{
		shared.Team1: []shared.ForageShareInfo(nil),
		shared.Team2: []shared.ForageShareInfo{receiveForagingInfo(52.7, 64, shared.Team1)},
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
		shared.Team1: {
			LifeStatus: shared.Alive,
		},
		shared.Team2: {
			LifeStatus: shared.Critical,
		},
		shared.Team3: {
			LifeStatus: shared.Dead,
		},
	}

	mockClient := map[shared.ClientID]*mockClientIIFO{
		shared.Team1: {},
		shared.Team2: {},
		shared.Team3: {},
	}

	clientMap := map[shared.ClientID]baseclient.Client{
		shared.Team1: mockClient[shared.Team1],
		shared.Team2: mockClient[shared.Team2],
		shared.Team3: mockClient[shared.Team3],
	}
	team1Prediction := shared.Prediction{
		CoordinateX: 0,
		CoordinateY: 1,
		Magnitude:   1,
		TimeLeft:    4,
		Confidence:  100,
	}
	input := shared.PredictionInfoDict{
		shared.Team1: makeDisasterPrediction(team1Prediction, []shared.ClientID{shared.Team2, shared.Team3}),
		shared.Team2: makeDisasterPrediction(shared.Prediction{}, []shared.ClientID{}),
		shared.Team3: makeDisasterPrediction(shared.Prediction{}, []shared.ClientID{shared.Team1, shared.Team2}),
	}
	want := map[shared.ClientID]shared.ReceivedDisasterPredictionsDict{
		shared.Team1: shared.ReceivedDisasterPredictionsDict(nil),
		shared.Team2: {shared.Team1: receiveDisasterPrediction(team1Prediction, shared.Team1)},
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

package server

import (
	"reflect"
	"sort"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// The response from any client to a gift-related query must be sanitised to have an entry for all alive clients.

type mockClientIITO struct {
	baseclient.BaseClient
	requests                shared.GiftRequestDict
	offers                  shared.GiftOfferDict
	responses               shared.GiftResponseDict
	receivedResponses       shared.GiftResponseDict
	otherIslandContribution shared.ReceivedIntendedContributionDict
}

func (c *mockClientIITO) GetGiftRequests() shared.GiftRequestDict {
	return c.requests
}

func (c *mockClientIITO) GetGiftOffers(requests shared.GiftRequestDict) shared.GiftOfferDict {
	return c.offers
}

func (c *mockClientIITO) GetGiftResponses(offers shared.GiftOfferDict) shared.GiftResponseDict {
	return c.responses
}

func (c *mockClientIITO) UpdateGiftInfo(responses shared.GiftResponseDict) {
	c.receivedResponses = responses
}

func (c *mockClientIITO) ReceiveIntendedContribution(receivedIntendedContribution shared.ReceivedIntendedContributionDict) {
	// You can check the other's common pool contributions like this
	// intededContributions := c.intendedContribution
	c.otherIslandContribution = receivedIntendedContribution

}
func shareIntendedContribution(contribution float64, shareTo []shared.ClientID) shared.IntendedContribution {
	if len(shareTo) > 0 {
		return shared.IntendedContribution{
			Contribution:   contribution,
			TeamsOfferedTo: shareTo,
		}
	}
	// People can be selfish and choose not to share their common pool intended contribution
	return shared.IntendedContribution{
		Contribution: contribution,
	}
}

func receiveIntendedContribution(contribution float64, sharedFrom shared.ClientID) shared.ReceivedIntendedContribution {
	return shared.ReceivedIntendedContribution{
		Contribution: contribution,
		SharedFrom:   sharedFrom,
	}
}

func (c *mockClientIITO) getOtherIslandsCommonPoolContribution() shared.ReceivedIntendedContributionDict {
	return c.otherIslandContribution
}

// Test that the server correctly forms the pipeline for IITO to run
func TestServerGetGiftRequests(t *testing.T) {
	// Mock a bunch of clients
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
		// Team 1 makes 1 valid request: 50 to team 2.
		shared.Team1: &mockClientIITO{requests: shared.GiftRequestDict{shared.Team1: 50, shared.Team2: 50, shared.Team3: 50}},
		// Team 2 makes no valid requests: a zero'ed entry, one to itself and one to a dead team.
		shared.Team2: &mockClientIITO{requests: shared.GiftRequestDict{shared.Team1: 0, shared.Team2: 50, shared.Team3: 50}},
		// Team 3 is dead boi
		shared.Team3: &mockClientIITO{},
	}

	want := map[shared.ClientID]shared.GiftRequestDict{
		shared.Team1: {shared.Team2: 50},
	}

	// Mock a server
	s := &SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
		clientMap: clientMap,
	}

	if !reflect.DeepEqual(want, s.getGiftRequests()) {
		t.Errorf("want '%v' got '%v'", want, s.getGiftRequests())
	}

}

func TestOfferKnapsackPacker(t *testing.T) {
	offers := shared.GiftOfferDict{
		shared.Team1: 200,
		shared.Team2: 500,
		shared.Team3: 500,
	}

	want := shared.GiftOffer(1000)
	wantCombi := []shared.ClientID{shared.Team2, shared.Team3}

	got, optimal := offersKnapsackSolver(1000, offers)
	sort.Sort(shared.SortClientByID(optimal))
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want '%v' got '%v'", want, got)
	}
	if !reflect.DeepEqual(wantCombi, optimal) {
		t.Errorf("want '%v' got '%v'", wantCombi, optimal)
	}
}

func TestServerGetGiftOffers(t *testing.T) {

	clientInfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Team1: {
			Resources:  500,
			LifeStatus: shared.Alive,
		},
		shared.Team2: {
			Resources:  500,
			LifeStatus: shared.Alive,
		},
		shared.Team3: {
			Resources:  50,
			LifeStatus: shared.Alive,
		},
		shared.Team4: {
			Resources:  50,
			LifeStatus: shared.Dead,
		},
	}

	clientMap := map[shared.ClientID]baseclient.Client{
		// Team 1 makes 1 valid offer: 500 to team 2.
		shared.Team1: &mockClientIITO{offers: shared.GiftOfferDict{shared.Team1: 50, shared.Team2: 500, shared.Team4: 50}},
		// TODO: Team 2 attempts to offer more than it has to team 1 and 4.
		shared.Team2: &mockClientIITO{offers: shared.GiftOfferDict{shared.Team1: 400, shared.Team3: 150}},
		// Team 3 makes no offers.
		shared.Team3: &mockClientIITO{},
		// Team 4 is dead, should not show up.
		shared.Team4: &mockClientIITO{offers: shared.GiftOfferDict{shared.Team1: 600}},
	}

	want := map[shared.ClientID]shared.GiftOfferDict{
		shared.Team1: {shared.Team2: 500},
		shared.Team2: {shared.Team1: 400},
	}

	// Mock a server
	s := &SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
		clientMap: clientMap,
	}
	offers := s.getGiftOffers(map[shared.ClientID]shared.GiftRequestDict{})
	if !reflect.DeepEqual(want, offers) {
		t.Errorf("want '%v' got '%v'", want, offers)
	}
}

// Test that the server caps accepted amount in responses to the amount offered.
// TODO: Test that the server makes a response for every offer, even if the client ignored it.
func TestServerGetGiftResponses(t *testing.T) {

	clientInfos := map[shared.ClientID]gamestate.ClientInfo{
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
			LifeStatus: shared.Alive,
		},
	}

	offers := map[shared.ClientID]shared.GiftOfferDict{
		shared.Team1: {shared.Team2: 500},
		shared.Team2: {shared.Team1: 500},
		shared.Team3: {shared.Team2: 100},
		shared.Team4: {shared.Team3: 100},
	}

	clientMap := map[shared.ClientID]baseclient.Client{
		// Team 1 accepts team 2's offer, and accepts an offer that was never made by team 3.
		shared.Team1: &mockClientIITO{
			responses: shared.GiftResponseDict{
				shared.Team2: {AcceptedAmount: 500, Reason: shared.Accept},
				shared.Team3: {AcceptedAmount: 700, Reason: shared.Accept},
			},
		},

		// Team 2 tries to accept more than it was offered, and ignores team 3's offer.
		shared.Team2: &mockClientIITO{
			responses: shared.GiftResponseDict{
				shared.Team1: {AcceptedAmount: 700, Reason: shared.Accept},
			},
		},
		// Team 3 has a malformed reply.
		shared.Team3: &mockClientIITO{
			responses: shared.GiftResponseDict{
				shared.Team4: {AcceptedAmount: 300, Reason: shared.DeclineDontLikeYou},
			},
		},
		shared.Team4: &mockClientIITO{},
	}

	want := map[shared.ClientID]shared.GiftResponseDict{
		shared.Team1: {
			shared.Team2: {AcceptedAmount: 500, Reason: shared.Accept},
		},
		shared.Team2: {
			shared.Team1: {AcceptedAmount: 500, Reason: shared.Accept},
			shared.Team3: {AcceptedAmount: 0, Reason: shared.Ignored},
		},
		shared.Team3: {
			shared.Team4: {AcceptedAmount: 0, Reason: shared.DeclineDontLikeYou},
		},
	}

	// Mock a server
	s := &SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
		clientMap: clientMap,
	}
	responses := s.getGiftResponses(offers)
	if !reflect.DeepEqual(want, responses) {
		t.Errorf("want '%v' got '%v'", want, responses)
	}

}

func TestDistributeGiftHistory(t *testing.T) {

	clientInfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Team1: {
			LifeStatus: shared.Alive,
		},
		shared.Team2: {
			LifeStatus: shared.Alive,
		},
		shared.Team3: {
			LifeStatus: shared.Alive,
		},
	}

	c1 := mockClientIITO{}
	c2 := mockClientIITO{}
	c3 := mockClientIITO{}
	clientMap := map[shared.ClientID]baseclient.Client{
		shared.Team1: &c1,
		shared.Team2: &c2,
		shared.Team3: &c3,
	}

	dontNeed := shared.GiftResponse{AcceptedAmount: 0, Reason: shared.DeclineDontNeed}
	responses := map[shared.ClientID]shared.GiftResponseDict{
		// Everybody loves team 1 and team 1 loves everybody.
		shared.Team1: {
			shared.Team2: {AcceptedAmount: 100, Reason: shared.Accept},
			shared.Team3: {AcceptedAmount: 100, Reason: shared.Accept},
		},
		// Team 2 hates receiving gifts.
		shared.Team2: {
			shared.Team1: {AcceptedAmount: 200, Reason: shared.Accept},
			shared.Team3: {AcceptedAmount: 0, Reason: shared.Ignored},
		},
		// Team 3
		shared.Team3: {
			shared.Team1: dontNeed,
			shared.Team2: dontNeed,
		},
	}

	want1 := shared.GiftResponseDict{
		shared.Team2: {AcceptedAmount: 200, Reason: shared.Accept},
		shared.Team3: dontNeed,
	}

	want2 := shared.GiftResponseDict{
		shared.Team1: {AcceptedAmount: 100, Reason: shared.Accept},
		shared.Team3: dontNeed,
	}

	want3 := shared.GiftResponseDict{
		shared.Team1: {AcceptedAmount: 100, Reason: shared.Accept},
		shared.Team2: {AcceptedAmount: 0, Reason: shared.Ignored},
	}

	// Mock a server
	s := &SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
		clientMap: clientMap,
	}
	s.distributeGiftHistory(responses)

	if !reflect.DeepEqual(want1, c1.receivedResponses) {
		t.Errorf("want '%v' got '%v'", want1, c1.receivedResponses)
	}

	if !reflect.DeepEqual(want2, c2.receivedResponses) {
		t.Errorf("want '%v' got '%v'", want2, c2.receivedResponses)
	}

	if !reflect.DeepEqual(want3, c3.receivedResponses) {
		t.Errorf("want '%v' got '%v'", want3, c3.receivedResponses)
	}
}

func TestDistributeContributions(t *testing.T) {
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

	mockClient := map[shared.ClientID]*mockClientIITO{
		shared.Team1: {},
		shared.Team2: {},
		shared.Team3: {},
	}

	clientMap := map[shared.ClientID]baseclient.Client{
		shared.Team1: mockClient[shared.Team1],
		shared.Team2: mockClient[shared.Team2],
		shared.Team3: mockClient[shared.Team3],
	}

	team1Contribution := float64(1)

	input := shared.IntendedContributionDict{
		shared.Team1: shareIntendedContribution(team1Contribution, []shared.ClientID{shared.Team2, shared.Team3}),
		shared.Team2: shareIntendedContribution(0.0, []shared.ClientID{}),
		shared.Team3: shareIntendedContribution(0.0, []shared.ClientID{shared.Team1, shared.Team2}),
	}
	want := map[shared.ClientID]shared.ReceivedIntendedContributionDict{
		shared.Team1: shared.ReceivedIntendedContributionDict(nil),
		shared.Team2: {shared.Team1: receiveIntendedContribution(team1Contribution, shared.Team1)},
	}

	server := &SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
		clientMap: clientMap,
	}

	//	distributeIntendedContributions and verify results
	server.distributeIntendedContributions(input)
	for id := range mockClient {
		got := mockClient[id].getOtherIslandsCommonPoolContribution()
		w := want[id]
		if !reflect.DeepEqual(w, got) {
			t.Errorf("want '%#v' got '%#v' for %#v", w, got, id)
		}
	}

}

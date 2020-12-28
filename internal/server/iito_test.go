package server

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// The response from any client to a gift-related query must be sanitised to have an entry for all alive clients.

type mockClientIITO struct {
	baseclient.BaseClient
	requests  shared.GiftRequestDict
	offers    shared.GiftOfferDict
	responses shared.GiftResponseDict
}

func (c *mockClientIITO) GetGiftRequests() shared.GiftRequestDict {
	return c.requests
}

func (c *mockClientIITO) GetGiftOffers(requests shared.GiftRequestDict) (shared.GiftOfferDict, error) {
	return c.offers, nil
}

func (c *mockClientIITO) GetGiftResponses(offers shared.GiftOfferDict) (shared.GiftResponseDict, error) {
	return c.responses, nil
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

func TestServerGetGiftOffers(t *testing.T) {

	clientInfos := map[shared.ClientID]gamestate.ClientInfo{
		shared.Team1: {
			Resources:  1000,
			LifeStatus: shared.Alive,
		},
		shared.Team2: {
			Resources:  500,
			LifeStatus: shared.Critical,
		},
		shared.Team3: {
			Resources:  50,
			LifeStatus: shared.Dead,
		},
	}

	clientMap := map[shared.ClientID]baseclient.Client{
		// Team 1 makes 1 valid request: 500 to team 2.
		shared.Team1: &mockClientIITO{offers: shared.GiftOfferDict{shared.Team1: 50, shared.Team2: 500, shared.Team3: 50}},
		// Team 2 attempts to offer more than it has to team 1.
		shared.Team2: &mockClientIITO{offers: shared.GiftOfferDict{shared.Team1: 600}},
		// Team 3 is dead boi
		shared.Team3: &mockClientIITO{},
	}

	want := map[shared.ClientID]shared.GiftOfferDict{
		shared.Team1: {shared.Team2: 500},
		shared.Team2: {shared.Team1: 500},
	}

	// Mock a server
	s := &SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
		clientMap: clientMap,
	}
	offers, _ := s.getGiftOffers(map[shared.ClientID]shared.GiftRequestDict{})
	if !reflect.DeepEqual(want, offers) {
		t.Errorf("want '%v' got '%v'", want, s.getGiftRequests())
	}

}

// Test that the server makes a response for every offer, even if the client ignored it.
func TestServerGetGiftResponses(t *testing.T) {

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

	offers := map[shared.ClientID]shared.GiftOfferDict{
		shared.Team1: {shared.Team2: 500},
		shared.Team2: {shared.Team1: 500},
	}

	clientMap := map[shared.ClientID]baseclient.Client{
		// Team 1 makes 1 valid request: 50 to team 2.
		shared.Team1: &mockClientIITO{
			responses: shared.GiftResponseDict{
				shared.Team2: {AcceptedAmount: 500, Reason: shared.Accept},
			},
		},

		// Team 2 tries to accept more than it was offered.
		shared.Team2: &mockClientIITO{
			responses: shared.GiftResponseDict{
				shared.Team1: {AcceptedAmount: 600, Reason: shared.Accept},
			},
		},
	}

	want := map[shared.ClientID]shared.GiftResponseDict{
		shared.Team1: {
			shared.Team2: {AcceptedAmount: 500, Reason: shared.Accept},
		},
		shared.Team2: {
			shared.Team1: {AcceptedAmount: 500, Reason: shared.Accept},
		},
	}

	// Mock a server
	s := &SOMASServer{
		gameState: gamestate.GameState{
			ClientInfos: clientInfos,
		},
		clientMap: clientMap,
	}
	responses, _ := s.getGiftResponses(offers)
	if !reflect.DeepEqual(want, responses) {
		t.Errorf("want '%v' got '%v'", want, responses)
	}

}

// A client's offer can't offer more than they have in a single offer.
func TestServerSanitisesOffers(t *testing.T) {}

// Test that the server caps accepted amount in responses to the amount offered.
func TestServerSanitisesResponses(t *testing.T) {
}

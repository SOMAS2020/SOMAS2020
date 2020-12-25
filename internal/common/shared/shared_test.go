package shared

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestSortClientByID(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	clients := []ClientID{
		Team2,
		Team4,
		Team6,
		Team1,
		Team3,
		Team5,
	}
	rand.Shuffle(len(clients), func(i, j int) { clients[i], clients[j] = clients[j], clients[i] })

	want := []ClientID{
		Team1,
		Team2,
		Team3,
		Team4,
		Team5,
		Team6,
	}

	sort.Sort(SortClientByID(clients))
	if !reflect.DeepEqual(want, clients) {
		t.Errorf("want '%v' got '%v'", want, clients)
	}
}

func TestSortGiftRequests(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	requests := []GiftRequest{
		{RequestFrom: Team1, RequestAmount: 0},
		{RequestFrom: Team2, RequestAmount: 0},
		{RequestFrom: Team3, RequestAmount: 0},
		{RequestFrom: Team4, RequestAmount: 0},
		{RequestFrom: Team5, RequestAmount: 0},
		{RequestFrom: Team6, RequestAmount: 0},
	}
	rand.Shuffle(len(requests), func(i, j int) { requests[i], requests[j] = requests[j], requests[i] })

	want := []GiftRequest{
		{RequestFrom: Team1, RequestAmount: 0},
		{RequestFrom: Team2, RequestAmount: 0},
		{RequestFrom: Team3, RequestAmount: 0},
		{RequestFrom: Team4, RequestAmount: 0},
		{RequestFrom: Team5, RequestAmount: 0},
		{RequestFrom: Team6, RequestAmount: 0},
	}

	sort.Sort(SortGiftRequestByTeam(requests))
	if !reflect.DeepEqual(want, requests) {
		t.Errorf("want '%v' got '%v'", want, requests)
	}
}

func TestSortGiftOffers(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	offers := []GiftOffer{
		{ReceivingTeam: Team2, OfferAmount: 0},
		{ReceivingTeam: Team4, OfferAmount: 0},
		{ReceivingTeam: Team3, OfferAmount: 0},
		{ReceivingTeam: Team6, OfferAmount: 0},
		{ReceivingTeam: Team5, OfferAmount: 0},
		{ReceivingTeam: Team1, OfferAmount: 0},
	}
	rand.Shuffle(len(offers), func(i, j int) { offers[i], offers[j] = offers[j], offers[i] })

	want := []GiftOffer{
		{ReceivingTeam: Team1, OfferAmount: 0},
		{ReceivingTeam: Team2, OfferAmount: 0},
		{ReceivingTeam: Team3, OfferAmount: 0},
		{ReceivingTeam: Team4, OfferAmount: 0},
		{ReceivingTeam: Team5, OfferAmount: 0},
		{ReceivingTeam: Team6, OfferAmount: 0},
	}

	sort.Sort(SortGiftOfferByTeam(offers))
	if !reflect.DeepEqual(want, offers) {
		t.Errorf("want '%v' got '%v'", want, offers)
	}
}

func TestSortGiftResponse(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	responses := []GiftResponse{
		{ResponseTo: Team5, AcceptedAmount: 0, Reason: 0},
		{ResponseTo: Team4, AcceptedAmount: 0, Reason: 0},
		{ResponseTo: Team1, AcceptedAmount: 0, Reason: 0},
		{ResponseTo: Team6, AcceptedAmount: 0, Reason: 0},
		{ResponseTo: Team3, AcceptedAmount: 0, Reason: 0},
		{ResponseTo: Team2, AcceptedAmount: 0, Reason: 0},
	}
	rand.Shuffle(len(responses), func(i, j int) { responses[i], responses[j] = responses[j], responses[i] })

	want := []GiftResponse{
		{ResponseTo: Team1, AcceptedAmount: 0, Reason: 0},
		{ResponseTo: Team2, AcceptedAmount: 0, Reason: 0},
		{ResponseTo: Team3, AcceptedAmount: 0, Reason: 0},
		{ResponseTo: Team4, AcceptedAmount: 0, Reason: 0},
		{ResponseTo: Team5, AcceptedAmount: 0, Reason: 0},
		{ResponseTo: Team6, AcceptedAmount: 0, Reason: 0},
	}

	sort.Sort(SortGiftResponseByTeam(responses))
	if !reflect.DeepEqual(want, responses) {
		t.Errorf("want '%v' got '%v'", want, responses)
	}
}

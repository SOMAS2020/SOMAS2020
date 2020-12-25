package server

import (
	"testing"
)

// The response from any client to a gift-related query must be sanitised to have an entry for all alive clients.

// TestGetGiftRequestsFromClients
func TestServerGetGiftRequests(t *testing.T) {
	// Mock a server
	// Mock a bunch of clients

	// totalRequests := testServer.getGiftRequests()

	// want := len(getNonDeadClientIDs())
	// if !reflect.DeepEqual(want, len(totalRequests)) {
	// 	t.Errorf("want '%v' got '%v'", want, totalRequests)
	// }
}

func TestServerGetGiftOffers(t *testing.T) {}

func TestServerGetGiftResponses(t *testing.T) {}

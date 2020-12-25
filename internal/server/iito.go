package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// runIITO : IITO makes recommendations about the optimal (and fairest) contributions this term
// to mitigate the common pool dilemma
func (s *SOMASServer) runIITO() error {
	s.logf("start runIITO")
	defer s.logf("finish runIITO")
	err := s.runGiftSession()
	if err != nil {
		return err
	}
	// TODO:- IITO team
	return nil
}

func (s *SOMASServer) runIITOEndOfTurn() error {
	s.logf("start runIITOEndOfTurn")
	defer s.logf("finish runIITOEndOfTurn")
	// TODO:- IITO team
	return nil
}

func (s *SOMASServer) runGiftSession() error {
	s.logf("start runGiftSession")
	defer s.logf("finish runGiftSession")

	giftRequestDict := s.getGiftRequests()
	giftOffersDict, err := s.getGiftOffers(giftRequestDict)
	if err != nil {
		return err
	}
	giftHistoryDict, err := s.getGiftResponses(giftOffersDict)
	if err != nil {
		return err
	}
	err = s.distributeGiftHistory(giftHistoryDict)
	if err != nil {
		return err
	}
	// Process actions
	for key, value := range giftHistoryDict {
		s.logf("Gifts from %s: %v\n", key, value)
	}
	return nil
}

// GetGiftRequests collects all gift requests from the clients in a map
func (s *SOMASServer) getGiftRequests() map[shared.ClientID][]shared.GiftRequest {
	giftRequestDict := map[shared.ClientID][]shared.GiftRequest{}
	for id, client := range s.clientMap {
		// TODO: Sort the requests by GiftRequest.RequestFrom first before storing them
		giftRequestDict[id] = client.GetGiftRequests()
	}
	return giftRequestDict
}

// getGiftOffers collects all responses from clients to their requests in a map
func (s *SOMASServer) getGiftOffers(totalRequests map[shared.ClientID][]shared.GiftRequest) (map[shared.ClientID][]shared.GiftOffer, error) {
	totalOffers := map[shared.ClientID][]shared.GiftOffer{}
	for id, client := range s.clientMap {
		// Gather all the requests made to this team
		requestsToThisTeam := []shared.GiftRequest{}
		for fromTeam, indivRequests := range totalRequests {
			// FIXME: Implement a sort so that this simple indexing method still works for lists.
			requestsToThisTeam[fromTeam] = indivRequests[id]
		}

		offer, err := client.GetGiftOffers(requestsToThisTeam)
		if err == nil {
			// totalResponses in this case may have a bogus or meaningless value
			return totalOffers, err
		}

		// TODO: Sort the offers by GiftOffer.ReceivingTeam first before storing them
		totalOffers[id] = offer
	}
	return totalOffers, nil
}

func (s *SOMASServer) getGiftResponses(totalOffers map[shared.ClientID][]shared.GiftOffer) (map[shared.ClientID][]shared.GiftResponse, error) {
	totalResponses := map[shared.ClientID][]shared.GiftResponse{}

	for id, client := range s.clientMap {
		offersToThisTeam := []shared.GiftOffer{}
		for fromTeam, indivOffers := range totalOffers {
			// FIXME: Implement a sort so that this simple indexing method still works for lists.
			offersToThisTeam[fromTeam] = indivOffers[id]
		}
		response, err := client.GetGiftResponses(offersToThisTeam)
		if err != nil {
			// totalResponses in this case may have a bogus or meaningless value
			return totalResponses, err
		}
		// TODO: Sort the responses by GiftResponse.ResponseTo first before storing them
		totalResponses[id] = response
	}
	return totalResponses, nil
}

// distributeGiftHistory collates all responses to a single client and calls that client to receive its responses.
func (s *SOMASServer) distributeGiftHistory(totalResponses map[shared.ClientID][]shared.GiftResponse) error {
	for id, client := range s.clientMap {
		responsesToThisTeam := []shared.GiftResponse{}
		for fromTeam, indivResponses := range totalResponses {
			// FIXME: Implement a sort so that this simple indexing method still works for lists.
			responsesToThisTeam[fromTeam] = indivResponses[id]
		}
		err := client.UpdateGiftInfo(responsesToThisTeam)
		if err != nil {
			return err
		}
	}
	return nil
}

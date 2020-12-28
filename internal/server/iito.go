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

func (s *SOMASServer) sanitiseTeamGiftRequests(requests []shared.GiftRequest) []shared.GiftRequest {
	return requests
}

// GetGiftRequests collects a map of gift requests from an individual client, for all clients, in a map
func (s *SOMASServer) getGiftRequests() map[shared.ClientID]shared.GiftRequestDict {
	giftRequestDict := map[shared.ClientID]shared.GiftRequestDict{}
	for _, id := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		giftRequestDict[id] = s.clientMap[id].GetGiftRequests()
	}
	return giftRequestDict
}

func (s *SOMASServer) sanitiseTeamGiftOffers(offers []shared.GiftOffer) []shared.GiftOffer {
	return offers
}

// getGiftOffers collects all responses from clients to their requests in a map
func (s *SOMASServer) getGiftOffers(totalRequests map[shared.ClientID]shared.GiftRequestDict) (map[shared.ClientID]shared.GiftOfferDict, error) {
	totalOffers := map[shared.ClientID]shared.GiftOfferDict{}
	for _, id := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		// Gather all the requests made to this team
		requestsToThisTeam := shared.GiftRequestDict{}
		for fromTeam, indivRequests := range totalRequests {
			requestsToThisTeam[fromTeam] = indivRequests[id]
		}

		offer, err := s.clientMap[id].GetGiftOffers(requestsToThisTeam)
		if err != nil {
			// totalResponses in this case may have a bogus or meaningless value
			return totalOffers, err
		}
		totalOffers[id] = offer
	}
	return totalOffers, nil
}

func (s *SOMASServer) sanitiseTeamGiftResponses(responses []shared.GiftResponse) []shared.GiftResponse {
	return responses
}

func (s *SOMASServer) getGiftResponses(totalOffers map[shared.ClientID]shared.GiftOfferDict) (map[shared.ClientID]shared.GiftResponseDict, error) {
	totalResponses := map[shared.ClientID]shared.GiftResponseDict{}

	for _, id := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		offersToThisTeam := shared.GiftOfferDict{}
		for fromTeam, indivOffers := range totalOffers {
			offersToThisTeam[fromTeam] = indivOffers[id]
		}
		response, err := s.clientMap[id].GetGiftResponses(offersToThisTeam)
		if err != nil {
			// totalResponses in this case may have a bogus or meaningless value
			return totalResponses, err
		}
		totalResponses[id] = response
	}
	return totalResponses, nil
}

// distributeGiftHistory collates all responses to a single client and calls that client to receive its responses.
func (s *SOMASServer) distributeGiftHistory(totalResponses map[shared.ClientID]shared.GiftResponseDict) error {
	for _, id := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		responsesToThisTeam := shared.GiftResponseDict{}
		for fromTeam, indivResponses := range totalResponses {
			// TODO: Will fail if the list does not have all 6 entries. Fix in a sanitisation function.
			responsesToThisTeam[fromTeam] = indivResponses[id]
		}
		err := s.clientMap[id].UpdateGiftInfo(responsesToThisTeam)
		if err != nil {
			return err
		}
	}
	return nil
}

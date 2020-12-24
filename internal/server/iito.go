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
func (s *SOMASServer) getGiftRequests() map[shared.ClientID]shared.GiftRequestDict {
	giftRequestDict := map[shared.ClientID]shared.GiftRequestDict{}
	for id, client := range s.clientMap {
		giftRequestDict[id] = client.GetGiftRequests()
	}
	return giftRequestDict
}

// GetGiftOffers collects all responses from clients to their requests in a map
func (s *SOMASServer) getGiftOffers(totalRequests map[shared.ClientID]shared.GiftRequestDict) (map[shared.ClientID]shared.GiftOfferDict, error) {
	totalOffers := map[shared.ClientID]shared.GiftOfferDict{}
	// Loop over each team
	for id, client := range s.clientMap {
		var offer map[shared.ClientID]shared.GiftOffer
		var err error
		// Gather all the requests made to this team
		requestsToThisTeam := shared.GiftRequestDict{}
		for fromTeam, indivRequests := range totalRequests {
			requestsToThisTeam[fromTeam] = indivRequests[id]
		}
		offer, err = client.GetGiftOffers(requestsToThisTeam)
		if err != nil {
			totalOffers[id] = offer
		}
	}
	return totalOffers, nil
}
func (s *SOMASServer) getGiftResponses(totalOffers map[shared.ClientID]shared.GiftOfferDict) (map[shared.ClientID]shared.GiftResponseDict, error) {
	totalResponses := map[shared.ClientID]shared.GiftResponseDict{}

	for id, client := range s.clientMap {
		var response map[shared.ClientID]shared.GiftResponse
		var err error

		offersToThisTeam := shared.GiftOfferDict{}
		for fromTeam, indivOffers := range totalOffers {
			offersToThisTeam[fromTeam] = indivOffers[id]
		}
		response, err = client.GetGiftResponses(offersToThisTeam)
		if err != nil {
			totalResponses[id] = response
		}
	}
	return totalResponses, nil
}

func (s *SOMASServer) distributeGiftHistory(totalResponses map[shared.ClientID]shared.GiftResponseDict) error {
	// Process acceptedGifts
	for id, client := range s.clientMap {
		responsesToThisTeam := shared.GiftResponseDict{}
		for fromTeam, indivResponses := range totalResponses {
			responsesToThisTeam[fromTeam] = indivResponses[id]
		}
		err := client.UpdateGiftInfo(responsesToThisTeam)
		if err != nil {
			return err
		}
	}
	return nil
}

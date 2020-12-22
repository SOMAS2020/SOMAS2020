package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// runIITO : IITO makes recommendations about the optimal (and fairest) contributions this term
// to mitigate the common pool dilemma
func (s *SOMASServer) runIITO() error {
	s.logf("start runIITO")
	defer s.logf("finish runIITO")
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
	giftHistoryDict, err := s.getGiftAcceptance(giftOffersDict)
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

func (s *SOMASServer) getGiftRequests() shared.GiftDict {
	giftRequestDict := shared.GiftDict{}
	for id, client := range s.clientMap {
		giftRequestDict[id] = client.RequestGift()
	}
	return giftRequestDict
}
func (s *SOMASServer) getGiftOffers(giftRequestDict shared.GiftDict) (map[shared.ClientID]shared.GiftDict, error) {
	giftOfferDict := map[shared.ClientID]shared.GiftDict{}
	var err error
	for id, client := range s.clientMap {
		giftOfferDict[id], err = client.OfferGifts(giftRequestDict)
		if err != nil {
			return giftOfferDict, err
		}
	}
	return giftOfferDict, nil
}
func (s *SOMASServer) getGiftAcceptance(giftOffersDict map[shared.ClientID]shared.GiftDict) (map[shared.ClientID]shared.GiftInfoDict, error) {
	acceptedGifts := map[shared.ClientID]shared.GiftInfoDict{}
	var err error

	received_by_client_dict := make(map[shared.ClientID]shared.GiftDict)

	//puts all the gifts received by a certain client accesible by the id of that client
	for idsend, _ := range giftOffersDict {
		for idto, _ := range giftOffersDict {
			received_by_client_dict[idsend][idto] = giftOffersDict[idto][idsend]
		}
	}

	for id, client := range s.clientMap {
		acceptedGifts[id], err = client.AcceptGifts(received_by_client_dict[id])
		if err != nil {
			return acceptedGifts, err
		}
	}
	return acceptedGifts, nil
}

func (s *SOMASServer) distributeGiftHistory(acceptedGifts map[shared.ClientID]shared.GiftInfoDict) error {
	//Process acceptedGifts
	for _, client := range s.clientMap {
		err := client.UpdateGiftInfo(acceptedGifts)
		if err != nil {
			return err
		}
	}
	return nil
}

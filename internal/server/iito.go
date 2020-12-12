package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// runIITO : IITO makes recommendations about the optimal (and fairest) contributions this term
// to mitigate the common pool dilemma
func (s *SOMASServer) runIITO() ([]common.Action, error) {
	s.logf("start runIITO")
	s.runGiftSession()
	defer s.logf("finish runIITO")
	// TOOD:- IITO team
	return nil, nil
}

func (s *SOMASServer) runGiftSession() ([]common.Action, error) {
	s.logf("start runGiftSession")
	giftRequestDict, err := s.getGiftRequests()
	if err != nil {
		return nil, err
	}
	giftOffersDict, err := s.getGiftOffers(giftRequestDict)
	if err != nil {
		return nil, err
	}
	giftHistoryDict, err := s.getGiftAcceptance(giftOffersDict)
	if err != nil {
		return nil, err
	}
	err = s.distributeGiftHistory(giftHistoryDict)
	if err != nil {
		return nil, err
	}
	// Process actions
	defer s.logf("finish runGiftSession")
	return nil, nil
}

func (s *SOMASServer) getGiftRequests() (shared.GiftDict, error) {
	giftRequestDict := shared.GiftDict{}
	for id, client := range s.clientMap {
		giftRequestDict[id] = client.RequestGift()
	}
	return giftRequestDict, nil
}
func (s *SOMASServer) getGiftOffers(giftRequestDict shared.GiftDict) (map[shared.ClientID]shared.GiftDict, error) {
	giftOfferDict := map[shared.ClientID]shared.GiftDict{}
	for id, client := range s.clientMap {
		giftOfferDict[id] = client.OfferGifts(giftRequestDict)
	}
	return giftOfferDict, nil
}
func (s *SOMASServer) getGiftAcceptance(giftOffersDict map[shared.ClientID]shared.GiftDict) (map[shared.ClientID]shared.GiftInfoDict, error) {
	acceptedGifts := map[shared.ClientID]shared.GiftInfoDict{}
	for id, client := range s.clientMap {
		acceptedGifts[id] = client.AcceptGifts(giftOffersDict[id])
	}
	return acceptedGifts, nil
}

func (s *SOMASServer) distributeGiftHistory(acceptedGifts map[shared.ClientID]shared.GiftInfoDict) error {
	//Process acceptedGifts
	for id, client := range s.clientMap {
		client.UpdateGiftInfo(acceptedGifts[id])
	}
	return nil
}

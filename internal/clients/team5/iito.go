package team5

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

//=============================================================
//=======GIFTS=================================================
//=============================================================

/*
	GetGiftRequests() shared.GiftRequestDict
	GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict
	GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict
	UpdateGiftInfo(receivedResponses shared.GiftResponseDict)
*/
func (c *client) GetGiftRequests() {}

// GiftAccceptance
func (c *client) GiftAcceptance(offer shared.g) shared.GiftResponse {

	if c.wealth() == JeffBezos {
		c.Logf("I don't need any gifts you peasants")
	} else if c.wealth() == MiddleClass {
		c.Logf("We only accept Gifts from Team 1, 2, 4, 6")
	} else if c.wealth() == ImperialStudent {
		c.Logf("We accept Gifts from everyone")
	} else {
		c.Logf("We accept Gifts from everyone")
	}

	GiftResponse := shared.GiftResponse{
		AcceptedAmount: 0,
		Reason:         shared.DeclineDontLikeYou,
	}
	return GiftResponse
}

/*Gift Requests*/
func (c *client) GiftRequests() shared.GiftRequest {
	var giftrequest shared.GiftRequest

	if c.wealth() == JeffBezos {
		c.Logf("I don't need any gifts you peasants")
	} else if c.wealth() == MiddleClass {
		c.Logf("We only accept Gifts from Team 1, 2, 4, 6")
	} else if c.wealth() == ImperialStudent {
		c.Logf("We accept Gifts from everyone")
	} else {
		c.Logf("We accept Gifts from everyone")
	}

	return giftrequest
}

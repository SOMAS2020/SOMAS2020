package test

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// TestBasicRuleEvaluatorPositive Checks whether rule we expect to evaluate as true actually evaluates as such
func TestRequestGift(t *testing.T) {
	client := common.NewClient(420)
	request, err := client.RequestGift()
	if err != nil {
		t.Errorf("An error occurred: '%v'", err)
	}
	if request < 0 {
		t.Errorf("Client is not allowed to request negative valued gift amounts")
	}
}

func TestOfferGift(t *testing.T) {
	client := common.NewClient(420)
	giftRequestDict := shared.GiftDict{
		shared.Team1: 10,
		shared.Team2: 2,
		shared.Team3: 4}
	giftOfferDict, err := client.OfferGifts(giftRequestDict)
	if err != nil {
		t.Errorf("An error occurred: '%v'", err)
	}
	for key, value := range giftRequestDict {
		if giftOfferDict[key] > value {
			t.Errorf("Client is not allowed to offer gifts greater than requested")
		}
	}
}

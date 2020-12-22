package baseclient

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// RequestGift allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
// COMPULSORY, you need to implement this method
func (c *BaseClient) RequestGift() (int, error) {
	return 0, nil
}

// OfferGifts allows clients to offer to give the gifts requested by other clients.
// It can offer multiple partial gifts
// COMPULSORY, you need to implement this method
func (c *BaseClient) OfferGifts(giftRequestDict shared.GiftDict) (shared.GiftDict, error) {
	return giftRequestDict, nil
}

// AcceptGifts allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
// COMPULSORY, you need to implement this method
func (c *BaseClient) AcceptGifts(receivedGiftDict shared.GiftDict) (shared.GiftInfoDict, error) {
	acceptedGifts := shared.GiftInfoDict{}
	for client, offer := range receivedGiftDict {
		acceptedGifts[client] = shared.GiftInfo{
			ReceivingTeam:  client,
			OfferingTeam:   c.GetID(),
			OfferAmount:    offer,
			AcceptedAmount: offer,
			Reason:         shared.Accept}
	}
	return acceptedGifts, nil
}

// UpdateGiftInfo gives information about the outcome from AcceptGifts.
// This allows for opinion formation.
// COMPULSORY, you need to implement this method
func (c *BaseClient) UpdateGiftInfo(acceptedGifts map[shared.ClientID]shared.GiftInfoDict) error {
	// PreviousGifts[count] = acceptedGifts
	// count++
	return nil
}

// SendGift is executed at the end of each turn and allows clients to
// send the gifts promised in the IITO
// COMPULSORY, you need to implement this method
func (c *BaseClient) SendGift(receivingClient shared.ClientID, amount int) error {
	return nil
}

// ReceiveGift is executed at the end of each turn and allows clients to
// receive the gifts promised in the IITO
// COMPULSORY, you need to implement this method
func (c *BaseClient) ReceiveGift(sendingClient shared.ClientID, amount int) error {
	return nil
}

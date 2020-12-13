package common

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// NewClient produces a new client with the BaseClient already implemented.
// BASE: Do not overwrite in team client.
func NewClient(id shared.ClientID) Client {
	return &BaseClient{id: id}
}

// BaseClient provides a basic implementation for all functions of the client interface and should always the interface fully.
// All clients should be based off of this BaseClient to ensure that all clients implement the interface,
// even when new features are added.
type BaseClient struct {
	id shared.ClientID
}

// Echo prints a message to show that the client exists
// BASE: Do not overwrite in team client.
func (c *BaseClient) Echo(s string) string {
	c.Logf("Echo: '%v'", s)
	return s
}

// GetID gets returns the id of the client.
// BASE: Do not overwrite in team client.
func (c *BaseClient) GetID() shared.ClientID {
	return c.id
}

// Logf is the client's logger that prepends logs with your ID. This makes
// it easier to read logs. DO NOT use other loggers that will mess logs up!
// BASE: Do not overwrite in team client.
func (c *BaseClient) Logf(format string, a ...interface{}) {
	log.Printf("[%v]: %v", c.id, fmt.Sprintf(format, a...))
}

// StartOfTurnUpdate is updates the gamestate of the client at the start of each turn.
// The gameState is served by the server.
// OPTIONAL. Base should be able to handle it but feel free to implement your own.
func (c *BaseClient) StartOfTurnUpdate(gameState GameState) {
	c.Logf("Received game state update: %v", gameState)
	// TODO
}

// EndOfTurnActions executes and returns the actions done by the client that turn.
// OPTIONAL. Base should be able to handle it but feel free to implement your own.
func (c *BaseClient) EndOfTurnActions() []Action {
	c.Logf("EndOfTurnActions")
	return nil
}

// RequestGift allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
// COMPULSORY, you need to implement this method
func (c *BaseClient) RequestGift() int {
	return 0
}

// OfferGifts allows clients to offer to give the gifts requested by other clients.
// It can offer multiple partial gifts
// COMPULSORY, you need to implement this method
func (c *BaseClient) OfferGifts(giftRequestDict shared.GiftDict) shared.GiftDict {
	return giftRequestDict
}

// AcceptGifts allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
// COMPULSORY, you need to implement this method
func (c *BaseClient) AcceptGifts(receivedGiftDict shared.GiftDict) shared.GiftInfoDict {
	acceptedGifts := shared.GiftInfoDict{}
	for client, offer := range receivedGiftDict {
		acceptedGifts[client] = shared.GiftInfo{
			ReceivingTeam:  client,
			OfferingTeam:   c.GetID(),
			OfferAmount:    offer,
			AcceptedAmount: offer,
			Reason:         shared.Accept}
	}
	return acceptedGifts
}

// UpdateGiftInfo gives information about the outcome from AcceptGifts.
// This allows for opinion formation.
// COMPULSORY, you need to implement this method
func (c *BaseClient) UpdateGiftInfo(acceptedGifts shared.GiftInfoDict) {

}

//Actions? Need to talk to LH and our team about this one:

// SendGift is executed at the end of each turn and allows clients to
// send the gifts promised in the IITO
// COMPULSORY, you need to implement this method
func (c *BaseClient) SendGift(receivingClient shared.ClientID, amount int) {

}

// ReceiveGift is executed at the end of each turn and allows clients to
// receive the gifts promised in the IITO
// COMPULSORY, you need to implement this method
func (c *BaseClient) ReceiveGift(sendingClient shared.ClientID, amount int) {

}

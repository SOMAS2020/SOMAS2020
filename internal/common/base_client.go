package common

import (
	"fmt"
	"log"
	"math/rand"

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
func (c *BaseClient) RequestGift() (int, error) {
	return 0, nil
}

// OfferGifts allows clients to offer to give the gifts requested by other clients.
// It can offer multiple partial gifts
// COMPULSORY, you need to implement this method
func (c *BaseClient) OfferGifts(giftRequestDict shared.GiftDict) (shared.GiftDict, error) {
	//******PSEUDO CODE**********

	// Assume we have knowledge about:
	//	a. Islands in Critical Condition
	//
	//	1. Check if our island is in Critical State:
	//		If yes then offer no gifts
	//		If no continue;
	//	2. Let x = number of islands in critical state
	//		If (x!=0)
	//			total no. of islands not in critical state = 6-x
	//			LOOP: for i=0; i<x; i++
	//				let y = resources requested by critical island i
	//				our contribution to island i (critical condition) = y/(6-x)
	//			end
	//	3. 	If (x==0)
	//			Contribute 20% of requested sum to all islands
	//
	//			NOTE: For agent implementation, take into consideration:
	//				a. possibility of disaster
	//				b. Confidence level or risk aversion parameters
	//				c. Agents current gamestate
	//				d. Rating of other agents based on previous turn
	//				e. variability of your own resources
	//			This is not part of MVP so should be handled by agents indiviually.
	//***************************

	//Setup Dummy Data for island conditions
	crit_cond_dict := shared.GiftDict{}
	//0 means not critical
	//1 means critical
	for client := range giftRequestDict {
		crit_cond_dict[client] = rand.Intn(2)
	}

	// We check ourselves if we are in critical
	if crit_cond_dict[c.id] == 1 {
		for client := range giftRequestDict {
			giftRequestDict[client] = 0
		}
		return giftRequestDict, nil
	}

	var critical_islands = 0

	for client := range giftRequestDict {
		if crit_cond_dict[client] == 1 {
			critical_islands++
			// sum_req_crit_isl += giftRequestDict[client]
		}
	}

	var amount_req_crit_isl = 0

	if critical_islands != 0 {
		var non_critical_islands = 6 - critical_islands
		for client := range giftRequestDict {
			if crit_cond_dict[client] == 1 {
				amount_req_crit_isl = giftRequestDict[client]
				giftRequestDict[client] = amount_req_crit_isl / non_critical_islands
			}
			giftRequestDict[client] = 0
		}
		return giftRequestDict, nil
	}

	//return 20% of requested value for all if no one is critical
	for client := range giftRequestDict {
		var amount = giftRequestDict[client]
		giftRequestDict[client] = amount / 5
	}

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
func (c *BaseClient) UpdateGiftInfo(acceptedGifts shared.GiftInfoDict) error {

	return nil
}

//Actions? Need to talk to LH and our team about this one:

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

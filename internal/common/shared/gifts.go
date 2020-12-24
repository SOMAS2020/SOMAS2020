package shared

// GiftRequest contains the details of a gift request from an island to another
type GiftRequest struct {
	RequestingTeam ClientID
	OfferingTeam   ClientID
	RequestAmount  Resources
}

// GiftRequestDict is a dictionary of offers that a client will make to other clients.
// A layer of redundancy exists, in that the key and the OfferingTeam encode the same information,
// but this may prove useful in implementation.
type GiftRequestDict = map[ClientID]GiftRequest

// GiftOffer contains the details of a gift offer from an island to another
type GiftOffer struct {
	ReceivingTeam ClientID
	OfferingTeam  ClientID
	OfferAmount   Resources
}

// GiftOfferDict is a dictionary of gifts that this client is going to receive from other clients.
// A layer of redundancy exists, in that the key and the ReceivingTeam encode the same information,
// but this may prove useful in implementation.
type GiftOfferDict = map[ClientID]GiftOffer

// AcceptReason is a just a logical wrapper on the int, a normal int declaration could easily have been used.
type AcceptReason int

const (
	// Accept ...
	Accept AcceptReason = 0
	// DeclineDontNeed ...
	DeclineDontNeed
	// DeclineDontLikeYou ...
	DeclineDontLikeYou
)

// GiftResponse represents an individual client's response to the server when presented with a GiftOffer.
type GiftResponse struct {
	OfferingTeam   ClientID
	AcceptedAmount Resources
	Reason         AcceptReason
}

// GiftResponseDict represents an individual client's response to all the other teams' GiftOffers.
// Similarly to GiftOfferDict, the key encodes the same information as the value's OfferingTeam member.
type GiftResponseDict map[ClientID]GiftResponse

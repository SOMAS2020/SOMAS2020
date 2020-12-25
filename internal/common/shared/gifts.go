package shared

// GiftRequest contains the details of a gift request from an island to another
type GiftRequest struct {
	RequestFrom   ClientID
	RequestAmount Resources
}

// GiftOffer contains the details of a gift offer from an island to another
type GiftOffer struct {
	ReceivingTeam ClientID
	OfferAmount   Resources
}

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
	AcceptedAmount Resources
	Reason         AcceptReason
}

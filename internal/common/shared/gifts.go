package shared

// GiftRequest contains the details of a gift request from an island to another
type GiftRequest Resources

// GiftRequestDict contains the details of an island's gift requests to everyone else.
type GiftRequestDict map[ClientID]GiftRequest

// GiftOffer contains the details of a gift offer from an island to another
type GiftOffer Resources

// GiftOfferDict contains the details of an island's gift offers to everyone else.
type GiftOfferDict map[ClientID]GiftOffer

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

// GiftResponseDict represents an island's responses to all the other island's gift offers.
type GiftResponseDict map[ClientID]GiftResponse

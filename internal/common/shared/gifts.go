package shared

// GiftRequest contains the details of a gift request from an island to another
type GiftRequest struct {
	RequestFrom   ClientID
	RequestAmount Resources
}

// SortGiftRequestByTeam implements a sorting interface for GiftOffer arrays / slices
type SortGiftRequestByTeam []GiftRequest

func (a SortGiftRequestByTeam) Len() int           { return len(a) }
func (a SortGiftRequestByTeam) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortGiftRequestByTeam) Less(i, j int) bool { return a[i].RequestFrom < a[j].RequestFrom }

// GiftOffer contains the details of a gift offer from an island to another
type GiftOffer struct {
	ReceivingTeam ClientID
	OfferAmount   Resources
}

// SortGiftOfferByTeam implements a sorting interface for GiftOffer arrays / slices
type SortGiftOfferByTeam []GiftOffer

func (a SortGiftOfferByTeam) Len() int           { return len(a) }
func (a SortGiftOfferByTeam) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortGiftOfferByTeam) Less(i, j int) bool { return a[i].ReceivingTeam < a[j].ReceivingTeam }

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
	ResponseTo     ClientID
	AcceptedAmount Resources
	Reason         AcceptReason
}

// SortGiftResponseByTeam implements a sorting interface for GiftOffer arrays / slices
type SortGiftResponseByTeam []GiftResponse

func (a SortGiftResponseByTeam) Len() int           { return len(a) }
func (a SortGiftResponseByTeam) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortGiftResponseByTeam) Less(i, j int) bool { return a[i].ResponseTo < a[j].ResponseTo }

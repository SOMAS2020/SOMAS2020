package shared

// GiftDict is a dictonary of gift
type GiftDict = map[ClientID]uint

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

// GiftInfoDict key is island index of the island who received your gift
type GiftInfoDict map[ClientID]GiftInfo

// GiftInfo is a struct containing the information describing a gift
// exchange between two islands
type GiftInfo struct {
	ReceivingTeam  ClientID
	OfferingTeam   ClientID
	OfferAmount    uint
	AcceptedAmount uint
	Reason         AcceptReason
}

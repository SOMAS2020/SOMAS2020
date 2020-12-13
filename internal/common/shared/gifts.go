package shared

// GiftDict is a dictonary of gift
type GiftDict = map[ClientID]int

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
type GiftInfoDict map[int]GiftInfo

// GiftInfo is a struct containing the information describing a gift
// exchange between two islands
type GiftInfo struct {
	receivingTeam  ClientID
	offeringTeam   ClientID
	offerAmount    int
	acceptedAmount int
	reason         AcceptReason
}

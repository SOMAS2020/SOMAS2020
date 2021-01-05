package team5

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

//================================================================
/*	Types */
//=================================================================
const id = shared.Team5

//wealthTier is how much money we have
type wealthTier int

//================ Common Pool =========================================

//cPRequestHistory history of CP Requests
type cPRequestHistory []shared.Resources

//cPAllocationHistory History of allocations
type cPAllocationHistory []shared.Resources

//================ Resource History =========================================

//resourceHistory OUR islands resources per turn
type resourceHistory map[uint]shared.Resources

//================ Foraging History =========================================

// forageOutcome records the ROI on a foraging session
type forageOutcome struct {
	turn   uint
	input  shared.Resources
	output shared.Resources
}

// forageHistory stores history of foraging outcomes
type forageHistory map[shared.ForageType][]forageOutcome

//================ Gifts ===========================================

// giftInfo holds information about the gifts
type giftInfo struct {
	requested      shared.GiftRequest
	offered        shared.GiftOffer
	response       shared.GiftResponse
	actualRecieved shared.Resources
}

//giftExchange Looks at how much they requested and we request
type giftExchange struct {
	// 							uint = turn
	TheirRequest map[uint]giftInfo
	OurRequest   map[uint]giftInfo
}

//	giftHistory is the history of our gifts according to which island sent it
type giftHistory map[shared.ClientID]giftExchange

// Client Information */

type clientConfig struct {

	// Initial non planned foraging
	InitialForageTurns      uint
	MinimumForagePercentage float64
	NormalForagePercentage  float64
	JBForagePercentage      float64

	// Normal foraging
	NormalRandomIncrease float64
	MaxForagePercentage  float64
	SkipForage           uint // Skip for X turns if no positive RoI

	// If resources go above this limit we are balling with money
	JBThreshold shared.Resources
	// Middle class:  Middle < Jeff bezos
	MiddleThreshold shared.Resources
	// Poor: Imperial student < Middle
	ImperialThreshold shared.Resources

	// How much to request when we are dying
	DyingGiftRequestAmount shared.Resources
	// How much to request when we are at Imperial
	ImperialGiftRequestAmount shared.Resources
	// How much to request when we are dying
	MiddleGiftRequestAmount shared.Resources
}

// Client is the island number
type client struct {
	*baseclient.BaseClient

	// History
	resourceHistory     resourceHistory
	forageHistory       forageHistory
	cpRequestHistory    cPRequestHistory
	cpAllocationHistory cPAllocationHistory

	giftHistory giftHistory

	taxAmount shared.Resources

	// allocation is the president's response to your last common pool resource request
	allocation shared.Resources

	config clientConfig
}

//================================================================
/*	Constants */
//=================================================================
// Wealth Tiers
const (
	dying           wealthTier = iota // Sets values = 0
	imperialStudent                   // iota sets the folloing values =1
	middleClass                       // = 2
	jeffBezos                         // = 3
)

const (
	// Accept ...
	Accept shared.AcceptReason = iota
	// DeclineDontNeed ...
	DeclineDontNeed
	// DeclineDontLikeYou ...
	DeclineDontLikeYou
	// Ignored ...
	Ignored
)

//================================================================
/*	Functions */
//=================================================================
//String converts string to number
func (wt wealthTier) String() string {
	strings := [...]string{"dying", "imperialStudent", "middleClass", "jeffBezos"}
	if wt >= 0 && int(wt) < len(strings) {
		return strings[wt]
	}
	return fmt.Sprintf("Unkown wealth state '%v'", int(wt))
}

// GoString implements GoStringer
func (wt wealthTier) GoString() string {
	return wt.String()
}

// MarshalText implements TextMarshaler
func (wt wealthTier) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(wt.String())
}

// MarshalJSON implements RawMessage
func (wt wealthTier) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(wt.String())
}

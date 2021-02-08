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
const startTurn = 1

//================ Common Pool =========================================

type resourceHistoryType map[uint]shared.Resources

//cpRequestHistory history of CP Requests
type cpRequestHistory resourceHistoryType

//cpAllocationHistory History of allocations
type cpAllocationHistory resourceHistoryType

//cpResourceHistory History of resource in common pool
type cpResourceHistory resourceHistoryType

//================ Resource History =========================================

//resourceHistory OUR islands resources per turn
type resourceHistory resourceHistoryType

//================ Foraging History =========================================

// forageOutcome records the ROI on a foraging session
type forageOutcome struct {
	turn   uint
	team   shared.ClientID
	input  shared.Resources
	output shared.Resources
	caught uint
}

// forageHistory stores history of foraging outcomes
type forageHistory map[shared.ForageType][]forageOutcome

//================ Gifts ===========================================

// giftInfo holds information about the gifts
type giftInfo struct {
	requested      shared.GiftRequest  // How much was requested
	offered        shared.GiftOffer    // How much offered
	response       shared.GiftResponse // Response to offer
	actualReceived shared.Resources    // How much was actually received
}

//giftExchange offers the two sides of gifting
type giftExchange struct {
	// 							uint = turn
	theirRequest map[uint]giftInfo
	ourRequest   map[uint]giftInfo
}

//	giftHistory is the history of our gifts according to which island sent it
type giftHistory map[shared.ClientID]giftExchange

type client struct {
	*baseclient.BaseClient

	// Roles
	team5President president

	// Roles
	team5Speaker speaker
	team5Judge   judge

	// History
	resourceHistory         resourceHistory
	forageHistory           forageHistory
	giftHistory             giftHistory
	cpRequestHistory        cpRequestHistory
	cpAllocationHistory     cpAllocationHistory
	cpResourceHistory       cpResourceHistory
	opinionHistory          opinionHistory
	forecastHistory         forecastHistory // history of forecasted disasters
	receivedForecastHistory receivedForecastHistory
	disasterHistory         disasterHistory // history of actual disasters

	// current states
	opinions               opinionMap // opinions of each team
	lastDisasterPrediction shared.DisasterPrediction

	clientsForecastSkill clientsForecastPerformance // record of forecasting perf of other agents
	disasterModel        disasterModel              // estimates of disaster parameters

	config clientConfig
}

//================================================================
/*	Constants */
//=================================================================
// Wealth Tiers
type wealthTier int

const (
	dying wealthTier = iota
	imperialStudent
	middleClass
	jeffBezos
)

type mindSet int

const (
	okBoomer mindSet = iota
	normal
	millennial
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
	return fmt.Sprintf("Unknown wealth state '%v'", int(wt))
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

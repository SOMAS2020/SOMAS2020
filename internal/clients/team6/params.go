package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// FriendshipLevel measures the friendship levels
type FriendshipLevel float64

// TrustRank indicates our trust rank on each other islands
type TrustRank map[shared.ClientID]float64

// Friendship is the friendship level between us and another island
type Friendship map[shared.ClientID]FriendshipLevel

// GiftsSentHistory stores what gifts we have sent
type GiftsSentHistory map[shared.ClientID]shared.Resources

// GiftsReceivedHistory stores what gifts we have received
type GiftsReceivedHistory map[shared.ClientID]shared.Resources

// GiftsRequestedHistory stores what gifts we have requested
type GiftsRequestedHistory map[shared.ClientID]shared.Resources

// DisastersHistory records the sequences of disastered happened
type DisastersHistory baseclient.PastDisastersList

// DisasterPredictions records the disasater disasters from other islands
type DisasterPredictions map[shared.ClientID]shared.DisasterPrediction

// Personality enumerate our personality
type Personality int

// ForageHistory stores our forage history
type ForageHistory map[shared.ForageType][]ForageResults

// enumerates personality
const (
	Selfish Personality = iota
	Normal
	Generous
)

//ForageResults stores the input and output for foraging
type ForageResults struct {
	forageIn     shared.Resources
	forageReturn shared.Resources
	turn         uint
}

// ClientConfig configures our island's initial state
type ClientConfig struct {
	minFriendship          FriendshipLevel
	maxFriendship          FriendshipLevel
	friendshipChangingRate FriendshipLevel
	selfishThreshold       shared.Resources
	normalThreshold        shared.Resources
	multiplier             float64
}

func getClientConfig() ClientConfig {
	return ClientConfig{
		minFriendship:          0.0,
		maxFriendship:          100.0,
		friendshipChangingRate: 20,
		selfishThreshold:       50.0,
		normalThreshold:        150.0,
		multiplier:             0.3,
	}
}

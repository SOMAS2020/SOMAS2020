package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// FriendshipLevel measures the friendship levels
type FriendshipLevel float64

// Friendship is the friendship level between us and another island
type Friendship map[shared.ClientID]FriendshipLevel

// GiftsSentHistory stores what gifts we have sent
type GiftsSentHistory map[shared.ClientID]shared.Resources

// GiftsReceivedHistory stores what gifts we have received
type GiftsReceivedHistory map[shared.ClientID]shared.Resources

// GiftsRequestedHistory stores what gifts we have requested
type GiftsRequestedHistory map[shared.ClientID]shared.Resources

// Personality enumerate our personality
type Personality int

// ForageHistory stores our forage history
type ForageHistory map[shared.ForageType][]ForageResults

type FavourRules []string

// enumerates personality
const (
	Selfish Personality = iota
	Normal
	Generous
)

type ForageResults struct {
	forageIn     shared.Resources
	forageReturn shared.Resources
}

// Config configures our island's initial state
type Config struct {
	minFriendship          FriendshipLevel
	maxFriendship          FriendshipLevel
	friendshipChangingRate FriendshipLevel
	selfishThreshold       shared.Resources
	normalThreshold        shared.Resources
	payingTax              shared.Resources
}

var (
	friendship = Friendship{
		shared.Team1: 50.0,
		shared.Team2: 50.0,
		shared.Team3: 50.0,
		shared.Team4: 50.0,
		shared.Team5: 50.0,
	}
	giftsSentHistory      = GiftsSentHistory{}
	giftsReceivedHistory  = GiftsReceivedHistory{}
	giftsRequestedHistory = GiftsRequestedHistory{}
	forageHistory         = ForageHistory{}
	favourRules           = FavourRules{}
	config                = Config{
		minFriendship:          0.0,
		maxFriendship:          100.0,
		friendshipChangingRate: 20.0,
		selfishThreshold:       50.0,
		normalThreshold:        150.0,
		payingTax:              15.0,
	}
)

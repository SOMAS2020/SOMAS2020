package team6

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

type Friendship map[shared.ClientID]float64

type GiftsOfferedHistory map[shared.ClientID]shared.Resources

type GiftsReceivedHistory map[shared.ClientID]shared.Resources

type Personality int

type ForageHistory map[shared.ForageType][]ForageResults

type FavourRules []string

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
	selfishThreshold shared.Resources
	normalThreshold  shared.Resources
}

var (
	friendship = Friendship{
		shared.Team1: 50.0,
		shared.Team2: 50.0,
		shared.Team3: 50.0,
		shared.Team4: 50.0,
		shared.Team5: 50.0,
	}
	giftsOfferedHistory  = GiftsOfferedHistory{}
	giftsReceivedHistory = GiftsReceivedHistory{}
	config               = Config{
		selfishThreshold: 50.0,
		normalThreshold:  150.0,
	}
)

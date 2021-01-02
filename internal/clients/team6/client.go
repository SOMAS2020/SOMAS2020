// Package team6 contains code for team 6's client implementation
package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const (
	id             = shared.Team6
	MIN_FRIENDSHIP = 0
	MAX_FRIENDSHIP = 100
)

// GiftsReceivedHistory is what gifts we have received from other islands
type GiftsReceivedHistory map[shared.ClientID]shared.Resources

// GiftsOfferedHistory is what gifts we have offered to other islands
type GiftsOfferedHistory map[shared.ClientID]shared.Resources

// Friendship indicates friendship level with other islands
type Friendship map[shared.ClientID]uint

// Config configures our island's initial state
type Config struct {
	friendship Friendship
}

type client struct {
	*baseclient.BaseClient

	config Config
}

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient: baseclient.NewClient(id),
			config:     configInfo,
		},
	)
}

// ########################
// ######  General  #######
// ########################

// ########################
// ###### Friendship ######
// ########################

// increases the friendship level with some other islands
func (c *client) raiseFriendshipLevel(clientID shared.ClientID) {
	currFriendship := c.config.friendship[clientID]

	if currFriendship == MAX_FRIENDSHIP {
		return
	}

	c.config.friendship[clientID]++
}

// decreases the friendship level with some other islands
func (c *client) lowerFriendshipLevel(clientID shared.ClientID) {
	currFriendship := c.config.friendship[clientID]

	if currFriendship == MIN_FRIENDSHIP {
		return
	}

	c.config.friendship[clientID]--
}

// ########################
// ######  Foraging  ######
// ########################

/*
------ TODO: COMPULSORY ------
func (c *client) DecideForage() (shared.ForageDecision, error)

func (c *client) ForageUpdate(shared.ForageDecision, shared.Resources)
*/

// ########################
// ######    IITO    ######
// ########################

/*
------ TODO: COMPULSORY ------
func (c *client) GetGiftRequests() shared.GiftRequestDict

func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict

func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict

func (c *client) UpdateGiftInfo(ReceivedResponses shared.GiftResponseDict)

func (c *client) SentGift(sent shared.Resources, to shared.ClientID)

func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID)
*/

// ########################
// ######    IIFO    ######
// ########################

/*
------ TODO: COMPULSORY ------
func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo

func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict)

------ TODO: OPTIONAL ------
func (c *client) MakeForageInfo() shared.ForageShareInfo

func (c *client) ReceiveForageInfo(neighbourForaging []shared.ForageShareInfo)
*/

// ########################
// ######    IIGO    ######
// ########################

func (c *client) GetClientPresidentPointer() roles.President {
	return &president{client: c}
}

func (c *client) GetClientJudgePointer() roles.Judge {
	return &judge{client: c}
}

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	return &speaker{client: c}
}

/*
------ TODO: COMPULSORY ------
func (c *client) MonitorIIGORole(shared.Role) bool

func (c *client) DecideIIGOMonitoringAnnouncement(bool) (bool, bool)
*/

// ########################
// ######  Disaster  ######
// ########################

/*
------ TODO: OPTIONAL ------
func (c *client) DisasterNotification(disasters.DisasterReport, map[shared.ClientID]shared.Magnitude)
*/

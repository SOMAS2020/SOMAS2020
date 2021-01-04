// Package team6 contains code for team 6's client implementation
package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const (
	id = shared.Team6
	// MinFriendship is the minimum friendship level
	MinFriendship = 0
	// MaxFriendship is the maximum friendship level
	MaxFriendship = 100
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
// ######  Foraging  ######
// ########################

// ------ TODO: COMPULSORY ------
func (c *client) DecideForage() (shared.ForageDecision, error) {
	return c.BaseClient.DecideForage()
}

func (c *client) ForageUpdate(forageDecision shared.ForageDecision, resources shared.Resources) {
	c.BaseClient.ForageUpdate(forageDecision, resources)
}

// ########################
// ######    IITO    ######
// ########################

// ------ TODO: COMPULSORY ------
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	return c.BaseClient.GetGiftRequests()
}

func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	return c.BaseClient.GetGiftOffers(receivedRequests)
}

func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
	return c.BaseClient.GetGiftResponses(receivedOffers)
}

func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {

}

func (c *client) SentGift(sent shared.Resources, to shared.ClientID) {

}

func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {

}

// ########################
// ######    IIFO    ######
// ########################

// ------ TODO: COMPULSORY ------
func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {
	return c.BaseClient.MakeDisasterPrediction()
}

func (c *client) ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict) {
	c.BaseClient.ReceiveDisasterPredictions(receivedPredictions)
}

// ------ TODO: OPTIONAL ------
func (c *client) MakeForageInfo() shared.ForageShareInfo {
	return c.BaseClient.MakeForageInfo()
}

func (c *client) ReceiveForageInfo(neighbourForaging []shared.ForageShareInfo) {
	c.BaseClient.ReceiveForageInfo(neighbourForaging)
}

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

// ------ TODO: COMPULSORY ------
func (c *client) MonitorIIGORole(roleName shared.Role) bool {
	return c.BaseClient.MonitorIIGORole(roleName)
}

func (c *client) DecideIIGOMonitoringAnnouncement(monitoringResult bool) (resultToShare bool, announce bool) {
	return c.BaseClient.DecideIIGOMonitoringAnnouncement(monitoringResult)
}

// ########################
// ######  Disaster  ######
// ########################

//------ TODO: OPTIONAL ------
func (c *client) DisasterNotification(dR disasters.DisasterReport, effects disasters.DisasterEffects) {
	c.BaseClient.DisasterNotification(dR, effects)
}

// ########################
// ######   Utils   #######
// ########################

// increases the friendship level with some other islands
func (c *client) raiseFriendshipLevel(clientID shared.ClientID) {
	currFriendship := c.config.friendship[clientID]

	if currFriendship == MaxFriendship {
		return
	}

	c.config.friendship[clientID]++
}

// decreases the friendship level with some other islands
func (c *client) lowerFriendshipLevel(clientID shared.ClientID) {
	currFriendship := c.config.friendship[clientID]

	if currFriendship == MinFriendship {
		return
	}

	c.config.friendship[clientID]--
}

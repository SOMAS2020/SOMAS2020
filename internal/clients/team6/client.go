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
	MinFriendship = 0.0
	// MaxFriendship is the maximum friendship level
	MaxFriendship = 100.0
)

type client struct {
	*baseclient.BaseClient

	friendship           map[shared.ClientID]float64
	giftsOfferedHistory  map[shared.ClientID]shared.Resources
	giftsReceivedHistory map[shared.ClientID]shared.Resources

	config Config
}

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient:           baseclient.NewClient(id),
			friendship:           friendship,
			giftsOfferedHistory:  giftsOfferedHistory,
			giftsReceivedHistory: giftsReceivedHistory,
			config:               config,
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

func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	ourStatus := c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	ourPersonality := c.getPersonality()
	friendshipCoffesOnOffer := make(map[shared.ClientID]shared.GiftOffer)

	for team, fsc := range c.getFriendshipCoeffs() {
		friendshipCoffesOnOffer[team] = shared.GiftOffer(fsc)
	}

	if ourStatus == shared.Critical {
		// rejects all offers if we are in critical status
		return offers
	}

	for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		amountOffer := shared.GiftOffer(0.0)

		if status == shared.Critical {
			// offers a minimum resources to all islands which are in critical status
			// TODO: what is the best amount to offer???
			offers[team] = shared.GiftOffer(10.0)
		} else if amountRequested, found := receivedRequests[team]; found {
			// otherwise gifts will be offered based on our friendship and the amount
			if shared.Resources(amountRequested) > ourResources/6.0 {
				amountOffer = shared.GiftOffer(ourResources / 6.0)
			} else {
				amountOffer = shared.GiftOffer(amountRequested)
			}

			if ourPersonality == Selfish {
				// introduces high penalty on the gift offers if we are selfish
				offers[team] = shared.GiftOffer(friendshipCoffesOnOffer[team] * friendshipCoffesOnOffer[team] * amountOffer)
			} else if ourPersonality == Normal {
				// introduces normal penalty
				offers[team] = shared.GiftOffer(friendshipCoffesOnOffer[team] * amountOffer)
			} else {
				// intorduces no penalty - we are rich!
				offers[team] = shared.GiftOffer(amountOffer)
			}

		}
	}

	// TODO: implement different wealth state for our island so we can decided --finished
	// whether to be generous or not
	return offers
}

// ------ TODO: COMPULSORY ------
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	return c.BaseClient.GetGiftRequests()
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

func (c *client) DecideGiftAmount(toTeam shared.ClientID, giftOffer shared.Resources) shared.Resources {
	return c.BaseClient.DecideGiftAmount(toTeam, giftOffer)
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
	currFriendship := c.friendship[clientID]

	if currFriendship == MaxFriendship {
		return
	}

	c.friendship[clientID]++
}

// decreases the friendship level with some other islands
func (c *client) lowerFriendshipLevel(clientID shared.ClientID) {
	currFriendship := c.friendship[clientID]

	if currFriendship == MinFriendship {
		return
	}

	c.friendship[clientID]--
}

// gets friendship cofficients
func (c client) getFriendshipCoeffs() map[shared.ClientID]float64 {
	friendshipCoeffs := make(map[shared.ClientID]float64)

	for team, fs := range c.friendship {
		friendshipCoeffs[team] = fs / MaxFriendship
	}

	return friendshipCoeffs
}

func (c client) getPersonality() Personality {
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources

	if ourResources <= c.config.selfishThreshold {
		return Selfish
	} else if ourResources <= c.config.normalThreshold {
		return Normal
	}

	return Generous
}

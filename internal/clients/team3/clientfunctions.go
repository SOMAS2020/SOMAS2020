package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// NewClient initialises the island state
func NewClient(clientID shared.ClientID) baseclient.Client {
	ourClient := client{
		// Initialise variables here
		BaseClient: baseclient.NewClient(clientID),
		params: islandParams{
			// Define parameter values here
			selfishness: 0.5,
		},
	}

	// Set trust scores
	for _, islandID := range shared.TeamIDs {
		ourClient.trustScore[islandID] = 50
		ourClient.theirTrustScore[islandID] = 50
	}
	// Set our trust in ourselves to 100
	ourClient.theirTrustScore[id] = 100

	return &ourClient
}

// getLocalResources retrieves our islands resrouces from server
func (c *client) getLocalResources() shared.Resources {
	currentState := c.BaseClient.ServerReadHandle.GetGameState()
	return currentState.ClientInfo.Resources
}

// getIslandsAlive retrives number of islands still alive
func (c *client) getIslandsAlive() int {
	var lifeStatuses map[shared.ClientID]shared.ClientLifeStatus
	var aliveCount int

	currentState := c.BaseClient.ServerReadHandle.GetGameState()
	lifeStatuses = currentState.ClientLifeStatuses
	for _, statusInfo := range lifeStatuses {
		if statusInfo == shared.Alive {
			aliveCount += 1
		}
	}
	return aliveCount
}

// func (c *client) GetClientPresidentPointer() roles.President {
// 	return c.presidentObj
// }

//func (c *client) Echo(s string) string { return "" }

//func (c *client) GetID() shared.ClientID { return id }

// func (c *client) Initialise(baseclient.ServerReadHandle) {}
// func (c *client) StartOfTurn()                           {}
// func (c *client) Logf(format string, a ...interface{})   {}

// func (c *client) GetVoteForRule(ruleName string) bool                          { return false }
// func (c *client) GetVoteForElection(roleToElect shared.Role) []shared.ClientID { return nil }
// func (c *client) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
// }
// func (c *client) GetCommunications() *map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent {
// 	return nil
// }

// func (c *client) CommonPoolResourceRequest() shared.Resources { return shared.Resources(0) }
// func (c *client) ResourceReport() shared.Resources            { return shared.Resources(0) }
// func (c *client) RuleProposal() string                        { return "" }
// func (c *client) GetClientPresidentPointer() roles.President  { return nil }
// func (c *client) GetClientJudgePointer() roles.Judge          { return nil }
// func (c *client) GetClientSpeakerPointer() roles.Speaker      { return nil }
// func (c *client) TaxTaken(shared.Resources)                   {}
// func (c *client) GetTaxContribution() shared.Resources        { return shared.Resources(0) }
// func (c *client) RequestAllocation() shared.Resources         { return shared.Resources(0) }

// //Foraging
// func (c *client) DecideForage() (shared.ForageDecision, error) {
// 	return shared.ForageDecision{}, nil
// }
// func (c *client) ForageUpdate(shared.ForageDecision, shared.Resources) {}

// //Disasters
// func (c *client) DisasterNotification(disasters.DisasterReport, map[shared.ClientID]shared.Magnitude) {
// }

// //IIFO: OPTIONAL
// func (c *client) MakeDisasterPrediction() shared.DisasterPredictionInfo {
// 	return shared.DisasterPredictionInfo{}
// }
// func (c *client) MakeForageInfo() shared.ForageShareInfo     { return shared.ForageShareInfo{} }
// func (c *client) ReceiveForageInfo([]shared.ForageShareInfo) {}

// //IITO: COMPULSORY
// func (c *client) GetGiftRequests() shared.GiftRequestDict { return nil }
// func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
// 	return nil
// }
// func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
// 	return nil
// }
// func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {}

// //IIGO: COMPULSORY
// func (c *client) MonitorIIGORole(shared.Role) bool                   { return false }
// func (c *client) DecideIIGOMonitoringAnnouncement(bool) (bool, bool) { return false, false }

// //TODO: THESE ARE NOT DONE yet, how do people think we should implement the actual transfer?
// func (c *client) SentGift(sent shared.Resources, to shared.ClientID)           {}
// func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {}

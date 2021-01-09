package baseclient

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Client is a base interface to be implemented by each client struct.
type Client interface {
	Echo(s string) string

	GetID() shared.ClientID
	Initialise(ServerReadHandle)
	StartOfTurn()
	Logf(format string, a ...interface{})

	VoteForRule(ruleMatrix rules.RuleMatrix) shared.RuleVoteType
	VoteForElection(roleToElect shared.Role, candidateList []shared.ClientID) []shared.ClientID
	ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent)
	GetCommunications() *map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent

	CommonPoolResourceRequest() shared.Resources
	ResourceReport() shared.ResourcesReport
	RuleProposal() rules.RuleMatrix
	GetClientPresidentPointer() roles.President
	GetClientJudgePointer() roles.Judge
	GetClientSpeakerPointer() roles.Speaker
	GetTaxContribution() shared.Resources
	GetSanctionPayment() shared.Resources
	RequestAllocation() shared.Resources
	ShareIntendedContribution() shared.IntendedContribution
	ReceiveIntendedContribution(receivedIntendedContributions shared.ReceivedIntendedContributionDict)

	//Foraging
	DecideForage() (shared.ForageDecision, error)
	ForageUpdate(shared.ForageDecision, shared.Resources, uint)

	//Disasters
	DisasterNotification(disasters.DisasterReport, disasters.DisasterEffects)

	//IIFO: OPTIONAL
	MakeDisasterPrediction() shared.DisasterPredictionInfo
	ReceiveDisasterPredictions(receivedPredictions shared.ReceivedDisasterPredictionsDict)
	MakeForageInfo() shared.ForageShareInfo
	ReceiveForageInfo([]shared.ForageShareInfo)

	//IITO: COMPULSORY
	GetGiftRequests() shared.GiftRequestDict
	GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict
	GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict
	UpdateGiftInfo(receivedResponses shared.GiftResponseDict)
	DecideGiftAmount(shared.ClientID, shared.Resources) shared.Resources

	//IIGO: COMPULSORY
	MonitorIIGORole(shared.Role) bool
	DecideIIGOMonitoringAnnouncement(bool) (bool, bool)

	//TODO: THESE ARE NOT DONE yet, how do people think we should implement the actual transfer?
	SentGift(sent shared.Resources, to shared.ClientID)
	ReceivedGift(received shared.Resources, from shared.ClientID)
}

// ServerReadHandle is a read-only handle to the game server, used for client to get up-to-date gamestate
type ServerReadHandle interface {
	GetGameState() gamestate.ClientGameState
	GetGameConfig() config.ClientConfig
}

// NewClient produces a new client with the BaseClient already implemented.
func NewClient(id shared.ClientID) *BaseClient {
	return &BaseClient{
		id:             id,
		Communications: map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent{},
	}
}

// BaseClient provides a basic implementation for all functions of the client interface and should always the interface fully.
// All clients should be based off of this BaseClient to ensure that all clients implement the interface,
// even when new features are added.
type BaseClient struct {
	id shared.ClientID

	predictionInfo       shared.DisasterPredictionInfo
	intendedContribution shared.IntendedContribution

	// exported variables are accessible by the client implementations
	LocalVariableCache map[rules.VariableFieldName]rules.VariableValuePair
	Communications     map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent
	ServerReadHandle   ServerReadHandle
}

// Echo prints a message to show that the client exists
// BASE: Do not overwrite in team client.
func (c *BaseClient) Echo(s string) string {
	c.Logf("Echo: '%v'", s)
	return s
}

// GetID gets returns the id of the client.
// BASE: Do not overwrite in team client.
func (c *BaseClient) GetID() shared.ClientID {
	return c.id
}

// Initialise initialises the base client.
// OPTIONAL: Overwrite, and make sure to keep the value of ServerReadHandle.
// You will need it to access the game state through its GetGameStateMethod.
func (c *BaseClient) Initialise(serverReadHandle ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle
	c.LocalVariableCache = rules.CopyVariableMap(c.ServerReadHandle.GetGameState().RulesInfo.VariableMap)
}

// StartOfTurn handles the start of a new turn.
// OPTIONAL: Use this method for any tasks you want to happen on the beginning
// of every turn (e.g. logging)
func (c *BaseClient) StartOfTurn() {}

// Logf is the client's logger that prepends logs with your ID. This makes
// it easier to read logs. DO NOT use other loggers that will mess logs up!
// BASE: Do not overwrite in team client.
func (c *BaseClient) Logf(format string, a ...interface{}) {
	log.Printf("[%v]: %v", c.id, fmt.Sprintf(format, a...))
}

// GetVoteForRule returns the client's vote in favour of or against a rule.
// COMPULSORY: vote to represent your island's opinion on a rule

func (c *BaseClient) VoteForRule(ruleMatrix rules.RuleMatrix) shared.RuleVoteType {
	// TODO implement decision on voting that considers the rule
	return shared.Approve
}

// GetVoteForElection returns the client's Borda vote for the role to be elected.
// COMPULSORY: use opinion formation to decide a rank for islands for the role
func (c *BaseClient) VoteForElection(roleToElect shared.Role, candidateList []shared.ClientID) []shared.ClientID {
	candidates := map[int]shared.ClientID{}
	for i := 0; i < len(candidateList); i++ {
		candidates[i] = candidateList[i]
	}
	// Recombine map, in shuffled order
	var returnList []shared.ClientID
	for _, v := range candidates {
		returnList = append(returnList, v)
	}
	return returnList
}

// ReceiveCommunication is a function called by IIGO to pass the communication sent to the client
// COMPULSORY: please override to save incoming communication relevant to your agent strategy
func (c *BaseClient) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	c.InspectCommunication(data)
	c.Communications[sender] = append(c.Communications[sender], data)
}

// InspectCommunication is a function which scans over Receive communications for any data
// that might be used to update the LocalVariableCache
func (c *BaseClient) InspectCommunication(data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	if c.LocalVariableCache == nil {
		return
	}
	for fieldName, dataPoint := range data {
		switch fieldName {
		case shared.IIGOTaxDecision:
			c.LocalVariableCache[rules.ExpectedTaxContribution] = dataPoint.IIGOValueData.Expected
			c.LocalVariableCache[rules.TaxDecisionMade] = dataPoint.IIGOValueData.DecisionMade
		case shared.IIGOAllocationDecision:
			c.LocalVariableCache[rules.ExpectedAllocation] = dataPoint.IIGOValueData.Expected
			c.LocalVariableCache[rules.AllocationMade] = dataPoint.IIGOValueData.DecisionMade
		case shared.SanctionAmount:
			c.LocalVariableCache[rules.SanctionExpected] = rules.VariableValuePair{
				VariableName: rules.SanctionExpected,
				Values:       []float64{float64(dataPoint.IntegerData)},
			}
		// TODO: Extend according to new rules added by Team 4
		default:
			return
		}

	}
}

// GetCommunications is used for testing communications
// BASE
func (c *BaseClient) GetCommunications() *map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent {
	return &c.Communications
}

//MonitorIIGORole decides whether to perform monitoring on a role
//COMPULOSRY: must be implemented
func (c *BaseClient) MonitorIIGORole(roleName shared.Role) bool {
	return true
}

//DecideIIGOMonitoringAnnouncement decides whether to share the result of monitoring a role and what result to share
//COMPULSORY: must be implemented
func (c *BaseClient) DecideIIGOMonitoringAnnouncement(monitoringResult bool) (resultToShare bool, announce bool) {
	resultToShare = monitoringResult
	announce = true
	return
}

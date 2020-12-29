// Package baseclient contains the Client interface as well as a base client implementation.
package baseclient

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Client is a base interface to be implemented by each client struct.
type Client interface {
	Echo(s string) string
	GetID() shared.ClientID
	Initialise(ServerReadHandle)
	StartOfTurn()
	Logf(format string, a ...interface{})

	GetVoteForRule(ruleName string) bool
	GetVoteForElection(roleToElect Role) []shared.ClientID
	ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent)
	GetCommunications() *map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent

	CommonPoolResourceRequest() shared.Resources
	ResourceReport() shared.Resources
	RuleProposal() string
	GetClientPresidentPointer() roles.President
	GetClientJudgePointer() roles.Judge
	GetClientSpeakerPointer() roles.Speaker
	TaxTaken(shared.Resources)
	GetTaxContribution() shared.Resources
	RequestAllocation() shared.Resources

	//Foraging
	DecideForage() (shared.ForageDecision, error)
	ForageUpdate(shared.ForageDecision, shared.Resources)

	//Disasters
	DisasterNotification(disasters.DisasterReport, map[shared.ClientID]shared.Magnitude)

	//IIFO: OPTIONAL
	MakePrediction() (shared.PredictionInfo, error)
	ReceivePredictions(receivedPredictions shared.PredictionInfoDict) error
	MakeForageInfo() shared.ForageShareInfo
	ReceiveForageInfo([]shared.ForageShareInfo)

	//IITO: COMPULSORY
	RequestGift() uint
	OfferGifts(giftRequestDict shared.GiftDict) (shared.GiftDict, error)
	AcceptGifts(receivedGiftDict shared.GiftDict) (shared.GiftInfoDict, error)
	UpdateGiftInfo(acceptedGifts map[shared.ClientID]shared.GiftInfoDict) error

	//TODO: THESE ARE NOT DONE yet, how do people think we should implement the actual transfer?
	SendGift(receivingClient shared.ClientID, amount int) error
	ReceiveGift(sendingClient shared.ClientID, amount int) error
}

// ServerReadHandle is a read-only handle to the game server, used for client to get up-to-date gamestate
type ServerReadHandle interface {
	GetGameState() gamestate.ClientGameState
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

	predictionInfo shared.PredictionInfo

	// exported variables are accessible by the client implementations
	Communications   map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent
	ServerReadHandle ServerReadHandle
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

// Role provides enumerated type for IIGO roles (President, Speaker and Judge)
type Role int

const (
	President Role = iota
	Speaker
	Judge
)

func (r Role) String() string {
	strs := [...]string{"President", "Speaker", "Judge"}
	if r >= 0 && int(r) < len(strs) {
		return strs[r]
	}
	return fmt.Sprintf("UNKNOWN Role '%v'", int(r))
}

// GoString implements GoStringer
func (r Role) GoString() string {
	return r.String()
}

// MarshalText implements TextMarshaler
func (r Role) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(r.String())
}

// MarshalJSON implements RawMessage
func (r Role) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(r.String())
}

// GetVoteForRule returns the client's vote in favour of or against a rule.
func (c *BaseClient) GetVoteForRule(ruleName string) bool {
	// TODO implement decision on voting that considers the rule
	return true
}

// GetVoteForElection returns the client's Borda vote for the role to be elected.
func (c *BaseClient) GetVoteForElection(roleToElect Role) []shared.ClientID {
	// Done ;)
	// Get all alive islands
	aliveClients := rules.VariableMap[rules.IslandsAlive]
	// Convert to ClientID type and place into unordered map
	aliveClientIDs := map[int]shared.ClientID{}
	for i, v := range aliveClients.Values {
		aliveClientIDs[i] = shared.ClientID(int(v))
	}
	// Recombine map, in shuffled order
	var returnList []shared.ClientID
	for _, v := range aliveClientIDs {
		returnList = append(returnList, v)
	}
	return returnList
}

// ReceiveCommunication is a function called by IIGO to pass the communication sent to the client
func (c *BaseClient) ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	c.Communications[sender] = append(c.Communications[sender], data)
}

// GetCommunications is used for testing communications
func (c *BaseClient) GetCommunications() *map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent {
	return &c.Communications
}

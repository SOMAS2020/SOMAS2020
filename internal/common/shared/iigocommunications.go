package shared

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

// CommunicationContentType provides type union for generic IIGO Inter-Island Communications
type CommunicationContentType int

const (
	CommunicationInt CommunicationContentType = iota
	CommunicationString
	CommunicationBool
	CommunicationIIGORole
	CommunicationIIGOValue
)

func (c CommunicationContentType) String() string {
	strs := [...]string{
		"CommunicationInt",
		"CommunicationString",
		"CommunicationBool",
		"CommunicationIIGORole",
		"CommunicationIIGOValue",
	}
	if c >= 0 && int(c) < len(strs) {
		return strs[c]
	}
	return fmt.Sprintf("UNKNOWN CommunicationContentType '%v'", int(c))
}

// GoString implements GoStringer
func (c CommunicationContentType) GoString() string {
	return c.String()
}

// MarshalText implements TextMarshaler
func (c CommunicationContentType) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(c.String())
}

// MarshalJSON implements RawMessage
func (c CommunicationContentType) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(c.String())
}

// ValueDecision is part of CommunicationContent and is used to send a tax decision from president to the client
type ValueDecision struct {
	Amount       Resources
	DecisionMade rules.VariableValuePair
	Expected     rules.VariableValuePair
}

// CommunicationContent is a general datastructure used for communications
type CommunicationContent struct {
	T              CommunicationContentType
	IntegerData    int
	TextData       string
	BooleanData    bool
	IIGORoleData   Role
	IIGOValueData  ValueDecision
	RuleMatrixData rules.RuleMatrix
}

type CommunicationFieldName int

const (
	BallotID CommunicationFieldName = iota
	SpeakerID
	RoleConducted
	ResAllocID
	PresidentID
	RuleName
	RuleVoteResult
	PardonClientID
	PardonTier
	SanctionAmount
	RoleMonitored
	MonitoringResult
	IIGOSanctionTier
	RuleSanctionPenalty
	IIGOTaxDecision
	IIGOAllocationDecision
	SanctionClientID
)

func (c CommunicationFieldName) String() string {
	strs := [...]string{
		"BallotID",
		"SpeakerID",
		"RoleConducted",
		"ResAllocID",
		"PresidentID",
		"RuleName",
		"RuleVoteResult",
		"PardonClientID",
		"PardonTier",
		"SanctionAmount",
		"RoleMonitored",
		"MonitoringResult",
		"IIGOSanctionTier",
		"RuleSanctionPenalty",
		"IIGOTaxDecision",
		"IIGOAllocationDecision",
		"SanctionClientID",
	}
	if c >= 0 && int(c) < len(strs) {
		return strs[c]
	}
	return fmt.Sprintf("UNKNOWN CommunicationFieldName '%v'", int(c))
}

// GoString implements GoStringer
func (c CommunicationFieldName) GoString() string {
	return c.String()
}

// MarshalText implements TextMarshaler
func (c CommunicationFieldName) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(c.String())
}

// MarshalJSON implements RawMessage
func (c CommunicationFieldName) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(c.String())
}

type Accountability struct {
	ClientID ClientID
	Pairs    []rules.VariableValuePair
}

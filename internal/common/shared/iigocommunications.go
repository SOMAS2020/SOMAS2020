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
	CommunicationTax
	CommunicationAllocation
	CommunicationIIGORole
)

func (c CommunicationContentType) String() string {
	strs := [...]string{
		"CommunicationInt",
		"CommunicationString",
		"CommunicationBool",
		"CommunicationTax",
		"CommunicationAllocation",
		"CommunicationIIGORole",
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

// TaxDecision is part of CommunicationContent and is used to send a tax decision from president to the client
type TaxDecision struct {
	TaxAmount   Resources
	TaxRule     rules.RuleMatrix
	ExpectedTax rules.VariableValuePair
	TaxDecided  bool
}

// AllocationDecision is part of CommunicationContent and is used to send an allocation decision from president to the client
type AllocationDecision struct {
	AllocationAmount   Resources
	AllocationRule     rules.RuleMatrix
	ExpectedAllocation rules.VariableValuePair
	AllocationDecided  bool
}

// CommunicationContent is a general datastructure used for communications
type CommunicationContent struct {
	T                  CommunicationContentType
	IntegerData        int
	TextData           string
	BooleanData        bool
	TaxDecision        TaxDecision
	AllocationDecision AllocationDecision
	IIGORoleData       Role
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
	Tax
	Allocation
	PardonClientID
	PardonTier
	SanctionAmount
	RoleMonitored
	MonitoringResult
	IIGOSanctionTier
	IIGOSanctionScore
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
		"Tax",
		"Allocation",
		"PardonClientID",
		"PardonTier",
		"SanctionAmount",
		"RoleMonitored",
		"MonitoringResult",
		"IIGOSanctionTier",
		"IIGOSanctionScore",
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

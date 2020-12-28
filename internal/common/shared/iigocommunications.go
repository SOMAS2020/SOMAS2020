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
)

func (c CommunicationContentType) String() string {
	strs := [...]string{
		"CommunicationInt",
		"CommunicationString",
		"CommunicationBool",
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

// CommunicationContent is a general datastructure used for communications
type CommunicationContent struct {
	T           CommunicationContentType
	IntegerData int
	TextData    string
	BooleanData bool
}

type CommunicationFieldName int

const (
	BallotID CommunicationFieldName = iota
	PresidentAllocationCheck
	SpeakerID
	RoleConducted
	ResAllocID
	SpeakerBallotCheck
	PresidentID
	RuleName
	RuleVoteResult
	TaxAmount
	AllocationAmount
)

func (c CommunicationFieldName) String() string {
	strs := [...]string{
		"BallotID",
		"PresidentAllocationCheck",
		"SpeakerID",
		"RoleConducted",
		"ResAllocID",
		"SpeakerBallotCheck",
		"PresidentID",
		"RuleName",
		"RuleVoteResult",
		"TaxAmount",
		"AllocationAmount",
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

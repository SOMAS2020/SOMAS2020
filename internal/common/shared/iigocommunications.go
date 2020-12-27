package shared

import "github.com/SOMAS2020/SOMAS2020/internal/common/rules"

type CommunicationContentType = int

const (
	CommunicationInt CommunicationContentType = iota
	CommunicationString
	CommunicationBool
)

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

type Accountability struct {
	ClientID ClientID
	Pairs    []rules.VariableValuePair
}

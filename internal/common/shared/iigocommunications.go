package shared

type CommunicationContentType = int

const (
	CommunicationInt CommunicationContentType = iota
	CommunicationString
	CommunicationBool
)

// Communication is a general datastructure used for communications
type Communication struct {
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

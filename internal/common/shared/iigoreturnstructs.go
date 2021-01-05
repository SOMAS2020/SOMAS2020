package shared

// PresidentReturnContentType is the index of the return content that should be retrieved if the action was taken.
type PresidentReturnContentType = int

// Enum to specify all the types of possible return values coming from president
const (
	PresidentAllocation PresidentReturnContentType = iota
	PresidentTaxation
	PresidentRuleProposal
	PresidentSpeakerSalary
)

// PresidentReturnContent is a general datastructure used for president return type
type PresidentReturnContent struct {
	ContentType   PresidentReturnContentType
	ResourceMap   map[ClientID]Resources
	ProposedRule  string
	SpeakerSalary Resources
	ActionTaken   bool
}

// SpeakerReturnContentType is the index of the return content that should be retrieved if the action was taken.
type SpeakerReturnContentType = int

// Enum to specify all the types of possible return values coming from the speaker
const (
	SpeakerAgenda SpeakerReturnContentType = iota
	SpeakerVote
	SpeakerAnnouncement
	SpeakerJudgeSalary
)

// SpeakerReturnContent is a general datastructure used for speaker return type
type SpeakerReturnContent struct {
	ContentType          SpeakerReturnContentType
	ParticipatingIslands []ClientID
	RuleID               string
	VotingResult         bool
	JudgeSalary          Resources
	ActionTaken          bool
}

// ResourcesReport is a struct returned by the Client when asked to report it's resources.
// The client can choose to report the resources by setting the Reported entry to true, along with ReportedAmount.
// If client doesn't want to share the information about its resources with president, it can set Reported to false.
type ResourcesReport struct {
	ReportedAmount Resources
	Reported       bool
}

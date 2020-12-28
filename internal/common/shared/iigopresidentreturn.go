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

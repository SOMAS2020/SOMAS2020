package roles

//President Object
type President interface {
	withdrawSpeakerSalary() error
	paySpeaker() error
	setTaxationAmount(int, int) error
	signalTaxation(int, int) error
	signalAllocationRequests(int) error
	evaluateAllocationRequest() map[int]int
	replyAllocationRequests(int) error
	pickRuleToVote() error
	sendRuleToSpeaker(int) error
	reset(string) error
}

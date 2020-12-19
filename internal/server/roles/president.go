package roles

//President Object
type President interface {
	paySpeaker() error
	setTaxationAmount(int, int) error
	evaluateAllocationRequest() map[int]int
	pickRuleToVote() error
	reset(string) error
}

	

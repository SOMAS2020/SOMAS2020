package roles

//President Object
type President interface {
	paySpeaker(common.GameState)
	setTaxationAmount(int, int) error
	evaluateAllocationRequest() map[int]int
	pickRuleToVote() string
	reset(string) error
}

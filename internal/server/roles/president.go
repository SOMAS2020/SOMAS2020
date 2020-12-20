package roles

//President Object
type President interface {
	paySpeaker() error
	setTaxationAmount(map[int]int) (map[int]int, error)
	evaluateAllocationRequests(map[int]int, int) (map[int]int, error)
	pickRuleToVote([]string) (string, error)
	reset(string) error
}

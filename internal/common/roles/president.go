package roles

//President Object
type President interface {
	PaySpeaker() (int, error)
	SetTaxationAmount(map[int]int) (map[int]int, error)
	EvaluateAllocationRequests(map[int]int, int) (map[int]int, error)
	PickRuleToVote([]string) (string, error)
	Reset(string) error
}

package roles

//base President Object
type basePresident struct {
	id                 int
	budget             int
	speakerSalary      int
	resourceRequests   map[int]int
	resourceAllocation map[int]int
	ruleToVote         int
	taxAmount          int
}

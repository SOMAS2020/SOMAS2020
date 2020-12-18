package roles

// base Judge object
type BaseJudge struct {
	id 	   int
	budget int
	presidentSalary int
	ballotID int
	resAlocID int
	actionLog map[int]string // not sure about this currently
}


func (j *BaseJudge) withdrawPresidentSalary() {
	// Withdraw president salary from the common pool
}

func (j *BaseJudge) payPresident() {
	// Pay the president
}

func (j *BaseJudge) inspectBallot() {
	// 1. Evaluate difference between newRules and oldRules to check
	//    rule changes are in line with ruleToVote in previous ballot
	// 2. Compare each ballot action adheres to rules in ruleSet matrix
}

func (j *BaseJudge) inspectAllocation() {
	// 1. Evaluate difference between commonPoolNew and commonPoolOld
	//    to check resource allocation changes are in line with resourceRequests
	//    in previous resourceAllocation
	// 2. Compare each resource allocation action adheres to rules in ruleSet
	//    matrix
}

func (j *BaseJudge) declareSpeakerPerformance() {
	// result := "pass/fail" based on outcome of audit in inspectBallot()
	// Broadcast (j.ballotID, result, "Speaker") to all islands
}

func (j *BaseJudge) declarePresidentPerformance() {
	// result := "pass/fail" based on outcome of audit in inspectAllocation()
	// Broadcast (j.resAlocID, result, "President") to all islands
}
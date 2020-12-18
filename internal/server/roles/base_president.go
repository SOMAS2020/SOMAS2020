package roles

//base President Object
type BasePresident struct {
	id                 int
	budget             int
	speakerSalary      int
	resourceRequests   map[int]int
	resourceAllocation map[int]int
	ruleToVote         int
	taxAmount          int
}

// Set taxation amount for all of the living islands
func (p *BasePresident) setTaxationAmount() {

}

// Obtain a rule proposal from each of the living islands
func (p *BasePresident) getRuleProposals() {

}

// Broadcast tax amount to each of the living islands
func (p *BasePresident) signalTaxationAmount() {

}

// Broadcast approved common pool resource request to each of the living islands
func (p *BasePresident) signalAllocationRequest() {

}

// Send rule to be voted on to Speaker
func (p *BasePresident) sendRuleToSpeaker() {

}

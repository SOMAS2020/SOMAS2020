package voting

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
)

type RuleVote struct {
	ruleToVote 		string
	islandsToVote 	[]shared.ClientID
	votes 			[]bool
}

// ProposeMotion is called by Speaker to set the rule to be voted on.
func (v *RuleVote) ProposeMotion(rule string) {
	ruleToVote := rule
}

// OpenBallot is called by Speaker to set the islands eligible to vote.
func (v *RuleVote) OpenBallot( clientIDs []shared.ClientID){
	s.GetVoteForRule(ruleToVote)
}

// Vote is called by Speaker to get votes from clients.
func (v *RuleVote) Vote() {

}

// CloseBallot is called by Speaker to count votes received.
func (v *RuleVote) CloseBallot() {

}
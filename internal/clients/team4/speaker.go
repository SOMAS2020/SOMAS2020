package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type speaker struct {
	*baseclient.BaseSpeaker
	parent *client
}

//SpeakerActionPriorities indicate speaker actions in order of priority
var SpeakerActionPriorities = []string{
	"SetVotingResult",
	"SetRuleToVote",
	"AnnounceVotingResult",
	"UpdateRules",
	"AppointNextJudge",
}

func (s *speaker) getSpeakerBudget() shared.Resources {
	return s.parent.ServerReadHandle.GetGameState().IIGORolesBudget[shared.Speaker]
}

func (s *speaker) getActionsCost(actions []string) shared.Resources {
	gameconfig := s.parent.ServerReadHandle.GetGameConfig().IIGOClientConfig
	costs := map[string]shared.Resources{
		"SetVotingResult":      gameconfig.SetVotingResultActionCost,
		"SetRuleToVote":        gameconfig.SetRuleToVoteActionCost,
		"AnnounceVotingResult": gameconfig.AnnounceVotingResultActionCost,
		"UpdateRules":          gameconfig.UpdateRulesActionCost,
		"AppointNextJudge":     gameconfig.AppointNextJudgeActionCost,
	}
	var SumOfCosts shared.Resources = 0
	for _, action := range actions {
		SumOfCosts += costs[action]
	}
	return SumOfCosts
}

func (s *speaker) getHigherPriorityActionsCost(baseaction string) shared.Resources {
	actionindex := len(SpeakerActionPriorities)
	for i, action := range SpeakerActionPriorities {
		if action == baseaction {
			actionindex = i
		}
	}
	return s.getActionsCost(SpeakerActionPriorities[:actionindex])
}

// PayJudge is used for paying judge for his service
func (s *speaker) PayJudge() shared.SpeakerReturnContent {
	var JudgeSalary shared.Resources = 0
	JudgeSalaryRule, ok := rules.RulesInPlay["salary_cycle_judge"]
	if ok {
		JudgeSalary = shared.Resources(JudgeSalaryRule.ApplicableMatrix.At(0, 1))
	}
	HigherPriorityActionsCost := s.getHigherPriorityActionsCost("")
	if s.getSpeakerBudget()-HigherPriorityActionsCost-JudgeSalary < 0 {
		//insufficient budget: JudgeSalary + action costs uses up budget
		JudgeSalary = s.getSpeakerBudget() - HigherPriorityActionsCost
		if JudgeSalary < 0 {
			JudgeSalary = 0
		}
	}
	s.parent.Logf("Judge salary: %v", JudgeSalary)
	return shared.SpeakerReturnContent{
		ContentType: shared.SpeakerJudgeSalary,
		JudgeSalary: JudgeSalary,
		ActionTaken: true,
	}
}

//DecideAgenda the interface implementation and example of a well behaved Speaker
//who sets the vote to be voted on to be the rule the President provided

//DecideVote is the interface implementation and example of a well behaved Speaker
//who calls a vote on the proposed rule and asks all available islands to vote.
//Return an empty string or empty []shared.ClientID for no vote to occur

//DecideAnnouncement is the interface implementation and example of a well behaved Speaker
//A well behaved speaker announces what had been voted on and the corresponding result
//Return "", _ for no announcement to occur

func (s *speaker) DecideAnnouncement(ruleMatrix rules.RuleMatrix, result bool) shared.SpeakerReturnContent {
	return shared.SpeakerReturnContent{
		ContentType:  shared.SpeakerAnnouncement,
		RuleMatrix:   ruleMatrix,
		VotingResult: result,
		ActionTaken:  true,
	}
}

// CallJudgeElection is called by the legislature to decide on power-transfer
// COMPULSORY: decide when to call an election following relevant rulesInPlay if you wish
/*func (s *speaker) CallJudgeElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	// example implementation calls an election if monitoring was performed and the result was negative
	// or if the number of turnsInPower exceeds 3
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.BordaCount,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}
	if monitoring.Performed && !monitoring.Result {
		electionsettings.HoldElection = true
	}
	if turnsInPower >= 2 {
		electionsettings.HoldElection = true
	}
	return electionsettings
}*/

// DecideNextJudge returns the ID of chosen next Judge
// OPTIONAL: override to manipulate the result of the election

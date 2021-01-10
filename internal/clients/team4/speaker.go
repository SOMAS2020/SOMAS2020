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

//SpeakerActionOrder indicates speaker actions in order of execution
var SpeakerActionOrder = []string{
	"SetRuleToVote",
	"SetVotingResult",
	"AnnounceVotingResult",
	"AppointNextJudge",
}

//SpeakerActionPriorities indicate speaker actions in order of priority
var SpeakerActionPriorities = []string{
	"SetVotingResult",
	"AnnounceVotingResult",
	"SetRuleToVote",
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
	actionindex := len(SpeakerActionOrder)
	for i, action := range SpeakerActionOrder {
		if action == baseaction {
			actionindex = i
		}
	}
	priorityindex := len(SpeakerActionPriorities)
	for i, action := range SpeakerActionPriorities {
		if action == baseaction {
			priorityindex = i
		}
	}
	var SAPcopy []string = SpeakerActionPriorities[:0]
	for _, action1 := range SpeakerActionPriorities[:priorityindex] {
		alreadyExecuted := false
		for _, action2 := range SpeakerActionOrder[:actionindex] {
			if action1 == action2 {
				alreadyExecuted = true
			}
		}
		if !alreadyExecuted {
			SAPcopy = append(SAPcopy, action1)
		}
	}
	return s.getActionsCost(SAPcopy)
}

// PayJudge is used for paying judge for his service
func (s *speaker) PayJudge() shared.SpeakerReturnContent {
	var JudgeSalary shared.Resources = 0
	JudgeSalaryRule, ok := s.parent.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay["salary_cycle_judge"]
	if ok {
		JudgeSalary = shared.Resources(JudgeSalaryRule.ApplicableMatrix.At(0, 1))
	}
	return shared.SpeakerReturnContent{
		ContentType: shared.SpeakerJudgeSalary,
		JudgeSalary: JudgeSalary,
		ActionTaken: true,
	}
}

//DecideAgenda the interface implementation and example of a well behaved Speaker
//who sets the vote to be voted on to be the rule the President provided
func (s *speaker) DecideAgenda(ruleMatrix rules.RuleMatrix) shared.SpeakerReturnContent {
	//there are more important things to do...
	s.parent.Logf("Action costs: %v, %v, %v, %v", s.getActionsCost([]string{"SetRuleToVote"}), s.getActionsCost([]string{"SetVotingResult"}), s.getActionsCost([]string{"AnnounceVotingResult"}), s.getActionsCost([]string{"AppointNextJudge"}))
	s.parent.Logf("Higher priority: %v, %v, %v, %v", s.getHigherPriorityActionsCost("SetRuleToVote"), s.getHigherPriorityActionsCost("SetVotingResult"), s.getHigherPriorityActionsCost("AnnounceVotingResult"), s.getHigherPriorityActionsCost("AppointNextJudge"))
	if s.getSpeakerBudget() < s.getHigherPriorityActionsCost("SetRuleToVote") {
		return shared.SpeakerReturnContent{
			ActionTaken: false,
		}
	}
	if ruleMatrix.RuleMatrixIsEmpty() {
		//president has not selected a rule
	}
	return shared.SpeakerReturnContent{
		ContentType: shared.SpeakerAgenda,
		RuleMatrix:  ruleMatrix,
		ActionTaken: true,
	}
}

//DecideVote is the interface implementation and example of a well behaved Speaker
//who calls a vote on the proposed rule and asks all available islands to vote.
//Return an empty string or empty []shared.ClientID for no vote to occur
func (s *speaker) DecideVote(ruleMatrix rules.RuleMatrix, aliveClients []shared.ClientID) shared.SpeakerReturnContent {
	//there are more important things to do...
	if s.getSpeakerBudget() < s.getHigherPriorityActionsCost("SetVotingResult") {
		return shared.SpeakerReturnContent{
			ActionTaken: false,
		}
	}
	return shared.SpeakerReturnContent{
		ContentType:          shared.SpeakerVote,
		ParticipatingIslands: aliveClients,
		RuleMatrix:           ruleMatrix,
		ActionTaken:          true,
	}
}

//DecideAnnouncement is the interface implementation and example of a well behaved Speaker
//A well behaved speaker announces what had been voted on and the corresponding result
//Return "", _ for no announcement to occur
func (s *speaker) DecideAnnouncement(ruleMatrix rules.RuleMatrix, result bool) shared.SpeakerReturnContent {
	//there are more important things to do...
	if s.getSpeakerBudget() < s.getHigherPriorityActionsCost("AnnounceVotingResult") {
		return shared.SpeakerReturnContent{
			ActionTaken: false,
		}
	}
	return shared.SpeakerReturnContent{
		ContentType:  shared.SpeakerAnnouncement,
		RuleMatrix:   ruleMatrix,
		VotingResult: result,
		ActionTaken:  true,
	}
}

// CallJudgeElection is called by the legislature to decide on power-transfer
// COMPULSORY: decide when to call an election following relevant rulesInPlay if you wish
func (s *speaker) CallJudgeElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	// example implementation calls an election if monitoring was performed and the result was negative
	// or if the number of turnsInPower exceeds 3
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.BordaCount,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}
	//there are more important things to do...
	if s.getSpeakerBudget() < s.getHigherPriorityActionsCost("AppointNextJudge") {
		return electionsettings
	}
	if monitoring.Performed && !monitoring.Result {
		electionsettings.HoldElection = true
	}
	if turnsInPower >= 2 {
		electionsettings.HoldElection = true
	}
	return electionsettings
}

// DecideNextJudge returns the ID of chosen next Judge
// OPTIONAL: override to manipulate the result of the election

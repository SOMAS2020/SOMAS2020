package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type speaker struct {
	*baseclient.BaseSpeaker
	parent                  *client
	SpeakerActionOrder      []string
	SpeakerActionPriorities []string
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
	//find where the action is in the order
	actionindex := len(s.SpeakerActionOrder)
	for i, action := range s.SpeakerActionOrder {
		if action == baseaction {
			actionindex = i
		}
	}
	//find where the action is in the priorities
	priorityindex := len(s.SpeakerActionPriorities)
	for i, action := range s.SpeakerActionPriorities {
		if action == baseaction {
			priorityindex = i
		}
	}
	//add to return slice if the higher priority action has not been executed yet
	var HigherPriorityActions = make([]string, 0)
	for _, action1 := range s.SpeakerActionPriorities[:priorityindex] {
		alreadyExecuted := false
		for _, action2 := range s.SpeakerActionOrder[:actionindex] {
			if action1 == action2 {
				alreadyExecuted = true
			}
		}
		if !alreadyExecuted {
			HigherPriorityActions = append(HigherPriorityActions, action1)
		}
	}
	return s.getActionsCost(HigherPriorityActions)
}

func (s *speaker) sendActionToBack(str string) {
	index := len(s.SpeakerActionPriorities) - 1
	for i, el := range s.SpeakerActionPriorities {
		if el == str {
			index = i
		}
	}
	if index != len(s.SpeakerActionPriorities)-1 {
		part1 := s.SpeakerActionPriorities[:index]
		part2 := s.SpeakerActionPriorities[(index + 1):]
		part3 := s.SpeakerActionPriorities[index]
		s.SpeakerActionPriorities = append(append(part1, part2...), part3)
	}
}

func (s *speaker) reorderPriorities() {
	copy(s.SpeakerActionPriorities, s.SpeakerActionOrder)
	_, ok := s.parent.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay["vote_called_rule"]
	if !ok {
		s.sendActionToBack("SetRuleToVote")
	}
	_, ok = s.parent.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay["islands_allowed_to_vote_rule"]
	if !ok {
		s.sendActionToBack("SetVotingResult")
	}
	_, ok = s.parent.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay["vote_result_rule"]
	if !ok {
		s.sendActionToBack("AnnounceVotingResult")
	}
	_, ok = s.parent.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay["roles_must_hold_election"]
	if !ok {
		s.sendActionToBack("AppointNextJudge")
	}
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
	s.reorderPriorities()
	//(there are more important things to do) or (there is no rule to put on the agenda)
	if (s.getSpeakerBudget() < s.getHigherPriorityActionsCost("SetRuleToVote")) || ruleMatrix.RuleMatrixIsEmpty() {
		return shared.SpeakerReturnContent{
			ActionTaken: false,
		}
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
	//(there are more important things to do) or (there is no rule to vote on)
	if s.getSpeakerBudget() < s.getHigherPriorityActionsCost("SetVotingResult") || ruleMatrix.RuleMatrixIsEmpty() {
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
	//(there are more important things to do) or (there is no result to announce)
	if s.getSpeakerBudget() < s.getHigherPriorityActionsCost("AnnounceVotingResult") || ruleMatrix.RuleMatrixIsEmpty() {
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
	//there are more important things to do
	if s.getSpeakerBudget() < s.getHigherPriorityActionsCost("AppointNextJudge") {
		return electionsettings
	}
	if (monitoring.Performed && !monitoring.Result) || turnsInPower >= 2 {
		electionsettings.HoldElection = true
	}
	return electionsettings
}

// DecideNextJudge returns the ID of chosen next Judge
// OPTIONAL: override to manipulate the result of the election

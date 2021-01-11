package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type speaker struct {
	*baseclient.BaseSpeaker
	c *client
}

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	c.Logf("Team 5 became speaker, drop the mic")
	return &speaker{c: c, BaseSpeaker: &baseclient.BaseSpeaker{GameState: c.ServerReadHandle.GetGameState()}}
}

// Pay Judge based on the status of our own wealth
// If we are not doing verywell, pay Judge less so we have more in the CP to take from
func (s *speaker) PayJudge() shared.SpeakerReturnContent {
	JudgeSalaryRule, ok := s.GameState.RulesInfo.CurrentRulesInPlay["salary_cycle_judge"]
	var JudgeSalary shared.Resources = 0
	if ok {
		JudgeSalary = shared.Resources(JudgeSalaryRule.ApplicableMatrix.At(0, 1))
	}
	if s.c.wealth() == jeffBezos {
		return shared.SpeakerReturnContent{
			ContentType: shared.SpeakerJudgeSalary,
			JudgeSalary: JudgeSalary,
			ActionTaken: true,
		}
	} else if s.c.wealth() == middleClass {
		JudgeSalary = JudgeSalary * 0.8
	} else {
		JudgeSalary = JudgeSalary * 0.5
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
	//TODO: disregard islands with sanctions
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
	if monitoring.Performed && !monitoring.Result {
		electionsettings.HoldElection = true
	}
	if turnsInPower >= 2 {
		electionsettings.HoldElection = true
	}
	return electionsettings
}

// if the real winner is on our bad side, then we choose our best friend
func (s *speaker) DecideNextJudge(winner shared.ClientID) shared.ClientID {
	aliveTeams := s.c.getAliveTeams(false) //not including us
	if s.c.opinions[winner].getScore() < 0 {
		ballot := s.c.VoteForElection(shared.President, aliveTeams)
		winner = ballot[0] //choose the first one in Borda Vote
	}
	return winner
}

package iigointernal

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
	"github.com/pkg/errors"
)

type legislature struct {
	gameState     *gamestate.GameState
	gameConf      *config.IIGOConfig
	SpeakerID     shared.ClientID
	ruleToVote    rules.RuleMatrix
	ballotBox     voting.BallotBox
	votingResult  bool
	clientSpeaker roles.Speaker
	iigoClients   map[shared.ClientID]baseclient.Client
	monitoring    *monitor
	logger        shared.Logger
}

func (l *legislature) Logf(format string, a ...interface{}) {
	l.logger("[LEGISLATURE]: %v", fmt.Sprintf(format, a...))
}

// loadClientSpeaker checks client pointer is good and if not panics
func (l *legislature) loadClientSpeaker(clientSpeakerPointer roles.Speaker) {
	if clientSpeakerPointer == nil {
		panic(fmt.Sprintf("Client '%v' has loaded a nil speaker pointer", l.SpeakerID))
	}
	l.clientSpeaker = clientSpeakerPointer
}

func (l *legislature) syncWithGame(gameState *gamestate.GameState, gameConf *config.IIGOConfig) {
	l.gameState = gameState
	l.gameConf = gameConf
}

// sendJudgeSalary conduct the transaction based on amount from client implementation
func (l *legislature) sendJudgeSalary() error {
	if l.clientSpeaker != nil {
		amountReturn := l.clientSpeaker.PayJudge()
		if amountReturn.ActionTaken && amountReturn.ContentType == shared.SpeakerJudgeSalary {
			// Subtract from common resources pool
			amountWithdraw, withdrawSuccess := WithdrawFromCommonPool(amountReturn.JudgeSalary, l.gameState)

			if withdrawSuccess {
				// Pay into the client private resources pool
				depositIntoClientPrivatePool(amountWithdraw, l.gameState.JudgeID, l.gameState)

				variablesToCache := []rules.VariableFieldName{rules.JudgePayment}
				valuesToCache := [][]float64{{float64(amountWithdraw)}}
				l.monitoring.addToCache(l.SpeakerID, variablesToCache, valuesToCache)
				return nil
			}
		}
		variablesToCache := []rules.VariableFieldName{rules.JudgePaid}
		valuesToCache := [][]float64{{boolToFloat(amountReturn.ActionTaken)}}
		l.monitoring.addToCache(l.SpeakerID, variablesToCache, valuesToCache)
	}
	return errors.Errorf("Cannot perform sendJudgeSalary")
}

// Receive a rule to call a vote on
func (l *legislature) setRuleToVote(ruleMatrix rules.RuleMatrix) error {
	if !CheckEnoughInCommonPool(l.gameConf.SetRuleToVoteActionCost, l.gameState) {
		return errors.Errorf("Insufficient Budget in common Pool: setRuleToVote")
	}

	agendaReturn := l.clientSpeaker.DecideAgenda(ruleMatrix)
	if agendaReturn.ActionTaken && agendaReturn.ContentType == shared.SpeakerAgenda {
		if !l.incurServiceCharge(l.gameConf.SetRuleToVoteActionCost) {
			return errors.Errorf("Insufficient Budget in common Pool: setRuleToVote")
		}
		l.ruleToVote = agendaReturn.RuleMatrix
	}
	return nil
}

//Asks islands to vote on a rule
//Called by orchestration
func (l *legislature) setVotingResult(clientIDs []shared.ClientID) (bool, error) {
	voteCalled := false
	if !CheckEnoughInCommonPool(l.gameConf.SetVotingResultActionCost, l.gameState) {
		return voteCalled, errors.Errorf("Insufficient Budget in common Pool: announceVotingResult")
	}

	returnVote := l.clientSpeaker.DecideVote(l.ruleToVote, clientIDs)
	returnedIslands := make([]shared.ClientID, len(returnVote.ParticipatingIslands))
	copy(returnedIslands, returnVote.ParticipatingIslands)
	sort.Sort(shared.SortClientByID(returnedIslands))
	sort.Sort(shared.SortClientByID(clientIDs))

	//log rule: All islands must participate in voting
	variablesToCache := []rules.VariableFieldName{rules.AllIslandsAllowedToVote}
	valuesToCache := [][]float64{{boolToFloat(reflect.DeepEqual(returnedIslands, clientIDs))}}
	l.monitoring.addToCache(l.SpeakerID, variablesToCache, valuesToCache)

	if returnVote.ActionTaken && returnVote.ContentType == shared.SpeakerVote {
		if !l.incurServiceCharge(l.gameConf.SetVotingResultActionCost) {
			return voteCalled, errors.Errorf("Insufficient Budget in common Pool: setVotingResult")
		}
		l.ballotBox = l.RunVote(returnVote.RuleMatrix, returnVote.ParticipatingIslands)

		l.votingResult = l.ballotBox.CountVotesMajority()
		voteCalled = true
	}
	return voteCalled, nil
}

//RunVote creates the voting object, returns votes by category (for, against) in BallotBox.
//Passing in empty ruleID or empty clientIDs results in no vote occurring
func (l *legislature) RunVote(ruleMatrix rules.RuleMatrix, clientIDs []shared.ClientID) voting.BallotBox {

	if ruleMatrix.RuleMatrixIsEmpty() || len(clientIDs) == 0 {
		return voting.BallotBox{}
	}
	l.Logf("Rule vote with islands %v allowed to vote", clientIDs)
	ruleVote := voting.RuleVote{
		Logger: l.logger,
	}

	//TODO: check if rule is valid, otherwise return empty ballot, raise error?
	ruleVote.SetRule(ruleMatrix)

	//TODO: intersection of islands alive and islands chosen to vote in case of client error
	//TODO: check if remaining slice is >0, otherwise return empty ballot, raise error?
	ruleVote.SetVotingIslands(clientIDs)

	ruleVote.GatherBallots(l.iigoClients)
	//TODO: log of vote occurring with ruleMatrix, clientIDs
	//TODO: log of clientIDs vs islandsAllowedToVote
	//TODO: log of ruleMatrix vs s.RuleToVote

	rulesEqual := false
	if reflect.DeepEqual(ruleMatrix, l.ruleToVote) {
		rulesEqual = true
	}

	variablesToCache := []rules.VariableFieldName{rules.SpeakerProposedPresidentRule}
	valuesToCache := [][]float64{{boolToFloat(rulesEqual)}}

	l.monitoring.addToCache(l.SpeakerID, variablesToCache, valuesToCache)

	return ruleVote.GetBallotBox()
}

//Speaker declares a result of a vote (see spec to see conditions on what this means for a rule-abiding speaker)
//Called by orchestration
func (l *legislature) announceVotingResult() (bool, error) {
	resultAnnounced := false
	if !CheckEnoughInCommonPool(l.gameConf.AnnounceVotingResultActionCost, l.gameState) {
		return resultAnnounced, errors.Errorf("Insufficient Budget in common Pool: announceVotingResult")
	}

	returnAnnouncement := l.clientSpeaker.DecideAnnouncement(l.ruleToVote, l.votingResult)

	if returnAnnouncement.ActionTaken && returnAnnouncement.ContentType == shared.SpeakerAnnouncement {
		//Deduct action cost
		if !l.incurServiceCharge(l.gameConf.AnnounceVotingResultActionCost) {
			return resultAnnounced, errors.Errorf("Insufficient Budget in common Pool: announceVotingResult")
		}

		//Perform announcement
		broadcastToAllIslands(l.iigoClients, shared.TeamIDs[l.SpeakerID], generateVotingResultMessage(returnAnnouncement.RuleMatrix, returnAnnouncement.VotingResult), *l.gameState)
		resultAnnounced = true

		//log rule "must announce what was called"
		announcementRuleMatchesVote := reflect.DeepEqual(returnAnnouncement.RuleMatrix, l.ruleToVote)
		announcementResultMatchesVote := returnAnnouncement.VotingResult == l.votingResult
		variablesToCache := []rules.VariableFieldName{rules.AnnouncementRuleMatchesVote, rules.AnnouncementResultMatchesVote}
		valuesToCache := [][]float64{{boolToFloat(announcementRuleMatchesVote)}, {boolToFloat(announcementResultMatchesVote)}}
		l.monitoring.addToCache(l.SpeakerID, variablesToCache, valuesToCache)
		l.Logf("Rule: %v , voted in by islands: %v , result heeded by speaker: %v", returnAnnouncement.RuleMatrix.RuleName, l.votingResult, returnAnnouncement.VotingResult)

	}
	return resultAnnounced, nil
}

func generateVotingResultMessage(ruleMatrix rules.RuleMatrix, result bool) map[shared.CommunicationFieldName]shared.CommunicationContent {
	returnMap := map[shared.CommunicationFieldName]shared.CommunicationContent{}

	returnMap[shared.RuleName] = shared.CommunicationContent{
		T:              shared.CommunicationString,
		RuleMatrixData: ruleMatrix,
	}
	returnMap[shared.RuleVoteResult] = shared.CommunicationContent{
		T:           shared.CommunicationBool,
		BooleanData: result,
	}

	return returnMap
}

// updateRules updates the rules in play according to the result of a vote.
func (l *legislature) updateRules(ruleMatrix rules.RuleMatrix, ruleIsVotedIn bool) error {
	if !l.incurServiceCharge(l.gameConf.UpdateRulesActionCost) {
		return errors.Errorf("Insufficient Budget in common Pool: updateRules")
	}
	//TODO: might want to log the errors as logging messages too?
	//notInRulesCache := errors.Errorf("Rule '%v' is not available in rules cache", ruleMatrix)
	if _, ok := l.gameState.RulesInfo.AvailableRules[ruleMatrix.RuleName]; !ok || reflect.DeepEqual(ruleMatrix, l.gameState.RulesInfo.AvailableRules[ruleMatrix.RuleName]) { //if the proposed ruleMatrix has the same content as the rule with the same name in AvailableRules, the proposal is for putting a rule in/out of play.
		if ruleIsVotedIn {
			err := l.gameState.PullRuleIntoPlay(ruleMatrix.RuleName)
			if ruleErr, ok := err.(*rules.RuleError); ok {
				if ruleErr.Type() == rules.RuleNotInAvailableRulesCache {
					return ruleErr
				}
			}
		} else {
			err := l.gameState.PullRuleOutOfPlay(ruleMatrix.RuleName)
			if ruleErr, ok := err.(*rules.RuleError); ok {
				if ruleErr.Type() == rules.RuleNotInAvailableRulesCache {
					return ruleErr
				}
			}

		}
	} else { //if the proposed ruleMatrix has different content to the rule with the same name in AvailableRules, the proposal is for modifying the rule in the rule caches. It doesn't put a rule in/out of play.
		if ruleIsVotedIn {
			err := l.gameState.ModifyRule(ruleMatrix.RuleName, ruleMatrix.ApplicableMatrix, ruleMatrix.AuxiliaryVector)
			return err
		}
	}

	return nil

}

// appointNextJudge returns the island ID of the island appointed to be Judge in the next turn
func (l *legislature) appointNextJudge(monitoring shared.MonitorResult, currentJudge shared.ClientID, allIslands []shared.ClientID) (shared.ClientID, error) {
	var election = voting.Election{
		Logger: l.logger,
	}
	var appointedJudge shared.ClientID
	allIslandsCopy1 := copyClientList(allIslands)
	electionSettings := l.clientSpeaker.CallJudgeElection(monitoring, int(l.gameState.IIGOTurnsInPower[shared.Judge]), allIslandsCopy1)

	//Log election rule
	termCondition := l.gameState.IIGOTurnsInPower[shared.Judge] > l.gameConf.IIGOTermLengths[shared.Judge]
	variablesToCache := []rules.VariableFieldName{rules.TermEnded, rules.ElectionHeld}
	valuesToCache := [][]float64{{boolToFloat(termCondition)}, {boolToFloat(electionSettings.HoldElection)}}
	l.monitoring.addToCache(l.SpeakerID, variablesToCache, valuesToCache)

	if electionSettings.HoldElection {
		if !l.incurServiceCharge(l.gameConf.AppointNextJudgeActionCost) {
			return l.gameState.JudgeID, errors.Errorf("Insufficient Budget in common Pool: appointNextJudge")
		}
		election.ProposeElection(shared.Judge, electionSettings.VotingMethod)
		allIslandsCopy2 := copyClientList(allIslands)
		election.OpenBallot(electionSettings.IslandsToVote, allIslandsCopy2)
		election.Vote(l.iigoClients)
		l.gameState.IIGOTurnsInPower[shared.Judge] = 0
		electedJudge := election.CloseBallot(l.iigoClients)
		appointedJudge = l.clientSpeaker.DecideNextJudge(electedJudge)

		//Log rule: Must appoint elected role
		appointmentMatchesVote := appointedJudge == electedJudge
		variablesToCache := []rules.VariableFieldName{rules.AppointmentMatchesVote}
		valuesToCache := [][]float64{{boolToFloat(appointmentMatchesVote)}}
		l.monitoring.addToCache(l.SpeakerID, variablesToCache, valuesToCache)
		l.Logf("Result of election for new Judge: %v", appointedJudge)
	} else {
		appointedJudge = currentJudge
	}
	l.gameState.IIGOElection = append(l.gameState.IIGOElection, election.GetVotingInfo())
	return appointedJudge, nil
}

func (l *legislature) incurServiceCharge(cost shared.Resources) bool {
	_, ok := WithdrawFromCommonPool(cost, l.gameState)
	if ok {
		l.gameState.IIGORolesBudget[shared.Speaker] -= cost
		if l.monitoring != nil {
			variablesToCache := []rules.VariableFieldName{rules.SpeakerLeftoverBudget}
			valuesToCache := [][]float64{{float64(l.gameState.IIGORolesBudget[shared.Speaker])}}
			l.monitoring.addToCache(l.SpeakerID, variablesToCache, valuesToCache)
		}
	}
	return ok
}

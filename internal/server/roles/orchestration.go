package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/pkg/errors"
)

var Base_judge = BaseJudge{
	id:                0,
	budget:            0,
	presidentSalary:   0,
	BallotID:          0,
	ResAllocID:        0,
	speakerID:         0,
	presidentID:       0,
	evaluationResults: nil,
}

var Base_speaker = baseSpeaker{
	id:          0,
	budget:      0,
	judgeSalary: 0,
	ruleToVote:  "",
}

var Base_President = basePresident{
	id:               0,
	budget:           0,
	speakerSalary:    0,
	resourceRequests: nil,
}

var SpeakerIDGlobal = 0
var JudgeIDGlobal = 0
var PresidentIDGlobal = 0

var judgePointer = Base_judge
var speakerPointer = Base_speaker
var presidentPointer = Base_President

func RunIIGO(g *gamestate.GameState) error {
	// Initialise IDs
	Base_judge.id = JudgeIDGlobal
	Base_speaker.id = SpeakerIDGlobal
	Base_President.id = PresidentIDGlobal

	// Initialise roles with their clientVersions
	Base_judge.clientJudge = &judgePointer
	Base_President.clientPresident = &presidentPointer
	Base_speaker.clientSpeaker = &speakerPointer

	// Withdraw the salaries
	errWithdrawPresident := judgePointer.withdrawPresidentSalary(g)
	errWithdrawJudge := speakerPointer.withdrawJudgeSalary(g)
	errWithdrawSpeaker := presidentPointer.withdrawSpeakerSalary(g)

	// Handle the lack of resources
	if errWithdrawPresident != nil {
		return errors.Errorf("Could not run IIGO since President has no resoruces to spend")
	}
	if errWithdrawJudge != nil {
		return errors.Errorf("Could not run IIGO since Judge has no resoruces to spend")
	}
	if errWithdrawSpeaker != nil {
		return errors.Errorf("Could not run IIGO since Speaker has no resoruces to spend")
	}

	// Pay salaries into budgets
	judgePointer.payPresident()
	speakerPointer.payJudge()
	presidentPointer.paySpeaker()

	// 1 Judge actions - inspect history
	_, judgeInspectingHistoryError := Base_judge.inspectHistory()

	// 2 President actions
	presidentPointer.requestAllocationRequest()
	presidentPointer.replyAllocationRequest(g.CommonPool)
	presidentPointer.requestRuleProposal()
	ruleToVote := presidentPointer.getRuleForSpeaker()

	// 3 Speaker actions
	speakerPointer.SetRuleToVote(ruleToVote)
	speakerPointer.setVotingResult()
	speakerPointer.announceVotingResult()

	// 4 Declare performance (Judge) (in future all the roles)
	if judgeInspectingHistoryError != nil {
		Base_judge.declarePresidentPerformanceWrapped()

		Base_judge.declareSpeakerPerformanceWrapped()
	}

	// Get new Judge ID
	JudgeIDGlobal = speakerPointer.appointNextJudge()
	// Get new Speaker ID
	SpeakerIDGlobal = presidentPointer.appointNextSpeaker()
	// Get new President ID
	PresidentIDGlobal = judgePointer.appointNextPresident()

	// Set judgePointer
	// https://imgur.com/a/HjVZIkh
	// Set speakerPointer

	// Set presidentPointer

	return nil
}

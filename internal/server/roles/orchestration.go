package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
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
	ruleToVote:  0,
}

var Base_President = basePresident{
	id:                 0,
	budget:             0,
	speakerSalary:      0,
	resourceRequests:   nil,
	resourceAllocation: nil,
	ruleToVote:         0,
	taxAmount:          0,
}

var SpeakerIDGlobal = 0
var JudgeIDGlobal = 0
var PresidentIDGlobal = 0

var judgePointer = Base_judge
var speakerPointer = Base_speaker
var presidentPointer = Base_President

func runIIGO() error {

	// Initialise IDs
	Base_judge.id = JudgeIDGlobal
	Base_speaker.id = SpeakerIDGlobal
	Base_President.id = PresidentIDGlobal

	// Initialise roles with their clientVersions
	Base_judge.clientJudge = &judgePointer

	// Pay the salaries
	errPayPresident := judgePointer.payPresident()
	errPayJudge := speakerPointer.PayJudge()
	errPaySpeaker := presidentPointer.paySpeaker()

	// Handle the lack of resources
	if errPayPresident == nil {
		return errors.Errorf("Could not run IIGO since President has no resoruces to spend")
	}

	if errPayJudge == nil {
		return errors.Errorf("Could not run IIGO since Judge has no resoruces to spend")
	}

	if errPaySpeaker == nil {
		return errors.Errorf("Could not run IIGO since Speaker has no resoruces to spend")
	}

	// 1 Judge actions - inspect history
	_, judgeInspectingHistoryError := Base_judge.inspectHistory()

	// 2 Speaker actions

	// 3 President actions

	// 4 Declare performance (Judge) (in future all the roles)
	if judgeInspectingHistoryError != nil {
		Base_judge.declarePresidentPerformanceWrapped()

		Base_judge.declareSpeakerPerformanceWrapped()
	}

	//TODO: Add election setting
	// Set JudgeIDGlobal
	// Set judgePointer
	// Set SpeakerIDGlobal
	// Set speakerPointer
	// Set PresidentIDGlobal
	// Set presidentPointer
	return nil
}

// callVote possible implementation of voting
func callVote(speakerID int, whateverIsBeingVotedOn string) {
	// Do voting

	noIslandAlive := rules.VariableValuePair{
		VariableName: "no_islands_alive",
		Values:       []float64{5},
	}
	noIslandsVoting := rules.VariableValuePair{
		VariableName: "no_islands_voted",
		Values:       []float64{5},
	}
	err := updateTurnHistory(speakerID, []rules.VariableValuePair{noIslandAlive, noIslandsVoting})
	if err != nil {
		// exit with error
	} else {
		// carry on
	}
}

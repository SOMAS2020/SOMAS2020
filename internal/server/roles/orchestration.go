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
		RID, result1, PID, checkRole1, judgeDeclaringPresidentPerformanceError := Base_judge.declarePresidentPerformance()
		if judgeDeclaringPresidentPerformanceError == nil {
			broadcastToAllIslands(Base_judge.id, generatePresidentPerformanceMessage(RID, result1, PID, checkRole1))
		}

		BID, result2, SID, checkRole2, judgeDeclaringSpeakerPerformanceError := Base_judge.declareSpeakerPerformance()
		if judgeDeclaringSpeakerPerformanceError == nil {
			broadcastToAllIslands(Base_judge.id, generateSpeakerPerformanceMessage(BID, result2, SID, checkRole2))
		}
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

func generateSpeakerPerformanceMessage(BID int, result bool, SID int, conductedRole bool) map[int]DataPacket {
	returnMap := map[int]DataPacket{}

	returnMap[BallotID] = DataPacket{
		integerData: BID,
	}
	returnMap[SpeakerBallotCheck] = DataPacket{
		booleanData: result,
	}
	returnMap[SpeakerID] = DataPacket{
		integerData: SID,
	}
	returnMap[RoleConducted] = DataPacket{
		booleanData: conductedRole,
	}
	return returnMap
}

func generatePresidentPerformanceMessage(RID int, result bool, PID int, conductedRole bool) map[int]DataPacket {
	returnMap := map[int]DataPacket{}

	returnMap[ResAllocID] = DataPacket{
		integerData: RID,
	}
	returnMap[PresidentAllocationCheck] = DataPacket{
		booleanData: result,
	}
	returnMap[PresidentID] = DataPacket{
		integerData: PID,
	}
	returnMap[RoleConducted] = DataPacket{
		booleanData: conductedRole,
	}
	return returnMap
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

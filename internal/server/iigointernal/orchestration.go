package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

// featureJudge is an instantiation of the Judge interface
// with both the Base Judge features and a reference to client judges
var featureJudge = BaseJudge{
	id:                0,
	budget:            0,
	presidentSalary:   0,
	BallotID:          0,
	ResAllocID:        0,
	speakerID:         0,
	presidentID:       0,
	evaluationResults: nil,
}

// featureSpeaker is an instantiation of the Speaker interface
// with both the baseSpeaker features and a reference to client speakers
var featureSpeaker = baseSpeaker{
	id:          0,
	budget:      0,
	judgeSalary: 0,
	ruleToVote:  "",
}

// featurePresident is an instantiation of the President interface
// with both the basePresident features and a reference to client presidents
var featurePresident = basePresident{
	id:               0,
	budget:           0,
	speakerSalary:    0,
	resourceRequests: nil,
}

// SpeakerIDGlobal is the single source of truth for speaker ID (MVP)
var SpeakerIDGlobal = 0

// JudgeIDGlobal is the single source of truth for judge ID (MVP)
var JudgeIDGlobal = 0

// PresidentIDGlobal is the single source of truth for president ID (MVP)
var PresidentIDGlobal = 0

// Pointers allow clients to customise implementations of mutable functions
var judgePointer roles.Judge = nil
var speakerPointer roles.Speaker = nil
var presidentPointer roles.President = nil

// iigoClients holds pointers to all the clients
var iigoClients map[shared.ClientID]baseclient.Client

// RunIIGO runs all iigo function in sequence
func RunIIGO(g *gamestate.GameState, clientMap *map[shared.ClientID]baseclient.Client) error {

	// TODO: Get Client pointers from gamestate https://imgur.com/a/HjVZIkh
	iigoClients = *clientMap

	// Initialise IDs
	featureJudge.id = JudgeIDGlobal
	featureSpeaker.id = SpeakerIDGlobal
	featurePresident.id = PresidentIDGlobal

	// Initialise iigointernal with their clientVersions
	featureJudge.clientJudge = judgePointer
	featurePresident.clientPresident = presidentPointer
	featureSpeaker.clientSpeaker = speakerPointer

	// Withdraw the salaries
	errWithdrawPresident := featureJudge.withdrawPresidentSalary(g)
	errWithdrawJudge := featureSpeaker.withdrawJudgeSalary(g)
	errWithdrawSpeaker := featurePresident.withdrawSpeakerSalary(g)

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
	featureJudge.sendPresidentSalary()
	featureSpeaker.sendJudgeSalary()
	featurePresident.sendSpeakerSalary()

	// 1 Judge actions - inspect history
	_, judgeInspectingHistoryError := featureJudge.inspectHistory()

	// 2 President actions
	resourceReports := map[int]int{}
	for _, v := range rules.VariableMap["islands_alive"].Values {
		resourceReports[int(v)] = iigoClients[shared.ClientID(int(v))].ResourceReport()
	}
	featurePresident.broadcastTaxation(resourceReports)
	featurePresident.requestAllocationRequest()
	featurePresident.replyAllocationRequest(g.CommonPool)
	featurePresident.requestRuleProposal()
	ruleToVote := featurePresident.getRuleForSpeaker()

	// 3 Speaker actions
	featureSpeaker.SetRuleToVote(ruleToVote)
	featureSpeaker.setVotingResult()
	featureSpeaker.announceVotingResult()

	// 4 Declare performance (Judge) (in future all the iigointernal)
	if judgeInspectingHistoryError != nil {
		featureJudge.declarePresidentPerformanceWrapped()

		featureJudge.declareSpeakerPerformanceWrapped()
	}

	// Get new Judge ID
	JudgeIDGlobal = featureSpeaker.appointNextJudge()
	// Get new Speaker ID
	SpeakerIDGlobal = featurePresident.appointNextSpeaker()
	// Get new President ID
	PresidentIDGlobal = featureJudge.appointNextPresident()

	// Set judgePointer
	judgePointer = iigoClients[shared.ClientID(JudgeIDGlobal)].GetClientJudgePointer()
	// Set speakerPointer
	speakerPointer = iigoClients[shared.ClientID(SpeakerIDGlobal)].GetClientSpeakerPointer()
	// Set presidentPointer
	presidentPointer = iigoClients[shared.ClientID(PresidentIDGlobal)].GetClientPresidentPointer()

	return nil
}

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
	ID:                0,
	budget:            0,
	presidentSalary:   0,
	BallotID:          0,
	ResAllocID:        0,
	speakerID:         0,
	presidentID:       0,
	EvaluationResults: nil,
}

// featureSpeaker is an instantiation of the Speaker interface
// with both the baseSpeaker features and a reference to client speakers
var featureSpeaker = baseSpeaker{
	Id:          0,
	budget:      0,
	judgeSalary: 0,
	RuleToVote:  "",
}

// featurePresident is an instantiation of the President interface
// with both the basePresident features and a reference to client presidents
var featurePresident = basePresident{
	ID:               0,
	budget:           0,
	speakerSalary:    0,
	ResourceRequests: nil,
}

// GetFeaturedRoles returns featured versions of the roles
func GetFeaturedRoles() (roles.Judge, roles.Speaker, roles.President) {
	return &featureJudge, &featureSpeaker, &featurePresident
}

// SpeakerIDGlobal is the single source of truth for speaker ID (MVP)
var SpeakerIDGlobal shared.ClientID

// JudgeIDGlobal is the single source of truth for judge ID (MVP)
var JudgeIDGlobal shared.ClientID

// PresidentIDGlobal is the single source of truth for president ID (MVP)
var PresidentIDGlobal shared.ClientID

// TaxAmountMapExport is a local tax amount cache for checking of rules
var TaxAmountMapExport map[shared.ClientID]shared.Resources

// AllocationAmountMapExport is a local allocation map for checking of rules
var AllocationAmountMapExport map[shared.ClientID]shared.Resources

// Pointers allow clients to customise implementations of mutable functions
var judgePointer roles.Judge = nil
var speakerPointer roles.Speaker = nil
var presidentPointer roles.President = nil

// iigoClients holds pointers to all the clients
var iigoClients map[shared.ClientID]baseclient.Client

var iigoRoleStates gamestate.IIGOBaseRoles

// RunIIGO runs all iigo function in sequence
func RunIIGO(g *gamestate.GameState, clientMap *map[shared.ClientID]baseclient.Client) error {

	// TODO: Get Client pointers from gamestate https://imgur.com/a/HjVZIkh
	iigoClients = *clientMap
	iigoRoleStates = g.IIGOInfo

	// Initialise IDs
	featureJudge.ID = JudgeIDGlobal
	featureSpeaker.Id = SpeakerIDGlobal
	featurePresident.ID = PresidentIDGlobal

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
		returnWithdrawnSalariesToCommonPool(g)
		return errors.Errorf("Could not run IIGO since President has no resoruces to spend")
	}
	if errWithdrawJudge != nil {
		returnWithdrawnSalariesToCommonPool(g)
		return errors.Errorf("Could not run IIGO since Judge has no resoruces to spend")
	}
	if errWithdrawSpeaker != nil {
		returnWithdrawnSalariesToCommonPool(g)
		return errors.Errorf("Could not run IIGO since Speaker has no resoruces to spend")
	}

	// Pay salaries into budgets
	featureJudge.sendPresidentSalary()
	featureSpeaker.sendJudgeSalary()
	featurePresident.sendSpeakerSalary()

	// 1 Judge actions - inspect history
	_, judgeInspectingHistoryError := featureJudge.InspectHistory()

	// 2 President actions
	resourceReports := map[shared.ClientID]shared.Resources{}
	var aliveClientIds []shared.ClientID
	for _, v := range rules.VariableMap["islands_alive"].Values {
		aliveClientIds = append(aliveClientIds, shared.ClientID(int(v)))
		resourceReports[shared.ClientID(int(v))] = iigoClients[shared.ClientID(int(v))].ResourceReport()
	}
	featurePresident.broadcastTaxation(resourceReports)
	featurePresident.requestAllocationRequest()
	featurePresident.replyAllocationRequest(g.CommonPool)
	featurePresident.requestRuleProposal()
	ruleToVote := featurePresident.getRuleForSpeaker()

	// 3 Speaker actions
	featureSpeaker.setRuleToVote(ruleToVote)
	featureSpeaker.setVotingResult(iigoClients)
	_ = featureSpeaker.announceVotingResult()

	// 4 Declare performance (Judge) (in future all the iigointernal)
	if judgeInspectingHistoryError != nil {
		featureJudge.declarePresidentPerformanceWrapped()

		featureJudge.declareSpeakerPerformanceWrapped()
	}

	// Get new Judge ID
	JudgeIDGlobal = featureSpeaker.appointNextJudge(aliveClientIds)
	// Get new Speaker ID
	SpeakerIDGlobal = featurePresident.appointNextSpeaker(aliveClientIds)
	// Get new President ID
	PresidentIDGlobal = featureJudge.appointNextPresident(aliveClientIds)

	// Set judgePointer
	judgePointer = iigoClients[shared.ClientID(JudgeIDGlobal)].GetClientJudgePointer()
	// Set speakerPointer
	speakerPointer = iigoClients[shared.ClientID(SpeakerIDGlobal)].GetClientSpeakerPointer()
	// Set presidentPointer
	presidentPointer = iigoClients[shared.ClientID(PresidentIDGlobal)].GetClientPresidentPointer()

	return nil
}

func returnWithdrawnSalariesToCommonPool(state *gamestate.GameState) {
	returnVal := featurePresident.returnSpeakerSalary() + featureSpeaker.returnJudgeSalary() + featureJudge.returnPresidentSalary()
	depositIntoCommonPool(returnVal, state)
}

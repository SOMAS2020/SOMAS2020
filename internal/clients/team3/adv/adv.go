package adv

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

type Spec int

const (
	NoAdv Spec = iota
	MaliceAdv
	TargetAdv
)

type Adv interface {
	Initialise(clientID shared.ClientID)
	ProposeRule(availableRules map[string]rules.RuleMatrix) (rules.RuleMatrix, bool)
	GetRuleViolationSeverity() (map[string]shared.IIGOSanctionsScore, bool)
	GetPardonedIslands(currentSanctions map[int][]shared.Sanction) (map[int][]bool, bool)
	CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) (shared.ElectionSettings, bool)
	DecideNextPresident(winner shared.ClientID) (shared.ClientID, bool)
	CallJudgeElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) (shared.ElectionSettings, bool)
	DecideNextJudge(winner shared.ClientID) (shared.ClientID, bool)
	CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) (shared.ElectionSettings, bool)
	DecideNextSpeaker(winner shared.ClientID) (shared.ClientID, bool)
	InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]shared.EvaluationReturn, bool, bool)
	SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) (shared.PresidentReturnContent, bool)
}

func copySingleRuleMatrix(inp rules.RuleMatrix) rules.RuleMatrix {
	return rules.RuleMatrix{
		RuleName:          inp.RuleName,
		RequiredVariables: copyRequiredVariables(inp.RequiredVariables),
		ApplicableMatrix:  *mat.DenseCopyOf(&inp.ApplicableMatrix),
		AuxiliaryVector:   *mat.VecDenseCopyOf(&inp.AuxiliaryVector),
		Mutable:           inp.Mutable,
		Link:              copyLink(inp.Link),
	}
}

func copyLink(inp rules.RuleLink) rules.RuleLink {
	return rules.RuleLink{
		Linked:     inp.Linked,
		LinkType:   inp.LinkType,
		LinkedRule: inp.LinkedRule,
	}
}

func copyRequiredVariables(inp []rules.VariableFieldName) []rules.VariableFieldName {
	targetList := make([]rules.VariableFieldName, len(inp))
	copy(targetList, inp)
	return targetList
}

package team6

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type president struct {
	*baseclient.BasePresident
	*client
}

func (p *president) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) shared.PresidentReturnContent {
	var requestSum shared.Resources
	var setStrategy string = "normal"
	resourceAllocation := make(map[shared.ClientID]shared.Resources)
	var otherIslandsRequests shared.Resources = 0.0
	var multiplier shared.Resources = 0.75
	var multiplierForOtherIslands shared.Resources = 0.75

	/*
		// When to use which strategy?
		if p.ServerReadHandle.GetGameState().ClientInfo.Resources < p.ServerReadHandle.GetGameConfig().MinimumResourceThreshold+p.ServerReadHandle.GetGameConfig().CostOfLiving {
			setStrategy := "egoistic"
		}
	*/

	if setStrategy == "normal" {
		for _, request := range resourceRequest {
			requestSum += request
		}

		if requestSum < multiplier*availCommonPool || requestSum == 0 {
			resourceAllocation = resourceRequest
		} else {
			for id, request := range resourceRequest {
				resourceAllocation[id] = multiplier * availCommonPool * request / requestSum
			}
		}
	} else if setStrategy == "egoistic" {
		for id, request := range resourceRequest {
			if id != shared.Team6 {
				otherIslandsRequests += request
			}

			multiplierForOtherIslands = otherIslandsRequests / (multiplier * (availCommonPool - resourceAllocation[shared.Team6]))

			for id, request := range resourceRequest {
				resourceAllocation[id] = multiplierForOtherIslands * availCommonPool * request / requestSum
				if id == shared.Team6 {
					resourceAllocation[id] = request
				}
			}
		}

	}

	return shared.PresidentReturnContent{
		ContentType: shared.PresidentAllocation,
		ResourceMap: resourceAllocation,
		ActionTaken: true,
	}
}

func (p *president) PickRuleToVote(rulesProposals []rules.RuleMatrix) shared.PresidentReturnContent {
	// DefaulContentType: No rules were proposed by the islands
	proposedRule := ""
	actionTaken := false

	// if some rules were proposed
	if len(rulesProposals) != 0 {
		proposedRule = ""
		// 		rulesProposals[rand.Intn(len(rulesProposals))]
		actionTaken = true
	}

	return shared.PresidentReturnContent{
		ContentType: shared.PresidentRuleProposal,
		ProposedRuleMatrix: rules.RuleMatrix{
			RuleName: proposedRule,
		},
		ActionTaken: actionTaken,
	}
}

// TODO
func (p *president) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {
	taxAmountMap := make(map[shared.ClientID]shared.Resources)

	for clientID, clientReport := range islandsResources {
		if clientReport.Reported {
			taxAmountMap[clientID] = shared.Resources(float64(clientReport.ReportedAmount) * rand.Float64())
		} else {
			taxAmountMap[clientID] = 15 //flat tax rate
		}
	}
	return shared.PresidentReturnContent{
		ContentType: shared.PresidentTaxation,
		ResourceMap: taxAmountMap,
		ActionTaken: true,
	}
}

func (p *president) PaySpeaker(salary shared.Resources) shared.PresidentReturnContent {
	return shared.PresidentReturnContent{
		ContentType:   shared.PresidentSpeakerSalary,
		SpeakerSalary: salary,
		ActionTaken:   true,
	}
}

func (p *president) CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
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

func (p *president) DecideNextSpeaker(winner shared.ClientID) shared.ClientID {
	return winner
}

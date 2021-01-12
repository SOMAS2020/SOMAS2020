package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type president struct {
	*baseclient.BasePresident
	*client
}

func (p *president) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) shared.PresidentReturnContent {
	resourceAllocation := make(map[shared.ClientID]shared.Resources)
	ourPersonality := p.client.getPersonality()
	evaluationCoeff := shared.Resources(0.5)

	friendCoeffsOnAllocation := make(map[shared.ClientID]shared.Resources)

	for team, fsc := range p.client.getFriendshipCoeffs() {
		friendCoeffsOnAllocation[team] = shared.Resources(fsc)
	}

	if ourPersonality != Selfish {
		requestSum := shared.Resources(0.0)

		for _, request := range resourceRequest {
			requestSum += request
		}

		if requestSum <= evaluationCoeff*availCommonPool || requestSum == 0 {
			for team, request := range resourceRequest {
				resourceAllocation[team] = friendCoeffsOnAllocation[team] * request
			}
		} else {
			for team, request := range resourceRequest {
				resourceAllocation[team] = friendCoeffsOnAllocation[team] * evaluationCoeff * availCommonPool * request / requestSum
			}
		}
	} else {
		// trying to be a selfish president
		commonPoolLeft := evaluationCoeff * availCommonPool
		otherRequestSum := shared.Resources(0.0)

		for team, request := range resourceRequest {
			if team == p.client.GetID() {
				if request <= evaluationCoeff*availCommonPool {
					resourceAllocation[team] = request
				} else {
					resourceAllocation[team] = evaluationCoeff * availCommonPool
				}

				commonPoolLeft = evaluationCoeff*availCommonPool - resourceAllocation[id]

				continue
			}

			otherRequestSum += request
		}

		if otherRequestSum <= commonPoolLeft || otherRequestSum == 0 {
			for team, request := range resourceRequest {
				if team == p.client.GetID() {
					continue
				}

				resourceAllocation[team] = friendCoeffsOnAllocation[team] * request
			}
		} else {
			for team, request := range resourceRequest {
				if team == p.client.GetID() {
					continue
				}

				resourceAllocation[team] = friendCoeffsOnAllocation[team] * evaluationCoeff * commonPoolLeft * request / otherRequestSum
			}
		}
	}

	return shared.PresidentReturnContent{
		ContentType: shared.PresidentAllocation,
		ResourceMap: resourceAllocation,
		ActionTaken: true,
	}
}

func (p *president) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {
	taxAmountMap := make(map[shared.ClientID]shared.Resources)
	friendshipCoffesOnTax := make(map[shared.ClientID]float64)

	for team, fsc := range p.client.getFriendshipCoeffs() {
		friendshipCoffesOnTax[team] = float64(1) - fsc
	}

	for clientID, clientReport := range islandsResources {
		if clientReport.Reported {
			taxRate := 0.2
			taxAmountMap[clientID] = shared.Resources(float64(clientReport.ReportedAmount) * taxRate * friendshipCoffesOnTax[clientID])
		} else {
			taxAmountMap[clientID] = shared.Resources(30 * friendshipCoffesOnTax[clientID]) //flat tax rate
		}
	}

	return shared.PresidentReturnContent{
		ContentType: shared.PresidentTaxation,
		ResourceMap: taxAmountMap,
		ActionTaken: true,
	}
}

func (p *president) CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	currSpeakerID := p.client.ServerReadHandle.GetGameState().SpeakerID
	otherIslands := []shared.ClientID{}

	for _, team := range allIslands {
		if currSpeakerID == team {
			continue
		}

		otherIslands = append(otherIslands, team)
	}

	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.Runoff,
		IslandsToVote: otherIslands,
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
	if winner == p.client.ServerReadHandle.GetGameState().SpeakerID {
		return p.client.GetID()
	}

	if p.client.friendship[winner] < p.clientConfig.maxFriendship/1.5 {
		return p.client.GetID()
	}

	return winner
}

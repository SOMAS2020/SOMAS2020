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
	evaluationCoeff := shared.Resources(0.75)

	if ourPersonality != Selfish {
		requestSum := shared.Resources(0.0)

		for _, request := range resourceRequest {
			requestSum += request
		}

		if requestSum <= evaluationCoeff*availCommonPool || requestSum == 0 {
			resourceAllocation = resourceRequest
		} else {
			for id, request := range resourceRequest {
				resourceAllocation[id] = evaluationCoeff * availCommonPool * request / requestSum
			}
		}
	} else {
		// trying to be a selfish president
		commonPoolLeft := availCommonPool
		otherRequestSum := shared.Resources(0.0)

		for id, request := range resourceRequest {
			if id == p.client.GetID() {
				if request <= evaluationCoeff*availCommonPool {
					resourceAllocation[id] = request
				} else {
					resourceAllocation[id] = evaluationCoeff * availCommonPool
				}

				commonPoolLeft = availCommonPool - resourceAllocation[id]

				continue
			}

			otherRequestSum += request
		}

		if otherRequestSum <= evaluationCoeff*commonPoolLeft || otherRequestSum == 0 {
			for id, request := range resourceRequest {
				if id == p.client.GetID() {
					continue
				}

				resourceAllocation[id] = request
			}
		} else {
			for id, request := range resourceRequest {
				if id == p.client.GetID() {
					continue
				}

				resourceAllocation[id] = evaluationCoeff * commonPoolLeft * request / otherRequestSum
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

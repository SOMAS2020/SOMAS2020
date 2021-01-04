package team3

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type president struct {
	// Base implementation
	*baseclient.BasePresident
	// Our client
	c *client
}

func (p *president) PaySpeaker(salary shared.Resources) (shared.Resources, bool) {
	// Use the base implementation
	return p.BasePresident.PaySpeaker(salary)
}

func (p *president) DecideNextSpeaker(winner shared.ClientID) shared.ClientID {
	p.c.clientPrint("choosing speaker")
	// Naively choose group 0
	return mostTrusted(p.c.trustMapAgg)
}

// Computes average request, excluding top and bottom
func findAvgNoTails(resourceRequest map[shared.ClientID]shared.Resources) shared.Resources {
	var sum shared.Resources
	minClient := shared.TeamIDs[0]
	maxClient := shared.TeamIDs[0]

	// Find min and max requests
	for island, request := range resourceRequest {
		if request < resourceRequest[minClient] {
			minClient = island
		}
		if request > resourceRequest[maxClient] {
			maxClient = island
		}
	}

	// Compute average ignoring highest and lowest
	for island, request := range resourceRequest {
		if island != minClient || island != maxClient {
			sum += request
		}
	}

	return shared.Resources(int(sum) / len(resourceRequest))
}

// EvaluateAllocationRequests sets allowed resource allocation based on each islands requests
func (p *president) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
	p.c.clientPrint("Evaluating allocations...")
	// var allocations, allocWeights map[shared.ClientID]float64
	var avgResource, avgRequest shared.Resources
	// var resources map[shared.ClientID]shared.Resources
	var allocSum, commonPoolThreshold, sumRequest float64
	// var finalAllocations map[shared.ClientID]shared.Resources

	// Make sure resource skew is greater than 1
	resourceSkew := math.Max(float64(p.c.params.resourcesSkew), 1)

	// TEMP
	p.c.declaredResources = map[shared.ClientID]shared.Resources{
		shared.Team1: 10,
		shared.Team2: 10,
		shared.Team3: 10,
		shared.Team4: 10,
		shared.Team5: 10,
		shared.Team6: 10,
	}

	resources := make(map[shared.ClientID]shared.Resources)
	allocations := make(map[shared.ClientID]float64)
	allocWeights := make(map[shared.ClientID]float64)
	finalAllocations := make(map[shared.ClientID]shared.Resources)

	for island, req := range resourceRequest {
		sumRequest += float64(req)
		resources[island] = shared.Resources(float64(p.c.declaredResources[island]) * math.Pow(resourceSkew, 1-p.c.trustScore[island]))
	}

	p.c.clientPrint("Resource requests: %+v\n", resourceRequest)
	p.c.clientPrint("Their resource estimation: %+v\n Our estimation: %+v\n", p.c.declaredResources, resources)

	avgRequest = findAvgNoTails(resourceRequest)
	avgResource = findAvgNoTails(resources)

	for island, resource := range resources {
		allocations[island] = float64(avgRequest) + p.c.params.equity*(float64(avgResource-resource)+float64(resourceRequest[island]-avgRequest))
		// p.c.clientPrint("Allocation for island %v: %f", island, allocations[island])
		if island == id {
			allocations[island] += math.Max(float64(resourceRequest[island])-allocations[island]*p.c.params.selfishness, 0)
		} else {
			allocations[island] = math.Min(float64(resourceRequest[island]), allocations[island]) // to prevent overallocating
			allocations[island] = math.Max(allocations[island], 0)
		}
	}

	// Collect weights
	for _, alloc := range allocations {
		allocSum += alloc
	}
	// Normalise
	for island, alloc := range allocations {
		allocWeights[island] = alloc / allocSum
	}
	// p.c.clientPrint("Allocation wieghts: %+v\n", allocWeights)

	commonPoolThreshold = float64(availCommonPool) * (1.0 - p.c.params.riskFactor)
	if p.c.params.saveCriticalIsland {
		for island := range resourceRequest {
			if resources[island] < p.c.criticalStatePrediction.lowerBound {
				finalAllocations[island] = shared.Resources(math.Max((allocWeights[island] * commonPoolThreshold), float64(p.c.criticalStatePrediction.lowerBound-resources[island])))
			} else {
				finalAllocations[island] = 0
			}
		}
	}

	for island := range resourceRequest {
		if finalAllocations[island] == 0 {
			if sumRequest < commonPoolThreshold {
				finalAllocations[island] = shared.Resources(allocWeights[island] * float64(sumRequest))
			} else {
				finalAllocations[island] = shared.Resources(allocWeights[island] * commonPoolThreshold)
			}
		}
	}

	p.c.clientPrint("Final allocations: %+v\n", finalAllocations)

	// Curently always evaluate, would there be a time when we don't want to?
	return finalAllocations, true
}

func (p *president) SetTaxationAmount(islandsResources map[shared.ClientID]shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
	//decide if we want to run SetTaxationAmount
	p.c.declaredResources = islandsResources
	gameState := p.c.BaseClient.ServerReadHandle.GetGameState()
	resourcesRequired := 100.0 - float64(gameState.CommonPool)
	if len(p.c.disasterPredictions) > int(p.c.ServerReadHandle.GetGameState().Turn) {
		if disaster, ok := p.c.disasterPredictions[int(p.c.ServerReadHandle.GetGameState().Turn)][p.c.GetID()]; ok {
			resourcesRequired = (disaster.Magnitude - float64(gameState.CommonPool)/float64(disaster.TimeLeft))
		}
	}
	AveTax := resourcesRequired / float64(len(islandsResources))
	var adjustedResources []float64
	adjustedResourcesMap := make(map[shared.ClientID]shared.Resources)
	for island, resource := range islandsResources {
		adjustedResource := resource * shared.Resources(math.Pow(p.c.params.resourcesSkew, (100-p.c.trustScore[island])/100))
		adjustedResources = append(adjustedResources, float64(adjustedResource))
		adjustedResourcesMap[island] = adjustedResource
	}
	AveAdjustedResources := getAverage(adjustedResources)
	taxationMap := make(map[shared.ClientID]shared.Resources)
	for island, resources := range adjustedResourcesMap {
		taxation := shared.Resources(AveTax) + shared.Resources(p.c.params.equity)*(resources-shared.Resources(AveAdjustedResources))
		if island == p.c.BaseClient.GetID() {
			taxation -= shared.Resources(p.c.params.selfishness) * taxation
		}
		taxation = shared.Resources(math.Max(float64(taxation), 0.0))
		taxationMap[island] = taxation
	}
	p.c.clientPrint("tax amounts : %v\n", taxationMap)
	return taxationMap, true

}

package team3

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type president struct {
	// Base implementation
	*baseclient.BasePresident
	// Our client
	c *client
}

// CallSpeakerElection is called by the executive to decide on power-transfer
// COMPULSORY: decide when to call an election following relevant rulesInPlay if you wish
func (p *president) CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	// example implementation calls an election if monitoring was performed and the result was negative
	// or if the number of turnsInPower exceeds 3
	if p.c.params.adv != nil {
		ret, done := p.c.params.adv.CallSpeakerElection(monitoring, turnsInPower, allIslands)
		if done {
			return ret
		}
	}

	return p.BasePresident.CallSpeakerElection(monitoring, turnsInPower, allIslands)
}

func (p *president) DecideNextSpeaker(winner shared.ClientID) shared.ClientID {
	if p.c.params.adv != nil {
		ret, done := p.c.params.adv.DecideNextSpeaker(winner)
		if done {
			return ret
		}
	}

	p.c.clientPrint("choosing speaker")
	// Naively choose group 0
	return mostTrusted(p.c.trustScore)

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
func (p *president) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) shared.PresidentReturnContent {
	p.c.clientPrint("Evaluating allocations...")
	var avgResource, avgRequest shared.Resources
	var allocSum, commonPoolThreshold, sumRequest float64

	// Make sure resource skew is greater than 1
	resourceSkew := math.Max(float64(p.c.params.resourcesSkew), 1)

	resources := make(map[shared.ClientID]shared.Resources)
	allocations := make(map[shared.ClientID]float64)
	allocWeights := make(map[shared.ClientID]float64)
	finalAllocations := make(map[shared.ClientID]shared.Resources)

	for island, req := range resourceRequest {
		sumRequest += float64(req)
		resources[island] = shared.Resources(float64(p.c.declaredResources[island]) * math.Pow(resourceSkew, ((100-p.c.trustScore[island])/100)))
	}

	p.c.clientPrint("Resource requests: %+v\n", resourceRequest)
	p.c.clientPrint("Their resource estimation: %+v\n", p.c.declaredResources)
	p.c.clientPrint("Our estimation: %+v\n", resources)
	p.c.clientPrint("Trust Scores: %+v\n", p.c.trustScore)

	avgRequest = findAvgNoTails(resourceRequest)
	avgResource = findAvgNoTails(resources)

	for island, resource := range resources {
		allocations[island] = float64(avgRequest) + p.c.params.equity*(float64(avgResource-resource)+float64(resourceRequest[island]-avgRequest))
		// p.c.clientPrint("Allocation for island %v: %f", island, allocations[island])
		if island == p.c.GetID() {
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

	commonPoolThreshold = math.Min(float64(availCommonPool)*(1.0-p.c.params.riskFactor), sumRequest)
	if p.c.params.saveCriticalIsland {
		for island := range resourceRequest {
			if resources[island] < p.c.criticalThreshold {
				finalAllocations[island] = shared.Resources(math.Max((allocWeights[island] * commonPoolThreshold), float64(p.c.criticalThreshold-resources[island])))
			} else {
				finalAllocations[island] = 0
			}
		}
	}

	for island := range resourceRequest {
		if finalAllocations[island] == 0 {
			if sumRequest < commonPoolThreshold {
				finalAllocations[island] = shared.Resources(math.Max(allocWeights[island]*float64(sumRequest), 0))
			} else {
				finalAllocations[island] = shared.Resources(math.Max(allocWeights[island]*commonPoolThreshold, 0))
			}
		}
	}

	p.c.clientPrint("Final allocations: %+v\n", finalAllocations)

	// Curently always evaluate, would there be a time when we don't want to?
	return shared.PresidentReturnContent{
		ContentType: shared.PresidentAllocation,
		ResourceMap: finalAllocations,
		ActionTaken: true,
	}
}

func (p *president) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {

	if p.c.params.adv != nil {
		val, done := p.c.params.adv.SetTaxationAmount(islandsResources)
		if done {
			return val
		}
	}

	//decide if we want to run SetTaxationAmount
	p.c.declaredResources = make(map[shared.ClientID]shared.Resources)
	for island, report := range islandsResources {
		if report.Reported {
			p.c.declaredResources[island] = report.ReportedAmount
		} else {
			p.c.declaredResources[island] = shared.Resources(p.c.ServerReadHandle.GetGameConfig().IIGOClientConfig.AssumedResourcesNoReport)
		}
	}
	gameState := p.c.BaseClient.ServerReadHandle.GetGameState()
	// Aim to have 100 in common pool after iigo run
	resourcesRequired := math.Max((float64(p.c.getIIGOCost())+100.0)-float64(gameState.CommonPool), 0)
	p.c.clientPrint("Resources required in common pool %f", resourcesRequired)

	if len(p.c.globalDisasterPredictions) > int(p.c.ServerReadHandle.GetGameState().Turn) {
		disaster := p.c.globalDisasterPredictions[int(p.c.ServerReadHandle.GetGameState().Turn)]
		resourcesRequired += disaster.Magnitude - safeDivFloat(float64(gameState.CommonPool), float64(disaster.TimeLeft+1))
	}

	length := math.Max(float64(len(p.c.declaredResources)), 1.0)
	AveTax := resourcesRequired / length
	var adjustedResources []float64
	adjustedResourcesMap := make(map[shared.ClientID]shared.Resources)
	for island, resource := range p.c.declaredResources {
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
		taxation = shared.Resources(math.Max(math.Round(float64(taxation)), 0.0))
		taxationMap[island] = taxation
	}
	p.c.clientPrint("tax amounts : %v\n", taxationMap)
	return shared.PresidentReturnContent{
		ContentType: shared.PresidentTaxation,
		ResourceMap: taxationMap,
		ActionTaken: true,
	}
}

// PickRuleToVote chooses a rule proposal from all the proposals
func (p *president) PickRuleToVote(rulesProposals []rules.RuleMatrix) shared.PresidentReturnContent {

	// Convert to map
	proposals := rmListToMap(rulesProposals)

	compliantProposal := p.c.generalRuleSelection(proposals)

	ret := shared.PresidentReturnContent{
		ContentType:        shared.PresidentRuleProposal,
		ProposedRuleMatrix: compliantProposal,
		ActionTaken:        true,
	}

	/*if p.c.params.complianceLevel > 0.5 {
		return ret
	}

	ret.ProposedRuleMatrix = p.c.RuleProposal()*/

	return ret
}

func rmListToMap(lst []rules.RuleMatrix) map[string]rules.RuleMatrix {
	ret := make(map[string]rules.RuleMatrix)
	for _, rule := range lst {
		ret[rule.RuleName] = rule
	}
	return ret
}

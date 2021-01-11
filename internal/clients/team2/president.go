package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"

	"sort"
)

type President struct {
	*baseclient.BasePresident
	c                 *client
	reportedResources map[shared.ClientID]shared.ResourcesReport
}

type IslandResources struct {
	Resources shared.Resources
	ID        shared.ClientID
}

type IslandResourceList []IslandResources

// Overwrite default sort implementation
func (x IslandResourceList) Len() int           { return len(x) }
func (x IslandResourceList) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x IslandResourceList) Less(i, j int) bool { return x[i].Resources < x[j].Resources }

func (p *President) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) shared.PresidentReturnContent {
	var modeMult, requestSum shared.Resources
	resourceAllocation := make(map[shared.ClientID]shared.Resources)

	for _, request := range resourceRequest {
		requestSum += request
	}

	// Scale resources that could be allocated from Common Pool according to AgentStrategy
	switch p.c.getAgentStrategy() {
	case Selfish:
		// pool is scarce - be cautious
		modeMult = 0.4
	case FairSharer:
		// default
		modeMult = 0.5
	case Altruist:
		// pool has a surplus - give more
		modeMult = 0.6
	}

	resourcesToGive := modeMult * availCommonPool

	// If there are sufficient resources, fulfill every request
	if requestSum <= resourcesToGive || requestSum == 0 {
		resourceAllocation = resourceRequest
	} else {
		// In this case we only give out resources to critical islands
		// order islands in terms of most critical to least critical
		rankedIslands := rankIslands(p, resourceRequest)

		ResourcesPerCritIsland := resourcesToGive
		count := 0

		for _, islandID := range rankedIslands {
			if p.c.gameState().ClientLifeStatuses[islandID] == shared.Critical {
				count++
			}
		}

		if count > 0 {
			ResourcesPerCritIsland = resourcesToGive / shared.Resources(count)
		}

		for _, islandID := range rankedIslands {
			if p.c.gameState().ClientLifeStatuses[islandID] == shared.Critical {
				if resourceRequest[islandID] < ResourcesPerCritIsland {
					resourceAllocation[islandID] = shared.Resources(resourceRequest[islandID])
				} else {
					resourceAllocation[islandID] = shared.Resources(ResourcesPerCritIsland)
				}
			} else {
				resourceAllocation[islandID] = shared.Resources(0)
			}
		}
	}

	return shared.PresidentReturnContent{
		ContentType: shared.PresidentAllocation,
		ResourceMap: resourceAllocation,
		ActionTaken: true,
	}
}

// rankIslands returns island IDs ranked from lowest to highest island resources
func rankIslands(p *President, resourceRequest map[shared.ClientID]shared.Resources) []shared.ClientID {
	sortedList := make(IslandResourceList, 0)
	notReportedList := make(IslandResourceList, 0)

	for id, reportedResource := range p.reportedResources {
		if reportedResource.Reported {
			sortedList = append(sortedList, IslandResources{Resources: reportedResource.ReportedAmount, ID: id})
		} else {
			notReportedList = append(notReportedList, IslandResources{Resources: reportedResource.ReportedAmount, ID: id})
		}
	}

	sort.Sort(sortedList)

	// Append those islands without reported resources to the end
	sortedList = append(sortedList, notReportedList...)

	// Return ranked list of IDs
	rankedList := make([]shared.ClientID, 0)
	for _, islandResources := range sortedList {
		rankedList = append(rankedList, islandResources.ID)
	}

	return rankedList
}

// TODO: We should always pick the rule WE want changed if we are the judge (but not sure how to do this)
//func (p *President) PickRuleToVote(rulesProposals []rules.RuleMatrix) shared.PresidentReturnContent {
//}
// ***************STRATEGY******************
// if islands declare resources we compute the average declared resources
// we can add a penalty (additional tax) to islands that do not declare tax
// and use the declared resources to compute the avg tax as a base
// so more like a fine - the fine can be whatever we want
// if we wanted to it could be 20% of the richest declared resources
// very little computation in this approach by comparison to the current
//*********************************************
func (p *President) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {

	p.reportedResources = islandsResources
	totalResourcesReported := shared.Resources(0)
	totalIslandsReported := shared.Resources(0)
	avgRepResources := shared.Resources(0)
	avgResourcesReq := shared.Resources(0)

	for _, report := range islandsResources {
		if report.Reported {
			totalResourcesReported += report.ReportedAmount
			totalIslandsReported++
		}
	}

	if totalIslandsReported != 0 {
		avgRepResources = totalResourcesReported / totalIslandsReported
	}

	// Check resourcesRequired to mitigate a disaster
	resourcesRequired := shared.Resources(0.0)
	disaster := p.c.MakeDisasterPrediction().PredictionMade

	if float64(disaster.TimeLeft) != 0.0 {
		resourcesRequired = shared.Resources(disaster.Magnitude - float64(p.c.gameState().CommonPool)/float64(disaster.TimeLeft))
	} else {
		resourcesRequired = shared.Resources(disaster.Magnitude - float64(p.c.gameState().CommonPool))
	}

	taxationMap := make(map[shared.ClientID]shared.Resources)

	if p.c.getNumAliveClients() != 0 {
		avgResourcesReq = resourcesRequired / shared.Resources(p.c.getNumAliveClients())
	}

	// assign average tax
	for island, report := range p.reportedResources {
		if report.Reported {
			taxationMap[island] = Max(0.2*report.ReportedAmount, avgResourcesReq)
		} else {
			taxationMap[island] = Max(avgRepResources, avgResourcesReq) * 1.2
		}

		// the best islands don't pay tax - tax evasion 101
		if island == p.c.BaseClient.GetID() {
			taxationMap[island] = 0
		}
	}

	return shared.PresidentReturnContent{
		ContentType: shared.PresidentTaxation,
		ResourceMap: taxationMap,
		ActionTaken: true,
	}
}

// example implementation calls an election if monitoring was performed and the result was negative
func (p *President) CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.InstantRunoff,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}

	if monitoring.Performed && !monitoring.Result {
		electionsettings.HoldElection = true
	}

	if turnsInPower >= 2 {
		//keep default because speaker is not too relevant
		electionsettings.HoldElection = true
	}

	return electionsettings
}

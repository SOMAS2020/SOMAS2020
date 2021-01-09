package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"

	"math"
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
	switch p.c.MethodOfPlay() {
	case 0:
		modeMult = 0.6 //pool is struggling, give less
	case 1:
		modeMult = 0.7 //default
	case 2:
		modeMult = 0.8 //pool has surplus, give a bit more than usual
	}

	if requestSum < modeMult*availCommonPool || requestSum == 0 { //enough resources to go round
		resourceAllocation = resourceRequest
	} else { // In this case we only give out resources to critical islands
		//order islands in terms of most critical to least critical
		resourcesToGive := modeMult * availCommonPool
		rankedIslands := rankIslands(p, resourceRequest)
		for _, islandID := range rankedIslands {
			if resourcesToGive-resourceRequest[islandID] > 0 {
				resourceAllocation[islandID] = shared.Resources(resourceRequest[islandID])
				resourcesToGive -= resourceRequest[islandID]
				continue
			}
			if resourcesToGive > 0 {
				resourceAllocation[islandID] = shared.Resources(resourcesToGive)
				resourcesToGive = 0
				break
			}
			break
		}
	}

	return shared.PresidentReturnContent{
		ContentType: shared.PresidentAllocation,
		ResourceMap: resourceAllocation,
		ActionTaken: true,
	}
}

// rankIslands returns a slice of island IDs from the lowest resource island (index 0) to highest
// resource island
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

	// Sort this list
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

// TODO by Eirik/Hardik
// We should always pick the rule WE want changed if we are the judge (but not sure how to do this)
//func (p *President) PickRuleToVote(rulesProposals []rules.RuleMatrix) shared.PresidentReturnContent {
//}

//Used Team3's implementation
func (p *President) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {
	p.reportedResources = islandsResources
	p.c.declaredResources = make(map[shared.ClientID]shared.Resources)
	for island, report := range islandsResources {
		if report.Reported {
			p.c.declaredResources[island] = report.ReportedAmount
		} else {
			//TODO: read from params.config file once its available now i just hope they die of satanic attacc <-Team3 comment
			p.c.declaredResources[island] = shared.Resources(666)
		}
	}
	// TODO: why are we not using our helper functions here?
	gameState := p.c.BaseClient.ServerReadHandle.GetGameState()
	resourcesRequired := 100.0 - float64(gameState.CommonPool)
	disaster := p.c.MakeDisasterPrediction().PredictionMade
	resourcesRequired = (disaster.Magnitude - float64(p.c.gameState().CommonPool)/float64(disaster.TimeLeft))

	AveTax := resourcesRequired / float64(len(p.c.declaredResources))

	var declaredResourcesFloat []float64
	declaredResourcesMapFloat := make(map[shared.ClientID]shared.Resources)
	for island, resource := range p.c.declaredResources {
		declaredResourceFloat := resource
		declaredResourcesFloat = append(declaredResourcesFloat, float64(declaredResourceFloat))
		declaredResourcesMapFloat[island] = declaredResourceFloat
	}

	AveDeclaredResources := getAverage(declaredResourcesFloat)

	taxationMap := make(map[shared.ClientID]shared.Resources)
	for island, resources := range declaredResourcesMapFloat {
		taxation := shared.Resources(AveTax) + (resources - shared.Resources(AveDeclaredResources)) //*shared.Resources(p.c.params.equity)
		if island == p.c.BaseClient.GetID() {
			taxation -= taxation //*shared.Resources(p.c.params.selfishness)
		}
		taxation = shared.Resources(math.Max(float64(taxation), 0.0))
		taxationMap[island] = taxation
	}
	return shared.PresidentReturnContent{
		ContentType: shared.PresidentTaxation,
		ResourceMap: taxationMap,
		ActionTaken: true,
	}
}

func (p *President) CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	// example implementation calls an election if monitoring was performed and the result was negative
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.InstantRunoff,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}
	if monitoring.Performed && !monitoring.Result {
		electionsettings.HoldElection = true
	}
	if turnsInPower >= 2 { //keep default because speaker is not too relevant
		electionsettings.HoldElection = true
	}
	return electionsettings
}

// getAverage returns the average of the list
func getAverage(lst []float64) float64 {
	if len(lst) == 0 {
		return 0.0
	}

	total := 0.0
	for _, val := range lst {
		total += val
	}

	return (float64(total) / float64(len(lst)))
}

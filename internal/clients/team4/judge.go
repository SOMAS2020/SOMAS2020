package team4

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func init() {
	os.Remove("./Historymap.txt")
}

type judge struct {
	*baseclient.BaseJudge
	parent *client
	t      *testing.T
}

// GetRuleViolationSeverity returns a custom map of named rules and how severe the sanction should be for transgressing them
// If a rule is not named here, the default sanction value added is 1
// OPTIONAL: override to set custom sanction severities for specific rules

// GetSanctionThresholds returns a custom map of sanction score thresholds for different sanction tiers
// For any unfilled sanction tiers will be filled with default values (given in judiciary.go)
// OPTIONAL: override to set custom sanction thresholds

// PayPresident pays the President a salary.
// OPTIONAL: override to pay the President less than the full amount.

// InspectHistory is the base implementation of evaluating islands choices the last turn.
// OPTIONAL: override if you want to evaluate the history log differently.
type valuePair struct {
	actual   shared.Resources
	expected shared.Resources
}

type judgeHistoryInfo struct {
	Resources  valuePair // amount of resources available and reported by the island
	Taxation   valuePair // amount of tax paid and expected
	Allocation valuePair // amount of allocation taken and granted
	Lied       int       // number of times the island has lied
}

func (j *judge) saveHistoryInfo(iigoHistory *[]shared.Accountability, lieCounts *map[shared.ClientID]int, turn uint) {
	accountabilityMap := map[shared.ClientID][]rules.VariableValuePair{}
	for _, clientID := range shared.TeamIDs {
		accountabilityMap[clientID] = []rules.VariableValuePair{}
	}

	for _, acc := range *iigoHistory {
		client := acc.ClientID
		accountabilityMap[client] = append(accountabilityMap[client], acc.Pairs...)
	}

	for client, pairs := range accountabilityMap {
		clientInfo, ok := buildHistoryInfo(pairs)
		if ok {
			clientLied := (*lieCounts)[client]
			clientInfo.Lied = clientLied
			savedHistory := *j.parent.savedHistory
			if savedHistory[turn] != nil {
				savedHistory[turn][client] = clientInfo
			} else {
				savedHistory[turn] = map[shared.ClientID]judgeHistoryInfo{client: clientInfo}
			}
			j.parent.savedHistory = &savedHistory
		}
	}
}

func (j *judge) InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]shared.EvaluationReturn, bool) {

	outputmap, state := j.BaseJudge.InspectHistory(iigoHistory, turnsAgo)

	turn := j.parent.getTurn() - uint(turnsAgo)
	lieCounts := map[shared.ClientID]int{}

	// Check number of times clients lied
	for client, eval := range outputmap {
		lieCount := 0
		for _, rulePasses := range eval.Evaluations {
			if !rulePasses {
				lieCount++
			}
		}
		lieCounts[client] = lieCount
	}

	j.saveHistoryInfo(&iigoHistory, &lieCounts, turn)

	return outputmap, state
}

// GetPardonedIslands decides which islands to pardon i.e. no longer impose sanctions on
// COMPULSORY: decide which islands, if any, to forgive

// HistoricalRetributionEnabled allows you to punish more than the previous turns transgressions
// OPTIONAL: override if you want to punish historical transgressions
func (j *judge) HistoricalRetributionEnabled() bool {
	return false
}

// CallPresidentElection is called by the judiciary to decide on power-transfer
// COMPULSORY: decide when to call an election following relevant rulesInPlay if you wish
func (j *judge) CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	// example implementation calls an election if monitoring was performed and the result was negative
	// or if the number of turnsInPower exceeds 3
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.Runoff,
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

// DecideNextPresident returns the ID of chosen next President
// OPTIONAL: override to manipulate the result of the election
func (j *judge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	// overloaded
	j.logf("hello world %v", winner)
	return id
}

func (j *judge) logf(format string, a ...interface{}) {
	if j.t != nil {
		j.t.Log(fmt.Sprintf(format, a...))
	}
}

func dump(filename string, format string, v ...interface{}) {
	//f, err := os.Create(filename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(fmt.Sprintf(format, v...))

	if err2 != nil {
		log.Fatal(err2)
	}

}

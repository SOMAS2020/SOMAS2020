package team4

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
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
// type judgeHistoryInfo struct {
// 	rules.VariableFieldName
// }

func (j *judge) InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]roles.EvaluationReturn, bool) {
	//outputMap := map[shared.ClientID]roles.EvaluationReturn{}

	// current turn
	// currentTurn := j.parent.getTurn()

	// store := map[int]map[shared.ClientID]map[rules.Variable]{}

	// j.logf(store)

	// Check who is lying about private resources
	// How many resources an island has

	// ExpectedTaxContribution vs IslandTaxContribution

	// Check which rules are in place
	outputmap, state := j.BaseJudge.InspectHistory(iigoHistory, turnsAgo)
	dump("./Historymap.txt", "iigoHistory: %v \n", outputmap)

	// Choose to lie about state
	// Track if any rules wore broken by looking at the output map

	// t := outputmap[shared.Team1]
	//dump("./Outputmap.txt", "Outputmap: %v \n", outputmap)

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

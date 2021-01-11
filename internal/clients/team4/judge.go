package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type judge struct {
	*baseclient.BaseJudge
	parent *client
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
	Resources   valuePair // amount of resources available and reported by the island
	Taxation    valuePair // amount of tax paid and expected
	Allocation  valuePair // amount of allocation taken and granted
	LawfulRatio float64   // ratio of truth/all in each island report
}

type accountabilityHistory struct {
	history     map[uint]map[shared.ClientID]judgeHistoryInfo // stores accountablity history of all turns
	updated     bool                                          // indicates whether a judge has updated the history
	updatedTurn uint                                          // indicates a turn which was updated
}

func (h *accountabilityHistory) getNewInfo() map[shared.ClientID]judgeHistoryInfo {
	if newInfo, ok := h.history[h.updatedTurn]; ok {
		return newInfo
	}
	return map[shared.ClientID]judgeHistoryInfo{}
}

func (h *accountabilityHistory) update(turn uint) {
	h.updated = true
	h.updatedTurn = turn
}

func (j *judge) saveHistoryInfo(iigoHistory *[]shared.Accountability, truthfulness *map[shared.ClientID]float64, turn uint) {
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
			truthfulRatio := (*truthfulness)[client]
			clientInfo.LawfulRatio = truthfulRatio
			if j.parent.savedHistory.history[turn] != nil {
				j.parent.savedHistory.history[turn][client] = clientInfo
			} else {
				j.parent.savedHistory.history[turn] = map[shared.ClientID]judgeHistoryInfo{client: clientInfo}
			}
			j.parent.savedHistory.update(turn)
		}
	}
}

func (j *judge) InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]shared.EvaluationReturn, bool) {
	outputmap, state := j.BaseJudge.InspectHistory(iigoHistory, turnsAgo)

	turn := j.parent.getTurn() - uint(turnsAgo)
	truthfulness := map[shared.ClientID]float64{}

	// Check number of times clients lied
	for client, eval := range outputmap {
		lieCount := 0
		for _, rulePasses := range eval.Evaluations {
			if !rulePasses {
				lieCount++
			}
		}

		if len(eval.Evaluations) == 0 {
			truthfulness[client] = 1.0
		} else {
			truthfulness[client] = 1.0 - float64(lieCount)/float64(len(eval.Evaluations))
		}
		// lieCounts[client] = float64(lieCount) / float64(len(eval.Evaluations))
	}

	j.saveHistoryInfo(&iigoHistory, &truthfulness, turn)

	return outputmap, state
}

// GetPardonedIslands decides which islands to pardon i.e. no longer impose sanctions on
// COMPULSORY: decide which islands, if any, to forgive
func (j *judge) GetPardonedIslands(currentSanctions map[int][]shared.Sanction) map[int][]bool {
	pardons := map[int][]bool{}

	for turn, sanctions := range currentSanctions {
		pardons[turn] = make([]bool, len(sanctions))
		for index, sanction := range sanctions {
			if sanction.SanctionTier != shared.NoSanction {
				considerTime := sanction.TurnsLeft <= j.parent.internalParam.maxPardonTime
				considerSeverity := sanction.SanctionTier <= j.parent.internalParam.maxTierToPardon
				considerTrust := j.parent.internalParam.minTrustToPardon <= j.parent.getTrust(sanction.ClientID)
				if considerTime && considerSeverity && considerTrust {
					pardons[turn][index] = true
				}
			}
		}
	}

	return pardons
}

// HistoricalRetributionEnabled allows you to punish more than the previous turns transgressions
// OPTIONAL: override if you want to punish historical transgressions
func (j *judge) HistoricalRetributionEnabled() bool {
	return false
}

// CallPresidentElection is called by the judiciary to decide on power-transfer
// COMPULSORY: decide when to call an election following relevant rulesInPlay if you wish
func (j *judge) CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {

	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.Runoff,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}

	// calculate whether the term has ended
	termEnded := uint(turnsInPower) > j.parent.getTermLength(shared.President)

	j.parent.LocalVariableCache[rules.TermEnded] = rules.VariableValuePair{
		VariableName: rules.TermEnded,
		Values:       []float64{boolToFloat(termEnded)},
	}

	// try not holding election
	j.parent.LocalVariableCache[rules.ElectionHeld] = rules.VariableValuePair{
		VariableName: rules.ElectionHeld,
		Values:       []float64{boolToFloat(false)},
	}

	compliantWithRules := j.parent.CheckCompliance(rules.ElectionHeld)

	if monitoring.Performed && !monitoring.Result {
		// If president didn't perform it's duties we can try to elect new president. Potentially decide based on trust or something
		electionsettings.HoldElection = true
	} else if compliantWithRules {
		// If we are compliant with rules ie. we don't have to (are not obliged) hold election
		electionsettings.HoldElection = false
	} else {
		// Now we know that we have to hold elections. We can check whether we want to hold election based on trust or something
		electionsettings.HoldElection = true
	}

	return electionsettings
}

// DecideNextPresident returns the ID of chosen next President
// OPTIONAL: override to manipulate the result of the election

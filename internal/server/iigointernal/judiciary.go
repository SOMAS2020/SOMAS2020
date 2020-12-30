package iigointernal

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
)

// to be moved to paramters
const sanctionCacheDepth = 3

// to be changed
const sanctionLength = 2

type judiciary struct {
	JudgeID               shared.ClientID
	budget                shared.Resources
	presidentSalary       shared.Resources
	EvaluationResults     map[shared.ClientID]roles.EvaluationReturn
	clientJudge           roles.Judge
	presidentTurnsInPower int
	sanctionRecord        map[shared.ClientID]roles.IIGOSanctionScore
	sanctionThresholds    map[roles.IIGOSanctionTier]roles.IIGOSanctionScore
	ruleViolationSeverity map[string]roles.IIGOSanctionScore
	localSanctionCache    map[int][]roles.Sanction
}

// Loads ruleViolationSeverity and sanction thresholds
func (j *judiciary) loadSanctionConfig() {
	j.sanctionThresholds = softMergeSanctionThresholds(j.clientJudge.GetSanctionThresholds())
	j.ruleViolationSeverity = j.clientJudge.GetRuleViolationSeverity()
}

// loadClientJudge checks client pointer is good and if not panics
func (j *judiciary) loadClientJudge(clientJudgePointer roles.Judge) {
	if clientJudgePointer == nil {
		panic(fmt.Sprintf("Client '%v' has loaded a nil judge pointer", j.JudgeID))
	}
	j.clientJudge = clientJudgePointer
}

// returnPresidentSalary returns the salary to the common pool.
func (j *judiciary) returnPresidentSalary() shared.Resources {
	x := j.presidentSalary
	j.presidentSalary = 0
	return x
}

// withdrawPresidentSalary withdraws the president's salary from the common pool.
func (j *judiciary) withdrawPresidentSalary(gameState *gamestate.GameState) bool {
	var presidentSalary = shared.Resources(rules.VariableMap[rules.PresidentSalary].Values[0])
	var withdrawAmount, withdrawSuccesful = WithdrawFromCommonPool(presidentSalary, gameState)
	j.presidentSalary = withdrawAmount
	return withdrawSuccesful
}

// sendPresidentSalary sends the president's salary to the president.
func (j *judiciary) sendPresidentSalary(executiveBranch *executive) {
	if j.clientJudge != nil {
		amount, payPresident := j.clientJudge.PayPresident(j.presidentSalary)
		if payPresident {
			executiveBranch.budget = amount
		}
		return
	}
	amount := j.PayPresident()
	executiveBranch.budget = amount
}

// PayPresident pays the president salary.
func (j *judiciary) PayPresident() shared.Resources {
	hold := j.presidentSalary
	j.presidentSalary = 0
	return hold
}

// InspectHistory checks all actions that happened in the last turn and audits them.
// This can be overridden by clients.
func (j *judiciary) inspectHistory(iigoHistory []shared.Accountability) (map[shared.ClientID]roles.EvaluationReturn, bool) {
	j.budget -= serviceCharge
	return j.clientJudge.InspectHistory(iigoHistory)
}

// searchForRule searches for a given rule in the RuleMatrix
func searchForRule(ruleName string, listOfRuleMatrices []rules.RuleMatrix) (int, bool) {
	for i, v := range listOfRuleMatrices {
		if v.RuleName == ruleName {
			return i, true
		}
	}
	return -1, false
}

// appointNextPresident returns the island ID of the island appointed to be President in the next turn
func (j *judiciary) appointNextPresident(currentPresident shared.ClientID, allIslands []shared.ClientID) shared.ClientID {
	var election voting.Election
	var nextPresident shared.ClientID
	electionsettings := j.clientJudge.CallPresidentElection(j.presidentTurnsInPower, allIslands)
	if electionsettings.HoldElection {
		// TODO: deduct the cost of holding an election
		election.ProposeElection(baseclient.President, electionsettings.VotingMethod)
		election.OpenBallot(electionsettings.IslandsToVote)
		election.Vote(iigoClients)
		j.presidentTurnsInPower = 0
		nextPresident = election.CloseBallot()
		nextPresident = j.clientJudge.DecideNextPresident(nextPresident)
	} else {
		j.presidentTurnsInPower++
		nextPresident = currentPresident
	}
	return nextPresident
}

// cycleSanctionCache rolls the sanction cahce one turn forward (effectively dropping any sanctions longer than the depth)
func (j *judiciary) cycleSanctionCache() {
	newMap := j.localSanctionCache
	delete(newMap, sanctionCacheDepth-1)
	newMapCache := newMap
	for i := 0; i < sanctionCacheDepth-1; i++ {
		newMap[i+1] = newMapCache[i]
	}
	newMap[0] = []roles.Sanction{}
	j.localSanctionCache = newMap
}

// Helper functions //

// getDefaultSanctionThresholds provides default thresholds for sanctions
func getDefaultSanctionThresholds() map[roles.IIGOSanctionTier]roles.IIGOSanctionScore {
	return map[roles.IIGOSanctionTier]roles.IIGOSanctionScore{
		roles.SanctionTier1: 1,
		roles.SanctionTier2: 5,
		roles.SanctionTier3: 10,
		roles.SanctionTier4: 20,
		roles.SanctionTier5: 30,
	}
}

// softMergeSanctionThresholds merges the default sanction thresholds with a (preferred) client version
func softMergeSanctionThresholds(clientSanctionMap map[roles.IIGOSanctionTier]roles.IIGOSanctionScore) map[roles.IIGOSanctionTier]roles.IIGOSanctionScore {
	defaultMap := getDefaultSanctionThresholds()
	for k := range defaultMap {
		if clientVal, ok := clientSanctionMap[k]; ok {
			defaultMap[k] = clientVal
		}
	}
	return defaultMap
}

// getIslandSanctionTier if statement based evaluator for which sanction tier a particular score is in
func getIslandSanctionTier(islandScore roles.IIGOSanctionScore, scoreMap map[roles.IIGOSanctionTier]roles.IIGOSanctionScore) roles.IIGOSanctionTier {
	if islandScore <= scoreMap[roles.SanctionTier1] {
		return roles.NoSanction
	} else if islandScore <= scoreMap[roles.SanctionTier2] {
		return roles.SanctionTier1
	} else if islandScore <= scoreMap[roles.SanctionTier3] {
		return roles.SanctionTier2
	} else if islandScore <= scoreMap[roles.SanctionTier4] {
		return roles.SanctionTier3
	} else if islandScore <= scoreMap[roles.SanctionTier5] {
		return roles.SanctionTier4
	} else {
		return roles.SanctionTier5
	}
}

// getTierSanctionMap basic mapping between snaciotn tier and rule that governs it
func getTierSanctionMap() map[roles.IIGOSanctionTier]string {
	return map[roles.IIGOSanctionTier]string{
		roles.SanctionTier1: "iigo_economic_sanction_1",
		roles.SanctionTier2: "iigo_economic_sanction_2",
		roles.SanctionTier3: "iigo_economic_sanction_3",
		roles.SanctionTier4: "iigo_economic_sanction_4",
		roles.SanctionTier5: "iigo_economic_sanction_5",
	}
}

// defaultInitLocalSanctionCache generates a blank sanction cache
func defaultInitLocalSanctionCache(depth int) map[int][]roles.Sanction {
	returnMap := map[int][]roles.Sanction{}
	for i := 0; i < depth; i++ {
		returnMap[i] = []roles.Sanction{}
	}
	return returnMap
}

// checkPardons checks the pardons issued by an island judge making sure they are valid, returns the remaining sanctions
func checkPardons(sanctionCache map[int][]roles.Sanction, pardons map[int]map[int]roles.Sanction) (pardonsValid bool, communications map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent, finalCache map[int][]roles.Sanction) {
	comms := map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent{}
	newSanctionCache := map[int][]roles.Sanction{}
	for k, v := range sanctionCache {
		newSanctionCache[k] = v
	}
	for i, v := range pardons {
		for iSan, vSan := range v {
			if sanctionCache[i][iSan] != vSan {
				return false, map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent{}, sanctionCache
			} else {
				comms[vSan.ClientID] = append(comms[vSan.ClientID], map[shared.CommunicationFieldName]shared.CommunicationContent{
					shared.PardonClientID: {
						T:           shared.CommunicationInt,
						IntegerData: int(vSan.ClientID),
					},
					shared.PardonTier: {
						T:           shared.CommunicationInt,
						IntegerData: int(vSan.SanctionTier),
					},
				})
				copyOfNewSanctionCache := make([]roles.Sanction, len(newSanctionCache[i]))
				copy(copyOfNewSanctionCache, newSanctionCache[i])
				newSanctionCache[i] = removeSanctions(copyOfNewSanctionCache, iSan-getDifferenceInLength(sanctionCache[i], copyOfNewSanctionCache))
			}
		}
	}
	return true, comms, newSanctionCache
}

// removeSanctions is a helper function to remoce a sanction element from a slice
func removeSanctions(slice []roles.Sanction, s int) []roles.Sanction {
	var output []roles.Sanction
	for i, v := range slice {
		if i != s {
			output = append(output, v)
		}
	}
	return output
}

// getDifferenceInLength helper function to get difference in length between two lists
func getDifferenceInLength(slice1 []roles.Sanction, slice2 []roles.Sanction) int {
	return len(slice1) - len(slice2)
}

package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) getTurn() uint {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().Turn
	}
	return 0
}

func (c *client) getMinimumThreshold() shared.Resources {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold
	}
	return 0
}

func (c *client) getCostOfLiving() shared.Resources {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameConfig().CostOfLiving
	}
	return 0
}

func (c *client) getSeason() uint {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().Season
	}
	return 0
}

func (c *client) getCommonPool() shared.Resources {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().CommonPool
	}
	return 0
}

func (c *client) getTermLength(role shared.Role) uint {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameConfig().IIGOClientConfig.IIGOTermLengths[role]
	}
	return 0
}

func (c *client) getResources() shared.Resources {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().ClientInfo.Resources
	}
	return 0
}

func (c *client) getLifeStatus() shared.ClientLifeStatus {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus
	}
	return 0
}

func (c *client) getAllLifeStatus() map[shared.ClientID]shared.ClientLifeStatus {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().ClientLifeStatuses
	}
	return make(map[shared.ClientID]shared.ClientLifeStatus)
}

func (c *client) getSafeResourceLevel() shared.Resources {
	if c.ServerReadHandle != nil {
		conf := c.ServerReadHandle.GetGameConfig()
		return conf.MinimumResourceThreshold + conf.CostOfLiving
	}
	return 0
}

func (c *client) getTrust(clientID shared.ClientID) float64 {
	if c.GetID() == clientID {
		return 0.4 + (0.6 * c.internalParam.selfishness)
	}
	return c.trustMatrix.GetClientTrust(clientID)
}

func (c *client) getTurnsInPower(role shared.Role) uint {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().IIGOTurnsInPower[role]
	}
	return 0
}

func (c *client) getRoleBudget(role shared.Role) shared.Resources {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().IIGORolesBudget[role]
	}
	return 0
}

func (c *client) getIIGOConfig() config.IIGOConfig {
	return c.ServerReadHandle.GetGameConfig().IIGOClientConfig
}

func buildHistoryInfo(pairs []rules.VariableValuePair) (retInfo judgeHistoryInfo, ok bool) {
	resourceOK := 0
	taxOK := 0
	allocationOK := 0
	for _, val := range pairs {
		switch val.VariableName {
		case rules.IslandActualPrivateResources:
			if len(val.Values) > 0 {
				retInfo.Resources.expected = shared.Resources(val.Values[0])
				resourceOK++
			}
		case rules.IslandReportedPrivateResources:
			if len(val.Values) > 0 {
				retInfo.Resources.actual = shared.Resources(val.Values[0])
				resourceOK++
			}
		case rules.ExpectedTaxContribution:
			if len(val.Values) > 0 {
				retInfo.Taxation.expected = shared.Resources(val.Values[0])
				taxOK++
			}
		case rules.IslandTaxContribution:
			if len(val.Values) > 0 {
				retInfo.Taxation.actual = shared.Resources(val.Values[0])
				taxOK++
			}
		case rules.ExpectedAllocation:
			if len(val.Values) > 0 {
				retInfo.Allocation.expected = shared.Resources(val.Values[0])
				allocationOK++
			}
		case rules.IslandAllocation:
			if len(val.Values) > 0 {
				retInfo.Allocation.actual = shared.Resources(val.Values[0])
				allocationOK++
			}
		default:
			//[exhaustive] reported by reviewdog 🐶
			//missing cases in switch of type rules.VariableFieldName: AllocationMade, AllocationRequestsMade, AnnouncementResultMatchesVote, AnnouncementRuleMatchesVote, AppointmentMatchesVote, ConstSanctionAmount, ElectionHeld, HasIslandReportPrivateResources, IslandReportedResources, IslandsAlive, IslandsAllowedToVote, IslandsProposedRules, JudgeBudgetIncrement, JudgeHistoricalRetributionPerformed, JudgeInspectionPerformed, JudgeLeftoverBudget, JudgePaid, JudgePayment, JudgeSalary, MaxSeverityOfSanctions, MonitorRoleAnnounce, MonitorRoleDecideToMonitor, MonitorRoleEvalResult, MonitorRoleEvalResultDecide, NumberOfAllocationsSent, NumberOfBallotsCast, NumberOfBrokenAgreements, NumberOfFailedForages, NumberOfIslandsAlive, NumberOfIslandsContributingToCommonPool, PresidentBudgetIncrement, PresidentLeftoverBudget, PresidentPaid, PresidentPayment, PresidentRuleProposal, PresidentSalary, RuleChosenFromProposalList, RuleSelected, SanctionExpected, SanctionPaid, SpeakerBudgetIncrement, SpeakerLeftoverBudget, SpeakerPaid, SpeakerPayment, SpeakerProposedPresidentRule, SpeakerSalary, TaxDecisionMade, TermEnded, TestVariable, TurnsLeftOnSanction, VoteCalled, VoteResultAnnounced (exhaustive)

		}
	}

	ok = resourceOK == 2 && taxOK == 2 && allocationOK == 2

	return retInfo, ok
}

/*func dump(filename string, format string, v ...interface{}) {
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
}*/

func (c *client) getPresident() shared.ClientID {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().PresidentID
	}
	return 0
}

func (c *client) getSpeaker() shared.ClientID {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().SpeakerID
	}
	return 0
}

func (c *client) getJudge() shared.ClientID {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().JudgeID
	}
	return 0
}

func boolToFloat(input bool) float64 {
	if input {
		return 1
	}
	return 0
}

func checkIfClientIsInList(lst []shared.ClientID, c shared.ClientID) bool {
	for _, e := range lst {
		if e == c {
			return true
		}
	}
	return false
}

func createClientSet(lst []shared.ClientID) []shared.ClientID {
	uniqueMap := make(map[shared.ClientID]bool)
	var uniqueLst []shared.ClientID
	for _, e := range lst {
		if _, ok := uniqueMap[e]; !ok {
			uniqueMap[e] = true
			uniqueLst = append(uniqueLst, e)
		}
	}
	return uniqueLst
}

func (c *client) getOurResources() shared.Resources {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().ClientInfo.Resources
	}
	return 0
}

func (c *client) getRole(role shared.Role) shared.ClientID {
	if c.ServerReadHandle != nil {
		switch role {
		case shared.Judge:
			return c.getJudge()
		case shared.President:
			return c.getPresident()
		case shared.Speaker:
			return c.getSpeaker()
		}
	}
	return 0
}

func (c *client) printConfig() {
	c.Logf("Client resources: %v", c.getResources())
	c.Logf("Internal Config: greediness %v, selfishness %v, colaboration %v, risk-taking %v",
		c.internalParam.greediness, c.internalParam.selfishness, c.internalParam.collaboration, c.internalParam.riskTaking)

	c.Logf("Trust: %v", c.trustMatrix.trustMap)
}

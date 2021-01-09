package team4

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) getTurn() uint {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().Turn
	}
	return 0
}

func (c *client) getSeason() uint {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameState().Season
	}
	return 0
}

func (c *client) getTurnLength(role shared.Role) uint {
	if c.ServerReadHandle != nil {
		return c.ServerReadHandle.GetGameConfig().IIGOClientConfig.IIGOTermLengths[role]
	}
	return 0
}

func (c *client) getTrust(clientID shared.ClientID) float64 {
	if int(clientID) < len(c.internalParam.agentsTrust) {
		return c.internalParam.agentsTrust[int(clientID)]
	}
	return 0
}

func newClient(id shared.ClientID, testing *testing.T) client {

	// have some config json file or something?
	internalConfig := internalParameters{
		greediness:       0,
		selfishness:      0,
		fairness:         0,
		collaboration:    0,
		riskTaking:       0,
		agentsTrust:      []float64{},
		minPardonTime:    3,
		maxTierToPardon:  shared.SanctionTier3,
		minTrustToPardon: 0.6,
	}

	iigoObs := iigoObservation{
		allocationGranted: shared.Resources(0),
		taxDemanded:       shared.Resources(0),
	}
	iifoObs := iifoObservation{}
	iitoObs := iitoObservation{}

	obs := observation{
		iigoObs: &iigoObs,
		iifoObs: &iifoObs,
		iitoObs: &iitoObs,
	}

	judgeHistory := map[uint]map[shared.ClientID]judgeHistoryInfo{}

	emptyRuleCache := map[string]rules.RuleMatrix{}

	newClient := client{
		BaseClient:         baseclient.NewClient(id),
		clientJudge:        judge{BaseJudge: &baseclient.BaseJudge{}, t: testing},
		clientSpeaker:      speaker{BaseSpeaker: &baseclient.BaseSpeaker{}},
		yes:                "",
		obs:                &obs,
		internalParam:      &internalConfig,
		idealRulesCachePtr: &emptyRuleCache,
		savedHistory:       &judgeHistory,
	}

	newClient.clientJudge.parent = &newClient
	newClient.clientSpeaker.parent = &newClient

	return newClient
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

// func dump(filename string, format string, v ...interface{}) {
// 	//f, err := os.Create(filename)
// 	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer f.Close()

// 	_, err2 := f.WriteString(fmt.Sprintf(format, v...))

// 	if err2 != nil {
// 		log.Fatal(err2)
// 	}

// }
func boolToFloat(input bool) float64 {
	if input {
		return 1
	}
	return 0
}

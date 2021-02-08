// Package team4 contains code for team 4's client implementation
package team4

import (
	"math"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

// DefaultClient creates the client that will be used for most simulations. All
// other personalities are considered alternatives. To give a different
// personality for your agent simply create another (exported) function with the
// same signature as "DefaultClient" that creates a different agent, and inform
// someone on the simulation team that you would like it to be included in
// testing
func DefaultClient(id shared.ClientID) baseclient.Client {
	return NewClient(id, honest)
}

// DishonestClient creates a dishonest client
func DishonestClient(id shared.ClientID) baseclient.Client {
	return NewClient(id, dishonest)
}

// ModerateClient creates a moderate client
func ModerateClient(id shared.ClientID) baseclient.Client {
	return NewClient(id, moderate)
}

func newClientInternal(clientID shared.ClientID, clientConfig ClientConfig) client {
	// have some config json file or something?
	internalConfig := configureClient(clientConfig)

	iigoObs := iigoObservation{
		allocationGranted: shared.Resources(0),
		taxDemanded:       shared.Resources(0),
		sanctionTiers:     make(map[shared.ClientID]shared.IIGOSanctionsTier),
	}
	iifoObs := iifoObservation{}
	iitoObs := iitoObservation{}

	obs := observation{
		iigoObs: &iigoObs,
		iifoObs: &iifoObs,
		iitoObs: &iitoObs,
	}

	judgeHistory := accountabilityHistory{
		history:     map[uint]map[shared.ClientID]judgeHistoryInfo{},
		updated:     false,
		updatedTurn: 0,
	}

	emptyRuleCache := map[string]rules.RuleMatrix{}
	trustMatrix := trust{
		trustMap: map[shared.ClientID]float64{},
	}
	trustMatrix.initialise()

	importancesMatrix := importances{
		requestAllocationImportance:                mat.NewVecDense(6, []float64{5.0, 1.0, -1.0, -1.0, 5.0, 1.0}),
		commonPoolResourceRequestImportance:        mat.NewVecDense(6, []float64{4.0, 1.0, -1.0, -1.0, 1.0, 1.0}),
		resourceReportImportance:                   mat.NewVecDense(6, []float64{5.0, 5.0, -5.0, -5.0, 1.0, 5.0}),
		getTaxContributionImportance:               mat.NewVecDense(4, []float64{-2.0, -2.0, 4.0, 1.0}),
		decideIIGOMonitoringAnnouncementImportance: mat.NewVecDense(3, []float64{1.0, -1.0, 1.0}),
		getGiftRequestsImportance:                  mat.NewVecDense(4, []float64{2.0, 1.0, -1.0, -1.0}),
	}

	baseJudge := baseclient.BaseJudge{}
	baseSpeaker := baseclient.BaseSpeaker{}
	basePresident := baseclient.BasePresident{}

	team4client := client{
		BaseClient:  baseclient.NewClient(clientID),
		clientJudge: judge{BaseJudge: &baseJudge},
		clientSpeaker: speaker{
			BaseSpeaker: &baseSpeaker,
			SpeakerActionOrder: []string{
				"SetVotingResult",
				"SetRuleToVote",
				"AnnounceVotingResult",
				"UpdateRules",
				"AppointNextJudge",
			},
			SpeakerActionPriorities: []string{
				"SetVotingResult",
				"SetRuleToVote",
				"AnnounceVotingResult",
				"UpdateRules",
				"AppointNextJudge",
			},
		},
		clientPresident:    president{BasePresident: &basePresident},
		obs:                &obs,
		internalParam:      &internalConfig,
		idealRulesCachePtr: &emptyRuleCache,
		savedHistory:       &judgeHistory,
		trustMatrix:        &trustMatrix,
		importances:        &importancesMatrix,
		forage: &forageStorage{
			forageHistory:      nil,
			receivedForageData: nil,
		},
	}

	team4client.updateParents()

	return team4client
}

// NewClient is a function that creates a new empty client
func NewClient(clientID shared.ClientID, clientConfig ClientConfig) baseclient.Client {
	team4client := newClientInternal(clientID, clientConfig)
	return &team4client
}

type client struct {
	*baseclient.BaseClient //client struct has access to methods and fields of the BaseClient struct which implements implicitly the Client interface.

	//custom fields
	clientJudge        judge
	clientSpeaker      speaker
	clientPresident    president
	obs                *observation        //observation is the raw input into our client
	internalParam      *internalParameters //internal parameter store the useful parameters for the our agent
	idealRulesCachePtr *map[string]rules.RuleMatrix
	savedHistory       *accountabilityHistory
	trustMatrix        *trust
	importances        *importances
	forage             *forageStorage
}

type importances struct {
	requestAllocationImportance                *mat.VecDense
	commonPoolResourceRequestImportance        *mat.VecDense
	resourceReportImportance                   *mat.VecDense
	getTaxContributionImportance               *mat.VecDense
	decideIIGOMonitoringAnnouncementImportance *mat.VecDense
	getGiftRequestsImportance                  *mat.VecDense
}

// Store extra information which is not in the server and is helpful for our client
type observation struct {
	iigoObs           *iigoObservation
	iifoObs           *iifoObservation
	iitoObs           *iitoObservation
	pastDisastersList baseclient.PastDisastersList
}

type iigoObservation struct {
	allocationGranted shared.Resources
	taxDemanded       shared.Resources
	sanctionTiers     map[shared.ClientID]shared.IIGOSanctionsTier
}

type iifoObservation struct {
	// receivedDisasterPredictions shared.ReceivedDisasterPredictionsDict
	ourDisasterPrediction   shared.DisasterPredictionInfo
	finalDisasterPrediction shared.DisasterPrediction
}

type iitoObservation struct {
}

// all parameters are from 0 to 1 and they determine the personality of the agent.
type internalParameters struct {
	//trustMatrix *mat.Dense //this shouldn't be in internal parameters
	greediness    float64
	selfishness   float64
	fairness      float64
	collaboration float64
	riskTaking    float64

	// Judge GetPardonIslands config
	// days left on the sanction after which we can even considering pardoning other islands
	maxPardonTime int
	// specifies the maximum sanction tier after which we will no longer consider pardoning others
	maxTierToPardon shared.IIGOSanctionsTier
	// we will only consider pardoning islands which we trust with at least this value
	minTrustToPardon float64

	// Trust config
	historyWeight                float64
	historyFullTruthfulnessBonus float64
	monitoringWeight             float64
	monitoringResultChange       float64
	giftExtra                    bool
}

// type personality struct {
// }

//Overriding and extending the Initialise method of the BaseClient to initilise our client. This function happens after the init() function. At this point server has just initialised and the ServerReadHandle is available.
func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.BaseClient.Initialise(serverReadHandle)

	//custom things below, trust matrix initilised to values of 0
	c.idealRulesCachePtr = deepCopyRulesCache(c.ServerReadHandle.GetGameState().RulesInfo.AvailableRules)
	c.updateParents()
}

func (c *client) updateParents() {
	c.clientJudge.parent = c
	c.clientSpeaker.parent = c
	c.clientPresident.parent = c
}

func deepCopyRulesCache(AvailableRules map[string]rules.RuleMatrix) *map[string]rules.RuleMatrix {
	idealRulesCache := map[string]rules.RuleMatrix{}
	for k, v := range AvailableRules {
		idealRulesCache[k] = v
	}
	return &idealRulesCache
}

//Overriding the StartOfTurn method of the BaseClient
// func (c *client) StartOfTurn() {
// }

// GetVoteForRule returns the client's vote in favour of or against a rule.
// COMPULSORY: vote to represent your island's opinion on a rule
func (c *client) VoteForRule(ruleMatrix rules.RuleMatrix) shared.RuleVoteType {
	ruleDistance := c.decideRuleDistance(ruleMatrix)
	if ruleDistance < 5 { // TODO: calibrate the distance ranges
		return shared.Reject
	} else if ruleDistance < 15 {
		return shared.Abstain
	} else if ruleDistance >= 15 {
		return shared.Approve
	}
	return shared.Abstain
}

// decideRuleDistance returns the evaluated distance for the rule given in the argument
func (c *client) decideRuleDistance(ruleMatrix rules.RuleMatrix) float64 {
	// link rules

	// find rule corresponding to the rule that you need to evaluate
	idealRuleMatrix := (*c.idealRulesCachePtr)[ruleMatrix.RuleName]

	// calculate a distance
	distance := 0.0
	for i := 0; i < ruleMatrix.AuxiliaryVector.Len(); i++ {
		currentAuxValue := ruleMatrix.AuxiliaryVector.AtVec(i)
		for j := range ruleMatrix.RequiredVariables {

			idealValue := idealRuleMatrix.ApplicableMatrix.At(i, j)
			actualValue := ruleMatrix.ApplicableMatrix.At(i, j)

			if currentAuxValue == 0 {
				// ==0 condition
				if idealValue > 0 {
					distance += math.Abs(idealValue-actualValue) / idealValue
				} else {
					distance += math.Abs(idealValue - actualValue)
				}
			} else if currentAuxValue == 1 {
				// TODO: ACTUALLY IMPLEMENT THESE CONDITIONS
				// >0 condition
				if idealValue > 0 {
					distance += math.Abs(idealValue-actualValue) / idealValue
				} else {
					distance += math.Abs(idealValue - actualValue)
				}
			} else if currentAuxValue == 2 {
				// <=0 condition
				if idealValue > 0 {
					distance += math.Abs(idealValue-actualValue) / idealValue
				} else {
					distance += idealValue - actualValue
				}
			} else if currentAuxValue == 3 {
				// !=0 condition
				if idealValue != 0 {
					distance += math.Abs(idealValue-actualValue) / idealValue
				} else {
					distance += idealValue - actualValue
				}
			} else if currentAuxValue == 4 {
				if idealValue > 0 {
					distance += math.Abs(idealValue-actualValue) / idealValue
				} else {
					distance += math.Abs(idealValue - actualValue)
				}
				// it returns the value of the calculation
			}
		}

	}

	return distance
}

// GetVoteForElection returns the client's Borda vote for the role to be elected.
// COMPULSORY: use opinion formation to decide a rank for islands for the role
func (c *client) VoteForElection(roleToElect shared.Role, candidateList []shared.ClientID) []shared.ClientID {

	trustToID := map[float64]shared.ClientID{}
	trustList := []float64{}
	returnList := []shared.ClientID{}
	for i := 0; i < len(candidateList); i++ {
		trustScore := c.trustMatrix.GetClientTrust(candidateList[i]) //c.internalParam.agentsTrust[candidateList[i]]
		_, ok := trustToID[trustScore]
		for ok {
			trustScore += 0.000001 // Tiny increment to make it unique
			_, ok = trustToID[trustScore]
		}
		trustToID[trustScore] = candidateList[i]
		trustList = append(trustList, trustScore)
	}

	sort.Float64s(trustList)

	for i := len(trustList) - 1; i >= 0; i-- {
		// The idea is to have the very untrusted island to split the points in order
		// to increase the gap with good islands that we include and that we want to be elected.
		if trustList[i] > 0.25 || (len(trustList)-1)-i < 2 { //TODO: calibrate the trustScore so we don't always not rank //currently the infra does not support not ranking someone
			returnList = append(returnList, trustToID[trustList[i]])
		}
	}

	return returnList
}

func (c *client) StartOfTurn() {
	c.updateTrustFromSavedHistory()
	c.printConfig()
}

func (c *client) updateTrustFromSavedHistory() {
	if c.savedHistory.updated {
		newInfo := c.savedHistory.getNewInfo()

		if len(newInfo) > 0 {
			var lawfulnessSum float64

			for _, history := range newInfo {
				lawfulnessSum += history.LawfulRatio
			}
			averageTruthfulness := lawfulnessSum / float64(len(newInfo))

			for clientID, history := range newInfo {
				lawfulness := history.LawfulRatio

				c.trustMatrix.ChangeClientTrust(clientID, c.internalParam.historyWeight*(lawfulness-averageTruthfulness)) //potentially add * historyWeight to scale the update

				if floatEqual(lawfulness, 1) { //bonus for being fully truthful
					c.trustMatrix.ChangeClientTrust(clientID, c.internalParam.historyFullTruthfulnessBonus)
				}
			}
		}
		c.savedHistory.updated = false
	}
}

func (c *client) updateTrustMonitoring(data map[shared.CommunicationFieldName]shared.CommunicationContent) {
	if roleMonitored, ok := data[shared.RoleMonitored]; ok && roleMonitored.T == shared.CommunicationIIGORole {
		if monitoringResult, ok := data[shared.MonitoringResult]; ok && monitoringResult.T == shared.CommunicationBool {
			roleID := c.getRole(roleMonitored.IIGORoleData)

			if monitoringResult.BooleanData {
				// Monitored role was truthful
				c.trustMatrix.ChangeClientTrust(roleID, c.internalParam.monitoringWeight*c.internalParam.monitoringResultChange) // config?
			} else {
				// Monitored role cheated
				c.trustMatrix.ChangeClientTrust(roleID, -c.internalParam.monitoringWeight*c.internalParam.monitoringResultChange)
			}
		}
	}
}

//MonitorIIGORole decides whether to perform monitoring on a role
//COMPULOSRY: must be implemented
func (c *client) MonitorIIGORole(roleName shared.Role) bool {

	presidentID := c.getPresident()
	speakerID := c.getSpeaker()
	judgeID := c.getJudge()
	ourResources := c.getOurResources()
	// TODO: Choose sensible thresholds!
	trustThreshold := 0.5
	resourcesThreshold := shared.Resources(100)
	monitoring := false
	switch c.GetID() {
	case presidentID:
		// If we are the president.
		monitoring = (c.getTrust(speakerID) < trustThreshold ||
			c.getTrust(judgeID) < trustThreshold) &&
			(ourResources > resourcesThreshold)

	case speakerID:
		// If we are the Speaker.
		monitoring = (c.getTrust(presidentID) < trustThreshold ||
			c.getTrust(judgeID) < trustThreshold) &&
			(ourResources > resourcesThreshold)
	case judgeID:
		// If we are the Judge.
		monitoring = (c.getTrust(speakerID) < trustThreshold ||
			c.getTrust(judgeID) < trustThreshold) &&
			(ourResources > resourcesThreshold)
	default:
		break
	}
	return monitoring
}

//DecideIIGOMonitoringAnnouncement decides whether to share the result of monitoring a role and what result to share
//COMPULSORY: must be implemented
func (c *client) DecideIIGOMonitoringAnnouncement(monitoringResult bool) (resultToShare bool, announce bool) {
	collaborationThreshold := 0.5
	importance := c.importances.decideIIGOMonitoringAnnouncementImportance

	parameters := mat.NewVecDense(3, []float64{
		c.internalParam.selfishness,
		c.internalParam.fairness,
		c.internalParam.collaboration,
	})
	// Initialise Return values.
	announce = false
	resultToShare = monitoringResult

	// Calculate collaborationLevel based on the current personality of the client.
	collaborationLevel := mat.Dot(importance, parameters)

	if collaborationLevel > collaborationThreshold {
		// announce only if we are collaborative enough.
		announce = true
	}
	return resultToShare, announce
}

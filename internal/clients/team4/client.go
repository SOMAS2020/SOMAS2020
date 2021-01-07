// Package team4 contains code for team 4's client implementation
package team4

import (
	"fmt"
	"math"
	"reflect"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team4

func init() {
	team4client := client{
		BaseClient:    baseclient.NewClient(id),
		clientJudge:   judge{BaseJudge: &baseclient.BaseJudge{}, t: nil},
		clientSpeaker: speaker{BaseSpeaker: &baseclient.BaseSpeaker{}},
		yes:           "",
		obs:           &observation{},
		internalParam: &internalParameters{},
		savedHistory:  map[uint]map[shared.ClientID]judgeHistoryInfo{},
	}
	team4client.clientSpeaker.parent = &team4client
	team4client.clientJudge.parent = &team4client

	baseclient.RegisterClient(id, &team4client)
}

type client struct {
	*baseclient.BaseClient //client struct has access to methods and fields of the BaseClient struct which implements implicitly the Client interface.
	clientJudge            judge
	clientSpeaker          speaker

	//custom fields
	yes                string              //this field is just for testing
	obs                *observation        //observation is the raw input into our client
	internalParam      *internalParameters //internal parameter store the useful parameters for the our agent
	idealRulesCachePtr *map[string]rules.RuleMatrix
	savedHistory       map[uint]map[shared.ClientID]judgeHistoryInfo
}

// Store extra information which is not in the server and is helpful for our client
type observation struct {
	iigoObs *iigoObservation
	iifoObs *iifoObservation
	iitoObs *iitoObservation
}

type iigoObservation struct {
	allocationGranted shared.Resources
	taxDemanded       shared.Resources
}

type iifoObservation struct {
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
	agentsTrust   []float64
}

type personality struct {
}

//Not our code
// func getislandParams() islandParams {
// 	return islandParams{
// 		giftingThreshold:            45,  // real number
// 		equity:                      0.1, // 0-1
// 		complianceLevel:             0.1, // 0-1
// 		resourcesSkew:               1.3, // >1
// 		saveCriticalIsland:          true,
// 		escapeCritcaIsland:          true,
// 		selfishness:                 0.3, //0-1
// 		minimumRequest:              30,  // real number
// 		disasterPredictionWeighting: 0.4, // 0-1
// 		recidivism:                  0.3, // real number
// 		riskFactor:                  0.2, // 0-1
// 		friendliness:                0.3, // 0-1
// 		anger:                       0.2, // 0-1
// 		aggression:                  0.4, // 0-1
// 	}
// }

// func (l *locator) calculateMetaStrategy() {
// 	newMetaStrategy := metaStrategy{}
// 	currentParams := l.islandParamsCache
// 	newMetaStrategy.conquest = !currentParams.saveCriticalIsland && currentParams.aggression > 0.7 && currentParams.friendliness < 0.5
// 	newMetaStrategy.saviour = currentParams.saveCriticalIsland && currentParams.friendliness >= 0.5 && currentParams.aggression <= 0.7
// 	newMetaStrategy.democrat = currentParams.complianceLevel > 0.5 && currentParams.selfishness < 0.5
// 	newMetaStrategy.generous = currentParams.selfishness < 0.5 && currentParams.saveCriticalIsland && currentParams.recidivism < 0.5 && currentParams.friendliness > 0.5
// 	newMetaStrategy.lawful = currentParams.complianceLevel > 0.5 && currentParams.selfishness < 0.5
// 	newMetaStrategy.legislative = currentParams.equity > 0.5
// 	newMetaStrategy.executive = currentParams.recidivism < 0.5
// 	newMetaStrategy.fury = !currentParams.saveCriticalIsland && currentParams.aggression > 0.6 && currentParams.anger > 0.8
// 	newMetaStrategy.independent = currentParams.recidivism > 0.5
// 	l.locationStrategy = newMetaStrategy
// }

/////////////////////////

//Overriding and extending the Initialise method of the BaseClient to initilise our client
func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	// c.BaseClient.Initialise(serverReadHandle)
	c.ServerReadHandle = serverReadHandle
	c.LocalVariableCache = rules.CopyVariableMap()

	//custom things below, trust matrix initilised to values of 1
	numClient := len(c.ServerReadHandle.GetGameState().ClientLifeStatuses)
	c.internalParam = &internalParameters{agentsTrust: make([]float64, numClient)}
	iigoObs := &iigoObservation{
		allocationGranted: shared.Resources(0),
		taxDemanded:       shared.Resources(0),
	}
	iifoObs := &iifoObservation{}
	iitoObs := &iitoObservation{}
	c.obs = &observation{
		iigoObs: iigoObs,
		iifoObs: iifoObs,
		iitoObs: iitoObs,
	}

	c.idealRulesCachePtr = deepCopyRulesCache(rules.AvailableRules)

	// numClient := len(shared.TeamIDs)
	// v := make([]float64, numClient*numClient)
	// for i := range v {
	// 	v[i] = 1
	// }
	// c.internalParam = &internalParameters{
	// 	trustMatrix: mat.NewDense(numClient, numClient, v),
	// }
}

func deepCopyRulesCache(AvailableRules map[string]rules.RuleMatrix) *map[string]rules.RuleMatrix {
	idealRulesCache := map[string]rules.RuleMatrix{}
	for k, v := range AvailableRules {
		idealRulesCache[k] = v
	}
	return &idealRulesCache
}

//Overriding the StartOfTurn method of the BaseClient
func (c *client) StartOfTurn() {
	c.yes = "yes"
	c.Logf(`what are you doing?
	=========================================
	==============================================
	================================================`)
	c.Logf("this is a %v for you ", c.yes)
	fmt.Println(reflect.TypeOf(c))
}

// GetVoteForRule returns the client's vote in favour of or against a rule.
// COMPULSORY: vote to represent your island's opinion on a rule
func (c *client) VoteForRule(ruleMatrix rules.RuleMatrix) shared.RuleVoteType {
	// TODO implement decision on voting that considers the rule
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
	// rules with 0(==) as auxiliary vector element(s)

	// find rule correspondent to the rule that you need to evaluate
	idealRuleMatrix := (*c.idealRulesCachePtr)[ruleMatrix.RuleName]

	// calculate a distance and a distance
	distance := 0.0
	for i := 0; i < ruleMatrix.AuxiliaryVector.Len(); i++ {
		currentAuxValue := ruleMatrix.AuxiliaryVector.AtVec(i)
		for j := range ruleMatrix.RequiredVariables {

			idealValue := idealRuleMatrix.ApplicableMatrix.At(i, j)
			actualValue := ruleMatrix.ApplicableMatrix.At(i, j)

			if currentAuxValue == 0 {
				// ==0 condition
				distance += math.Abs(idealValue-actualValue) / idealValue
			} else if currentAuxValue == 1 {
				// TODO: ACTUALLY IMPLEMENT THESE CONDITIONS
				// >0 condition
				distance += 10000
			} else if currentAuxValue == 2 {
				// >=0 condition
				distance += 10000
			} else if currentAuxValue == 3 {
				// !=0 condition
				distance += math.Abs(idealValue-actualValue) / idealValue
			} else if currentAuxValue == 4 {
				distance += 10000
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
		trustScore := c.internalParam.agentsTrust[candidateList[i]]
		trustToID[trustScore] = candidateList[i]
		trustList = append(trustList, trustScore)
	}
	sort.Float64s(trustList)

	for i := len(trustList) - 1; i >= 0; i-- {
		// The idea is to have the very untrusted island to split the points in order
		// to increase the gap with good islands that we include and that we want to be elected.
		if trustList[i] > 0.25 { //TODO: calibrate the trustScore so we don't always not rank
			returnList = append(returnList, trustToID[trustList[i]])
		}
	}

	return returnList
}

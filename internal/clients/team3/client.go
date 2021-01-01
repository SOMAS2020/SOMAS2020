// Package team3 contains code for team 3's client implementation
package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team3

func init() {
	ourClient := &client{BaseClient: baseclient.NewClient(id)}

	baseclient.RegisterClient(id, ourClient)
}

type client struct {
	*baseclient.BaseClient
	// ## Gifting ##

	acceptedGifts        map[shared.ClientID]int
	requestedGiftAmounts map[shared.ClientID]int
	// receivedResponses // todo

	// ## Trust ##

	theirTrustScore map[shared.ClientID]uint
	trustScore      map[shared.ClientID]uint

	// ## Role performance ##

	judgePerformance     map[shared.ClientID]int
	speakerPerformance   map[shared.ClientID]int
	presidentPerformance map[shared.ClientID]int

	// ## Game state & History ##

	// unused or replaced by getter functions
	// currentIteration iterationInfo
	// islandsAlive uint
	// localPool int

	// declaredResources is a map of all declared island resources
	declaredResources map[shared.ClientID]shared.Resources
	//disasterPredictions gives a list of predictions by island for each turn
	disasterPredictions []map[shared.ClientID]shared.DisasterPrediction

	// ## Compliance ##

	timeSinceCaught uint
	numTimeCaught   uint
	compliance      float32

	// allVotes [rule][shared.ClientID] TODO

	// params is list of island wide function parameters
	params islandParams

	presidentObj president
}

type islandParams struct {
	giftingThreshold            int
	equity                      float32
	complianceLevel             float32
	resourcesSkew               float32
	saveCriticalIsland          bool
	escapeCritcaIsland          bool
	selfishness                 float32
	minimumRequest              int
	disasterPredictionWeighting float32
	DesiredRuleSet              []string
	recidivism                  float32
	riskFactor                  float32
	friendliness                float32
	anger                       float32
	aggression                  float32
}

func (c *client) DemoEvaluation() {
	evalResult, err := rules.BasicBooleanRuleEvaluator("Kinda Complicated Rule")
	if err != nil {
		panic(err.Error())
	}
	c.Logf("Rule Eval: %t", evalResult)
}

// Package team3 contains code for team 3's client implementation
package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
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
	receivedResponses    []shared.GiftResponse

	// ## Trust ##

	theirTrustScore  map[shared.ClientID]float64
	trustScore       map[shared.ClientID]float64
	trustMapAgg      map[shared.ClientID][]float64
	theirTrustMapAgg map[shared.ClientID][]float64

	// ## Role performance ##

	judgePerformance     map[shared.ClientID]int
	speakerPerformance   map[shared.ClientID]int
	presidentPerformance map[shared.ClientID]int

	// ## Game state & History ##
	criticalStatePrediction critcalStatePrediction

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
	compliance      float64

	// allVotes stores the votes of each island for/against each rule
	allVotes map[string]map[shared.ClientID]bool

	// params is list of island wide function parameters
	params islandParams
}

type critcalStatePrediction struct {
	upperBound shared.Resources
	lowerBound shared.Resources
}

type islandParams struct {
	giftingThreshold            int
	equity                      float64
	complianceLevel             float64
	resourcesSkew               float64
	saveCriticalIsland          bool
	escapeCritcaIsland          bool
	selfishness                 float64
	minimumRequest              int
	disasterPredictionWeighting float64
	DesiredRuleSet              []string
	recidivism                  float64
	riskFactor                  float64
	friendliness                float64
	anger                       float64
	aggression                  float64
}

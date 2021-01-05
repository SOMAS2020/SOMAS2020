// Package team4 contains code for team 4's client implementation
package team4

import (
	"fmt"
	"reflect"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team4

func init() {
	baseclient.RegisterClient(id, &client{BaseClient: baseclient.NewClient(id)})
}

type client struct {
	*baseclient.BaseClient //client struct has access to methods and fields of the BaseClient struct which implements implicitly the Client interface.

	//custom fields
	yes           string              //this field is just for testing
	obs           *observation        //observation is the raw input into our client
	internalParam *internalParameters //internal parameter store the useful parameters for the our agent

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

//Overriding the Initialise method of the BaseClient to initilise the trust matrix too
func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle

	//custom things below, trust matrix initilised to values of 1

	// numClient := len(shared.TeamIDs)
	// v := make([]float64, numClient*numClient)
	// for i := range v {
	// 	v[i] = 1
	// }
	// c.internalParam = &internalParameters{
	// 	trustMatrix: mat.NewDense(numClient, numClient, v),
	// }
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

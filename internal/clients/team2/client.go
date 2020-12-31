// Package team2 contains code for team 2's client implementation
package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team2

// Histories for what we want our agent to remember

type CommonPoolState struct {
	turn   int
	amount float64
}

type CommonPoolHistory map[shared.Resources][]CommonPoolState

// Type for Empathy level assigned to each other team
type EmpathyLevel int

const (
	Selfish EmpathyLevel = iota
	FairSharer
	Altruist
)

type IslandEmpathies map[shared.ClientID]EmpathyLevel

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient: baseclient.NewClient(id),
			// add other things we want to remember (Histories)
			// commonpoolHistory: CommonPoolHistory{},
			// need to init to initially assume other islands are fair
			// IslandEmpathies: IslandEmpathies{},
		},
	)
}

// we have to initialise our client somehow
type client struct {
	// should this have a * in front?
	*baseclient.BaseClient

	islandEmpathies   IslandEmpathies
	commonPoolHistory CommonPoolHistory
}

// After declaring the struct we have to actually make an object for the client
func NewClient(clientID shared.ClientID) baseclient.Client {
	// return a reference to the client struct variable's memory address
	return &client{
		BaseClient: baseclient.NewClient(clientID),
		// commonpoolHistory: CommonPoolHistory{},
		// we could experiment with how being more/less trustful affects agent performance
		// i.e. start with assuming all islands selfish, normal, altruistic
		islandEmpathies: IslandEmpathies{},
	}
}

func (c *client) islandEmpathyLevel() EmpathyLevel {
	clientInfo := c.gameState().ClientInfo

	// switch statement to toggle between three levels
	// change our state based on these cases
	switch {
	case clientInfo.LifeStatus == shared.Critical:
		return Selfish
		// replace with some expression
	case (true):
		return Altruist
	default:
		return FairSharer
	}
}

// func (c client) functionName() {

// }

//this determines our internal threshold for survival, allocationrec is the output of the function AverageCommonPool which determines which role we will be
// func InternalThreshold(DaysUntilDisaster uint, allocationrec float64) float64 {
// 	var initialThreshold float64      //figure this out
// 	var defaultdisasterday uint       //set this for when we do not know when a disaster will occur
// 	var reverseDisasterCountdown uint //tune this variable

// 	return initialThreshold + (float64(reverseDisasterCountdown)-float64(DaysUntilDisaster))*allocationrec
// }

// func (c *client) criticalStatus() bool {
// 	clientInfo := c.gameState().ClientInfo
// 	if clientInfo.LifeStatus == shared.Critical { //not sure about shared.Critical
// 		return true
// 	}

// 	return false
// }

//CommonPoolResourceRequest() shared.Resources
//this determines whether we need or can give resources based on the gamestate
// func (c *client) TakeRequestResources()
// 	if criticalStatus() {

// 	} else if

// }

// //this determines how we spend our resources given specific situations
// func (c *client) handleResources(){
// 	clientInfo := c.gameState().ClientInfo
// 	if criticalStatus()==1{  //if we are critical then take resources from the pool
// 		TakeRequestResources()
// 	}
// 	if determineTax()>clientInfo.Resources { //if we cannot pay our tax
// 		TakeRequestResources()
// 	}
// 	if CheckOthersCrit(){ //if another island is critical
// 		SelectResourcesToGive()
// 	}
// }

// //will find out how much tax we have to pay
// func (c *client) determineTax() float64{

// }

// //will loop through all agents and check if anyone is critical, CHANGE FOR CRIT RATHER THAN DEATH
// func (c *client) CheckOthersCrit(){
// 	for clientID, status := range c.gameState().ClientLifeStatuses {
// 		if status != shared.Dead && clientID != c.GetID() {
// 			return
// 		}
// 	}

// }

// //determine if an agent is worthy of being helped out of critical
// func (c *client) SelectResourcesToGive(){

// }

func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

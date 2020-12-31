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

func (c *client) StartTurn() shared.Resources {
	
}

//TODO: change if else statements to switch case

//determines how much to request from the pool
func (c *client) CommonPoolResourceRequest() shared.Resources {
	clientInfo := c.gameState().ClientInfo
	if criticalStatus() {  //if we are critical then take resources from the pool
		resources = 20 //change to (threshold - clientInfo.Resources) but need to work out threshold amount
		return resources
		}
	if determineTax()>clientInfo.Resources { //if we cannot pay our tax
		return 15 //change to (tax-clientInfo.Resources)
	}
	if clientInfo.shared.Resources < InternalThreshold() { //if we are below our own internal threshold
		return 10
	}
		standardGamePlay()
}

//determines how many resources you actually take
func (c *client) RequestAllocation() shared.Resources {
	if criticalStatus() {  //if we are critical then take resources from the pool
		//code for when critical
	}
	if determineTax()>clientInfo.Resources { //if we cannot pay our tax
		//
	}
	if CheckOthersCrit(){ //if another island is critical
		//
	}
	if clientInfo.shared.Resources < InternalThreshold() { //if we are below our own internal threshold
		//
	}
	//standardGamePlay()
	return allocation
}

//determines how much we put into pool
func (c *client) GetTaxContribution() shared.Resources {
	var allocation shared.Resources
	if criticalStatus() {  //if we are critical then take resources from the pool
		return 0     //tax evasion
	}
	if determineTax()>clientInfo.Resources { //if we cannot pay our tax
		return 0     //no choice but tax evasion
	}
	if clientInfo.shared.Resources < InternalThreshold() { //if we are below our own internal threshold
		return 0     //make this TaxMin, which will be the minimum contribution
	}
	//This is our default allocation, this determines how much to give based off of previous common pool level
	return AverageCommonPoolDilemma()
}



//this determines our internal threshold for survival, allocationrec is the output of the function AverageCommonPool which determines which role we will be
func InternalThreshold(DaysUntilDisaster uint, allocationrec float64) float64 {
	var initialThreshold float64      //figure this out
	var defaultdisasterday uint       //set this for when we do not know when a disaster will occur
	var reverseDisasterCountdown uint //tune this variable

	return initialThreshold + (float64(reverseDisasterCountdown)-float64(DaysUntilDisaster))*allocationrec
}

func (c *client) criticalStatus() bool {
	clientInfo := c.gameState().ClientInfo
	if clientInfo.LifeStatus == shared.Critical { //not sure about shared.Critical
		return true
	}
	return false
}


//will find out how much tax we have to pay
func (c *client) determineTax() shared.Resources {
	return c.taxAmount
}

//will loop through all agents and check if anyone is critical
func (c *client) CheckOthersCrit(){
	for clientID, status := range c.gameState().ClientInfo.LifeStatus {
		if status != shared.Critical && clientID != c.GetID() {
			return
		}
	}

}


func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

//this function determines how much to contribute to the common pool depending on whether other agents are altruists,fair sharers etc
//it only needs the current resource level and the current turn as inputs
//the output will be an integer which is a recommendation on how much to add to the pool, with this recommendation there will be a weighting of how important it is we contribute that exact amount
//this will be a part of other decision making functions which will have their own weights 

//tunable parameters:
//how much to give the pool on our first turn: default_strat
//after how many rounds of struggling pool to intervene and become altruist: intervene
//the number of turns at the beginning we cannot free ride: no_freeride
//the factor in which the common pool increases by to decide if we should free ride: freeride
//the factor which we multiply the fair_sharer average by: tune_average
//the factor which we multiply the altruist value by: tune_alt

//NEED TO CHANGE RESOURCE AND TURN VARIABLE INTO THEIR C *CLIENT EQUIVALENTS

func (c *client)AverageCommonPoolDilemma() float64{

	var default_strat float64=x    //this parameter will determine how much we contribute on the first turn when there is no data to make a decision 
	var ResourceHistory [10]float64    //store all previous resource levels to find an average increase or decrease, currently stores 10 turns

	var fair_sharer float64          //this is how much we contribute when we are a fair sharer and altruist
	var altruist float64

	var decreasing_pool float64           //records for how many turns the common pool is decreasing
	var intervene float64=y                 //when the pool is struggling for this number of rounds we intervene
	var no_freeride float64=z             //how many turns at the beginning we cannot free ride for
	var freeride float64=k 		//what factor the common pool must increase by for us to considered free riding
	var tune_average float64=v            //what factor of the average to contribute when fair sharing, default is 1 to give the average
	var tune_alt float64=l                //what factor of the average to contribute when being altruistic, will be much higher than fair sharing

	if ResourceHistory==nil {        //if there is no historical data then use default strategy
		return default_strat
		decreasing_pool=0 
	}

	altruist=determine_altruist(ResourceHistory,turn)  //determines altruist amount
	fair_sharer=determine_fair(ResourceHistory,turn)   //determines fair sharer amount
	
	if ResourceHistory[turn-1]=>ResourceHistory[turn] {  //decreasing common pool means consider altruist
		if decreasing_pool=>intervene {
			decreasing_pool=0            //once we have contriubuted a lot we reset
			return altruist             
		}	
	}

	if turn>no_freeride {          //we will not allow ourselves to use free riding at the start of the game
		if ResourceHistory[turn-1]<(ResourceHistory[turn]*freeride) {    
			if ResourceHistory[turn-2]<(ResourceHistory[turn-1]*freeride {   //two large jumps then we free ride
				return 0   //change this to be tax
			}    
		}
	}
	return fair_sharer    //by default we contribute a fair share
}	

func determine_altruist(ResourceHistory [10]float64,turn uint) int{   //identical to fair sharing but a larger factor to multiple the average contribution by
	for j:=turn; j>0; j-- {               //we are trying to find the most recent instance of the common pool increasing and then use that value
		if ResourceHistory[j-1]-ResourceHistory[j]>0 {  
			return ((ResourceHistory[j-1]-ResourceHistory[j])/6)*tune_alt
		}
	}	
}

func determine_fair(ResourceHistory,turn uint) int{     //can make more sophisticated! Right now just contribute the average, default matters the most
	for j:=turn; j>0; j-- {               //we are trying to find the most recent instance of the common pool increasing and then use that value
		if ResourceHistory[j-1]-ResourceHistory[j]>0 {  
			return ((ResourceHistory[j-1]-ResourceHistory[j])/6)*tune_average   //make 6 variable for no of agents
		}
	}
}


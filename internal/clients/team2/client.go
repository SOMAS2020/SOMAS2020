// Package team2 contains code for team 2's client implementation
package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team2

// Histories for what we want our agent to remember

//type CommonPoolState struct {
//	season uint
//	amount float64
//}

//type Turn uint

// type ServerReadHandle interface {
// 	GetGameState() gamestate.ClientGameState
// }

func (c *client) StartOfTurn() {
	CommonPoolUpdate(c, c.commonPoolHistory) //add the common pool level
}

func CommonPoolUpdate(c *client, commonPoolHistory CommonPoolHistory) {
	currentPool := c.gameState().CommonPool
	c.commonPoolHistory[c.gameState().Turn] = float64(currentPool)
}

type CommonPoolHistory map[uint]float64

// Type for Empathy level assigned to each other team
type EmpathyLevel int

const (
	Selfish EmpathyLevel = iota
	FairSharer
	Altruist
)

type IslandEmpathies map[shared.ClientID]EmpathyLevel

// Could be used to store our expectation on an island's behaviour (about anything)
// vs what they actually did
type ExpectationReality struct {
	exp  int
	real int
}

type Situation string

const (
	President   Situation = "President"
	Judge       Situation = "Judge"
	Speaker     Situation = "Speaker"
	Foraging    Situation = "Foraging"
	GiftRequest Situation = "GiftRequest"
	GiftGiven   Situation = "GiftGiven"
)

type Opinion struct {
	Histories    map[Situation][]int
	Performances map[Situation]ExpectationReality
}

type ForageInfo struct {
	DecisionMade      shared.ForageDecision
	ResourcesObtained shared.Resources
}

type GiftInfo struct {
	requested shared.GiftRequest
	gifted    shared.GiftOffer
	reason    shared.AcceptReason
}

type GiftExchange struct {
	IslandRequest map[uint]GiftInfo
	OurRequest    map[uint]GiftInfo
}

type OpinionHist map[shared.ClientID]Opinion
type PredictionsHist map[shared.ClientID][]shared.DisasterPrediction
type ForagingReturnsHist map[shared.ClientID][]ForageInfo
type GiftHist map[shared.ClientID]GiftExchange

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

	islandEmpathies     IslandEmpathies
	commonPoolHistory   CommonPoolHistory
	opinionHist         OpinionHist
	predictionsHist     PredictionsHist
	foragingReturnsHist ForagingReturnsHist
	giftHist            GiftHist
}

//NewClient After declaring the struct we have to actually make an object for the client
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

func criticalStatus(c *client) bool {
	clientInfo := c.gameState().ClientInfo
	if clientInfo.LifeStatus == shared.Critical { //not sure about shared.Critical
		return true
	}
	return false
}

//TODO: maybe change if else statements to switch case
//determines how much to request from the pool
func (c *client) CommonPoolResourceRequest() shared.Resources {
	clientInfo := c.gameState().ClientInfo
	if criticalStatus(c) {
		return 20 //change to (threshold - clientInfo.Resources) but need to work out threshold amount
	}
	if determineTax(c) > clientInfo.Resources { //if we cannot pay our tax
		return 15 //change to (tax-clientInfo.Resources)
	}
	if c.gameState().ClientInfo.Resources < internalThreshold(c) { //if we are below our own internal threshold
		return 10
	}
	//TODO: maybe separate standard gameplay when no-one is critical vs when others are critical
	//standardGamePlay or others are critical
	return 0
}

//determines how many resources you actually take - currrently going to take however much we say (playing nicely)
func (c *client) RequestAllocation() shared.Resources {
	if criticalStatus(c) {
		//change to (threshold - clientInfo.Resources) but need to work out threshold amount
		return 20
	}
	if determineTax(c) > c.gameState().ClientInfo.Resources { //if we cannot pay our tax
		return 15 //change to (tax-clientInfo.Resources)
	}
	if c.gameState().ClientInfo.Resources < internalThreshold(c) { //if we are below our own internal threshold
		return 10
	}
	//TODO: maybe separate standard gameplay when no-one is critical vs when others are critical
	//standardGamePlay or others are critical
	return 0
}

//GetTaxContribution determines how much we put into pool
func (c *client) GetTaxContribution() shared.Resources {
	var ourResources = c.gameState().ClientInfo.Resources
	var Taxmin shared.Resources = determineTax(c)                          //minimum tax contribution to not break law
	var allocation shared.Resources = AverageCommonPoolDilemma(c) + Taxmin //This is our default allocation, this determines how much to give based off of previous common pool level
	if criticalStatus(c) {
		return 0 //tax evasion
	}
	if determineTax(c) > ourResources {
		return 0 //no choice but tax evasion
	}
	if ourResources < internalThreshold(c) { //if we are below our own internal threshold
		return Taxmin
	}
	if checkOthersCrit(c) {
		return (ourResources - internalThreshold(c) - Taxmin) / 2 //TODO: tune this
	}

	allocation = AverageCommonPoolDilemma(c) + Taxmin
	return allocation
}

//internalThreshold determines our internal threshold for survival, allocationrec is the output of the function AverageCommonPool which determines which role we will be
func internalThreshold(c *client) shared.Resources {
	var initialThreshold shared.Resources = 30 //TODO: Fix a number ourselves when we don't know threshold, else find out how to access this
	var allocationrec = AverageCommonPoolDilemma(c)
	var disasterEffectOnUs float64 = 3            //TODO: call function from Hamish's part to get map clientID: effect
	var disasterPredictionConfidence float64 = 50 //TODO: call function from Hamish's oart to get this confidence level
	var turnsLeftUntilDisaster uint = 3           //TODO: call function from Hamish's oart to get number of turns
	if disasterEffectOnUs > 4 {
		return (initialThreshold + allocationrec) * shared.Resources(disasterPredictionConfidence/10) //TODO: tune divisor
	}
	if turnsLeftUntilDisaster < 3 {
		return 3
	}
	return initialThreshold + allocationrec
}

//determineTaxreturns how much tax we have to pay
func determineTax(c *client) shared.Resources {
	return shared.Resources(shared.TaxAmount) //TODO: not sure if this is correct tax amount to use
}

//checkOthersCrit checks if anyone else is critical
func checkOthersCrit(c *client) bool {
	for clientID, status := range c.gameState().ClientLifeStatuses {
		if status == shared.Critical && clientID != c.GetID() {
			return true
		}
	}
	return false
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

func AverageCommonPoolDilemma(c *client) shared.Resources {
	ResourceHistory := c.commonPoolHistory
	turn := c.gameState().Turn
	var default_strat float64 = 50 //this parameter will determine how much we contribute on the first turn when there is no data to make a decision

	var fair_sharer float64 //this is how much we contribute when we are a fair sharer and altruist
	var altruist float64

	var decreasing_pool float64 //records for how many turns the common pool is decreasing
	var intervene float64 = 3   //when the pool is struggling for this number of rounds we intervene
	var no_freeride float64 = 3 //how many turns at the beginning we cannot free ride for
	var freeride float64 = 5    //what factor the common pool must increase by for us to considered free riding

	if turn == 1 { //if there is no historical data then use default strategy
		decreasing_pool = 0
		return shared.Resources(default_strat)
	}

	altruist = c.determine_altruist(turn) //determines altruist amount
	fair_sharer = c.determine_fair(turn)  //determines fair sharer amount

	prevTurn := turn - 1
	prevTurn2 := turn - 2
	if ResourceHistory[prevTurn] > ResourceHistory[turn] { //decreasing common pool means consider altruist
		decreasing_pool++ //increment decreasing pool every turn the pool decreases
		if decreasing_pool > intervene {
			decreasing_pool = 0 //once we have contributed a lot we reset
			return shared.Resources(altruist)
		}
	}

	if float64(turn) > no_freeride { //we will not allow ourselves to use free riding at the start of the game
		if ResourceHistory[prevTurn] < (ResourceHistory[turn] * freeride) {
			if ResourceHistory[prevTurn2] < (ResourceHistory[prevTurn] * freeride) { //two large jumps then we free ride
				return 0
			}
		}
	}
	return shared.Resources(fair_sharer) //by default we contribute a fair share
}

func (c *client) determine_altruist(turn uint) float64 { //identical to fair sharing but a larger factor to multiple the average contribution by
	ResourceHistory := c.commonPoolHistory
	var tune_alt float64 = 2    //what factor of the average to contribute when being altruistic, will be much higher than fair sharing
	for j := turn; j > 0; j-- { //we are trying to find the most recent instance of the common pool increasing and then use that value
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn] > 0 {
			return ((ResourceHistory[j] - ResourceHistory[prevTurn]) / 6) * tune_alt
		}
	}
	return 0
}

func (c *client) determine_fair(turn uint) float64 { //can make more sophisticated! Right now just contribute the average, default matters the most
	ResourceHistory := c.commonPoolHistory
	var tune_average float64 = 1 //what factor of the average to contribute when fair sharing, default is 1 to give the average
	for j := turn; j > 0; j-- {  //we are trying to find the most recent instance of the common pool increasing and then use that value
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn] > 0 {
			return ((ResourceHistory[j] - ResourceHistory[prevTurn]) / 6) * tune_average //make 6 variable for no of agents
		}
	}
	return 0
}

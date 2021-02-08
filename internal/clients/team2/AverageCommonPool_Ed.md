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

```
Func (c *client)AverageCommonPoolDilemma() float64{

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
				return 0 
			}    
		}
	}
	return fair_sharer    //by default we contribute a fair share
}	

Func determine_altruist(ResourceHistory [10]float64,turn uint) int{   //identical to fair sharing but a larger factor to multiple the average contribution by
	for j:=turn; j>0; j-- {               //we are trying to find the most recent instance of the common pool increasing and then use that value
		if ResourceHistory[j-1]-ResourceHistory[j]>0 {  
			return ((ResourceHistory[j-1]-ResourceHistory[j])/6)*tune_alt
		}
	}	
}

Func determine_fair(ResourceHistory,turn uint) int{     //can make more sophisticated! Right now just contribute the average, default matters the most
	for j:=turn; j>0; j-- {               //we are trying to find the most recent instance of the common pool increasing and then use that value
		if ResourceHistory[j-1]-ResourceHistory[j]>0 {  
			return ((ResourceHistory[j-1]-ResourceHistory[j])/6)*tune_average   //make 6 variable for no of agents
		}
	}
}
```

```
///////////////////////////  TESTBENCH ////////////////////////////////
package main
import "fmt"
func main(){
 fmt.Println(AverageCommonPoolDilemma())
 return
}

func AverageCommonPoolDilemma() float64 {
	ResourceHistory := make(map[uint]float64)
	ResourceHistory[0] = 0
	ResourceHistory[1] = 100
	ResourceHistory[2] = 200
	ResourceHistory[3] = 300
	
	var turn uint=3
	var default_strat float64 = 20 
	var fair_sharer float64 
	var altruist float64

	var decreasing_pool float64
	var intervene float64 = 3   
	var no_freeride float64 = 1
	var freeride float64 = 5  

	if turn==0 { 
		decreasing_pool = 0
		return default_strat
	}

	altruist = determine_altruist(turn,ResourceHistory)
	fair_sharer = determine_fair(turn,ResourceHistory)  

	prevTurn := turn - 1
	prevTurn2 := turn -2
	if ResourceHistory[prevTurn] > ResourceHistory[turn] {
		decreasing_pool++
		if decreasing_pool > intervene {
			decreasing_pool = 0
			return altruist
		}
	}

	if float64(turn) > no_freeride { 
		if ResourceHistory[prevTurn] < (ResourceHistory[turn] * freeride) {
			if ResourceHistory[prevTurn2] < (ResourceHistory[prevTurn] * freeride) {
				return 0
			}
		}
	}
	return fair_sharer
}

func determine_altruist(turn uint, ResourceHistory map[uint]float64 ) float64 { 
	var tune_alt float64 = 2    
	for j := turn; j > 0; j-- { 
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn] > 0 {
			return ((ResourceHistory[j] - ResourceHistory[prevTurn]) / 6) * tune_alt
		}
	}
	return 0
}

func determine_fair(turn uint,ResourceHistory map[uint]float64) float64 {
	var tune_average float64 = 1
	for j := turn; j > 0; j-- {  
		prevTurn := j - 1
		if ResourceHistory[j]-ResourceHistory[prevTurn]> 0 {
			return ((ResourceHistory[j]- ResourceHistory[prevTurn]) / 6) * tune_average 
		}
	}
	return 0
}
```
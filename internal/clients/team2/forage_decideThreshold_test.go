package team2

import "fmt"

func main() {
	fmt.Printl(Otheragentinfo())
	fmt.Println(decideThreshold())
	return
}

func decideThreshold() float64 {
	if Otheragentinfo() == 1 {
		return 0.95
	} else if Otheragentinfo() > 1 {
		return 0.95 - (Otheragentinfo() * 0.15)
	} else {
		return 0.1
	}
}

func Otheragentinfo() float64 { 
	HuntNum := 0.00   
	TotalHuntNum := 0.00
	totalDecisions := 0.00   
	for _, id := range clientInfo {​​​​​​​    
		for index, forageInfo := range c.foragingReturnsHist[id] {  
			HuntNum = (HuntNum + forageInfo.DecisionMade.ForageType)/totalDecisions
			totalDecisions++
			TotalHuntNum=TotalHuntNum + HuntNum
		}
	totalDecisions=0
	HuntNum=0
	}​​​​
	return TotalHuntNum
}

//TEST CASE: 
//INPUT: 
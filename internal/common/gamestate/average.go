package gamestate

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

func AverageGameStates(currentAvg []GameState, newEntry []GameState) []GameState {
	if len(currentAvg) >= len(newEntry) {
		newStates := make([]GameState, len(currentAvg))
		for i := 0; i < len(newEntry); i++ {
			newStates[i] = averageSingleGameState(currentAvg[i], newEntry[i])
		}
		for i := len(newEntry); i < len(currentAvg); i++ {
			newStates[i] = currentAvg[i]
		}
		return newStates
	}
	newStates := make([]GameState, len(newEntry))
	for i := 0; i < len(currentAvg); i++ {
		newStates[i] = averageSingleGameState(currentAvg[i], newEntry[i])
	}
	for i := len(currentAvg); i < len(newEntry); i++ {
		newStates[i] = newEntry[i]
	}
	return newStates
}

func checkStatesLineUp(state1 GameState, state2 GameState) bool {
	return state1.Turn == state2.Turn
}

func averageSingleGameState(currentAvg GameState, newVal GameState) GameState {
	copyOfCurrent := currentAvg.Copy()
	copyOfNewVal := newVal.Copy()
	if checkStatesLineUp(copyOfCurrent, copyOfNewVal) {
		copyOfCurrent.CommonPool += (copyOfNewVal.CommonPool - copyOfCurrent.CommonPool) / 2
		copyOfCurrent.ClientInfos = averageClientInfos(copyOfCurrent.ClientInfos, copyOfNewVal.ClientInfos)
		copyOfCurrent.IIGOAllocationMap = averageIIGOPayments(copyOfCurrent.IIGOAllocationMap, copyOfNewVal.IIGOAllocationMap)
		copyOfCurrent.IIGOSanctionMap = averageIIGOPayments(copyOfCurrent.IIGOSanctionMap, copyOfNewVal.IIGOSanctionMap)
		copyOfCurrent.IIGOTaxAmount = averageIIGOPayments(copyOfCurrent.IIGOTaxAmount, copyOfNewVal.IIGOTaxAmount)
	}
	return copyOfCurrent
}

func averageIIGOPayments(pay1 map[shared.ClientID]shared.Resources, pay2 map[shared.ClientID]shared.Resources) map[shared.ClientID]shared.Resources {
	newMap := make(map[shared.ClientID]shared.Resources)
	for key, val := range pay1 {
		if val2, ok := pay2[key]; ok {
			newMap[key] = (val + val2) / 2
		}
	}
	return newMap
}

func averageClientInfos(cl1 map[shared.ClientID]ClientInfo, cl2 map[shared.ClientID]ClientInfo) map[shared.ClientID]ClientInfo {
	newMap := make(map[shared.ClientID]ClientInfo)
	for client, info := range cl1 {
		temp := info.Copy()
		cl2Info := cl2[client]
		temp.Resources = (cl2Info.Resources + info.Resources) / 2
		newMap[client] = temp
	}
	return newMap
}

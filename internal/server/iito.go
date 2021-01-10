package server

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// runIITO : IITO makes recommendations about the optimal (and fairest) contributions this term
// to mitigate the common pool dilemma
func (s *SOMASServer) runIITO() error {
	s.logf("start runIITO")
	defer s.logf("finish runIITO")
	s.gameState.IITOTransactions = s.runGiftSession()

	// This is for sharing an island's intended contributions to the common pool
	s.runIntendedContributionSession()
	// TODO:- IITO team
	return nil
}

func (s *SOMASServer) runIITOEndOfTurn() error {
	s.logf("start runIITOEndOfTurn")
	defer s.logf("finish runIITOEndOfTurn")
	s.executeTransactions(s.gameState.IITOTransactions)
	return nil
}

func (s *SOMASServer) runGiftSession() map[shared.ClientID]shared.GiftResponseDict {
	s.logf("start runGiftSession")
	defer s.logf("finish runGiftSession")

	requests := s.getGiftRequests()
	offers := s.getGiftOffers(requests)
	responses := s.getGiftResponses(offers)

	// Clean all rejected / ignored responses so the remaining are only the transactions
	transactions := s.distributeGiftHistory(responses)
	for key := range transactions {
		for fromTeam, response := range transactions[key] {
			s.logf("[IITO]: Gifts to %v from %v: %v\n", fromTeam, key, response.AcceptedAmount)
			if response.Reason != shared.Accept {
				delete(transactions[key], fromTeam)
			}
		}
	}

	return transactions
}

func (s *SOMASServer) sanitiseTeamGiftRequests(requests shared.GiftRequestDict, thisTeam shared.ClientID) shared.GiftRequestDict {
	for team, request := range requests {
		if s.gameState.ClientInfos[team].LifeStatus == shared.Dead || team == thisTeam || request == 0 {
			delete(requests, team)
			// s.logf("%v violated request conventions. To %v, requested %v", thisTeam, team, request)
		}
	}
	return requests
}

// GetGiftRequests collects a map of gift requests from an individual client, for all clients, in a map
func (s *SOMASServer) getGiftRequests() map[shared.ClientID]shared.GiftRequestDict {
	totalRequests := map[shared.ClientID]shared.GiftRequestDict{}
	for _, id := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		totalRequests[id] = s.sanitiseTeamGiftRequests(s.clientMap[id].GetGiftRequests(), id)

		if len(totalRequests[id]) == 0 {
			delete(totalRequests, id)
		}
	}
	return totalRequests
}

func offersKnapsackSolver(capacity shared.GiftOffer, offers shared.GiftOfferDict) (shared.GiftOffer, []shared.ClientID) {
	// Base case
	if len(offers) == 1 {
		for team, offer := range offers {
			if offer <= capacity {
				return offer, []shared.ClientID{team}
			}
			return 0, []shared.ClientID{}
		}
	}

	bestOffer := shared.GiftOffer(0)
	bestCombination := []shared.ClientID{}

	for team, offer := range offers {
		currentOffer := shared.GiftOffer(0)
		var currentCombi []shared.ClientID

		lessThisOffer := shared.GiftOfferDict{}
		for team, offer := range offers {
			lessThisOffer[team] = offer
		}
		delete(lessThisOffer, team)

		// You might need this if you're debugging <3
		// fmt.Printf("Looping with: %v \n", lessThisOffer)

		if offer <= capacity {
			// Pack and find optimal of remaining capacity
			optimalWithThisOffer, remainingCombination := offersKnapsackSolver(capacity-offer, lessThisOffer)
			currentOffer += (offer + optimalWithThisOffer)
			currentCombi = append(remainingCombination, team)

		} else {
			// Find optimal of remaining capacity
			optimalWithoutThisOffer, remainingCombination := offersKnapsackSolver(capacity, lessThisOffer)
			currentOffer += optimalWithoutThisOffer
			currentCombi = remainingCombination
		}

		if currentOffer > bestOffer {
			bestOffer = currentOffer
			bestCombination = currentCombi
		}
	}

	return bestOffer, bestCombination
}

func (s *SOMASServer) sanitiseTeamGiftOffers(offers shared.GiftOfferDict, thisTeam shared.ClientID) shared.GiftOfferDict {
	totalOffers := shared.GiftOffer(0)
	for team, offer := range offers {
		totalOffers += offer
		if s.gameState.ClientInfos[team].LifeStatus == shared.Dead || team == thisTeam || offer == 0 {
			delete(offers, team)
			// s.logf("%v made an invalid offer", thisTeam)
		}
	}

	// Find the optimal combination of offers if the sum of their offers exceeds their capacity.
	totalResources := shared.GiftOffer(s.gameState.ClientInfos[thisTeam].Resources)
	if totalOffers > totalResources {
		s.logf("[IITO]: Total offerings exceed total resources for %v", thisTeam)
		// Yay, a knapsack problem!
		_, bestCombination := offersKnapsackSolver(totalResources, offers)
		newOffers := shared.GiftOfferDict{}
		for _, team := range bestCombination {
			newOffers[team] = offers[team]
		}
		offers = newOffers
	}

	return offers
}

// getGiftOffers collects all responses from clients to their requests in a map
func (s *SOMASServer) getGiftOffers(totalRequests map[shared.ClientID]shared.GiftRequestDict) map[shared.ClientID]shared.GiftOfferDict {
	totalOffers := map[shared.ClientID]shared.GiftOfferDict{}
	for _, thisTeam := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		// Gather all the requests made to this team
		requestsToThisTeam := shared.GiftRequestDict{}
		for fromTeam, indivRequests := range totalRequests {
			if request, ok := indivRequests[thisTeam]; ok {
				requestsToThisTeam[fromTeam] = request
			}
		}

		offers := s.sanitiseTeamGiftOffers(s.clientMap[thisTeam].GetGiftOffers(requestsToThisTeam), thisTeam)
		if len(offers) > 0 {
			totalOffers[thisTeam] = offers
		}
	}
	return totalOffers
}

func (s *SOMASServer) sanitiseTeamGiftResponses(responses shared.GiftResponseDict, offers shared.GiftOfferDict, thisTeam shared.ClientID) shared.GiftResponseDict {
	for team, response := range responses {
		// If the reason isn't "Accept", the accepted amount should be 0. Otherwise,
		// cap each response so that the an island can't accept more than it was offered.
		if response.Reason != shared.Accept {
			response.AcceptedAmount = 0
			s.logf("[IITO]: %v had a malformed response. Accepted: %v, Reason: %v, Offered: %v ", thisTeam, response.AcceptedAmount, response.Reason, offers[team])
		} else if response.AcceptedAmount > shared.Resources(offers[team]) {
			response.AcceptedAmount = shared.Resources(offers[team])
			s.logf("[IITO]: %v tried to accept more than it was offered. Accepted: %v, Offered: %v ", thisTeam, response.AcceptedAmount, offers[team])
		}
		responses[team] = response

		// Can't respond to an offer that was not given.
		if _, ok := offers[team]; !ok {
			delete(responses, team)
			s.logf("[IITO]: %v tried to accept a non-existent offer. Accepted: %v, Offered: %v ", thisTeam, response.AcceptedAmount, team)
		}
	}
	// Pad the responses so that each offer is responded to, even if ignored.
	for offerFrom := range offers {
		if _, ok := responses[offerFrom]; !ok {
			responses[offerFrom] = shared.GiftResponse{AcceptedAmount: 0.0, Reason: shared.Ignored}
		}
	}
	return responses
}

func (s *SOMASServer) getGiftResponses(totalOffers map[shared.ClientID]shared.GiftOfferDict) map[shared.ClientID]shared.GiftResponseDict {
	totalResponses := map[shared.ClientID]shared.GiftResponseDict{}

	for _, id := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		offersToThisTeam := shared.GiftOfferDict{}
		for fromTeam, indivOffers := range totalOffers {
			if offer, ok := indivOffers[id]; ok {
				offersToThisTeam[fromTeam] = offer
			}
		}
		response := s.sanitiseTeamGiftResponses(s.clientMap[id].GetGiftResponses(offersToThisTeam), offersToThisTeam, id)
		if len(response) > 0 {
			totalResponses[id] = response
		}
	}
	return totalResponses
}

// distributeGiftHistory collates all responses to a single client and calls that client to receive its responses.
func (s *SOMASServer) distributeGiftHistory(totalResponses map[shared.ClientID]shared.GiftResponseDict) map[shared.ClientID]shared.GiftResponseDict {
	s.updateIITOGameState(totalResponses)
	responsesToClients := map[shared.ClientID]shared.GiftResponseDict{}
	for _, id := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		responsesToThisTeam := shared.GiftResponseDict{}
		for fromTeam, indivResponses := range totalResponses {
			if response, ok := indivResponses[id]; ok {
				responsesToThisTeam[fromTeam] = response
			}
		}
		responsesToClients[id] = responsesToThisTeam
		s.clientMap[id].UpdateGiftInfo(responsesToThisTeam)
	}
	return responsesToClients
}

// executeTransactions runs all the accepted responses from the gift session.
// TODO: UNTESTED
func (s *SOMASServer) executeTransactions(transactions map[shared.ClientID]shared.GiftResponseDict) {
	for fromTeam, responses := range transactions {
		for toTeam, indivResponse := range responses {
			giftAmount := s.clientMap[fromTeam].DecideGiftAmount(toTeam, indivResponse.AcceptedAmount)
			if giftAmount < 0 {
				s.logf("[IITO]: Negative resources received in executeTransactions() from %v. Nice Try", fromTeam)
				continue
			}
			transactionMsg := fmt.Sprintf("[IITO]: %v received gift from %v: %v", toTeam, fromTeam, giftAmount)
			errTake := s.takeResources(fromTeam, giftAmount, "TAKE: "+transactionMsg)
			if errTake != nil {
				s.logf("[IITO]: Error deducting amount: %v", errTake)
			} else {
				err := s.giveResources(toTeam, giftAmount, "GIVE: "+transactionMsg)
				if err != nil {
					s.logf("Ignoring failure to give resources in executeTransactions: %v", err)
				}
				s.clientMap[toTeam].ReceivedGift(giftAmount, fromTeam)
				s.clientMap[fromTeam].SentGift(giftAmount, toTeam)
			}
		}
	}
}

func (s *SOMASServer) runIntendedContributionSession() {
	s.logf("start runIntendedContributionSession")
	defer s.logf("finish runIntendedContributionSession")
	islandContributionDict := s.getIntendedContribution()
	s.distributeIntendedContributions(islandContributionDict)
}

func (s *SOMASServer) getIntendedContribution() shared.IntendedContributionDict {
	islandPredictionsDict := shared.IntendedContributionDict{}
	for _, id := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		islandPredictionsDict[id] = s.clientMap[id].ShareIntendedContribution()
	}
	return islandPredictionsDict
}

func (s *SOMASServer) distributeIntendedContributions(islandPredictionDict shared.IntendedContributionDict) {
	reorderDictionary := make(map[shared.ClientID]shared.ReceivedIntendedContributionDict)
	// Add the predictions/sources to the dict containing which predictions each island should receive
	// Don't allow teams to know who else these predictions were shared with in MVP
	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)

	for idSource, info := range islandPredictionDict {
		if clientArrayContains(nonDeadClients, idSource) {
			for _, idShare := range info.TeamsOfferedTo {
				if idShare == idSource {
					continue
				}
				if reorderDictionary[idShare] == nil {
					reorderDictionary[idShare] = make(shared.ReceivedIntendedContributionDict)
				}
				reorderDictionary[idShare][idSource] = shared.ReceivedIntendedContribution{Contribution: info.Contribution, SharedFrom: idSource}
			}
		}
	}

	// Now distribute these predictions to the islands
	for _, id := range nonDeadClients {
		c := s.clientMap[id]
		c.ReceiveIntendedContribution(reorderDictionary[id])
	}
}

func (s *SOMASServer) updateIITOGameState(totalResponses map[shared.ClientID]shared.GiftResponseDict) {
	s.gameState.IITOTransactions = totalResponses
}

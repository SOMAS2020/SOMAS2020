package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// runIITO : IITO makes recommendations about the optimal (and fairest) contributions this term
// to mitigate the common pool dilemma
func (s *SOMASServer) runIITO() error {
	s.logf("start runIITO")
	defer s.logf("finish runIITO")
	err := s.runGiftSession()
	if err != nil {
		return err
	}
	// TODO:- IITO team
	return nil
}

func (s *SOMASServer) runIITOEndOfTurn() error {
	s.logf("start runIITOEndOfTurn")
	defer s.logf("finish runIITOEndOfTurn")
	// TODO:- IITO team
	return nil
}

func (s *SOMASServer) runGiftSession() error {
	s.logf("start runGiftSession")
	defer s.logf("finish runGiftSession")

	giftRequestDict := s.getGiftRequests()
	giftOffersDict, err := s.getGiftOffers(giftRequestDict)
	if err != nil {
		return err
	}
	giftHistoryDict, err := s.getGiftResponses(giftOffersDict)
	if err != nil {
		return err
	}
	err = s.distributeGiftHistory(giftHistoryDict)
	if err != nil {
		return err
	}
	// Process actions
	for key, value := range giftHistoryDict {
		s.logf("Gifts from %s: %v\n", key, value)
	}
	return nil
}

func (s *SOMASServer) sanitiseTeamGiftRequests(requests shared.GiftRequestDict, thisTeam shared.ClientID) (shared.GiftRequestDict, error) {
	for team, request := range requests {
		// Delete the request if it's to yourself or a dead team.
		if s.gameState.ClientInfos[team].LifeStatus == shared.Dead || team == thisTeam || request == 0 {
			delete(requests, team)
		}
	}

	// TODO: Maybe return some kind of helpful message if any of the above cases break?
	return requests, nil
}

// GetGiftRequests collects a map of gift requests from an individual client, for all clients, in a map
func (s *SOMASServer) getGiftRequests() map[shared.ClientID]shared.GiftRequestDict {
	totalRequests := map[shared.ClientID]shared.GiftRequestDict{}
	for _, id := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		totalRequests[id], _ = s.sanitiseTeamGiftRequests(s.clientMap[id].GetGiftRequests(), id)

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
		currentCombi := []shared.ClientID{}

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

// TODO: Return an error?
func (s *SOMASServer) sanitiseTeamGiftOffers(offers shared.GiftOfferDict, thisTeam shared.ClientID) (shared.GiftOfferDict, error) {
	totalOffers := shared.GiftOffer(0)
	for team, offer := range offers {
		totalOffers += offer
		if s.gameState.ClientInfos[team].LifeStatus == shared.Dead || team == thisTeam || offer == 0 {
			delete(offers, team)
		}
	}

	// Find the optimal combination of offers if the sum of their offers exceeds their capacity.
	totalResources := shared.GiftOffer(s.gameState.ClientInfos[thisTeam].Resources)
	if totalOffers > totalResources {
		// Yay, a knapsack problem!
		_, bestCombination := offersKnapsackSolver(totalResources, offers)
		newOffers := shared.GiftOfferDict{}
		for _, team := range bestCombination {
			newOffers[team] = offers[team]
		}
		offers = newOffers
	}

	// TODO: Maybe return some kind of helpful message if any of the above cases break?
	return offers, nil
}

// getGiftOffers collects all responses from clients to their requests in a map
func (s *SOMASServer) getGiftOffers(totalRequests map[shared.ClientID]shared.GiftRequestDict) (map[shared.ClientID]shared.GiftOfferDict, error) {
	totalOffers := map[shared.ClientID]shared.GiftOfferDict{}
	for _, thisTeam := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		// Gather all the requests made to this team
		requestsToThisTeam := shared.GiftRequestDict{}
		for fromTeam, indivRequests := range totalRequests {
			if request, ok := indivRequests[thisTeam]; ok {
				requestsToThisTeam[fromTeam] = request
			}
		}

		offers, err := s.sanitiseTeamGiftOffers(s.clientMap[thisTeam].GetGiftOffers(requestsToThisTeam), thisTeam)
		if err != nil {
			// totalResponses in this case may have a bogus or meaningless value
			return totalOffers, err
		}

		if len(offers) > 0 {
			totalOffers[thisTeam] = offers
		}
	}
	return totalOffers, nil
}

func (s *SOMASServer) sanitiseTeamGiftResponses(responses shared.GiftResponseDict, offers shared.GiftOfferDict) (shared.GiftResponseDict, error) {
	for team, response := range responses {
		// Cap each response so that the an island can't accept more than it was offered.
		if response.AcceptedAmount > shared.Resources(offers[team]) {
			response.AcceptedAmount = shared.Resources(offers[team])
			responses[team] = response
		}

		// Can't respond to an offer that was not given.
		if _, ok := offers[team]; !ok {
			delete(responses, team)
		}
	}
	// TODO: Pad the responses so that each offer is responded to, even if ignored.
	// TODO: Maybe return some kind of helpful message if any of the above cases break?
	return responses, nil
}

func (s *SOMASServer) getGiftResponses(totalOffers map[shared.ClientID]shared.GiftOfferDict) (map[shared.ClientID]shared.GiftResponseDict, error) {
	totalResponses := map[shared.ClientID]shared.GiftResponseDict{}

	for _, id := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		offersToThisTeam := shared.GiftOfferDict{}
		for fromTeam, indivOffers := range totalOffers {
			if offer, ok := indivOffers[id]; ok {
				offersToThisTeam[fromTeam] = offer
			}
		}
		response, err := s.sanitiseTeamGiftResponses(s.clientMap[id].GetGiftResponses(offersToThisTeam), offersToThisTeam)
		if err != nil {
			// totalResponses in this case may have a bogus or meaningless value
			return totalResponses, err
		}
		totalResponses[id] = response
	}
	return totalResponses, nil
}

// distributeGiftHistory collates all responses to a single client and calls that client to receive its responses.
func (s *SOMASServer) distributeGiftHistory(totalResponses map[shared.ClientID]shared.GiftResponseDict) error {
	for _, id := range getNonDeadClientIDs(s.gameState.ClientInfos) {
		responsesToThisTeam := shared.GiftResponseDict{}
		for fromTeam, indivResponses := range totalResponses {
			if response, ok := indivResponses[id]; ok {
				responsesToThisTeam[fromTeam] = response
			}
		}
		err := s.clientMap[id].UpdateGiftInfo(responsesToThisTeam)
		if err != nil {
			return err
		}
	}
	return nil
}

package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

// runIIFO : IIFO allows sharing of disaster predictions between islands
func (s *SOMASServer) runIIFO() error {
	s.logf("start runIIFO")
	defer s.logf("finish runIIFO")

	// This is for Disaster prediction
	err := s.runPredictionSession()
	if err != nil {
		return errors.Errorf("Error running prediction session: %v", err)
	}

	s.runForageSharing()

	// TODO:- IIFO team
	return nil
}

func (s *SOMASServer) runIIFOEndOfTurn() error {
	s.logf("start runIIFOEndOfTurn")
	defer s.logf("finish runIIFOEndOfTurn")
	// TODO:- IIFO team
	return nil
}

func (s *SOMASServer) runPredictionSession() error {
	s.logf("start runPredictionSession")
	defer s.logf("finish runPredictionSession")
	islandPredictionDict, err := s.getPredictions()
	if err != nil {
		return err
	}

	err = s.distributePredictions(islandPredictionDict)
	if err != nil {
		return err
	}

	return nil
}

func (s *SOMASServer) runForageSharing() {
	s.logf("Run Forage Predictions")
	defer s.logf("Finish Running Forage Predictions")

	otherIslandInfo := s.getForageSharing()
	s.distributeForageSharing(otherIslandInfo)
}

func (s *SOMASServer) getPredictions() (shared.PredictionInfoDict, error) {
	islandPredictionsDict := shared.PredictionInfoDict{}
	var err error
	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)
	for _, id := range nonDeadClients {
		c := s.clientMap[id]
		islandPredictionsDict[id], err = c.MakePrediction()
		if err != nil {
			return islandPredictionsDict, errors.Errorf("Failed to get prediction from %v: %v", id, err)
		}
	}
	return islandPredictionsDict, nil
}

func (s *SOMASServer) distributePredictions(islandPredictionDict shared.PredictionInfoDict) error {
	receivedPredictionsDict := make(shared.ReceivedPredictionsDict)
	var err error
	// Add the predictions/sources to the dict containing which predictions each island should receive
	// Don't allow teams to know who else these predictions were shared with in MVP
	for idSource, info := range islandPredictionDict {
		for _, idShare := range info.TeamsOfferedTo {
			if idShare == idSource {
				continue
			}
			if receivedPredictionsDict[idShare] == nil {
				receivedPredictionsDict[idShare] = make(shared.PredictionInfoDict)
			}
			info.TeamsOfferedTo = nil
			receivedPredictionsDict[idShare][idSource] = info
		}
	}

	// Now distribute these predictions to the islands
	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)
	for _, id := range nonDeadClients {
		c := s.clientMap[id]
		err = c.ReceivePredictions(receivedPredictionsDict[id])
		if err != nil {
			return errors.Errorf("Failed to receive prediction from client %v: %v", id, err)
		}
	}

	return nil
}

// getForageSharing will ping each nonDeadClient and will save what their ForagingDecision,
// their ResourceObtained from that decision, and which ClientID they want share this info
// with.
func (s *SOMASServer) getForageSharing() shared.ForagingOfferDict {
	s.logf("Getting Forage Information")
	islandShareForageDict := shared.ForagingOfferDict{}
	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)
	for _, id := range nonDeadClients {
		c := s.clientMap[id]
		islandShareForageDict[id] = c.MakeForageInfo()
	}
	return islandShareForageDict
}

// distributeForageSharing sends the collected ForageDecisions and ResourceObtained to specified ClientID
func (s *SOMASServer) distributeForageSharing(otherIslandInfo shared.ForagingOfferDict) {
	s.logf("Distributing Forage Information")
	islandForagingDict := shared.ForagingReceiptDict{}
	for islandID, foragingInfo := range otherIslandInfo {
		for _, shareID := range foragingInfo.ShareTo {
			if islandID == shareID {
				continue
			}
			islandForagingDict[shareID] = append(
				islandForagingDict[shareID],
				shared.ForageShareInfo{
					DecisionMade:     foragingInfo.DecisionMade,
					ResourceObtained: foragingInfo.ResourceObtained,
					SharedFrom:       islandID})
		}
	}
	nonDeadClients := getNonDeadClientIDs(s.gameState.ClientInfos)
	for _, id := range nonDeadClients {
		c := s.clientMap[id]

		c.ReceiveForageInfo(islandForagingDict[id])
	}
}

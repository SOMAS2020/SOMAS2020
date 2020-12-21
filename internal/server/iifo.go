package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// runIIFO : IIFO allows sharing of disaster predictions between islands
func (s *SOMASServer) runIIFO() error {
	s.logf("start runIIFO")
	s.runPredictionSession()
	defer s.logf("finish runIIFO")
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
	islandPredictionDict, err := s.getPredictions()
	if err != nil {
		return err
	}

	err = s.distributePredictions(islandPredictionDict)
	if err != nil {
		return err
	}

	defer s.logf("finish runPredictionSession")
	return nil
}

func (s *SOMASServer) getPredictions() (shared.PredictionInfoDict, error) {
	islandPredictionsDict := shared.PredictionInfoDict{}
	tempPredictionInfo := shared.PredictionInfo{}
	var err error
	for id, client := range s.clientMap {
		tempPredictionInfo, err = client.MakePrediction()
		islandPredictionsDict[id] = tempPredictionInfo
		if err != nil {
			return islandPredictionsDict, err
		}
	}
	return islandPredictionsDict, nil
}
func (s *SOMASServer) distributePredictions(islandPredictionDict shared.PredictionInfoDict) error {
	recievedPredictionsDict := make(shared.RecievedPredictionsDict)
	var err error
	// Add the predictions/sources to the dict containing which predictions each island should recieve
	// Don't allow teams to know who else these predictions were shared with in MVP
	for idSource, info := range islandPredictionDict {
		for _, idShare := range info.TeamsOfferedTo {
			if idShare == idSource {
				continue
			}
			if recievedPredictionsDict[idShare] == nil {
				recievedPredictionsDict[idShare] = make(shared.PredictionInfoDict)
			}
			info.TeamsOfferedTo = nil
			recievedPredictionsDict[idShare][idSource] = info
		}
	}

	// Now distribute these predictions to the islands
	for id, client := range s.clientMap {
		err = client.RecievePredictions(recievedPredictionsDict[id])
		if err != nil {
			return err
		}
	}
	return nil
}

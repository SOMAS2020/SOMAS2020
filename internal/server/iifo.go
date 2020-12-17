package server

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// runIIFO : IIFO allows sharing of disaster predictions between islands
func (s *SOMASServer) runIIFO() ([]common.Action, error) {
	s.logf("start runIIFO")
	s.runPredictionSession()
	defer s.logf("finish runIIFO")
	// TODO:- IIFO team
	return nil, nil
}

func (s *SOMASServer) runPredictionSession() ([]common.Action, error) {
	s.logf("start runPredictionSession")
	islandPredictionDict, err := s.getPredictions()
	if err != nil {
		return nil, err
	}

	err = s.distributePredictions(islandPredictionDict)
	if err != nil {
		return nil, err
	}

	defer s.logf("finish runPredictionSession")
	return nil, nil
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
	// Add the predictions/sources to the dict which determines which predictions each island should recieve
	// Don't allow teams to know who these predictions were shared with in MVP
	for idSource, info := range islandPredictionDict {
		for _, idShare := range info.TeamsOfferedTo {
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

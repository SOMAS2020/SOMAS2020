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

func (s *SOMASServer) getPredictions() (shared.IslandPredictionDict, error) {
	islandPredictionsDict := shared.IslandPredictionDict{}
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
func (s *SOMASServer) distributePredictions(islandPredictionDict shared.IslandPredictionDict) error {
	recievedPredictionsDict := shared.RecievedPredictionsDict{}
	tempPredictionSlice := []shared.Prediction{}
	tempClientIDSlice := make([]shared.ClientID, len(s.clientMap))
	var err error

	// Add the predictions/sources to the dict which determines which predictions each island should recieve
	for idSource, info := range islandPredictionDict {
		for _, idShare := range info.TeamsOfferedTo {
			//s.logf("Gifts from %v\n", &recievedPredictionsDict[idShare].Predictions)
			tempPredictionSlice = append(recievedPredictionsDict[idShare].Predictions, islandPredictionDict[idSource].PredictionMade)
			tempClientIDSlice = append(recievedPredictionsDict[idShare].SourceIslands, idSource)
			recievedPredictionsDict[idShare].Predictions = tempPredictionSlice
			recievedPredictionsDict[idShare].SourceIslands = tempClientIDSlice
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

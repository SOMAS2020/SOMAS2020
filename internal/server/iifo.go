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
	var err error
	for id, client := range s.clientMap {
		islandPredictionsDict[id].PredictionMade, islandPredictionsDict[id].TeamsOfferedTo, err = client.MakePrediction()
		if err != nil {
			return islandPredictionsDict, err
		}
	}
	return islandPredictionsDict, nil
}
func (s *SOMASServer) distributePredictions(islandPredictionDict shared.IslandPredictionDict) error {
	recievedPredictionsDict := shared.RecievedPredictionsDict{}
	var err error

	// Add the predictions/sources to the dict which determines which predictions each island should recieve
	for idSource, info := range islandPredictionDict {
		for _, idShare := range info.TeamsOfferedTo {
			recievedPredictionsDict[idShare].Predictions = append(recievedPredictionsDict[idShare].Predictions, islandPredictionDict[idSource].PredictionMade)
			recievedPredictionsDict[idShare].SourceIslands = append(recievedPredictionsDict[idShare].SourceIslands, idSource)
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

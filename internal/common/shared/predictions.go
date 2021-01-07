package shared

// DisasterPrediction is a struct containing the necessary parameters for an island to
// make a prediction about a disaster
type DisasterPrediction struct {
	CoordinateX Coordinate
	CoordinateY Coordinate
	Magnitude   Magnitude
	TimeLeft    uint
	Confidence  PredictionConfidence
}

// PredictionConfidence is bounded between 0 and 100
type PredictionConfidence = float64

// DisasterPredictionInfo is a struct containing the information describing a prediction
// made by a given island
type DisasterPredictionInfo struct {
	PredictionMade DisasterPrediction
	TeamsOfferedTo []ClientID
}

// ReceivedDisasterPredictionInfo is a struct containing the information describing a prediction
// made by a given island
type ReceivedDisasterPredictionInfo struct {
	PredictionMade DisasterPrediction
	SharedFrom     ClientID
}

// DisasterPredictionInfoDict is a dictionary of DisasterPredictionInfo
type DisasterPredictionInfoDict = map[ClientID]DisasterPredictionInfo

// ReceivedDisasterPredictionsDict is a dictionary of ReceivedDisasterPredictionInfo
type ReceivedDisasterPredictionsDict = map[ClientID]ReceivedDisasterPredictionInfo

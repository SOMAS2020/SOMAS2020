package shared

// Prediction is a struct containing the necessary parameters for an island to
// make a prediction about a disaster
type Prediction struct {
	CoordinateX float64
	CoordinateY float64
	Magnitude   float64
	TimeLeft    int
	Confidence  float64
}

// PredictionInfo is a struct containing the information describing a prediction
// made by a given island
type PredictionInfo struct {
	PredictionMade Prediction
	TeamsOfferedTo []ClientID
}

// PredictionInfoDict is a dictionary of PredictionInfo
type PredictionInfoDict = map[ClientID]PredictionInfo

// RecievedPredictionsDict is a dictionary of PredictionInfoDict
type RecievedPredictionsDict = map[ClientID]PredictionInfoDict

package shared

// Prediction is a struct containing the necessary parameters for an island to
// make a prediction about a disaster
type Prediction struct {
	CoordinateX float32
	CoordinateY float32
	Magnitude   int
	TimeLeft    int
	Confidence  int
}

// PredictionInfo is a struct containing the information describing a prediction
// made by a given island
type PredictionInfo struct {
	PredictionMade Prediction
	TeamsOfferedTo []ClientID
}

// IslandPredictionDict is a dictonary of a pointer to PredictionInfo
type IslandPredictionDict = map[ClientID]PredictionInfo

// RecievedPredictions is a struct containing the predictions an island should recieve and
// the IDs of the islands that these predictions came from
type RecievedPredictions struct {
	Predictions   []Prediction
	SourceIslands []ClientID
}

// RecievedPredictionsDict is a dictonary of pointer to RecievedPredictions
type RecievedPredictionsDict = map[ClientID]RecievedPredictions

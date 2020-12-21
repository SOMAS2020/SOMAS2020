// Package baseclient contains the Client interface as well as a base client implementation.
package baseclient

import (
	"fmt"
	"log"
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Assume that after each disaster, the disaster location and magnitude are stored for use in this function.
// These disasters will be stored in PastDisastersDict (maps round number to disaster that occurred)
type Disaster struct {
	CoordinateX float64
	CoordinateY float64
	Magnitude   int
	TurnNumber  int
}
type PastDisastersDict = map[int]Disaster

var ourPredictionInfo shared.PredictionInfo

// NewClient produces a new client with the BaseClient already implemented.
// BASE: Do not overwrite in team client.
func NewClient(id shared.ClientID) Client {
	return &BaseClient{id: id}
}

// BaseClient provides a basic implementation for all functions of the client interface and should always the interface fully.
// All clients should be based off of this BaseClient to ensure that all clients implement the interface,
// even when new features are added.
type BaseClient struct {
	id              shared.ClientID
	clientGameState gamestate.ClientGameState
}

// Echo prints a message to show that the client exists
// BASE: Do not overwrite in team client.
func (c *BaseClient) Echo(s string) string {
	c.Logf("Echo: '%v'", s)
	return s
}

// GetID gets returns the id of the client.
// BASE: Do not overwrite in team client.
func (c *BaseClient) GetID() shared.ClientID {
	return c.id
}

// Logf is the client's logger that prepends logs with your ID. This makes
// it easier to read logs. DO NOT use other loggers that will mess logs up!
// BASE: Do not overwrite in team client.
func (c *BaseClient) Logf(format string, a ...interface{}) {
	log.Printf("[%v]: %v", c.id, fmt.Sprintf(format, a...))
}

// StartOfTurnUpdate is updates the gamestate of the client at the start of each turn.
// The gameState is served by the server.
// OPTIONAL. Base should be able to handle it but feel free to implement your own.
func (c *BaseClient) StartOfTurnUpdate(gameState gamestate.ClientGameState) {
	c.Logf("Received start of turn game state update: %#v", gameState)
	c.clientGameState = gameState
	// TODO
}

// GameStateUpdate updates game state mid-turn.
// The gameState is served by the server.
// OPTIONAL. Base should be able to handle it but feel free to implement your own.
func (c *BaseClient) GameStateUpdate(gameState gamestate.ClientGameState) {
	c.Logf("Received game state update: %#v", gameState)
	c.clientGameState = gameState
}

// RequestGift allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
// COMPULSORY, you need to implement this method
func (c *BaseClient) RequestGift() (int, error) {
	return 0, nil
}

// OfferGifts allows clients to offer to give the gifts requested by other clients.
// It can offer multiple partial gifts
// COMPULSORY, you need to implement this method
func (c *BaseClient) OfferGifts(giftRequestDict shared.GiftDict) (shared.GiftDict, error) {
	return giftRequestDict, nil
}

// AcceptGifts allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
// COMPULSORY, you need to implement this method
func (c *BaseClient) AcceptGifts(receivedGiftDict shared.GiftDict) (shared.GiftInfoDict, error) {
	acceptedGifts := shared.GiftInfoDict{}
	for client, offer := range receivedGiftDict {
		acceptedGifts[client] = shared.GiftInfo{
			ReceivingTeam:  client,
			OfferingTeam:   c.GetID(),
			OfferAmount:    offer,
			AcceptedAmount: offer,
			Reason:         shared.Accept}
	}
	return acceptedGifts, nil
}

// UpdateGiftInfo gives information about the outcome from AcceptGifts.
// This allows for opinion formation.
// COMPULSORY, you need to implement this method
func (c *BaseClient) UpdateGiftInfo(acceptedGifts shared.GiftInfoDict) error {

	return nil
}

//Actions? Need to talk to LH and our team about this one:

// SendGift is executed at the end of each turn and allows clients to
// send the gifts promised in the IITO
// COMPULSORY, you need to implement this method
func (c *BaseClient) SendGift(receivingClient shared.ClientID, amount int) error {
	return nil
}

// ReceiveGift is executed at the end of each turn and allows clients to
// receive the gifts promised in the IITO
// COMPULSORY, you need to implement this method
func (c *BaseClient) ReceiveGift(sendingClient shared.ClientID, amount int) error {
	return nil
}

// MakePrediction is called on each client for them to make a prediction about a disaster
// Prediction includes location, magnitude, confidence etc
// COMPULSORY, you need to implement this method
func (c *BaseClient) MakePrediction() (shared.PredictionInfo, error) {
	// Set up dummy disasters for testing purposes
	pastDisastersDict := make(PastDisastersDict)
	tempDisaster := Disaster{}
	for i := 0; i < 4; i++ {
		tempDisaster.CoordinateX = float64(i)
		tempDisaster.CoordinateY = float64(i)
		tempDisaster.Magnitude = i
		tempDisaster.TurnNumber = i
		pastDisastersDict[i] = tempDisaster
	}

	// Use the sample mean of each field as our prediction
	meanDisaster, err := getMeanDisaster(pastDisastersDict)
	if err != nil {
		return shared.PredictionInfo{}, err
	}

	prediction := shared.Prediction{
		CoordinateX: meanDisaster.CoordinateX,
		CoordinateY: meanDisaster.CoordinateY,
		Magnitude:   meanDisaster.Magnitude,
		TimeLeft:    meanDisaster.TurnNumber,
	}

	// Use (variance limit - mean(sample variance)), where the mean is taken over each field, as confidence
	// Use a variance limit of 100 for now
	varianceLimit := 100.0
	confidence, err := determineConfidence(pastDisastersDict, meanDisaster, varianceLimit)
	if err != nil {
		return shared.PredictionInfo{}, err
	}
	prediction.Confidence = confidence

	/*
		if c.GetID() == 1 {
			c.Logf("CoordinateX Prediction:[%v]\n", prediction.CoordinateX)
			c.Logf("CoordinateY Prediction:[%v]\n", prediction.CoordinateY)
			c.Logf("Magnitude Prediction:[%v]\n", prediction.Magnitude)
			c.Logf("Time Left Prediction:[%v]\n", prediction.TimeLeft)
			c.Logf("Confidence:[%v]\n", prediction.Confidence)
		}*/

	// For MVP, share this prediction with all islands since trust has not yet been implemented
	trustedIslands := make([]shared.ClientID, 6)
	for index, id := range shared.TeamIDs {
		trustedIslands[index] = id
	}

	// Return all prediction info and store our own island's prediction in global variable
	predictionInfo := shared.PredictionInfo{
		PredictionMade: prediction,
		TeamsOfferedTo: trustedIslands,
	}
	ourPredictionInfo = predictionInfo
	return predictionInfo, nil
}

func getMeanDisaster(pastDisastersDict PastDisastersDict) (Disaster, error) {
	totalCoordinateX, totalCoordinateY, totalMagnitude, totalTurnNumber := 0.0, 0.0, 0.0, 0.0
	numberDisastersPassed := float64(len(pastDisastersDict))

	for _, disaster := range pastDisastersDict {
		totalCoordinateX += disaster.CoordinateX
		totalCoordinateY += disaster.CoordinateY
		totalMagnitude += float64(disaster.Magnitude)
		totalTurnNumber += float64(disaster.TurnNumber)
	}

	meanDisaster := Disaster{
		CoordinateX: totalCoordinateX / numberDisastersPassed,
		CoordinateY: totalCoordinateY / numberDisastersPassed,
		Magnitude:   int(totalMagnitude/numberDisastersPassed + 0.5),
		TurnNumber:  int(totalTurnNumber/numberDisastersPassed + 0.5),
	}
	return meanDisaster, nil
}
func determineConfidence(pastDisastersDict PastDisastersDict, meanDisaster Disaster, varianceLimit float64) (int, error) {
	totalCoordinateX, totalCoordinateY, totalMagnitude, totalTurnNumber := 0.0, 0.0, 0.0, 0.0
	numberDisastersPassed := float64(len(pastDisastersDict))

	// Find the sum of the square of the difference between the actual and mean, for each field
	for _, disaster := range pastDisastersDict {
		totalCoordinateX += math.Pow(disaster.CoordinateX-meanDisaster.CoordinateX, 2)
		totalCoordinateY += math.Pow(disaster.CoordinateY-meanDisaster.CoordinateY, 2)
		totalMagnitude += math.Pow(float64(disaster.Magnitude-meanDisaster.Magnitude), 2)
		totalTurnNumber += math.Pow(float64(disaster.TurnNumber-meanDisaster.TurnNumber), 2)
	}

	// Find the sum of the variances and the average variance
	varianceSum := (totalCoordinateX + totalCoordinateY + totalMagnitude + totalTurnNumber) / numberDisastersPassed
	averageVariance := varianceSum / 4

	// Implement the variance cap chosen
	if averageVariance > varianceLimit {
		averageVariance = varianceLimit
	}

	// Return the confidence of the prediction
	return int(varianceLimit - averageVariance + 0.5), nil
}

// RecievePredictions provides each client with the prediction info, in addition to the source island,
// that they have been granted access to see
// COMPULSORY, you need to implement this method
func (c *BaseClient) RecievePredictions(recievedPredictions shared.PredictionInfoDict) error {
	// If we assume that we trust each island equally (including ourselves), then take the final prediction
	// of disaster as being the weighted mean of predictions according to confidence
	numberOfPredictions := float64(len(recievedPredictions) + 1)
	selfConfidence := ourPredictionInfo.PredictionMade.Confidence

	// Initialise running totals using our own island's predictions
	totalCoordinateX := float64(selfConfidence) * ourPredictionInfo.PredictionMade.CoordinateX
	totalCoordinateY := float64(selfConfidence) * ourPredictionInfo.PredictionMade.CoordinateY
	totalMagnitude := selfConfidence * ourPredictionInfo.PredictionMade.Magnitude
	totalTimeLeft := selfConfidence * ourPredictionInfo.PredictionMade.TimeLeft
	totalConfidence := float64(selfConfidence)

	// Add other island's predictions using their confidence values
	for _, prediction := range recievedPredictions {
		totalCoordinateX += float64(prediction.PredictionMade.Confidence) * prediction.PredictionMade.CoordinateX
		totalCoordinateY += float64(prediction.PredictionMade.Confidence) * prediction.PredictionMade.CoordinateY
		totalMagnitude += prediction.PredictionMade.Confidence * prediction.PredictionMade.Magnitude
		totalTimeLeft += prediction.PredictionMade.Confidence * prediction.PredictionMade.TimeLeft
		totalConfidence += float64(prediction.PredictionMade.Confidence)
	}

	// Finally get the final prediction generated by considering predictions from all islands that we have available
	// This result is currently unused but would be used in decision making in full implementation
	finalPrediction := shared.Prediction{
		CoordinateX: totalCoordinateX / totalConfidence,
		CoordinateY: totalCoordinateY / totalConfidence,
		Magnitude:   int((float64(totalMagnitude) / totalConfidence) + 0.5),
		TimeLeft:    int((float64(totalTimeLeft) / totalConfidence) + 0.5),
		Confidence:  int((totalConfidence / numberOfPredictions) + 0.5),
	}

	c.Logf("Final Prediction: [%v]\n", finalPrediction)
	return nil
}

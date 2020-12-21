package gamestate

// ClientGameState contains game state only for a specific client.
type ClientGameState struct {
	// Season represents the current (1-index) season of the game.
	Season uint

	// Turn represents the current (1-index) Turn of the game.
	Turn uint

	// ClientInfo
	ClientInfo ClientInfo
}

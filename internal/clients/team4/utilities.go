package team4

func (c *client) getTurn() uint {
	return c.ServerReadHandle.GetGameState().Turn
}

func (c *client) getSeason() uint {
	return c.ServerReadHandle.GetGameState().Season
}

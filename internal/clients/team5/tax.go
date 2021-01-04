package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) GetTaxContribution() shared.Resources {
	c.Logf("Current tax amount: %v", c.taxAmount)
	if c.wealth() == ImperialStudent || c.wealth() == Dying {
		return 0
	}
	if c.gameState().Turn == 1 {
		return 0.5 * c.taxAmount
	}
	return c.taxAmount
}

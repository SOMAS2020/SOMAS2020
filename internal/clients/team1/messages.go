package team1

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

func (c *client) ReceiveCommunication(
	sender shared.ClientID,
	data map[shared.CommunicationFieldName]shared.CommunicationContent,
) {
	for field, content := range data {
		switch field {
		case shared.TaxAmount:
			// if !__isPresident__(sender) {
			// }
			c.taxAmount = shared.Resources(content.IntegerData)
		case shared.AllocationAmount:
			// if !__isPresident__(sender) {
			// }
			c.allocation = shared.Resources(content.IntegerData)
		}
	}
}

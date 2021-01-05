package baseclient

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
)

// DisasterNotification is an event handler for disasters. Server will notify client
// of the ramifications of a disaster via this method.
// OPTIONAL: Use this method for any tasks you want to happen when a disaster occurs
func (c *BaseClient) DisasterNotification(
	dR disasters.DisasterReport,
	effects disasters.DisasterEffects) { // effects contain abs magnitude, prop. mag relative to other islands and CP mitigated mag.
}

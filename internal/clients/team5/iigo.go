package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)
//MonitorIIGORole decides whether to perform monitoring on a role
//COMPULOSRY: must be implemented
//always monitor a role
func (c *client) MonitorIIGORole(roleName shared.Role) bool {
	return true 
}

//DecideIIGOMonitoringAnnouncement decides whether to share the result of monitoring a role and what result to share
//COMPULSORY: must be implemented
// always broadcast monitoring result
func (c *client) DecideIIGOMonitoringAnnouncement(monitoringResult bool) (resultToShare bool, announce bool) {
	resultToShare = monitoringResult
	announce = true
	return
}

// Cannot evade sanction :c
func (c *client) GetSanctionPayment() shared.Resources {
	return c.sanctionAmount
}
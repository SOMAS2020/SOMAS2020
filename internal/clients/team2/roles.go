package team2

import (
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	return &Speaker{c: c, BaseSpeaker: &baseclient.BaseSpeaker{GameState: c.ServerReadHandle.GetGameState()}}
}

func (c *client) GetClientJudgePointer() roles.Judge {
	return &Judge{c: c, BaseJudge: &baseclient.BaseJudge{GameState: c.ServerReadHandle.GetGameState()}}
}

func (c *client) GetClientPresidentPointer() roles.President {
	return &President{c: c, BasePresident: &baseclient.BasePresident{GameState: c.ServerReadHandle.GetGameState()}}
}

func (c *client) VoteForElection(roleToElect shared.Role, candidateList []shared.ClientID) []shared.ClientID {
	var situation Situation
	switch roleToElect {
	case shared.President:
		situation = "President"
	case shared.Judge:
		situation = "Judge"
	default:
		situation = "Gifts"
	}

	var trustRank IslandTrustList
	for _, candidate := range candidateList {
		islandConf := IslandTrust{
			island: candidate,
			trust:  c.confidence(situation, candidate),
		}
		trustRank = append(trustRank, islandConf)
	}

	if situation == "Judge" {
		c.updateJudgeTrust()
	}

	sort.Sort(trustRank)
	bordaList := make([]shared.ClientID, 0)

	for _, islandTrust := range trustRank {
		bordaList = append(bordaList, islandTrust.island)
	}

	return bordaList
}

//MonitorIIGORole decides whether to perform monitoring on a role
//COMPULOSRY: must be implemented
// TODO: This function is not implemented
func (c *client) MonitorIIGORole(roleName shared.Role) bool {
	return true
}

// DecideIIGOMonitoringAnnouncement decides whether to share the result of monitoring a role and what result to share
// COMPULSORY: must be implemented
// TODO: This function is not implemented
func (c *client) DecideIIGOMonitoringAnnouncement(monitoringResult bool) (resultToShare bool, announce bool) {
	resultToShare = monitoringResult
	announce = true
	return
}

// Returns a ResourceReport - if Agent Strategy is selfish lies about resources
// Otherwise accurately shares resources
func (c *client) ResourceReport() shared.ResourcesReport {
	mood := c.getAgentStrategy()
	switch mood {
	case Selfish:
		return shared.ResourcesReport{
			ReportedAmount: 0.5 * c.gameState().ClientInfo.Resources,
			Reported:       true,
		}
	default:
		return shared.ResourcesReport{
			ReportedAmount: c.gameState().ClientInfo.Resources,
			Reported:       true,
		}
	}
}

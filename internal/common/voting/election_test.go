package voting

import (
	"fmt"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestElection(t *testing.T) {
	var ele Election
	ele.roleToElect = 0
	ele.votingMethod = 0
	ele.candidateList = []shared.ClientID{shared.Team1, shared.Team2, shared.Team3, shared.Team4, shared.Team5, shared.Team6}
	ele.voterList = ele.candidateList
	ele.votes = [][]shared.ClientID{{shared.Team3, shared.Team2, shared.Team1, shared.Team6, shared.Team5, shared.Team4},
		{shared.Team4, shared.Team6, shared.Team5, shared.Team2, shared.Team3, shared.Team1},
		{shared.Team3, shared.Team6, shared.Team1, shared.Team4, shared.Team5, shared.Team2},
		{shared.Team2, shared.Team5, shared.Team6, shared.Team4, shared.Team3, shared.Team1},
		{shared.Team6, shared.Team4, shared.Team1, shared.Team5, shared.Team2, shared.Team3},
		{shared.Team5, shared.Team2, shared.Team3, shared.Team6, shared.Team1, shared.Team4}}
	//fmt.Println("Initializing struct Election done. ")
	//clientMap := make(map[shared.ClientID]baseclient.Client)

	//ele.bordaCountResult()
	//ele.runOffResult(clientMap)
	//ele.instantRunoffResult(clientMap)
	//ele.approvalResult()

}

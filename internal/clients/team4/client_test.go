package team4

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestGetJudgePointer(t *testing.T) {
	c := client{
		BaseClient:    baseclient.NewClient(id),
		clientJudge:   judge{BaseJudge: &baseclient.BaseJudge{}, t: t},
		clientSpeaker: speaker{},
	}

	j := c.GetClientJudgePointer()
	winner := j.DecideNextPresident(shared.Team1)

	if winner != shared.Team1 {
		t.Errorf("Got wrong judge pointer. Winner is %v", winner)
	}

}

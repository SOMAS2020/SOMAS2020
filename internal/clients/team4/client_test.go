package team4

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestGetJudgePointer(t *testing.T) {
	c := client{
		BaseClient:    baseclient.NewClient(id),
		clientJudge:   Judge{BaseJudge: &baseclient.BaseJudge{}, t: t},
		clientSpeaker: ourSpeaker{},
	}

	judge := c.GetClientJudgePointer()
	winner := judge.DecideNextPresident(shared.Team1)

	if winner != id {
		t.Errorf("Got wrong judge pointer. Winner is %v", winner)
	}

}

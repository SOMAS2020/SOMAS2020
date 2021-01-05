package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
)

type judge struct {
    *baseclient.BaseJudge
    //c *client comment out right now because we'll just use the baseSpeaker implementation 
}

func (c *client) GetClientJudgePointer() roles.Judge {
	//c.Logf("Team 5 became judge")
	return &c.team5Judge
}
package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
)

type speaker struct {
    *baseclient.BaseSpeaker
    //c *client comment out right now because we'll just use the baseSpeaker implementation 
}

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	c.Logf("Team 5 became speaker")
	return &c.team5Speaker
}
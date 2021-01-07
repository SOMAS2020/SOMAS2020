package team2

import "github.com/SOMAS2020/SOMAS2020/internal/common/roles"

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	return &c.currSpeaker
}

func (c *client) GetClientJudgePointer() roles.Judge {
	return &c.currJudge
}

func (c *client) GetClientPresidentPointer() roles.President {
	return &c.currPresident
}

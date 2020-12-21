package shared

import "math/rand"

type baseCommonPool struct {
	amount int
}

func (cp *baseCommonPool) checkEnoughInCommonPool() int {
	return rand.Intn(1)
}

func (cp *baseCommonPool) withdrawFromCommonPool() int {
	return rand.Intn(1000)
}

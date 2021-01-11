package dynamics

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

type Account struct {
	Id           shared.ClientID
	TargetVal    shared.Resources
	Coeff        float64
	personalPool shared.Resources
	turnTarget   shared.Resources
}

func (a *Account) UpdatePersonalPool(val shared.Resources) {
	a.personalPool = val
}

func (a *Account) Cycle() {
	def := a.TargetVal - a.personalPool
	a.turnTarget = def * shared.Resources(a.Coeff)
}

func (a *Account) LoadTaxation(calcVal shared.Resources) {
	a.turnTarget += calcVal
}

func (a *Account) GetAllocMin() shared.Resources {
	if a.turnTarget < 0 {
		a.TargetVal += -1 * a.turnTarget
		return a.turnTarget * -1
	}
	return a.turnTarget
}

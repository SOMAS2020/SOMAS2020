package roles

import (

)

type BaseJudge struct {
	int id
	int budget
	int president_salary
	int ballotID
	int resAlocID
	map actionLog
}
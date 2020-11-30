// Package team3 contains code for team 3's client implementation
package team3

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"gonum.org/v1/gonum/mat"
)

const id = common.Team3

func init() {
	common.RegisterClient(id, &client{id: id})
}

type client struct {
	id common.ClientID
}

func BasicRuleEvaluator(ruleName string) bool {
	rm := rules.AvailableRules[ruleName]

	variables := rm.RequiredVariables

	var variableVect []float64

	for _, v := range variables {
		if rules.VariableMap[v].Multivalued {
			variableVect = append(variableVect, rules.VariableMap[v].MultiValue...)
		} else {
			variableVect = append(variableVect, rules.VariableMap[v].SingleValue)
		}

	}

	variableVect = append(variableVect, 1)

	variableFormalVect := mat.NewVecDense(len(variables)+1, variableVect)

	rows, _ := rm.ApplicableMatrix.Dims()

	actual := make([]float64, rows)

	c := mat.NewVecDense(rows, actual)

	c.MulVec(&rm.ApplicableMatrix, variableFormalVect)
	aux := rm.AuxiliaryVector

	var resultVect []bool

	for i := 0; i < rows; i++ {
		switch interpret := aux.AtVec(i); interpret {
		case 0:
			resultVect = append(resultVect, c.AtVec(i) == 0)
		case 1:
			resultVect = append(resultVect, c.AtVec(i) > 0)
		case 2:
			resultVect = append(resultVect, c.AtVec(i) >= 0)
		case 3:
			resultVect = append(resultVect, c.AtVec(i) != 0)
		}
	}

	var finalBool = true

	for _, v := range resultVect {
		finalBool = finalBool && v
	}

	return finalBool
}

func (c *client) Echo(s string) string {

	//c.DemoEvaluation()
	c.Logf("Echo: '%v'", s)

	return s
}

func (c *client) DemoEvaluation() {
	getEval := BasicRuleEvaluator("Kinda Complicated Rule")
	c.Logf("Rule Eval: %t", getEval)
}

func (c *client) GetID() common.ClientID {
	return c.id
}

// Logf is the client's logger that prepends logs with your ID. This makes
// it easier to read logs. DO NOT use other loggers that will mess logs up!
func (c *client) Logf(format string, a ...interface{}) {
	log.Printf("[%v]: %v", c.id, fmt.Sprintf(format, a...))
}

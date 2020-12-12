package common

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
)

// Client is a base interface to be implemented by each client struct.
type Client interface {
	Echo(s string) string
	GetID() shared.ClientID

	// StartOfTurnUpdate is where SOMASServer.updateIsland sends the game state over
	// at start of turn. Do whatever you like here :).
	StartOfTurnUpdate(gameState GameState)

	// EndOfTurnActions should return all end of turn actions.
	EndOfTurnActions() []Action
}

// RegisteredClients contain all registered clients, exposed for the server.
var RegisteredClients = map[shared.ClientID]Client{}

// RegisterClient registers clients into RegisteredClients
func RegisterClient(id shared.ClientID, c Client) {
	// prevent double registrations
	if _, ok := RegisteredClients[id]; ok {
		// OK to panic here, as this is a _crucial_ step.
		panic(fmt.Sprintf("Duplicate client ID %v in RegisterClient!", id))
	}
	RegisteredClients[id] = c
}

func createVarList(variables []string) ([]float64, error){
	var variableVect []float64

    for _, v := range variables {
        if val, varOk := rules.VariableMap[v]; varOk {
            variableVect = append(variableVect, val.Values...)
        } else {
            return nil, errors.New(fmt.Sprintf("Variable: '%v' not found in global variable cache", v))
        }
    }

    variableVect = append(variableVect, 1)
    return variableVect, nil

}

func ruleMul(variableVect []float64, ApplicableMatrix mat.Dense)(*mat.VecDense){
    nRows,_ := ApplicableMatrix.Dims()

    variableFormalVect := mat.NewVecDense(len(variableVect), variableVect)

    actual := make([]float64, nRows)

    c := mat.NewVecDense(nRows, actual)

    c.MulVec(&ApplicableMatrix, variableFormalVect)

    return c

}

func genResult(aux mat.VecDense, c *mat.VecDense)([]bool, error){
    var resultVect []bool
    nRows, _ := c.Dims()

    for i := 0; i < nRows; i++ {
        res := false
        switch interpret := aux.AtVec(i); interpret {
        case 0:
            res = c.AtVec(i) == 0
        case 1:
            res = c.AtVec(i) > 0
        case 2:
            res = c.AtVec(i) >= 0
        case 3:
            res = c.AtVec(i) != 0
        default:
            return nil, errors.New(fmt.Sprintf("At auxillary vector entry: '%v' aux value outside of 0-3: '%v' was found", i, interpret))
        }
        resultVect = append(resultVect, res)
    }

    return resultVect, nil

}

func checkForFalse(resultVect []bool)(bool){

    var finalBool = true

    for _, v := range resultVect {
        finalBool = finalBool && v
    }

    return finalBool

}




// BasicRuleEvaluator implements a basic version of the Matrix rule evaluator, provides single boolean output (and error if present)
func BasicRuleEvaluator(ruleName string) (bool, error) {

	if rm, ok := rules.AvailableRules[ruleName]; ok {
        variableVect, err := createVarList(rm.RequiredVariables)
        if err != nil{
            return false,  err
        }

		//Checking dimensions line up
		_, nCols := rm.ApplicableMatrix.Dims()

		if nCols != len(variableVect) {
			return false, errors.Errorf(
				"dimension mismatch in evaluating rule: '%v' rule matrix has '%v' columns, while we sourced '%v' variables",
				ruleName,
				nCols,
				len(variableVect),
			)
		}

        c := ruleMul(variableVect, rm.ApplicableMatrix)

		resultVect, err := genResult(rm.AuxiliaryVector, c)
        if err != nil{
            return false,  err
        }

		return checkForFalse(resultVect), nil
	} else {
		return false, errors.Errorf("rule name: '%v' provided doesn't exist in global rule list", ruleName)
	}
}

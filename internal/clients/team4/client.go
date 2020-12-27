// Package team4 contains code for team 4's client implementation
package team4

import (
	"fmt"
	"reflect"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team4

func init() {
	baseclient.RegisterClient(id, &client{Client: baseclient.NewClient(id)})
}

type client struct {
	baseclient.Client
	yes string //this field is just for testing
}

func (c *client) StartOfTurn() {
	c.yes = "for you"
	c.Logf(`what are you doing?
	=========================================
	==============================================
	================================================`)
	c.Logf("this is yes " + c.yes)
	fmt.Println(reflect.TypeOf(c))
}

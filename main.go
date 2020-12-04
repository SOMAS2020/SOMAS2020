// Package main is the main entrypoint of the program.
package main

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/server"
)

func main() {
	s := server.SOMASServerFactory()
	gameStates := s.EntryPoint()
	for _, st := range gameStates {
		fmt.Printf("DAY: %v\n", st.Day)
		fmt.Printf("%#v\n", st)
	}
}

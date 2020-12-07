// Package main is the main entrypoint of the program.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SOMAS2020/SOMAS2020/internal/server"
)

func main() {
	s := server.SOMASServerFactory()
	if gameStates, err := s.EntryPoint(); err != nil {
		log.Printf("Run failed with: %v", err)
		os.Exit(1)
	} else {
		for _, st := range gameStates {
			fmt.Printf("===== START OF TURN %v (END OF TURN %v) =====\n", st.Turn, st.Turn-1)
			// this is fine for now, we shall visualise the data later on
			fmt.Printf("%#v\n", st)
		}
	}
}

// Package main is the main entrypoint of the program.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/server"
	"github.com/SOMAS2020/SOMAS2020/pkg/fileutils"
)

// output represents what is output into the output.json file
type output struct {
	GameStates []gamestate.GameState
	Config     config.Config
}

const outputJSONFileName = "output.json"

var cwd = fileutils.GetCurrFileDir()
var outputDir = path.Join(cwd, "output")
var outputJSONFilePath = path.Join(outputDir, outputJSONFileName)

func init() {
	// cleanup output
	err := fileutils.RemovePathIfExists(outputDir)
	if err != nil {
		panic("Cannot remove output directory.")
	}
	// make output directory
	err = os.Mkdir(outputDir, 0644)
	if err != nil {
		panic("Cannot make output directory.")
	}
}

func main() {
	s := server.SOMASServerFactory()
	if gameStates, err := s.EntryPoint(); err != nil {
		log.Printf("Run failed with: %+v", err)
		os.Exit(1)
	} else {
		fmt.Printf("===== GAME CONFIGURATION =====\n")
		fmt.Printf("%#v\n", config.GameConfig())
		for _, st := range gameStates {
			fmt.Printf("===== START OF TURN %v (END OF TURN %v) =====\n", st.Turn, st.Turn-1)
			fmt.Printf("%#v\n", st)
		}

		outputJSON(output{
			GameStates: gameStates,
			Config:     config.GameConfig(),
		})
	}
}

func outputJSON(o output) {
	log.Printf("Writing JSON output to '%v'\n", outputJSONFilePath)
	jsonBuf, err := json.MarshalIndent(o, "", "\t")
	if err != nil {
		log.Printf("Failed to Marshal gameStates: %v", err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(outputJSONFilePath, jsonBuf, 0644)
	if err != nil {
		log.Printf("Failed to write file: %v", err)
		os.Exit(1)
	}
	log.Printf("Finished writing JSON output to '%v'", outputJSONFilePath)
}

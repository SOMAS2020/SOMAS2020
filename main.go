// +build !js

// Package main is the main entrypoint of the program.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/SOMAS2020/SOMAS2020/internal/server"
	"github.com/SOMAS2020/SOMAS2020/pkg/fileutils"
	"github.com/SOMAS2020/SOMAS2020/pkg/gitinfo"
	"github.com/SOMAS2020/SOMAS2020/pkg/logger"
)

const outputJSONFileName = "output.json"
const outputLogFileName = "log.txt"

var cwd = fileutils.GetCurrFileDir()
var outputDir = path.Join(cwd, "output")
var outputJSONFilePath = path.Join(outputDir, outputJSONFileName)
var outputLogFilePath = path.Join(outputDir, outputLogFileName)

func init() {
	// cleanup output
	err := fileutils.RemovePathIfExists(outputDir)
	if err != nil {
		panic(err)
	}
	// make output directory
	err = os.Mkdir(outputDir, 0777)
	if err != nil {
		panic(err)
	}
	initLogger()
}

func initLogger() {
	f, err := os.OpenFile(outputLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(fmt.Sprintf("Unable to open log file, try running using sudo: %v", err))
	}
	log.SetOutput(
		logger.NewLogWriter([]io.Writer{os.Stderr, f}),
	)
}

func main() {
	timeStart := time.Now()
	flag.Parse()
	gameConfig, err := parseConfig()
	if err != nil {
		log.Printf("Flag parse error: %v\nUse --help.", err)
		os.Exit(1)
	}
	s := server.NewSOMASServer(gameConfig)
	if gameStates, err := s.EntryPoint(); err != nil {
		log.Printf("Run failed with: %+v", err)
		os.Exit(1)
	} else {
		fmt.Printf("===== GAME CONFIGURATION =====\n")
		fmt.Printf("%#v\n", gameConfig)
		for _, st := range gameStates {
			fmt.Printf("===== START OF TURN %v (END OF TURN %v) =====\n", st.Turn, st.Turn-1)
			fmt.Printf("%#v\n", st)
		}
		timeEnd := time.Now()
		outputJSON(output{
			GameStates: gameStates,
			Config:     gameConfig,
			GitInfo:    getGitInfo(),
			RunInfo: runInfo{
				TimeStart:       timeStart,
				TimeEnd:         timeEnd,
				DurationSeconds: timeEnd.Sub(timeStart).Seconds(),
				Version:         runtime.Version(),
				GOOS:            runtime.GOOS,
				GOARCH:          runtime.GOARCH,
			},
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
	err = ioutil.WriteFile(outputJSONFilePath, jsonBuf, 0777)
	if err != nil {
		log.Printf("Failed to write file: %v", err)
		os.Exit(1)
	}
	log.Printf("Finished writing JSON output to '%v'", outputJSONFilePath)
}

func getGitInfo() gitinfo.GitInfo {
	gitInfo, err := gitinfo.GetGitInfo(cwd)
	if err != nil {
		log.Printf("Ignoring error in getting git info--are you running this in a valid git repo? Error: %v", err)
	}
	return gitInfo
}

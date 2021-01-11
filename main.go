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
	"math/rand"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/SOMAS2020/SOMAS2020/internal/server"
	"github.com/SOMAS2020/SOMAS2020/pkg/fileutils"
	"github.com/SOMAS2020/SOMAS2020/pkg/gitinfo"
	"github.com/SOMAS2020/SOMAS2020/pkg/logger"
	"github.com/pkg/errors"
)

const outputJSONFileName = "output.json"
const outputLogFileName = "log.txt"

// non-WASM flags.
// see `params.go` for shared flags.
var (
	outputFolderName = flag.String(
		"output",
		"output",
		"The relative path (to the current working directory) to store output.json and logs in.\n"+
			"WARNING: This folder will be removed prior to running!",
	)
	logLevel = flag.Uint(
		"logLevel",
		0,
		"Logging verbosity level. Note that output artifacts will remain the same.\n"+
			"0: No logs at all\n"+
			"1: Game logs (identical to logs.txt) (to stderr)\n"+
			"2: As in 1 plus game states (similar to output.json) (to stdout)\n",
	)
)

func main() {
	timeStart := time.Now()
	rand.Seed(timeStart.UTC().UnixNano())

	flag.Parse()

	var err error

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("%v", err)
	}

	absOutputDir := path.Join(wd, *outputFolderName)

	err = prepareOutputFolder(absOutputDir)
	if err != nil {
		log.Fatalf("Failed to prepare output folder: %v", err)
	}
	err = prepareLogger(absOutputDir)
	if err != nil {
		log.Fatalf("Failed to prepare logger: %v", err)
	}
	gameConfig, err := parseConfig()
	if err != nil {
		log.Fatalf("Flag parse error: %v\nUse --help.", err)
	}

	s, err := server.NewSOMASServer(gameConfig)
	if err != nil {
		log.Fatalf("Failed to initial SOMASServer: %v", err)
	}
	if gameStates, err := s.EntryPoint(); err != nil {
		log.Fatalf("Run failed with: %+v", err)
	} else {
		if *logLevel >= 2 {
			fmt.Printf("===== GAME CONFIGURATION =====\n")
			fmt.Printf("%#v\n", gameConfig)
			for _, st := range gameStates {
				fmt.Printf("===== START OF TURN %v (END OF TURN %v) =====\n", st.Turn, st.Turn-1)
				fmt.Printf("%#v\n", st)
			}
		}
		timeEnd := time.Now()
		err = outputJSON(output{
			GameStates: gameStates,
			Config:     gameConfig,
			GitInfo:    getGitInfo(),
			AuxInfo:    getAuxInfo(),
			RunInfo: runInfo{
				TimeStart:       timeStart,
				TimeEnd:         timeEnd,
				DurationSeconds: timeEnd.Sub(timeStart).Seconds(),
				Version:         runtime.Version(),
				GOOS:            runtime.GOOS,
				GOARCH:          runtime.GOARCH,
			},
		}, absOutputDir)
		if err != nil {
			log.Fatalf("Failed to output JSON: %v", err)
		}
	}
}

func prepareOutputFolder(absOutputDir string) error {
	// cleanup output
	err := fileutils.RemovePathIfExists(absOutputDir)
	if err != nil {
		return err
	}
	// make output directory
	err = os.Mkdir(absOutputDir, 0777)
	if err != nil {
		return err
	}
	return nil
}

func prepareLogger(absOutputDir string) error {
	outputLogFilePath := path.Join(absOutputDir, outputLogFileName)

	f, err := os.OpenFile(outputLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return errors.Errorf("Unable to open log file, try running using sudo: %v", err)
	}

	writers := []io.Writer{f}

	if *logLevel >= 1 {
		writers = append(writers, os.Stderr)
	}

	log.SetOutput(
		logger.NewLogWriter(writers),
	)

	return nil
}

func outputJSON(o output, absOutputDir string) error {
	outputJSONFilePath := path.Join(absOutputDir, outputJSONFileName)

	log.Printf("Writing JSON output to '%v'\n", outputJSONFilePath)
	jsonBuf, err := json.MarshalIndent(o, "", "\t")
	if err != nil {
		return errors.Errorf("Failed to Marshal gameStates: %v", err)
	}
	err = ioutil.WriteFile(outputJSONFilePath, jsonBuf, 0777)
	if err != nil {
		return errors.Errorf("Failed to write file: %v", err)
	}

	log.Printf("Finished writing JSON output to '%v'", outputJSONFilePath)
	return nil
}

func getGitInfo() gitinfo.GitInfo {
	repoRootPath := fileutils.GetCurrFileDir()
	gitInfo, err := gitinfo.GetGitInfo(repoRootPath)
	if err != nil {
		log.Printf("Ignoring error in getting git info--are you running this in a valid git repo? Error: %v", err)
	}
	return gitInfo
}

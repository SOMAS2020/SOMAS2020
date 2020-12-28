// +build js,wasm

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"syscall/js"

	"github.com/SOMAS2020/SOMAS2020/internal/server"
	"github.com/SOMAS2020/SOMAS2020/pkg/logger"
	"github.com/pkg/errors"
)

func main() {
	js.Global().Set(
		"RunGame", js.FuncOf(RunGame),
	)
	select {}
}

func RunGame(this js.Value, args []js.Value) interface{} {
	// log into a buffer
	logBuf := bytes.Buffer{}
	log.SetOutput(logger.NewLogWriter([]io.Writer{&logBuf}))

	gameConfig := parseConfig()
	s := server.NewSOMASServer(gameConfig)

	var o output
	var outputJSON string
	var err error
	gameStates, err := s.EntryPoint()

	// === IF NO ERROR ===
	if err == nil {
		o = output{
			GameStates: gameStates,
			Config:     gameConfig,
			// no git info
		}
		outputJSON, err = getOutputJSON(o)
	}

	return js.ValueOf(map[string]interface{}{
		"output": outputJSON,
		"logs":   logBuf.String(),
		"error":  convertError(err),
	})
}

func getOutputJSON(o output) (string, error) {
	jsonBuf, err := json.Marshal(o)
	if err != nil {
		return "", errors.Errorf("Failed to Marshal output: %v", err)
	}
	return string(jsonBuf), nil
}

func convertError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

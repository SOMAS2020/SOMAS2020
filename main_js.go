// +build js,wasm

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"reflect"
	"syscall/js"

	"github.com/SOMAS2020/SOMAS2020/internal/server"
	"github.com/SOMAS2020/SOMAS2020/pkg/logger"
	"github.com/pkg/errors"
)

type flagInfo struct {
	Name     string
	Usage    string
	DefValue string
	Type     string
}

func main() {
	js.Global().Set(
		"RunGame", js.FuncOf(RunGame),
		"GetFlagsFormats", js.FuncOf(getFlagsFormats),
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

func getFlagsFormats(this js.Value, args []js.Value) interface{} {
	flagInfos := []flagInfo{}
	flagVisitor := flagVisitOuter(&flagInfos)
	flag.VisitAll(flagVisitor)

	output, err := getFlagsFormatsJSON(flagInfos)

	return js.ValueOf(map[string]interface{}{
		"output": output,
		"error":  convertError(err),
	})
}

// getFlagBaseType processes the String of the reflect.Type of a flag to
// only return the bare type.
// e.g. *flag.float64Value -> float64
func getFlagBaseType(f *flag.Flag) string {
	s := reflect.TypeOf(f.Value).String()
	return s[6 : len(s)-5]
}

func flagVisitOuter(fis *[]flagInfo) func(*flag.Flag) {
	return func(f *flag.Flag) {
		*fis = append(*fis, flagInfo{
			Name:     f.Name,
			Usage:    f.Usage,
			DefValue: f.DefValue,
			Type:     getFlagBaseType(f),
		})
		return
	}
}
func getFlagsFormatsJSON(fis []flagInfo) (string, error) {
	jsonBuf, err := json.Marshal(fis)
	if err != nil {
		return "", errors.Errorf("Failed to Marshal flags formats: %v", err)
	}
	return string(jsonBuf), nil
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

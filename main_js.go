// +build js,wasm

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"reflect"
	"runtime"
	"strings"
	"syscall/js"
	"time"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
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
	)
	js.Global().Set(
		"GetFlagsFormats", js.FuncOf(GetFlagsFormats),
	)
	select {}
}

// RunGame runs the game.
// args[0] are optional arguments to set flag values
// The format is `arg1=value,arg2=value,...`
func RunGame(this js.Value, args []js.Value) interface{} {
	// log into a buffer
	logBuf := bytes.Buffer{}
	log.SetOutput(logger.NewLogWriter([]io.Writer{&logBuf}))

	timeStart := time.Now()
	gameConfig, err := getConfigFromArgs(args)
	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"error": convertError(err),
		})
	}

	s := server.NewSOMASServer(gameConfig)

	var o output
	var outputJSON string
	gameStates, err := s.EntryPoint()
	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"error": convertError(err),
		})
	}
	timeEnd := time.Now()
	o = output{
		GameStates: gameStates,
		Config:     gameConfig,
		// no git info
		RunInfo: runInfo{
			TimeStart:       timeStart,
			TimeEnd:         timeEnd,
			DurationSeconds: timeEnd.Sub(timeStart).Seconds(),
			Version:         runtime.Version(),
			GOOS:            runtime.GOOS,
			GOARCH:          runtime.GOARCH,
		},
	}
	outputJSON, err = getOutputJSON(o)

	return js.ValueOf(map[string]interface{}{
		"output": outputJSON,
		"logs":   logBuf.String(),
		"error":  convertError(err),
	})
}

// GetFlagsFormats returns the format of the flags.
// See flagInfo type for more information.
func GetFlagsFormats(this js.Value, args []js.Value) interface{} {
	flagInfos := []flagInfo{}
	flagVisitor := flagVisitOuter(&flagInfos)
	flag.VisitAll(flagVisitor)

	output, err := getFlagsFormatsJSON(flagInfos)
	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"error": convertError(err),
		})
	}

	return js.ValueOf(map[string]interface{}{
		"output": output,
		"error":  convertError(err),
	})
}

func getConfigFromArgs(jsArgs []js.Value) (config.Config, error) {
	args := strings.Split(jsArgs[0].String(), ",")

	if len(args) == 1 && args[0] == "" {
		// handle empty input
		args = []string{}
	}

	flag.Parse()

	for _, arg := range args {
		nameValuePair := strings.Split(arg, "=")
		if len(nameValuePair) != 2 {
			return config.Config{}, errors.Errorf("Invalid arg: %v", arg)
		}
		name := nameValuePair[0]
		val := nameValuePair[1]
		err := flag.Set(name, val)
		if err != nil {
			return config.Config{}, errors.Errorf("Cannot set flag name '%v' with value '%v': %v", name, val, err)
		}
	}

	conf, err := parseConfig()
	if err != nil {
		return conf, errors.Errorf("Flag parse error: %v", err)
	}

	return conf, nil
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

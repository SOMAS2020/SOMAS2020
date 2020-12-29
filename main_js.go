// +build js,wasm

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"reflect"
	"strings"
	"syscall/js"

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

func RunGame(this js.Value, args []js.Value) interface{} {
	gameConfig, err := getConfigFromArgs(args)
	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"error": convertError(err),
		})
	}

	// log into a buffer
	logBuf := bytes.Buffer{}
	log.SetOutput(logger.NewLogWriter([]io.Writer{&logBuf}))

	s := server.NewSOMASServer(gameConfig)

	var o output
	var outputJSON string
	gameStates, err := s.EntryPoint()
	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"error": convertError(err),
		})
	}

	o = output{
		GameStates: gameStates,
		Config:     gameConfig,
		// no git info
	}
	outputJSON, err = getOutputJSON(o)

	return js.ValueOf(map[string]interface{}{
		"output": outputJSON,
		"logs":   logBuf.String(),
		"error":  convertError(err),
	})
}

func getConfigFromArgs(jsArgs []js.Value) (config.Config, error) {
	args := make([]string, len(jsArgs))

	for i, jsArg := range jsArgs {
		args[i] = jsArg.String()
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

	return parseConfig(), nil
}

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

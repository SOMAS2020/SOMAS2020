package main

import (
	"encoding/json"
	"flag"
	"log"
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/pkg/jstype"
	"github.com/pkg/errors"
)

type flagInfo struct {
	Name     string
	Usage    string
	DefValue string
	Type     string
}

func flagVisitOuter(fis *[]flagInfo) func(*flag.Flag) {
	return func(f *flag.Flag) {
		jst, err := jstype.GoToJSType(reflect.TypeOf(f.Value).Kind())
		if err != nil {
			log.Printf("Unable to get type for %v: %v", f.Name, err)
			jst = jstype.JSInvalid
		}
		*fis = append(*fis, flagInfo{
			Name:     f.Name,
			Usage:    f.Usage,
			DefValue: f.DefValue,
			Type:     jst.String(),
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

func getFlagsFormats() interface{} {
	flagInfos := []flagInfo{}
	flagVisitor := flagVisitOuter(&flagInfos)
	flag.VisitAll(flagVisitor)

	output, err := getFlagsFormatsJSON(flagInfos)

	return (map[string]interface{}{
		"output": output,
		"error":  convertError(err),
	})
}

func convertError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func TestGetFlagsFormats(t *testing.T) {
	t.Errorf("%v", getFlagsFormats())
}

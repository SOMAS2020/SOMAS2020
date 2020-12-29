// Package jstype provides utilities for WASM interop between Golang and JS.
// Author: lhl2617.
package jstype

import (
	"fmt"
	"reflect"

	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
	"github.com/pkg/errors"
)

// JSType is an enum of JavaScript types
type JSType int

// JS types
const (
	JSNull JSType = iota
	JSUndefined
	JSBoolean
	JSNumber
	JSString
	JSObject
	JSSymbol
	JSBigInt
	JSArray
	JSMap
	JSFunction
	JSInvalid
)

var strs = [...]string{
	"null",
	"undefined",
	"boolean",
	"number",
	"string",
	"object",
	"symbol",
	"bigint",
	"array",
	"map",
	"function",
	"invalid",
}

func (t JSType) String() string {
	if t >= 0 && int(t) < len(strs) {
		return strs[t]
	}
	return fmt.Sprintf("UNKNOWN JSType '%v'", int(t))
}

// GoString implements GoStringer
func (t JSType) GoString() string {
	return t.String()
}

// MarshalText implements TextMarshaler
func (t JSType) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(t.String())
}

// MarshalJSON implements RawMessage
func (t JSType) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(t.String())
}

// goTypeToJSTypeMap attempts to map golang types to their JS counterparts.
// This is strictly by personal choice. Copy and change for your use-case.
var goTypeToJSTypeMap = map[reflect.Kind]JSType{
	reflect.Invalid:       JSInvalid,
	reflect.Bool:          JSBoolean,
	reflect.Int:           JSBigInt,
	reflect.Int8:          JSBigInt,
	reflect.Int16:         JSBigInt,
	reflect.Int32:         JSBigInt,
	reflect.Int64:         JSBigInt,
	reflect.Uint:          JSBigInt,
	reflect.Uint8:         JSBigInt,
	reflect.Uint16:        JSBigInt,
	reflect.Uint32:        JSBigInt,
	reflect.Uint64:        JSBigInt,
	reflect.Uintptr:       JSBigInt,
	reflect.Float32:       JSNumber,
	reflect.Float64:       JSNumber,
	reflect.Complex64:     JSInvalid,
	reflect.Complex128:    JSInvalid,
	reflect.Array:         JSArray,
	reflect.Chan:          JSInvalid,
	reflect.Func:          JSFunction,
	reflect.Interface:     JSObject,
	reflect.Map:           JSMap,
	reflect.Ptr:           JSBigInt,
	reflect.Slice:         JSArray,
	reflect.String:        JSString,
	reflect.Struct:        JSObject,
	reflect.UnsafePointer: JSInvalid,
}

// GoToJSType gets the JS type from the Go type.
func GoToJSType(t reflect.Kind) (JSType, error) {
	jst, ok := goTypeToJSTypeMap[t]
	if !ok || jst == JSInvalid {
		return jst, errors.Errorf("Unknown type '%v'", t)
	}
	return jst, nil
}

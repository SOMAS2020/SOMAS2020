// Package logger provides a Writer implementation to be used for log.SetOutput
// to enable logging to multiple streams.
// Based on https://github.com/facebook/openbmc/tools/flashy/lib/logger/logger.go
// Original author: lhl2617, Facebook, (GPLv2)
package logger

import (
	"fmt"
	"io"
)

// LogWriter is an implementation of io.Writer that will call write
// to all streams defined in Streams.
type LogWriter struct {
	Streams []io.Writer
}

// LogWriter implements io.Writer
func (w LogWriter) Write(p []byte) (n int, err error) {
	for _, stream := range w.Streams {
		n, err = stream.Write(p)
		if err != nil {
			// don't panic, just log to stderr that this has failed
			println(fmt.Sprintf("%v", err))
		}
	}
	return
}

// NewLogWriter returns a new LogWriter that logs to the given streams.
func NewLogWriter(streams []io.Writer) LogWriter {
	return LogWriter{
		Streams: streams,
	}
}

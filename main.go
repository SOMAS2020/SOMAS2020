// Package main is the main entrypoint of the program.
package main

import "github.com/SOMAS2020/SOMAS2020/internal/server"

func main() {
	s := server.SOMASServerFactory()
	s.GetEcho("Hello World")
}

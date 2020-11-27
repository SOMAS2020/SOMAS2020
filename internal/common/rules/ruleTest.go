package rules

import (
	"log"
)

var testLog = ""

func init() {
	testLog = "Hello from a rules test"
}

func CheckInitialisation() {
	log.Print(testLog)
}

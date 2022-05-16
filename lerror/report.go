package lerror

import "log"

var HadError = false
var HadRuntimeError = false

func Report(line int, where string, message string) {
	log.Printf("[line %d ] Error %s : %s\n", line, where, message)
	HadError = true
}

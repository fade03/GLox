package loxerror

import "log"

func Report(line int, where string, message string) {
	log.Printf("[line %d ] Error %s : %s\n", line, where, message)
	HadError = true
}

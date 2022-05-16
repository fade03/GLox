package scanner

import e "LoxGo/lerror"

func serror(line int, message string) {
	e.Report(line, "", message)
}

package scanner

import e "GLox/lerror"

func serror(line int, message string) {
	e.Report(line, "", message)
}

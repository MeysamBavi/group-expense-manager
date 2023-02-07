package log

import (
	"fmt"
	"log"
	"path"
	"runtime"
)

func init() {
	log.SetFlags(0)
}

func Error(err error) {
	log.Print(getLineInfo(0), err)
}

func FatalError(err error) {
	log.Fatal(getLineInfo(0), err)
}

func FatalErrorByCaller(err error) {
	log.Fatal(getLineInfo(1), err)
}

func getLineInfo(skipOffset int) string {
	pc, fileName, lineNumber, ok := runtime.Caller(2 + skipOffset)
	return formatLineInfo(fileName, runtime.FuncForPC(pc).Name(), lineNumber, ok)
}

func formatLineInfo(fileName string, funcName string, lineNumber int, ok bool) string {
	if !ok {
		return "unknown error"
	}
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	return fmt.Sprintf("at %s %s:%d: ", funcName, fileName, lineNumber)
}

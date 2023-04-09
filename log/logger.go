package log

import (
	"fmt"
	"log"
)

var (
	logger Logger = NewLogger(true)
)

type Logger struct {
	shortFlag bool //是否使用简写
}

func NewLogger(shortFlag bool) Logger {
	return Logger{shortFlag: shortFlag}
}

func (l Logger) logger(format string, v ...any) {

	if l.shortFlag {
		fmt.Printf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}

func Info(msg string) {
	logger.logger("%s\r\n", msg)
}

func InfoF(format string, v ...any) {
	logger.logger(format+"\r\n", v...)
}

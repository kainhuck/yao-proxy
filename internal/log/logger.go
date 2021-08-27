package log

import (
	"fmt"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

var loggerColor map[string]string = map[string]string{
	DEBUG: "\x1b[34m",
	INFO:  "\x1b[32m",
	WARN:  "\x1b[33m",
	ERROR: "\x1b[31m",
}

type logger bool

func (l logger) log(logLevel string, format *string, args ...interface{}) {
	if l {
		color := loggerColor[logLevel]
		if format != nil {
			msg := fmt.Sprintf(*format, args...)
			fmt.Printf("%s[%s]\x1b[0m %s\n", color, logLevel, msg)
		} else {
			fmt.Printf("%s[%s]\x1b[0m ", color, logLevel)
			fmt.Println(args...)
		}
	}
}

func (l logger) Debug(args ...interface{}) {
	if l {
		l.log(DEBUG, nil, args...)
	}
}

func (l logger) Info(args ...interface{}) {
	if l {
		l.log(INFO, nil, args...)
	}
}

func (l logger) Warn(args ...interface{}) {
	if l {
		l.log(WARN, nil, args...)
	}
}

func (l logger) Error(args ...interface{}) {
	if l {
		l.log(ERROR, nil, args...)
	}
}

func (l logger) Debugf(format string, args ...interface{}) {
	if l {
		l.log(DEBUG, &format, args...)
	}
}

func (l logger) Infof(format string, args ...interface{}) {
	if l {
		l.log(INFO, &format, args...)
	}
}

func (l logger) Warnf(format string, args ...interface{}) {
	if l {
		l.log(WARN, &format, args...)
	}
}

func (l logger) Errorf(format string, args ...interface{}) {
	if l {
		l.log(ERROR, &format, args...)
	}
}

func NewLogger(debug bool) Logger {
	return logger(debug)
}

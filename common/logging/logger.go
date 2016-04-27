package logging

import (
	"io"
	"log"
	"fmt"
)

const (
	DebugLevel LogLevel = iota
	WarningLevel
	InfoLevel
	ErrorLevel
	FatalLevel
)

type LogLevel int
type LogTargets []*log.Logger

var loggers map[LogLevel]LogTargets

func AddTarget(target io.Writer, minLevel LogLevel) error {
	if target == nil {
		// Do nothing.
		return fmt.Errorf("Nil cannot be added as a log target!")
	}

	for level := DebugLevel; level <= FatalLevel; level++ {
		// Add the target to the selected levels
		if minLevel <= level {
			var levelPrefix string
			switch level {
			case DebugLevel:
				levelPrefix = "DEBUG:   "
			case InfoLevel:
				levelPrefix = "INFO:    "
			case WarningLevel:
				levelPrefix = "WARNING: "
			case ErrorLevel:
				levelPrefix = "ERROR:   "
			case FatalLevel:
				levelPrefix = "FATAL:   "
			default:
				return fmt.Errorf("Invalid level specified while adding log target! Requested log level: %v", minLevel)
			}

			// Lazy init of loggers map
			if loggers == nil {
				loggers = make(map[LogLevel]LogTargets)
			}

			loggers[level] = append(loggers[level],
				log.New(target, levelPrefix, log.LstdFlags|log.LUTC|log.Lshortfile|log.Lmicroseconds))
		}
	}

	return nil
}

//// Debug
func Debug(v ...interface{}) {
	printAll(DebugLevel, v...)
}

func Debugf(format string, v ...interface{}) {
	printAllf(DebugLevel, format, v...)
}

//// Info
func Info(v ...interface{}) {
	printAll(InfoLevel, v...)
}

func Infof(format string, v ...interface{}) {
	printAllf(InfoLevel, format, v...)
}

//// Warning
func Warning(v ...interface{}) {
	printAll(WarningLevel, v...)
}

func Warningf(format string, v ...interface{}) {
	printAllf(WarningLevel, format, v...)
}

//// Error
func Error(v ...interface{}) {
	printAll(ErrorLevel, v...)
}

func Errorf(format string, v ...interface{}) {
	printAllf(ErrorLevel, format, v...)
}

//// Fatal
func Fatal(v ...interface{}) {
	printAll(FatalLevel, v...)
}

func Fatalf(format string, v ...interface{}) {
	printAllf(FatalLevel, format, v...)
}

// Private implementations
func printAll(level LogLevel, v ...interface{}) {
	if loggers == nil {
		// Loggers can be nil, if no logger has been added yet.
		return
	}

	targets := loggers[level]
	for _, target := range targets {
		if target == nil {
			continue
		}
		if level == FatalLevel {
			target.Fatal(v...)
		} else {
			target.Print(v...)
		}
	}
}

func printAllf(level LogLevel, format string, v ...interface{}) {
	if loggers == nil {
		// Loggers can be nil, if no logger has been added yet.
		return
	}

	targets := loggers[level]
	for _, target := range targets {
		if target == nil {
			continue
		}
		if level == FatalLevel {
			target.Fatalf(format, v...)
		} else {
			target.Printf(format, v...)
		}
	}
}

package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	LevelDebug LogLevel = iota
	LevelWarning
	LevelInfo
	LevelError
	LevelFatal
)

type LogLevel int
type LogTargets []*log.Logger

var loggers map[LogLevel]LogTargets

func AddTarget(target io.Writer, minLevel LogLevel) error {
	if target == nil {
		// Do nothing.
		return fmt.Errorf("Nil cannot be added as a log target!")
	}

	for level := LevelDebug; level <= LevelFatal; level++ {
		// Add the target to the selected levels
		if minLevel <= level {
			var levelPrefix string
			switch level {
			case LevelDebug:
				levelPrefix = "DEBUG:   "
			case LevelInfo:
				levelPrefix = "INFO:    "
			case LevelWarning:
				levelPrefix = "WARNING: "
			case LevelError:
				levelPrefix = "ERROR:   "
			case LevelFatal:
				levelPrefix = "FATAL:   "
			default:
				return fmt.Errorf("Invalid level specified while adding log target! Requested log level: %v", minLevel)
			}

			// Lazy init of loggers map
			if loggers == nil {
				loggers = make(map[LogLevel]LogTargets)
			}

			loggers[level] = append(loggers[level],
				log.New(target, levelPrefix, log.LstdFlags|log.LUTC|log.Lmicroseconds))
		}
	}

	return nil
}

//// Debug
func Debug(v ...interface{}) {
	printAll(LevelDebug, v...)
}

func Debugf(format string, v ...interface{}) {
	printAllf(LevelDebug, format, v...)
}

//// Info
func Info(v ...interface{}) {
	printAll(LevelInfo, v...)
}

func Infof(format string, v ...interface{}) {
	printAllf(LevelInfo, format, v...)
}

//// Warning
func Warning(v ...interface{}) {
	printAll(LevelWarning, v...)
}

func Warningf(format string, v ...interface{}) {
	printAllf(LevelWarning, format, v...)
}

//// Error
func Error(v ...interface{}) {
	printAll(LevelError, v...)
}

func Errorf(format string, v ...interface{}) {
	printAllf(LevelError, format, v...)
}

//// Fatal
func Fatal(v ...interface{}) {
	printAll(LevelFatal, v...)
}

func Fatalf(format string, v ...interface{}) {
	printAllf(LevelFatal, format, v...)
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
		target.Print(v...)

		if level == LevelFatal {
			os.Exit(1)
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
		target.Printf(format, v...)

		if level == LevelFatal {
			os.Exit(1)
		}
	}
}

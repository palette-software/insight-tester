package logging

import (
	"io"
	"log"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	WarningLevel
	InfoLevel
	ErrorLevel
	FatalLevel
)

var (
	debugTargets   	[]*log.Logger
	infoTargets    	[]*log.Logger
	warningTargets 	[]*log.Logger
	errorTargets   	[]*log.Logger
	fatalTargets   	[]*log.Logger

)
func Init() {
	// Let's init them with the size of 5, that should be enough.
	// But append will resize these arrays, if necessary.
	debugTargets   = make([]*log.Logger, 0, 5)
	infoTargets    = make([]*log.Logger, 0, 5)
	warningTargets = make([]*log.Logger, 0, 5)
	errorTargets   = make([]*log.Logger, 0, 5)
	fatalTargets   = make([]*log.Logger, 0, 5)
}

func AddTarget(target io.Writer, minLevel LogLevel) {
	if target == nil {
		// Do nothing.
		return
	}

	if minLevel <= DebugLevel {
		debugTargets = append(debugTargets,
			log.New(target, "DEBUG:   ", log.LstdFlags|log.LUTC|log.Lshortfile|log.Lmicroseconds))
	}
	if minLevel <= InfoLevel {
		infoTargets = append(infoTargets,
			log.New(target, "INFO:    ", log.LstdFlags|log.LUTC|log.Lshortfile|log.Lmicroseconds))
	}
	if minLevel <= WarningLevel {
		warningTargets = append(warningTargets,
			log.New(target, "WARNING: ", log.LstdFlags|log.LUTC|log.Lshortfile|log.Lmicroseconds))
	}
	if minLevel <= ErrorLevel {
		errorTargets = append(errorTargets,
			log.New(target, "ERROR:   ", log.LstdFlags|log.LUTC|log.Lshortfile|log.Lmicroseconds))
	}
	if minLevel <= FatalLevel {
		fatalTargets = append(fatalTargets,
			log.New(target, "FATAL:   ", log.LstdFlags|log.LUTC|log.Lshortfile|log.Lmicroseconds))
	}
}

//// Debug
func Debug(v ...interface{}) {
	printAll(debugTargets, v...)
}

func Debugf(format string, v ...interface{}) {
	printAllf(debugTargets, format, v...)
}

//// Info
func Info(v ...interface{}) {
	printAll(infoTargets, v...)
}

func Infof(format string, v ...interface{}) {
	printAllf(infoTargets, format, v...)
}

//// Warning
func Warning(v ...interface{}) {
	printAll(warningTargets, v...)
}

func Warningf(format string, v ...interface{}) {
	printAllf(warningTargets, format, v...)
}

//// Error
func Error(v ...interface{}) {
	printAll(errorTargets, v...)
}

func Errorf(format string, v ...interface{}) {
	printAllf(errorTargets, format, v...)
}

//// Fatal
func Fatal(v ...interface{}) {
	for _, target := range fatalTargets {
		if target == nil {
			continue
		}
		target.Fatal(v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	for _, target := range fatalTargets {
		if target == nil {
			continue
		}
		target.Fatalf(format, v...)
	}
}

// Private implementations
func printAll(targets []*log.Logger, v ...interface{}) {
	for _, target := range targets {
		if target == nil {
			continue
		}
		target.Print(v...)
	}
}

func printAllf(targets []*log.Logger, format string, v ...interface{}) {
	for _, target := range targets {
		if target == nil {
			continue
		}
		target.Printf(format, v...)
	}
}

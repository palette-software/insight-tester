package logging

import (
    "log"
    "io"
)

type FatalLogger struct {
	innerLogger *log.Logger
}

func (fl *FatalLogger) Fatal(v ...interface{}) {
	fl.innerLogger.Fatal(v...)
}

func (fl *FatalLogger) Fatalf(format string, v ...interface{}) {
	fl.innerLogger.Fatalf(format, v...)
}

func (fl *FatalLogger) Fatalln(v ...interface{}) {
	fl.innerLogger.Fatalln(v...)
}

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
	Fatal   FatalLogger
)

func InitLog(
        traceHandle		io.Writer,
        infoHandle		io.Writer,
        warningHandle	io.Writer,
        errorHandle 	io.Writer,
		fatalHandle  	io.Writer)  {
    Trace 	= log.New(traceHandle,   "TRACE: ",   log.Ldate|log.Ltime|log.Lshortfile)
    Info 	= log.New(infoHandle, 	 "INFO: ",    log.Ldate|log.Ltime|log.Lshortfile)
    Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    Error 	= log.New(errorHandle, 	 "ERROR: ",   log.Ldate|log.Ltime|log.Lshortfile)
	fatal  := log.New(fatalHandle,	 "FATAL: ",	  log.Ldate|log.Ltime|log.Lshortfile)
	Fatal   = FatalLogger{fatal}
}

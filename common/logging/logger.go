package logging

import (
    "log"
    "io"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	fatal	*log.Logger
)

func Fatal(v ...interface{}) {
	fatal.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	fatal.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	fatal.Fatalln(v...)
}

func InitLog(
        traceHandle		io.Writer,
        infoHandle		io.Writer,
        warningHandle	io.Writer,
        errorHandle 	io.Writer,
		fatalHandle  	io.Writer)  {
    Trace 	= log.New(traceHandle,   "TRACE:   ", log.Ldate|log.Ltime|log.Lshortfile)
    Info 	= log.New(infoHandle, 	 "INFO:    ", log.Ldate|log.Ltime|log.Lshortfile)
    Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    Error 	= log.New(errorHandle, 	 "ERROR:   ", log.Ldate|log.Ltime|log.Lshortfile)
	fatal   = log.New(fatalHandle,	 "FATAL:   ", log.Ldate|log.Ltime|log.Lshortfile)
}

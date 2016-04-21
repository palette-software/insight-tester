package logging

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	fatal   *log.Logger
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
	debugHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	fatalHandle io.Writer) {
	Debug = log.New(debugHandle, "DEBUG:   ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile|log.Lmicroseconds)
	Info = log.New(infoHandle, "INFO:    ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile|log.Lmicroseconds)
	Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile|log.Lmicroseconds)
	Error = log.New(errorHandle, "ERROR:   ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile|log.Lmicroseconds)
	fatal = log.New(fatalHandle, "FATAL:   ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile|log.Lmicroseconds)
}

func Init() {
	InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
}

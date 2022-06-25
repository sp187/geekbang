package mid

import (
	"log"
	"os"

	"github.com/urfave/negroni"
)

var _DefaultLogger *negroni.Logger

func init() {
	_DefaultLogger = &negroni.Logger{ALogger: log.New(os.Stdout, os.Getenv("APP_NAME")+" ", 0)}
	_DefaultLogger.SetDateFormat(negroni.LoggerDefaultDateFormat)
	_DefaultLogger.SetFormat("{{.StartTime}} {{.Method}} {{.Status}} {{.Duration}} {{.Hostname}} {{.Path}}")
}

func DefaultLogger() *negroni.Logger {
	return _DefaultLogger
}

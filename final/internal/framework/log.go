package fw

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"text/template"
	"time"
)

// ALogger interface
type ALogger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	SetOutput(w io.Writer)
}

type LoggerEntry struct {
	Time     string
	Level    string
	Location string
	Message  string
}

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	ALogger
	dateFormat string
	template   *template.Template
	fd         *os.File
	level      LogLevel
}

func (l *Logger) SetFormat(format string) {
	l.template = template.Must(template.New("parser").Parse(format))
}

func (l *Logger) SetDateFormat(format string) {
	l.dateFormat = format
}

func (l *Logger) SaveLogToFile(filepath string) {
	if l == nil {
		return
	}
	logfd, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.Warn(filepath + " open log file err")
		return
	}
	l.fd = logfd

	l.ALogger.SetOutput(io.MultiWriter(os.Stdout, logfd))
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level > DEBUG {
		return
	}
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	log := LoggerEntry{
		Time:     time.Now().Format(l.dateFormat),
		Level:    "DEBUG",
		Location: file + ":" + strconv.Itoa(line) + ":",
		Message:  fmt.Sprintf(format, v...),
	}

	buff := &bytes.Buffer{}
	l.template.Execute(buff, log)
	l.Println(buff.String())
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level > INFO {
		return
	}
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	log := LoggerEntry{
		Time:     time.Now().Format(l.dateFormat),
		Level:    "INFO",
		Location: file + ":" + strconv.Itoa(line) + ":",
		Message:  fmt.Sprintf(format, v...),
	}

	buff := &bytes.Buffer{}
	l.template.Execute(buff, log)
	l.Println(buff.String())
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level > WARN {
		return
	}
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	log := LoggerEntry{
		Time:     time.Now().Format(l.dateFormat),
		Level:    "WARN",
		Location: file + ":" + strconv.Itoa(line) + ":",
		Message:  fmt.Sprintf(format, v...),
	}

	buff := &bytes.Buffer{}
	l.template.Execute(buff, log)
	l.Println(buff.String())
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.level > ERROR {
		return
	}
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	log := LoggerEntry{
		Time:     time.Now().Format(l.dateFormat),
		Level:    "ERROR",
		Location: file + ":" + strconv.Itoa(line) + ":",
		Message:  fmt.Sprintf(format, v...),
	}

	buff := &bytes.Buffer{}
	l.template.Execute(buff, log)
	l.Println(buff.String())
}

var (
	logOnce sync.Once
	logger  *Logger
)

func GetLogger() *Logger {
	logOnce.Do(func() {
		name := "geekbang"
		logger = &Logger{ALogger: log.New(os.Stdout, name+" ", 0), dateFormat: time.RFC3339}
		logger.SetFormat("{{.Time}} {{.Location}} {{.Level}} {{.Message}}")
	})
	return logger
}

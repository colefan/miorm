package miorm

import (
	"io"
	"log"
)

type LogLevel int

const (
	LOG_UNKNOWN = iota
	LOG_OFF
	LOG_DEBUG
	LOG_INFO
	LOG_WARNING
	LOG_ERROR
)

const (
	DEFAULT_LOG_PREFIX = "[miorm]"
	DEFAULT_LOG_FLAG   = log.Ldate | log.Lmicroseconds
	DEFAULT_LOG_LEVEL  = LOG_DEBUG
)

type ILogger interface {
	Debug(v ...interface{}) (err error)
	Debugf(format string, v ...interface{}) (err error)
	Err(v ...interface{}) (err error)
	Errf(format string, v ...interface{}) (err error)
	Info(v ...interface{}) (err error)
	Infof(format string, v ...interface{}) (err error)
	Warning(v ...interface{}) (err error)
	Warningf(format string, v ...interface{}) (err error)
	Level() LogLevel
	SetLevel(l LogLevel)
}

type DefaultLogger struct {
	log      *log.Logger
	level    LogLevel
	prefix   string
	levelfmt string
}

func levelfmt(l LogLevel) string {
	switch l {
	case LOG_DEBUG:
		return " [D]"
	case LOG_WARNING:
		return " [W]"
	case LOG_INFO:
		return " [I]"
	case LOG_ERROR:
		return " [E]"
	}

	return ""

}

func NewDefaultLogger(out io.Writer) *DefaultLogger {
	return NewDefaultLogger2(out, DEFAULT_LOG_PREFIX, DEFAULT_LOG_FLAG, DEFAULT_LOG_LEVEL)
}

func NewDefaultLogger2(out io.Writer, prefix string, flag int, level LogLevel) *DefaultLogger {
	return &DefaultLogger{
		log:      log.New(out, prefix+levelfmt(level), flag),
		level:    level,
		prefix:   prefix,
		levelfmt: levelfmt(level)}
}

func (d *DefaultLogger) Level() LogLevel {
	return d.level
}

func (d *DefaultLogger) SetLevel(l LogLevel) {
	d.level = l
	d.levelfmt = levelfmt(l)
	d.log.SetPrefix(d.prefix + d.levelfmt)
}

func (d *DefaultLogger) Debug(v ...interface{}) (err error) {
	if d.level >= LOG_DEBUG {
		d.log.Println(v...)
	}
	return
}

func (d *DefaultLogger) Debugf(format string, v ...interface{}) (err error) {
	if d.level >= LOG_DEBUG {
		d.log.Printf(format, v...)
	}
	return
}

func (d *DefaultLogger) Info(v ...interface{}) (err error) {
	if d.level >= LOG_INFO {
		d.log.Println(v...)
	}
	return
}

func (d *DefaultLogger) Infof(format string, v ...interface{}) (err error) {
	if d.level >= LOG_INFO {
		d.log.Printf(format, v...)
	}
	return

}

func (d *DefaultLogger) Warning(v ...interface{}) (err error) {
	if d.level >= LOG_WARNING {
		d.log.Println(v...)
	}
	return
}

func (d *DefaultLogger) Waringf(format string, v ...interface{}) (err error) {
	if d.level >= LOG_WARNING {
		d.log.Printf(format, v...)
	}
	return

}

func (d *DefaultLogger) Err(v ...interface{}) (err error) {
	if d.level >= LOG_ERROR {
		d.log.Println(v...)
	}
	return
}

func (d *DefaultLogger) Errf(format string, v ...interface{}) (err error) {
	if d.level >= LOG_ERROR {
		d.log.Printf(format, v...)
	}
	return
}

package logrusw

import (
	"sync"

	"github.com/go-interfaces/log"
	"github.com/sirupsen/logrus"
)

const loggerLabel = "logger"

var _ log.Logger = (*Logger)(nil)

// nolog is the reference noop logger.
// It will be returned when V(int) is called with a disabled log level.
// It is safe to use it as a global variable that might be used by multiple threads
// because no variables are modified withing its implementation
var nolog = &log.NoLog{}

// Logger wraps logrus logger in go-interfaces logger
type Logger struct {
	m         sync.Mutex
	name      string
	log       *logrus.Logger
	verbosity int
	logName   bool
}

// NewLogger returns an encapsulated logrus in logr interface
// It accepts
// - an already initialized logrus logger object
// - a verbosity level
// - a logger name to be added to all log lines (use empty string to avoid)
func NewLogger(l *logrus.Logger, verbosity int, optLogName string) *Logger {
	// logVerbosity = verbosity
	return &Logger{
		log:       l,
		name:      optLogName,
		verbosity: verbosity,
		logName:   len(optLogName) != 0,
	}
}

// Info logs non error messages
func (l *Logger) Info(msg string, kv ...log.KV) {
	if l.logName {
		kv = append(kv, log.KV{K: loggerLabel, V: l.name})
	}
	fl := logrus.FieldLogger(l.log)
	if len(kv) != 0 {
		newFields := l.fieldsFromKV(kv...)
		fl = fl.WithFields(*newFields)
	}
	fl.Info(msg)
}

// Error logs error messages
func (l *Logger) Error(err error, msg string, kv ...log.KV) {
	if l.logName {
		kv = append(kv, log.KV{K: loggerLabel, V: l.name})
	}
	fl := logrus.FieldLogger(l.log)
	if len(kv) != 0 {
		newFields := l.fieldsFromKV(kv...)
		fl = fl.WithFields(*newFields)
	}
	if err != nil {
		fl = fl.WithError(err)
	}
	fl.Error(msg)
}

// V will return a reference to the logger if the verbosity level is supported,
// otherwise a NoLog logger is returned
func (l *Logger) V(level int) log.InfoWriter {
	if level > l.verbosity {
		return nolog
	}
	return l
}

// SetLevel wont do anything
func (l *Logger) SetLevel(level int) {
	l.m.Lock()
	l.verbosity = level
	l.m.Unlock()
}

func (l *Logger) fieldsFromKV(kv ...log.KV) *logrus.Fields {
	newFields := logrus.Fields{}
	for i := range kv {
		newFields[kv[i].K] = kv[i].V
	}
	return &newFields
}

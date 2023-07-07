package logs

import (
	"github.com/sirupsen/logrus"
)

var l = logrus.New()

func SetNewLog(log *logrus.Logger)  {
	l = log
}

type Level logrus.Level

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

func SetLevel(level Level) {
	l.Level = logrus.Level(level)
}

func Debug(args ...interface{}) {
	l.Debug(args)
}

func Info(args ...interface{}) {
	l.Info(args)
}

func Warn(args ...interface{}) {
	l.Warn(args)
}

func Error(args ...interface{}) {
	l.Error(args)
}

func Fatal(args ...interface{}) {
	l.Fatal(args)
}

func Debugf(format string, args ...interface{}) {
	l.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	l.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	l.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	l.Fatalf(format, args...)
}

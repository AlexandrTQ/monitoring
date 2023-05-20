package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type logger struct {
	debugLogger *logrus.Logger
	errorLogger *logrus.Logger
	sqlLogger   *logrus.Logger
	httpLogger  *logrus.Logger
}

var loggerImpl logger

func InitMockLogger() {
	l := logrus.New()
	l.SetOutput(io.Discard)

	loggerImpl = logger{debugLogger: l, errorLogger: l, sqlLogger: l, httpLogger: l}
}

func Init(logPath string) error {
	err := os.MkdirAll(logPath, 0755)
	if err != nil {
		return err
	}

	debugPath := filepath.Join(logPath, "debug.log")
	debugFile, err := os.OpenFile(debugPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	errorPath := filepath.Join(logPath, "error.log")
	errorFile, err := os.OpenFile(errorPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	sqlPath := filepath.Join(logPath, "sql.log")
	sqlFile, err := os.OpenFile(sqlPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	httpPath := filepath.Join(logPath, "http.log")
	httpFile, err := os.OpenFile(httpPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	debugLogger := logrus.New()
	debugLogger.SetOutput(io.MultiWriter(os.Stdout, debugFile))
	debugLogger.SetLevel(logrus.DebugLevel)
	debugLogger.SetFormatter(&customFormatter{})

	sqlLogger := logrus.New()
	sqlLogger.SetOutput(io.MultiWriter(os.Stdout, sqlFile))
	sqlLogger.SetLevel(logrus.DebugLevel)
	sqlLogger.SetFormatter(&customFormatter{})

	errorLogger := logrus.New()
	errorLogger.SetOutput(io.MultiWriter(os.Stdout, errorFile))
	errorLogger.SetLevel(logrus.ErrorLevel)
	errorLogger.SetFormatter(&customFormatter{})

	httpLogger := logrus.New()
	httpLogger.SetOutput(io.MultiWriter(os.Stdout, httpFile))
	httpLogger.SetLevel(logrus.DebugLevel)
	httpLogger.SetFormatter(&customFormatter{})

	loggerImpl = logger{
		debugLogger: debugLogger,
		errorLogger: errorLogger,
		sqlLogger:   sqlLogger,
		httpLogger:  httpLogger,
	}

	return nil
}

type customFormatter struct{}

func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.UTC().Format("2006-01-02T15:04:05")
	level := strings.ToUpper(entry.Level.String())
	msg := entry.Message

	formattedMsg := fmt.Sprintf("%s %s: %s\n", level, timestamp, msg)
	return []byte(formattedMsg), nil
}

func Debugf(format string, args ...interface{}) {
	if loggerImpl.debugLogger == nil {
		panic("default logger not init initialized ")
	}
	loggerImpl.debugLogger.Debugf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	if loggerImpl.errorLogger == nil {
		panic("error logger not init initialized ")
	}
	loggerImpl.errorLogger.Errorf(format, args...)
}

func Sqlf(format string, args ...interface{}) {
	if loggerImpl.errorLogger == nil {
		panic("sql logger not init initialized ")
	}
	loggerImpl.sqlLogger.Debugf(format, args...)
}

func Httpf(format string, args ...interface{}) {
	if loggerImpl.httpLogger == nil {
		panic("http logger not init initialized ")
	}
	loggerImpl.httpLogger.Debugf(format, args...)
}

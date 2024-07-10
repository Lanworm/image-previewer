package logger

import (
	"errors"
	"fmt"
	"io"
)

var ErrUnknownLogLevel = errors.New("unexpected log level")

type Logger struct {
	Level LogLevel
	out   io.Writer
}

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

func New(
	level string,
	out io.Writer,
) (*Logger, error) {
	lv, err := parseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	return &Logger{
		Level: lv,
		out:   out,
	}, nil
}

func parseLevel(level string) (LogLevel, error) {
	switch level {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARNING":
		return WARNING, nil
	case "ERROR":
		return ERROR, nil
	}

	return 0, ErrUnknownLogLevel
}

func (l *Logger) Debug(msg interface{}) {
	if l.Level > DEBUG {
		return
	}

	fmt.Fprint(l.out, "DEBUG ", msg, "\n")
}

func (l *Logger) Info(msg interface{}) {
	if l.Level > INFO {
		return
	}

	fmt.Fprint(l.out, "INFO ", msg, "\n")
}

func (l *Logger) Warning(msg interface{}) {
	if l.Level > WARNING {
		return
	}

	fmt.Fprint(l.out, "WARNING ", msg, "\n")
}

func (l *Logger) Error(msg interface{}) {
	if l.Level > ERROR {
		return
	}

	fmt.Fprint(l.out, "ERROR ", msg, "\n")
}

func (l *Logger) ServerLog(msg interface{}) {
	fmt.Fprint(l.out, "SERVER ", msg, "\n")
}

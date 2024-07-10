package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerDebug(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	l := Logger{
		Level: DEBUG,
		out:   buf,
	}

	l.Debug("debug")
	l.Info("info")
	l.Warning("warning")
	l.Error("error")

	expected := `DEBUG debug
INFO info
WARNING warning
ERROR error
`

	assert.Equal(t, expected, buf.String())
}

func TestLoggerInfo(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	l := Logger{
		Level: INFO,
		out:   buf,
	}

	l.Debug("debug")
	l.Info("info")
	l.Warning("warning")
	l.Error("error")

	expected := `INFO info
WARNING warning
ERROR error
`

	assert.Equal(t, expected, buf.String())
}

func TestLoggerWarning(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	l := Logger{
		Level: WARNING,
		out:   buf,
	}

	l.Debug("debug")
	l.Info("info")
	l.Warning("warning")
	l.Error("error")

	expected := `WARNING warning
ERROR error
`

	assert.Equal(t, expected, buf.String())
}

func TestLoggerError(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	l := Logger{
		Level: ERROR,
		out:   buf,
	}

	l.Debug("debug")
	l.Info("info")
	l.Warning("warning")
	l.Error("error")

	expected := `ERROR error
`

	assert.Equal(t, expected, buf.String())
}

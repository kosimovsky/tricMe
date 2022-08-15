package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	*logrus.Logger
	file *os.File
}

func NewLogger(filename string) (*Logger, error) {
	l := logrus.New()
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	l.SetOutput(file)
	l.SetFormatter(new(logrus.JSONFormatter))
	return &Logger{l, file}, nil
}

func (l *Logger) SetWarnLevel() {
	l.SetLevel(logrus.WarnLevel)
}

func (l *Logger) Close() error {
	return l.file.Close()
}

package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type logger struct {
	*logrus.Logger
}

func NewLogger() *logger {
	return &logger{logrus.New()}
}

func (l *logger) Default(filename string) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error occured opening file : %s, %s", err.Error(), filename)
	}
	l.SetOutput(file)
	l.SetFormatter(new(logrus.JSONFormatter))
}

func (l *logger) SetWarnLevel() {
	l.SetLevel(logrus.WarnLevel)
}

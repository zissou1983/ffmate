package sev

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Format("15:04:05.000")

	// Color the log level based on the level
	var levelColor func(a ...interface{}) string
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = color.New(color.FgBlue).Sprint
	case logrus.WarnLevel:
		levelColor = color.New(color.FgYellow).Sprint
	case logrus.ErrorLevel:
		levelColor = color.New(color.FgRed).Sprint
	default:
		levelColor = color.New(color.FgWhite).Sprint
	}
	logLine := fmt.Sprintf("%s %s %s\n", timestamp, levelColor(entry.Level.String()), entry.Message)

	return []byte(logLine), nil
}

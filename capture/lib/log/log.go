package log

import (
	"github.com/OhYee/rainbow/color"
	"github.com/OhYee/rainbow/log"
)

var (
	// Error logger
	Error = log.New().SetColor(color.New().SetFrontRed()).SetOutputToStdout()
	// Info logger
	Info = log.New().SetOutputToStdout()
	// Debug logger
	Debug = log.New().SetColor(color.New().SetFrontYellow()).SetOutputToStdout()
)

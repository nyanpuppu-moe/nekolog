package log

import (
	"fmt"
	"os"
	"time"
)

const (
	colorReset = "\033[0m"

	colorGray   = "\033[90m"
	colorBlue   = "\033[34m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
)

func output(level string, msg string, args ...any) {
	fmt.Fprintf(
		os.Stdout,
		"[%s] [%s] %s\n",
		time.Now().Format("15:04:05"),
		level,
		fmt.Sprintf(msg, args...),
	)
}

func Debug(msg string, args ...any) {
	output(colorBlue+"Debug"+colorReset, msg, args...)
}

func Info(msg string, args ...any) {
	output(colorGray+"Info"+colorReset, msg, args...)
}

func Warn(msg string, args ...any) {
	output(colorYellow+"Warn"+colorReset, msg, args...)
}

func Error(msg string, args ...any) {
	output(colorRed+"Error"+colorReset, msg, args...)
}

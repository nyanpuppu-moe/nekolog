package log

import (
	"fmt"
	"os"
	"time"
)

const (
	reset = "\033[0m"

	gray   = "\033[90m"
	blue   = "\033[34m"
	yellow = "\033[33m"
	red    = "\033[31m"
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
	output(blue+"Debug"+reset, msg, args...)
}

func Info(msg string, args ...any) {
	output(gray+"Info"+reset, msg, args...)
}

func Warn(msg string, args ...any) {
	output(yellow+"Warn"+reset, msg, args...)
}

func Error(msg string, args ...any) {
	output(red+"Error"+reset, msg, args...)
}

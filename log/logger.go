package log

import (
	"fmt"
	"github.com/fatih/color"
)

func InfoF(format string, args ...any) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func Info(msg string) {
	fmt.Println(msg)
}

func SuccessF(format string, args ...any) {
	fmt.Println(fmt.Sprintf("%s\t\t%s", fmt.Sprintf(format, args...), color.GreenString("✔")))
}

func Success(msg string) {
	fmt.Println(fmt.Sprintf("%s\t\t%s", msg, color.GreenString("✔")))
}

func WarnF(format string, args ...any) {
	fmt.Println(color.YellowString(format, args...))
}

func Warn(msg string) {
	fmt.Println(color.YellowString(msg))
}

func ErrorF(format string, args ...any) {
	fmt.Println(color.RedString(format, args...))
}

func Error(msg string) {
	fmt.Println(color.RedString(msg))
}

func InfoShell(cmd string) {
	fmt.Printf("%s %s\r\n", color.GreenString("~"), cmd)
}

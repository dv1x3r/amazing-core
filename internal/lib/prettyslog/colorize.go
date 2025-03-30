package prettyslog

import (
	"fmt"
	"strconv"
)

type Color int

const (
	Black        Color = 30
	Red          Color = 31
	Green        Color = 32
	Yellow       Color = 33
	Blue         Color = 34
	Magenta      Color = 35
	Cyan         Color = 36
	LightGray    Color = 37
	DarkGray     Color = 90
	LightRed     Color = 91
	LightGreen   Color = 92
	LightYellow  Color = 93
	LightBlue    Color = 94
	LightMagenta Color = 95
	LightCyan    Color = 96
	White        Color = 97
)

func Colorize(color Color, v string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", strconv.Itoa(int(color)), v)
}

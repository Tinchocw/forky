package common

// Color type for ANSI color codes
type Color string

// ANSI color escape codes
const (
	COLOR_RESET   Color = "\x1b[0m"
	COLOR_RED     Color = "\x1b[31m"
	COLOR_YELLOW  Color = "\x1b[33m"
	COLOR_BLUE    Color = "\x1b[34m"
	COLOR_MAGENTA Color = "\x1b[35m"
	COLOR_GREEN   Color = "\x1b[32m"
	COLOR_CYAN    Color = "\x1b[36m"
	COLOR_WHITE   Color = "\x1b[37m"
	COLOR_BLACK   Color = "\x1b[30m"
)

// Colorize applies a color to a string and returns it with reset
func Colorize(text string, color Color) string {
	return string(color) + text + string(COLOR_RESET)
}

// Title creates a formatted title with "==" decoration
func Title(text string) string {
	return Colorize("== "+text+" ==", COLOR_CYAN)
}

package common

import "fmt"

// ANSI color escape codes (exported for reuse across packages)
const (
	COLOR_RESET   = "\x1b[0m"
	COLOR_RED     = "\x1b[31m"
	COLOR_YELLOW  = "\x1b[33m"
	COLOR_BLUE    = "\x1b[34m"
	COLOR_MAGENTA = "\x1b[35m"
	COLOR_GREEN   = "\x1b[32m"
	COLOR_CYAN    = "\x1b[36m"
	COLOR_WHITE   = "\x1b[37m"
	COLOR_BLACK   = "\x1b[30m"
)

// ANSI style codes
const (
	STYLE_BOLD      = "\x1b[1m"
	STYLE_DIM       = "\x1b[2m"
	STYLE_ITALIC    = "\x1b[3m"
	STYLE_UNDERLINE = "\x1b[4m"
	STYLE_BLINK     = "\x1b[5m"
	STYLE_REVERSE   = "\x1b[7m"
	STYLE_STRIKE    = "\x1b[9m"
)

// Background colors
const (
	BG_BLACK   = "\x1b[40m"
	BG_RED     = "\x1b[41m"
	BG_GREEN   = "\x1b[42m"
	BG_YELLOW  = "\x1b[43m"
	BG_BLUE    = "\x1b[44m"
	BG_MAGENTA = "\x1b[45m"
	BG_CYAN    = "\x1b[46m"
	BG_WHITE   = "\x1b[47m"
)

// ColorString applies a color to a string and returns it with reset
func ColorString(text, color string) string {
	return color + text + COLOR_RESET
}

// Colorize applies multiple ANSI codes to a string
func Colorize(text string, codes ...string) string {
	result := ""
	for _, code := range codes {
		result += code
	}
	return result + text + COLOR_RESET
}

// Predefined color functions for convenience
func Red(text string) string     { return ColorString(text, COLOR_RED) }
func Green(text string) string   { return ColorString(text, COLOR_GREEN) }
func Yellow(text string) string  { return ColorString(text, COLOR_YELLOW) }
func Blue(text string) string    { return ColorString(text, COLOR_BLUE) }
func Magenta(text string) string { return ColorString(text, COLOR_MAGENTA) }
func Cyan(text string) string    { return ColorString(text, COLOR_CYAN) }
func White(text string) string   { return ColorString(text, COLOR_WHITE) }
func Black(text string) string   { return ColorString(text, COLOR_BLACK) }

// Style functions
func Bold(text string) string      { return ColorString(text, STYLE_BOLD) }
func Dim(text string) string       { return ColorString(text, STYLE_DIM) }
func Italic(text string) string    { return ColorString(text, STYLE_ITALIC) }
func Underline(text string) string { return ColorString(text, STYLE_UNDERLINE) }
func Blink(text string) string     { return ColorString(text, STYLE_BLINK) }
func Reverse(text string) string   { return ColorString(text, STYLE_REVERSE) }
func Strike(text string) string    { return ColorString(text, STYLE_STRIKE) }

// Combined style functions - the most useful ones!
func BoldRed(text string) string     { return Colorize(text, STYLE_BOLD, COLOR_RED) }
func BoldGreen(text string) string   { return Colorize(text, STYLE_BOLD, COLOR_GREEN) }
func BoldYellow(text string) string  { return Colorize(text, STYLE_BOLD, COLOR_YELLOW) }
func BoldBlue(text string) string    { return Colorize(text, STYLE_BOLD, COLOR_BLUE) }
func BoldMagenta(text string) string { return Colorize(text, STYLE_BOLD, COLOR_MAGENTA) }
func BoldCyan(text string) string    { return Colorize(text, STYLE_BOLD, COLOR_CYAN) }

func UnderlineRed(text string) string   { return Colorize(text, STYLE_UNDERLINE, COLOR_RED) }
func UnderlineGreen(text string) string { return Colorize(text, STYLE_UNDERLINE, COLOR_GREEN) }
func UnderlineBlue(text string) string  { return Colorize(text, STYLE_UNDERLINE, COLOR_BLUE) }

// Fun rainbow function 🌈
func Rainbow(text string) string {
	colors := []string{COLOR_RED, COLOR_YELLOW, COLOR_GREEN, COLOR_CYAN, COLOR_BLUE, COLOR_MAGENTA}
	result := ""
	for i, char := range text {
		color := colors[i%len(colors)]
		result += color + string(char)
	}
	return result + COLOR_RESET
}

// Gradient function - transitions between two colors
func Gradient(text, startColor, endColor string) string {
	// For simplicity, just alternate between the two colors
	result := ""
	for i, char := range text {
		if i%2 == 0 {
			result += startColor + string(char)
		} else {
			result += endColor + string(char)
		}
	}
	return result + COLOR_RESET
}

// Box function - creates a colored box around text
func Box(text, color, bgColor string) string {
	border := "═"

	lines := []string{text} // For simplicity, assume single line
	maxWidth := len(text)

	result := ""
	// Top border
	result += color + bgColor + "╔" + repeatString(border, maxWidth+2) + "╗" + COLOR_RESET + "\n"
	// Content
	for _, line := range lines {
		result += color + bgColor + "║ " + line + " ║" + COLOR_RESET + "\n"
	}
	// Bottom border
	result += color + bgColor + "╚" + repeatString(border, maxWidth+2) + "╝" + COLOR_RESET

	return result
}

// Helper function to repeat a string
func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// Printf-style functions with colors
func PrintfRed(format string, args ...interface{}) {
	fmt.Printf(Red(format), args...)
}

func PrintfGreen(format string, args ...interface{}) {
	fmt.Printf(Green(format), args...)
}

func PrintfYellow(format string, args ...interface{}) {
	fmt.Printf(Yellow(format), args...)
}

func PrintfBlue(format string, args ...interface{}) {
	fmt.Printf(Blue(format), args...)
}

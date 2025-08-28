package prompt

import (
	"fmt"
	"os"
	"strings"
)

// ANSI color codes
const (
	// Reset
	Reset = "\033[0m"

	// Colors
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Gray    = "\033[90m"

	// Bold colors
	BoldRed     = "\033[1;31m"
	BoldGreen   = "\033[1;32m"
	BoldYellow  = "\033[1;33m"
	BoldBlue    = "\033[1;34m"
	BoldMagenta = "\033[1;35m"
	BoldCyan    = "\033[1;36m"
	BoldWhite   = "\033[1;37m"

	// Background colors
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

// ColorFunc represents a function that applies color to text
type ColorFunc func(string) string

// Color functions
var (
	RedText     = makeColorFunc(Red)
	GreenText   = makeColorFunc(Green)
	YellowText  = makeColorFunc(Yellow)
	BlueText    = makeColorFunc(Blue)
	MagentaText = makeColorFunc(Magenta)
	CyanText    = makeColorFunc(Cyan)
	WhiteText   = makeColorFunc(White)
	GrayText    = makeColorFunc(Gray)

	BoldRedText     = makeColorFunc(BoldRed)
	BoldGreenText   = makeColorFunc(BoldGreen)
	BoldYellowText  = makeColorFunc(BoldYellow)
	BoldBlueText    = makeColorFunc(BoldBlue)
	BoldMagentaText = makeColorFunc(BoldMagenta)
	BoldCyanText    = makeColorFunc(BoldCyan)
	BoldWhiteText   = makeColorFunc(BoldWhite)
)

// makeColorFunc creates a color function
func makeColorFunc(color string) ColorFunc {
	return func(text string) string {
		if !isTerminalSupportsColor() {
			return text
		}
		return color + text + Reset
	}
}

// isTerminalSupportsColor checks if the terminal supports color
func isTerminalSupportsColor() bool {
	term := os.Getenv("TERM")
	if term == "" {
		return false
	}

	// Check for common terminals that support color
	colorTerms := []string{"xterm", "xterm-256color", "screen", "tmux", "rxvt", "ansi"}
	for _, colorTerm := range colorTerms {
		if strings.Contains(term, colorTerm) {
			return true
		}
	}

	return false
}

// Themed formatting functions
func SuccessText(text string) string {
	return BoldGreenText("✅ " + text)
}

func ErrorText(text string) string {
	return BoldRedText("❌ " + text)
}

func WarningText(text string) string {
	return BoldYellowText("⚠️ " + text)
}

func InfoText(text string) string {
	return BoldBlueText("ℹ️ " + text)
}

func PromptText(text string) string {
	return BoldCyanText("🔤 " + text)
}

func PasswordText(text string) string {
	return BoldMagentaText("🔒 " + text)
}

func ConfirmText(text string) string {
	return BoldYellowText("❓ " + text)
}

func SelectText(text string) string {
	return BoldBlueText("📋 " + text)
}

func HighlightText(text string) string {
	return BoldWhiteText(text)
}

func DimText(text string) string {
	return GrayText(text)
}

// FormatOption formats an option with optional highlighting
func FormatOption(index int, text string, isSelected bool) string {
	marker := " "
	if isSelected {
		marker = BoldGreenText(">")
	}

	indexStr := fmt.Sprintf("%d)", index+1)
	if isSelected {
		indexStr = BoldWhiteText(indexStr)
		text = BoldWhiteText(text)
	} else {
		indexStr = DimText(indexStr)
	}

	return fmt.Sprintf("  %s %s %s", marker, indexStr, text)
}

// FormatDefault formats the default value display
func FormatDefault(value string) string {
	return DimText(fmt.Sprintf("[%s]", value))
}

// FormatValidationError formats validation error messages
func FormatValidationError(err error) string {
	return ErrorText(err.Error() + " Try again.")
}

package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

type prompter struct {
	reader *bufio.Reader
}

// NewPrompter creates a new interactive prompter
func NewPrompter() Prompter {
	return &prompter{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (p *prompter) Text(message string, defaultValue ...string) (string, error) {
	var defaultStr string
	if len(defaultValue) > 0 && defaultValue[0] != "" {
		defaultStr = " " + FormatDefault(defaultValue[0])
	}

	fmt.Printf("%s%s: ", PromptText(message), defaultStr)

	input, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	if input == "" && len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return input, nil
}

func (p *prompter) Password(message string) (string, error) {
	fmt.Printf("%s: ", PasswordText(message))

	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println() // New line after password input

	return string(bytePassword), nil
}

func (p *prompter) Confirm(message string, defaultValue ...bool) (bool, error) {
	var defaultStr string
	var defaultBool bool

	if len(defaultValue) > 0 {
		defaultBool = defaultValue[0]
		if defaultBool {
			defaultStr = " " + FormatDefault("Y/n")
		} else {
			defaultStr = " " + FormatDefault("y/N")
		}
	} else {
		defaultStr = " " + FormatDefault("y/n")
	}

	fmt.Printf("%s%s: ", ConfirmText(message), defaultStr)

	input, err := p.reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" && len(defaultValue) > 0 {
		return defaultBool, nil
	}

	return input == "y" || input == "yes" || input == "sim" || input == "s", nil
}

func (p *prompter) Select(message string, options []string, defaultIndex ...int) (string, error) {
	fmt.Printf("%s:\n", SelectText(message))

	for i, option := range options {
		isSelected := len(defaultIndex) > 0 && i == defaultIndex[0]
		fmt.Println(FormatOption(i, option, isSelected))
	}

	var defaultStr string
	if len(defaultIndex) > 0 && defaultIndex[0] >= 0 && defaultIndex[0] < len(options) {
		defaultStr = " " + FormatDefault(fmt.Sprintf("%d", defaultIndex[0]+1))
	}

	fmt.Printf("%s%s: ", HighlightText("Choose option"), defaultStr)

	input, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)

	if input == "" && len(defaultIndex) > 0 && defaultIndex[0] >= 0 && defaultIndex[0] < len(options) {
		return options[defaultIndex[0]], nil
	}

	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(options) {
		return "", fmt.Errorf("invalid selection: choose between 1 and %d", len(options))
	}

	return options[index-1], nil
}

func (p *prompter) MultiSelect(message string, options []string) ([]string, error) {
	fmt.Printf("%s %s:\n", SelectText(message), DimText("(separate multiple choices with commas)"))

	for i, option := range options {
		fmt.Println(FormatOption(i, option, false))
	}

	fmt.Printf("%s: ", HighlightText("Choose options (e.g., 1,3,5)"))

	input, err := p.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return []string{}, nil
	}

	var selected []string
	for _, indexStr := range strings.Split(input, ",") {
		indexStr = strings.TrimSpace(indexStr)
		index, err := strconv.Atoi(indexStr)
		if err != nil || index < 1 || index > len(options) {
			return nil, fmt.Errorf("invalid selection: %s", indexStr)
		}
		selected = append(selected, options[index-1])
	}

	return selected, nil
}

// TextWithValidation prompts for text input with validation
func (p *prompter) TextWithValidation(message string, validator ValidationFunc, defaultValue ...string) (string, error) {
	for {
		input, err := p.Text(message, defaultValue...)
		if err != nil {
			return "", err
		}

		if validator != nil {
			if err := validator(input); err != nil {
				fmt.Println(FormatValidationError(err))
				continue
			}
		}

		return input, nil
	}
}

// PasswordWithValidation prompts for password input with validation
func (p *prompter) PasswordWithValidation(message string, validator ValidationFunc) (string, error) {
	for {
		input, err := p.Password(message)
		if err != nil {
			return "", err
		}

		if validator != nil {
			if err := validator(input); err != nil {
				fmt.Println(FormatValidationError(err))
				continue
			}
		}

		return input, nil
	}
}

// Email prompts for email input with built-in validation
func (p *prompter) Email(message string, defaultValue ...string) (string, error) {
	return p.TextWithValidation(message, ValidateEmail, defaultValue...)
}

// Domain prompts for domain input with built-in validation
func (p *prompter) Domain(message string, defaultValue ...string) (string, error) {
	return p.TextWithValidation(message, ValidateDomain, defaultValue...)
}

// InteractiveMultiSelect provides an interactive multi-selection interface
func (p *prompter) InteractiveMultiSelect(message string, options []Option, preSelected ...int) ([]string, error) {
	selected := make(map[int]bool)

	// Pre-select options if provided
	for _, idx := range preSelected {
		if idx >= 0 && idx < len(options) {
			selected[idx] = true
		}
	}

	for {
		// Clear screen and show current state
		fmt.Printf("\033[2J\033[H") // Clear screen and move cursor to top

		fmt.Printf("%s:\n", SelectText(message))
		fmt.Println(DimText("Use numbers to select/deselect. Type 'done' to finish."))
		fmt.Println()

		// Show options with current selection state
		for i, option := range options {
			checkbox := "[ ]"
			if selected[i] {
				checkbox = BoldGreenText("[✓]")
			}

			status := ""
			if selected[i] {
				status = BoldGreenText(" (selected)")
			}

			fmt.Printf("  %s %s %s%s\n",
				DimText(fmt.Sprintf("%d)", i+1)),
				checkbox,
				option.Label,
				status)
		}

		fmt.Println()

		// Show current selection count
		selectionCount := len(getSelectedValues(selected))
		if selectionCount > 0 {
			fmt.Printf("%s: %d options selected\n",
				InfoText("Status"),
				selectionCount)
		} else {
			fmt.Println(DimText("No options selected"))
		}

		fmt.Println()
		fmt.Printf("%s: ", HighlightText("Enter option number or 'done' to finish"))

		input, err := p.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		input = strings.TrimSpace(strings.ToLower(input))

		if input == "done" || input == "d" {
			break
		}

		// Try to parse as number
		if idx, err := strconv.Atoi(input); err == nil {
			if idx >= 1 && idx <= len(options) {
				idx--                          // Convert to 0-based index
				selected[idx] = !selected[idx] // Toggle selection
			} else {
				fmt.Printf("%s\n", ErrorText("Invalid number. Press Enter to continue..."))
				p.reader.ReadString('\n')
			}
		} else {
			fmt.Printf("%s\n", ErrorText("Invalid input. Use numbers or 'done'. Press Enter to continue..."))
			p.reader.ReadString('\n')
		}
	}

	// Build result
	var result []string
	for i, option := range options {
		if selected[i] {
			result = append(result, option.Value)
		}
	}

	return result, nil
}

// Helper function to get selected values
func getSelectedValues(selected map[int]bool) []int {
	var values []int
	for idx, isSelected := range selected {
		if isSelected {
			values = append(values, idx)
		}
	}
	return values
}

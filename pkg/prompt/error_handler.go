package prompt

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// ErrorHandler provides enhanced error handling for interactive prompts
type ErrorHandler struct {
	prompter Prompter
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(prompter Prompter) *ErrorHandler {
	return &ErrorHandler{
		prompter: prompter,
	}
}

// WithGracefulInterrupt wraps a function with graceful interrupt handling
func (eh *ErrorHandler) WithGracefulInterrupt(ctx context.Context, fn func() error) error {
	// Setup interrupt handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Channel to capture function result
	resultCh := make(chan error, 1)

	// Run the function in a goroutine
	go func() {
		resultCh <- fn()
	}()

	select {
	case <-sigCh:
		fmt.Println(WarningText("\n⚠️ Operation interrupted by user"))

		// Ask if user wants to save partial progress
		save, err := eh.prompter.Confirm("Do you want to save partial progress?", false)
		if err != nil {
			return fmt.Errorf("error confirming save: %w", err)
		}

		if save {
			return fmt.Errorf("operation interrupted with save requested")
		}

		return fmt.Errorf("operation cancelled by user")

	case err := <-resultCh:
		return err

	case <-ctx.Done():
		return ctx.Err()
	}
}

// RetryOnError retries an operation on specific errors with user confirmation
func (eh *ErrorHandler) RetryOnError(operation func() error, maxRetries int, errorMsg string) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		if attempt == maxRetries {
			break
		}

		// Show error and ask for retry
		fmt.Println(ErrorText(fmt.Sprintf("Error: %v", err)))

		retry, promptErr := eh.prompter.Confirm(
			fmt.Sprintf("%s. Try again? (attempt %d/%d)", errorMsg, attempt+1, maxRetries+1),
			true,
		)
		if promptErr != nil {
			return fmt.Errorf("error confirming retry: %w", promptErr)
		}

		if !retry {
			return fmt.Errorf("operation cancelled by user after %d attempts", attempt+1)
		}

		fmt.Println(InfoText("Retrying..."))
	}

	return fmt.Errorf("operation failed after %d attempts: %w", maxRetries+1, lastErr)
}

// ValidateWithRetry validates input with retry mechanism
func (eh *ErrorHandler) ValidateWithRetry(
	promptMsg string,
	validator ValidationFunc,
	maxRetries int,
	defaultValue ...string,
) (string, error) {
	var result string
	var lastInput string

	err := eh.RetryOnError(func() error {
		input, err := eh.prompter.TextWithValidation(promptMsg, validator, defaultValue...)
		if err != nil {
			return err
		}
		lastInput = input
		return nil
	}, maxRetries, "Invalid input")

	if err != nil {
		return "", err
	}

	result = lastInput
	return result, nil
}

// SafeOperation wraps an operation with comprehensive error handling
func (eh *ErrorHandler) SafeOperation(ctx context.Context, name string, operation func() error) error {
	fmt.Println(InfoText("🚀 Starting: " + name))

	return eh.WithGracefulInterrupt(ctx, func() error {
		return eh.RetryOnError(operation, 2, "Operation failed: "+name)
	})
}

// RecoverFromPanic recovers from panics and provides user-friendly error messages
func (eh *ErrorHandler) RecoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Println(ErrorText("❌ Critical error detected!"))
		fmt.Printf("%s: %v\n", DimText("Details"), r)

		// Ask user if they want to report the issue
		report, err := eh.prompter.Confirm("Do you want to report this error to the developers?", true)
		if err == nil && report {
			fmt.Println(InfoText("Please open an issue at: https://github.com/leandrodaf/harborctl/issues"))
			fmt.Println(DimText("Include the error details above in your report."))
		}

		os.Exit(1)
	}
}

// ConfirmDestructiveOperation confirms dangerous operations with extra verification
func (eh *ErrorHandler) ConfirmDestructiveOperation(operationName string, target string) error {
	fmt.Println(WarningText("⚠️ WARNING: Destructive operation detected!"))
	fmt.Printf("%s: %s\n", BoldYellowText("Operation"), operationName)
	fmt.Printf("%s: %s\n", BoldYellowText("Target"), target)
	fmt.Println()

	// First confirmation
	confirm1, err := eh.prompter.Confirm("Are you sure you want to continue?", false)
	if err != nil {
		return err
	}
	if !confirm1 {
		return fmt.Errorf("operation cancelled by user")
	}

	// Second confirmation with typing
	confirmText, err := eh.prompter.Text(
		fmt.Sprintf("Type '%s' to confirm", BoldRedText("CONFIRM")),
	)
	if err != nil {
		return err
	}

	if confirmText != "CONFIRM" {
		return fmt.Errorf("invalid confirmation. Operation cancelled")
	}

	fmt.Println(SuccessText("✅ Operation confirmed"))
	return nil
}

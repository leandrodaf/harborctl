package main

import (
	"context"
	"fmt"
	"log"

	"github.com/leandrodaf/harborctl/pkg/prompt"
)

// Simple example demonstrating HarborCtl prompt features
func main() {
	ctx := context.Background()

	// Create prompter and error handler
	prompter := prompt.NewPrompter()
	errorHandler := prompt.NewErrorHandler(prompter)

	fmt.Println("🚀 HarborCtl Prompt System Example")
	fmt.Println()

	// Example 1: Basic prompts with colors
	fmt.Println("=== 1. Basic prompts with colors ===")

	name, err := prompter.Text("What's your name?", "User")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Hello, %s!\n\n", name)

	// Example 2: Email validation
	fmt.Println("=== 2. Email validation ===")

	email, err := prompter.Email("Enter your email")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Valid email: %s\n\n", email)

	// Example 3: Domain validation
	fmt.Println("=== 3. Domain validation ===")

	domain, err := prompter.Domain("Enter your domain", "example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Valid domain: %s\n\n", domain)

	// Example 4: Selection
	fmt.Println("=== 4. Environment selection ===")

	environment, err := prompter.Select("Choose environment", []string{
		"Local Development",
		"Production",
	}, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Selected environment: %s\n\n", environment)

	// Example 5: Confirmation
	fmt.Println("=== 5. Confirmation ===")

	includeServices, err := prompter.Confirm("Include observability services?", true)
	if err != nil {
		log.Fatal(err)
	}

	if includeServices {
		fmt.Println("✅ Observability services will be included")
	} else {
		fmt.Println("❌ Observability services will be skipped")
	}

	// Example 6: Error handling
	fmt.Println("\n=== 6. Error handling ===")

	err = errorHandler.SafeOperation(ctx, "Example operation", func() error {
		confirm, err := prompter.Confirm("Confirm operation?", true)
		if err != nil {
			return err
		}

		if !confirm {
			return fmt.Errorf("operation not confirmed")
		}

		fmt.Println("✅ Operation completed successfully!")
		return nil
	})

	if err != nil {
		log.Printf("Operation error: %v", err)
	}

	fmt.Println()
	fmt.Println("🎉 Example completed!")
	fmt.Println()
	fmt.Println("Features demonstrated:")
	fmt.Println("• ✅ Colored prompts with emojis")
	fmt.Println("• ✅ Automatic email validation")
	fmt.Println("• ✅ Automatic domain validation")
	fmt.Println("• ✅ Enhanced selection interface")
	fmt.Println("• ✅ Confirmation prompts")
	fmt.Println("• ✅ Error handling with retry")
	fmt.Println("• ✅ Graceful interruption (Ctrl+C)")
}

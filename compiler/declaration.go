package compiler

import (
	"fmt"
	"regexp"
	"strings"
)

func NewEntry() string {
	return "func main() {\n"
}

func Entry(msg string) *CompilerCommand {
	msg = strings.TrimSpace(msg)

	format := `margiang utama`
	re := regexp.MustCompile(format)
	match := re.FindStringSubmatch(msg)
	if match == nil {
		return nil
	}

	return &CompilerCommand{
		Syntax: "package main",
		Entry:  true,
	}
}

// Var parses the input and returns the corresponding CompilerCommand
func Var(msg string) *CompilerCommand {
	// Define the pattern to match the input format
	format := `teges ([a-zA-Z_]+[a-zA-Z0-9_]*) ne (.+)`
	re := regexp.MustCompile(format)
	match := re.FindStringSubmatch(msg)
	if match == nil {
		return nil // Return nil if the format doesn't match
	}

	// Extract variable name and value
	varName := match[1]
	value := match[2]

	// Transform the value if necessary
	transformedValue := valueTransform(value)

	// Create the assignment expression
	syntax := fmt.Sprintf("var %s = %s", varName, transformedValue)

	// Return the result in a CompilerCommand struct
	return &CompilerCommand{
		Syntax: syntax,
	}
}

// If parses the input string and returns a Go if condition
func If(msg string) *CompilerCommand {
	msg = strings.TrimSpace(msg)

	if strings.HasPrefix(msg, "tiosan") {
		return ElseIf(msg)
	}

	// Define the format pattern to match:
	// "yen urutan ne <value>:", "yen urutan ne sing <value>:",
	// "yen urutan gedenan ken <value>:", "yen urutan cenikan ken <value>:"
	format := regexp.MustCompile(`yen ([a-zA-Z0-9_]+) (ne|ne sing|gedenan ken|cenikan ken) ([^\n]+)`)
	match := format.FindStringSubmatch(msg)
	if match == nil {
		return nil
	}

	// Determine the operator based on the match
	operator := mapComparable[match[2]]

	// Apply value transformation to the third part (value)
	conditionValue := strings.TrimSuffix(valueTransform(match[3]), ":")

	// Construct Go if condition
	goIf := fmt.Sprintf("if %s %s %s", match[1], operator, conditionValue)
	return &CompilerCommand{
		Syntax:    goIf,
		OpenGroup: true,
	}
}

// ElseIf parses the input string and returns a Go if condition
func ElseIf(msg string) *CompilerCommand {
	// Trim any leading or trailing spaces
	msg = strings.TrimSpace(msg)

	// Adjusting the format to capture the phrase "tiosan yen"
	format := regexp.MustCompile(`tiosan yen ([a-zA-Z0-9]+) ([a-zA-Z ]+) ([^\[\]\(\)\n]+)`)
	match := format.FindStringSubmatch(msg)
	if match == nil {
		return nil
	}

	// Transform the second captured group (operator) using mapCompare
	operator := mapComparable[match[2]]
	if operator == "" {
		// If no operator found, just return without transformation
		operator = match[2]
	}

	// Construct the Go code for "else if" statement
	goLog := fmt.Sprintf("else if %s %s %s", match[1], operator, strings.TrimSuffix(valueTransform(match[3]), ":"))

	return &CompilerCommand{
		Syntax:     goLog,
		OpenGroup:  true,
		CloseGroup: true,
	}
}

// Else parses the input string and returns a Go if condition
func Else(msg string) *CompilerCommand {
	// Trim any leading or trailing spaces
	msg = strings.TrimSpace(msg)

	// Define the regular expression pattern to match the input string "tiosan"
	format := `tiosan:`
	re := regexp.MustCompile(format)
	match := re.FindStringSubmatch(msg)

	if match == nil {
		return nil
	}

	return &CompilerCommand{
		Syntax:     "else",
		OpenGroup:  true,
		CloseGroup: true,
	}
}

func CloseIf(msg string) *CompilerCommand {
	// Define the regular expression pattern to match the input string "suud"
	format := `suud$`
	re := regexp.MustCompile(format)
	match := re.FindStringSubmatch(msg)

	if match == nil {
		return nil
	}

	return &CompilerCommand{
		Syntax:     "",
		CloseGroup: true,
	}
}

func Print(msg string) *CompilerCommand {
	// Trim any leading or trailing spaces
	msg = strings.TrimSpace(msg)

	// Define the format pattern to match "pesuang <message>"
	format := regexp.MustCompile(`pesuang(.*)`)
	match := format.FindStringSubmatch(msg)
	if match == nil {
		return nil
	}

	// Construct Go print statement using fmt.Println for correct output
	goLog := fmt.Sprintf(`print%s`, match[1])

	// Return the compiled command
	return &CompilerCommand{
		Syntax: goLog,
	}
}

func ReturnEmpty(msg string) *CompilerCommand {
	msg = strings.TrimSpace(msg)

	format := regexp.MustCompile(`uliang$`)
	match := format.FindStringSubmatch(msg)
	if match == nil {
		return nil
	}

	goCmd := "return"

	return &CompilerCommand{
		Syntax:    goCmd,
		Returning: true,
	}
}

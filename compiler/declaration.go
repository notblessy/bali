package compiler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/notblessy/bali/anggenan"
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

func Import(msg string) *CompilerCommand {
	msg = strings.TrimSpace(msg)

	// Regex to match `anggen "package_name"`
	format := `^anggen\s+"([^"]+)"$`
	re := regexp.MustCompile(format)
	match := re.FindStringSubmatch(msg)

	if match == nil {
		return nil
	}

	packageName := match[1]

	return &CompilerCommand{
		Syntax:      fmt.Sprintf("\"%s\"", packageName),
		IsImporting: true,
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

	var syntaxValue string

	pkgName, identifier, bracketValue, isPackageUsage := getPotentialPackageUsage(value)
	if isPackageUsage {
		syntaxValue = fmt.Sprintf("%s.%s%s", pkgName, anggenan.PackageInterpreter[pkgName][identifier], bracketValue)
	} else {
		// Transform the value if necessary
		syntaxValue = valueTransform(value)
	}

	// Create the assignment expression
	syntax := fmt.Sprintf("var %s = %s", varName, syntaxValue)

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

	format := regexp.MustCompile(`yen ([^:]+):`)
	match := format.FindStringSubmatch(msg)
	if match == nil {
		return nil
	}

	// Extract and parse condition expression
	condition := strings.TrimSpace(match[1])

	// Handle logical operators
	condition = strings.ReplaceAll(condition, "lan", "&&")
	condition = strings.ReplaceAll(condition, "utawi", "||")

	// Transform each part of the condition if needed
	parts := regexp.MustCompile(`([a-zA-Z0-9_]+) (ne|ne sing|gedenan ken|cenikan ken) ([^\s]+)`)
	condition = parts.ReplaceAllStringFunc(condition, func(part string) string {
		matches := parts.FindStringSubmatch(part)
		if matches == nil {
			return part
		}
		varName, comparator, value := matches[1], matches[2], matches[3]
		operator := mapComparable[comparator]
		return fmt.Sprintf("%s %s %s", varName, operator, valueTransform(value))
	})

	// Construct the Go if condition
	goIf := fmt.Sprintf("if %s", condition)

	return &CompilerCommand{
		Syntax:    goIf,
		OpenGroup: true,
	}
}

// ElseIf parses the input string and returns a Go else if condition
func ElseIf(msg string) *CompilerCommand {
	// Trim any leading or trailing spaces
	msg = strings.TrimSpace(msg)

	// Adjusting the format to capture the phrase "tiosan yen"
	format := regexp.MustCompile(`tiosan yen ([^:]+):`)
	match := format.FindStringSubmatch(msg)
	if match == nil {
		return nil
	}

	// Extract and parse condition expression
	condition := strings.TrimSpace(match[1])

	// Handle logical operators
	condition = strings.ReplaceAll(condition, "lan", "&&")
	condition = strings.ReplaceAll(condition, "utawi", "||")

	// Handle comparisons and transform condition parts
	parts := regexp.MustCompile(`([a-zA-Z0-9_]+) (ne|ne sing|gedenan ken|cenikan ken) ([^\s]+)`)
	condition = parts.ReplaceAllStringFunc(condition, func(part string) string {
		matches := parts.FindStringSubmatch(part)
		if matches == nil {
			return part
		}
		varName, comparator, value := matches[1], matches[2], matches[3]
		operator := mapComparable[comparator]
		return fmt.Sprintf("%s %s %s", varName, operator, valueTransform(value))
	})

	// Construct the Go else if condition
	goLog := fmt.Sprintf("else if %s", condition)

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

	// Construct Go print statement using print for correct output
	goLog := fmt.Sprintf(`fmt.Println%s`, match[1])

	// Return the compiled command
	return &CompilerCommand{
		Syntax:     goLog,
		IsPrinting: true,
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
		Syntax:      goCmd,
		IsReturning: true,
	}
}

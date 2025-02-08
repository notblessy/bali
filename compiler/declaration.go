package compiler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/notblessy/bali/anggenan"
)

// NewEntry :nodoc:
func NewEntry() string {
	return "func main() {\n"
}

// Entry parses the input string and returns a Go entry point
func Entry(msg string) CompilerCommand {
	msg = strings.TrimSpace(msg)

	format := `margiang utama`
	re := regexp.MustCompile(format)
	match := re.FindStringSubmatch(msg)
	if match == nil {
		return nil
	}

	cmd := NewCommand("package main")
	cmd.Toggle("entry", true)

	return cmd
}

// Import parses the input string and returns a Go import statement
func Import(msg string) CompilerCommand {
	msg = strings.TrimSpace(msg)

	format := `^anggen\s+"([^"]+)"$`
	re := regexp.MustCompile(format)
	match := re.FindStringSubmatch(msg)

	if match == nil {
		return nil
	}

	packageName := match[1]

	cmd := NewCommand(fmt.Sprintf("\"%s\"", packageName))
	cmd.Toggle("importing", true)

	return cmd
}

// Var parses the input and returns the corresponding CompilerCommand
func Var(msg string) CompilerCommand {
	format := `teges ([a-zA-Z_]+[a-zA-Z0-9_]*) ne (.+)`

	re := regexp.MustCompile(format)
	match := re.FindStringSubmatch(msg)
	if match == nil {
		return nil
	}

	varName := match[1]
	value := match[2]

	var syntaxValue string

	pkgName, identifier, bracketValue, isPackageUsage := getPotentialPackageUsage(value)
	if isPackageUsage {
		syntaxValue = fmt.Sprintf("%s.%s%s", pkgName, anggenan.PackageInterpreter[pkgName][identifier], bracketValue)
	} else {
		syntaxValue = valueTransform(value)
	}

	return NewCommand(fmt.Sprintf("var %s = %s", varName, syntaxValue))
}

// If parses the input string and returns a Go if condition
func If(msg string) CompilerCommand {
	msg = strings.TrimSpace(msg)

	if strings.HasPrefix(msg, "tiosan") {
		return ElseIf(msg)
	}

	format := regexp.MustCompile(`yen ([^:]+):`)
	match := format.FindStringSubmatch(msg)
	if match == nil {
		return nil
	}

	condition := strings.TrimSpace(match[1])

	condition = strings.ReplaceAll(condition, "lan", "&&")
	condition = strings.ReplaceAll(condition, "utawi", "||")

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

	goIf := fmt.Sprintf("if %s", condition)

	cmd := NewCommand(goIf)
	cmd.Toggle("opengroup", true)

	return cmd
}

// ElseIf parses the input string and returns a Go else if condition
func ElseIf(msg string) CompilerCommand {
	msg = strings.TrimSpace(msg)

	format := regexp.MustCompile(`tiosan yen ([^:]+):`)
	match := format.FindStringSubmatch(msg)
	if match == nil {
		return nil
	}

	condition := strings.TrimSpace(match[1])

	condition = strings.ReplaceAll(condition, "lan", "&&")
	condition = strings.ReplaceAll(condition, "utawi", "||")

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

	goLog := fmt.Sprintf("else if %s", condition)
	cmd := NewCommand(goLog)
	cmd.Toggle("opengroup", true)
	cmd.Toggle("closegroup", true)

	return cmd
}

// Else parses the input string and returns a Go if condition
func Else(msg string) CompilerCommand {
	msg = strings.TrimSpace(msg)

	format := `tiosan:`
	re := regexp.MustCompile(format)
	match := re.FindStringSubmatch(msg)

	if match == nil {
		return nil
	}

	cmd := NewCommand("else")
	cmd.Toggle("opengroup", true)
	cmd.Toggle("closegroup", true)

	return cmd
}

// CloseStatement parses the input string and returns a Go closing statement
func CloseStatement(msg string) CompilerCommand {
	format := `suud$`

	re := regexp.MustCompile(format)
	match := re.FindStringSubmatch(msg)

	if match == nil {
		return nil
	}

	cmd := NewCommand("")
	cmd.Toggle("closegroup", true)

	return cmd
}

// Print parses the input string and returns a Go print statement
func Print(msg string) CompilerCommand {
	msg = strings.TrimSpace(msg)

	format := regexp.MustCompile(`pesuang(.*)`)
	match := format.FindStringSubmatch(msg)
	if match == nil {
		return nil
	}

	goLog := fmt.Sprintf(`fmt.Println%s`, match[1])

	cmd := NewCommand(goLog)
	cmd.Toggle("printing", true)

	return cmd
}

// ReturnEmpty parses the input string and returns a Go return statement
func ReturnEmpty(msg string) CompilerCommand {
	msg = strings.TrimSpace(msg)

	format := regexp.MustCompile(`uliang$`)
	match := format.FindStringSubmatch(msg)
	if match == nil {
		return nil
	}

	goCmd := "return"

	cmd := NewCommand(goCmd)
	cmd.Toggle("returning", true)

	return cmd
}

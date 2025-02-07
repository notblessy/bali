package compiler

import (
	"regexp"
)

var mapComparable = map[string]string{
	"gadenan ken": ">",
	"cenikan ken": "<",
	"ne":          "==",
	"ne sing":     "!=",
	"lan":         "&&",
	"utawi":       "||",
}

type CompilerCommand struct {
	Syntax      string
	OpenGroup   bool
	CloseGroup  bool
	Entry       bool
	IsReturning bool
	IsImporting bool
	IsPrinting  bool
}

// valueTransform simulates value transformation (e.g., removing quotes for strings)
func valueTransform(value string) string {
	return value
}

// getPotentialPackageUsage parses the input string and returns the package name, identifier, and bracket value
func getPotentialPackageUsage(value string) (string, string, string, bool) {
	pattern := `^([a-zA-Z_][a-zA-Z0-9_]*)\.([a-zA-Z_][a-zA-Z0-9_]*)(\((.*?)\)|\[(.*?)\])?`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(value)

	if len(match) > 2 {
		packageName := match[1]
		identifier := match[2]
		bracketValue := match[3]

		return packageName, identifier, bracketValue, true
	}

	return "", "", "", false
}

package compiler

var mapComparable = map[string]string{
	"gadenan ken": ">",
	"cenikan ken": "<",
	"ne":          "==",
	"ne sing":     "!=",
}

type CompilerCommand struct {
	Syntax     string
	OpenGroup  bool
	CloseGroup bool
	Entry      bool
	Returning  bool
}

// valueTransform simulates value transformation (e.g., removing quotes for strings)
func valueTransform(value string) string {
	// Example: If value is a number, keep it as is. You can expand this as needed.
	return value
}

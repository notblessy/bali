package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/notblessy/bali/compiler"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go [run|build] <filename>")
		os.Exit(1)
	}

	command := os.Args[1] // "run" or "build"
	filename := os.Args[2]

	tempDir := ".tempbuilds"
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		fmt.Println("Error creating temp directory:", err)
		os.Exit(1)
	}

	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Parse the input and convert it to Go code
	output := compiler.ToGolang(extract(string(content)))

	tempFile := filepath.Join(tempDir, "output.go")
	err = os.WriteFile(tempFile, []byte(output), 0644)
	if err != nil {
		fmt.Println("Error writing temp Go file:", err)
		os.Exit(1)
	}

	defer cleanup()

	switch command {
	case "run":
		cmd := exec.Command("go", "run", tempFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error running generated Go code:", err)
			os.Exit(1)
		}

	case "build":
		outputBinary := strings.TrimSuffix(filename, ".bali")

		cmd := exec.Command("go", "build", "-o", outputBinary, tempFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error building binary:", err)
			os.Exit(1)
		}

		fmt.Println("Build successfully:", outputBinary)

	default:
		fmt.Println("Unknown command. Use 'run' or 'build'.")
		os.Exit(1)
	}
}

func extract(input string) []compiler.CompilerCommand {
	// Step 1: Replace \r\n with \n and split the input into lines
	cmdLines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")

	// Step 2: Filter out empty lines
	var nonEmptyCmds []string
	for _, cmd := range cmdLines {
		if cmd != "" {
			nonEmptyCmds = append(nonEmptyCmds, cmd)
		}
	}

	cmds := compiler.GetCompilerCommand(nonEmptyCmds)

	return cmds
}

func cleanup() {
	tempDir := ".tempbuilds"
	err := os.RemoveAll(tempDir)
	if err != nil {
		fmt.Println("Warning: Failed to delete temp files:", err)
	}
}

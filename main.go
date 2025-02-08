package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/notblessy/bali/compiler"
)

const baliVersion = "0.0.3"

func main() {
	versionFlag := flag.Bool("version", false, "Print the Bali language version")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("Basa Bali Version: %s\n", baliVersion)
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		fmt.Println("USAGE: go run main.go [run|build] <filename>")
		os.Exit(1)
	}

	// accepts "run" or "build" command
	command := os.Args[1]
	filename := os.Args[2]

	tempDir := filepath.Join(os.Getenv("HOME"), ".local", "bin", "tempbuilds")
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		fmt.Println("[DIR] Error creating temp directory:", err)
		os.Exit(1)
	}

	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("[DIR] Error reading file:", err)
		os.Exit(1)
	}

	// Parse the input and convert it to Go code
	output := compiler.ToGolang(extract(string(content)))

	tempFile := filepath.Join(tempDir, "output.go")
	err = os.WriteFile(tempFile, []byte(output), 0644)
	if err != nil {
		fmt.Println("[DIR] Error writing temp Go file:", err)
		os.Exit(1)
	}

	defer cleanup(tempDir)

	switch command {
	case "run":
		cmd := exec.Command("go", append([]string{"run", tempFile}, os.Args[3:]...)...)
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

func cleanup(tempDir string) {
	err := os.RemoveAll(tempDir)
	if err != nil {
		fmt.Println("Warning: Failed to delete temp files:", err)
	}
}

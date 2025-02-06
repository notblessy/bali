package compiler

import "fmt"

func GetCompilerCommand(lines []string) []CompilerCommand {
	commands := []func(string) *CompilerCommand{
		Entry,
		Var,
		If,
		Else,
		CloseIf,
		Print,
		ReturnEmpty,
	}

	var compilerCommands []CompilerCommand

	for _, line := range lines {
		for _, command := range commands {
			if c := command(line); c != nil {
				compilerCommands = append(compilerCommands, *c)
			}
		}
	}
	return compilerCommands
}

func ToGolang(commands []CompilerCommand) string {
	var goCommands string
	isOpenGroup := false
	hasEntry := false

	for _, cmd := range commands {
		currSyntax := cmd.Syntax

		if cmd.Entry {
			currSyntax = currSyntax + "\nfunc main() {"
			hasEntry = true
		}

		if cmd.CloseGroup {
			currSyntax = "} " + currSyntax
			isOpenGroup = false
		}

		if cmd.OpenGroup {
			currSyntax = currSyntax + " {"
			isOpenGroup = true
		}

		if cmd.Returning {
			currSyntax = fmt.Sprintf("%s\n}", currSyntax)
		}

		goCommands += currSyntax + "\n"
	}

	if isOpenGroup {
		goCommands += " }"
	}

	if hasEntry {
		goCommands = goCommands + "}"
	}

	return goCommands
}

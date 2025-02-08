package compiler

import "fmt"

func GetCompilerCommand(lines []string) []CompilerCommand {
	commands := []func(string) CompilerCommand{
		Entry,
		Import,
		Var,
		If,
		Else,
		CloseStatement,
		Print,
		ReturnEmpty,
	}

	var compilerCommands []CompilerCommand

	for _, line := range lines {
		for _, command := range commands {
			if c := command(line); c != nil {
				compilerCommands = append(compilerCommands, c)
			}
		}
	}
	return compilerCommands
}

func ToGolang(commands []CompilerCommand) string {
	var goCommands string
	isOpenGroup := false
	hasEntry := false

	isPrinting := func() bool {
		for _, cmd := range commands {
			if cmd.IsPrinting() {
				return true
			}
		}

		return false
	}

	for _, cmd := range commands {
		currSyntax := cmd.Syntax()

		if cmd.IsImporting() {
			continue
		}

		if cmd.IsEntry() {
			for i, importcmd := range commands {
				if i == 0 {
					continue
				}

				if !commands[i-1].IsImporting() && !importcmd.IsImporting() {
					continue
				}

				if !commands[i-1].IsImporting() && importcmd.IsImporting() {
					currSyntax = currSyntax + "\nimport (\n" + importcmd.Syntax()
				}

				if commands[i-1].IsImporting() && importcmd.IsImporting() {
					currSyntax = currSyntax + "\n" + importcmd.Syntax()
				}

				if commands[i-1].IsImporting() && !importcmd.IsImporting() {
					if isPrinting() {
						currSyntax = currSyntax + "\n" + "\"fmt\""
					}

					currSyntax = currSyntax + "\n)"
				}
			}

			currSyntax = currSyntax + "\nfunc main() {"
			hasEntry = true
		}

		if cmd.IsCloseGroup() {
			currSyntax = "} " + currSyntax
			isOpenGroup = false
		}

		if cmd.IsOpenGroup() {
			currSyntax = currSyntax + " {"
			isOpenGroup = true
		}

		if cmd.IsReturning() {
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

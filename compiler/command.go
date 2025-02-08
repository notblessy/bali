package compiler

// CompilerCommand :nodoc:
type CompilerCommand interface {
	Toggle(label string, flag bool)
	Syntax() string
	IsOpenGroup() bool
	IsCloseGroup() bool
	IsEntry() bool
	IsReturning() bool
	IsImporting() bool
	IsPrinting() bool
}

type cmd struct {
	syntax       string
	isOpenGroup  bool
	isCloseGroup bool
	isEntry      bool
	isReturning  bool
	isImporting  bool
	isPrinting   bool
}

// NewCommand :nodoc:
func NewCommand(syntax string) CompilerCommand {
	return &cmd{
		syntax: syntax,
	}
}

// Toggle toggles the flag of the command
func (s *cmd) Toggle(label string, flag bool) {
	switch label {
	case "opengroup":
		s.isOpenGroup = flag
	case "closegroup":
		s.isCloseGroup = flag
	case "entry":
		s.isEntry = flag
	case "returning":
		s.isReturning = flag
	case "importing":
		s.isImporting = flag
	case "printing":
		s.isPrinting = flag
	default:
		return
	}
}

// Syntax :nodoc:
func (s *cmd) Syntax() string {
	return s.syntax
}

// IsOpenGroup :nodoc:
func (s *cmd) IsOpenGroup() bool {
	return s.isOpenGroup
}

// IsCloseGroup :nodoc:
func (s *cmd) IsCloseGroup() bool {
	return s.isCloseGroup
}

// IsEntry :nodoc:
func (s *cmd) IsEntry() bool {
	return s.isEntry
}

// IsReturning :nodoc:
func (s *cmd) IsReturning() bool {
	return s.isReturning
}

// IsImporting :nodoc:
func (s *cmd) IsImporting() bool {
	return s.isImporting
}

// IsPrinting :nodoc:
func (s *cmd) IsPrinting() bool {
	return s.isPrinting
}

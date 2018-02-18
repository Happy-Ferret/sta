package commands

// types
type Command struct {
	Cmd  string
	Args []string
}

type UnknownCommandError struct {
	Cmd Command
}

func (err UnknownCommandError) Error() string {
	return "Unknown command: " + err.Cmd.Cmd
}


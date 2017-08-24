package commands

import (
	"testing"
	"fmt"
)

func TestTerminalCommand_NoArgs(t *testing.T) {
	text := []string{"cli"}

	command := TerminalCommand{
		func(params map[string]string, args []string) {
			assert(0, len(params), t)
			assert(1, len(args), t)
		},
	}

	command.Execute(text)
}

func TestTerminalCommand_HandlerExecuted(t *testing.T) {
	text := []string{"command"}

	executed := false
	command := TerminalCommand{
		func(params map[string]string, args []string) {
			executed = true
		},
	}

	command.Execute(text)
	if !executed {
		t.Error("handler is not executed")
	}
}

func TestTerminalCommand_OneArg(t *testing.T) {
	text := []string{"command", "arg1"}

	command := TerminalCommand{
		func(params map[string]string, args []string) {
			assert(0, len(params), t)
			assert(2, len(args), t)
			assert("command", args[0], t)
			assert("arg1", args[1], t)
		},
	}

	command.Execute(text)
}

func TestTerminalCommand_TwoArgs(t *testing.T) {
	text := []string{"com", "arg1", "arg2"}

	command := TerminalCommand{
		func(params map[string]string, args []string) {
			assert(0, len(params), t)
			assert(3, len(args), t)
			assert("com", args[0], t)
			assert("arg1", args[1], t)
			assert("arg2", args[2], t)
		},
	}

	command.Execute(text)
}

func TestTerminalCommand_OneParamWithTwoDashes(t *testing.T) {
	text := []string{"cli", "--version"}

	command := TerminalCommand{
		func(params map[string]string, args []string) {
			assert(1, len(args), t)
			assert(1, len(params), t)
			_, ok := params["version"]
			assert(true, ok, t)
		},
	}

	command.Execute(text)
}

func TestTerminalCommand_OneParamWithOneDash(t *testing.T) {
	text := []string{"cli", "-h"}

	command := TerminalCommand{
		func(params map[string]string, args []string) {
			assert(1, len(args), t)
			assert(1, len(params), t)
			_, ok := params["h"]
			assert(true, ok, t)
		},
	}

	command.Execute(text)
}

func TestTerminalCommand_ParamWithValue(t *testing.T) {
	text := []string{"cli", "--location.storage=/dev/null"}

	command := TerminalCommand{
		func(params map[string]string, args []string) {
			assert(1, len(args), t)
			assert(1, len(params), t)
			assert("/dev/null", params["location.storage"], t)
		},
	}

	command.Execute(text)
}

func TestTerminalCommand_ParamWithValueParamWithoutValue(t *testing.T) {
	text := []string{"cli", "--location.storage=/dev/null", "--version"}

	command := TerminalCommand{
		func(params map[string]string, args []string) {
			assert(1, len(args), t)
			assert(2, len(params), t)
			assert("/dev/null", params["location.storage"], t)
			_, ok := params["version"]
			assert(true, ok, t)
		},
	}

	command.Execute(text)
}

func TestTerminalCommand_ParamAndArg(t *testing.T) {
	text := []string{"cli", "--location.storage=/dev/null", "copy"}

	command := TerminalCommand{
		func(params map[string]string, args []string) {
			assert(2, len(args), t)
			assert(1, len(params), t)
			assert("/dev/null", params["location.storage"], t)
			assert("copy", args[1], t)
		},
	}

	command.Execute(text)
}

func TestTerminalCommand_LonelyParamAndArg(t *testing.T) {
	text := []string{"cli", "--turbo", "flash"}

	command := TerminalCommand{
		func(params map[string]string, args []string) {
			assert(2, len(args), t)
			assert(1, len(params), t)
			_, ok := params["turbo"]
			assert(true, ok, t)
			assert("flash", args[1], t)
		},
	}

	command.Execute(text)
}

func TestIntermediateCommand_NoArgs(t *testing.T) {
	text := []string{"cli", "upgrade"}

	invoked := false
	var tc Command
	tc = &TerminalCommand{
		func(params map[string]string, args []string) {
			assert(1, len(args), t)
			assert(0, len(params), t)
			invoked = true
		},
	}

	command := IntermediateCommand{
		func(params map[string]string, args []string) *Command {
			assert(2, len(args), t)
			assert(0, len(params), t)
			return &tc
		},
	}

	command.Execute(text)
	assert(true, invoked, t)
}

func TestIntermediateCommand_ParamInIntermedArgInTerminal(t *testing.T) {
	text := []string{"cli", "--global", "upgrade", "repo"}

	invoked := false
	var tc Command
	tc = &TerminalCommand{
		func(params map[string]string, args []string) {
			assert(2, len(args), t)
			assert(0, len(params), t)
			invoked = true
		},
	}

	command := IntermediateCommand{
		func(params map[string]string, args []string) *Command {
			assert(3, len(args), t)
			assert(1, len(params), t)
			_, ok := params["global"]
			assert(true, ok, t)
			return &tc
		},
	}

	command.Execute(text)
	assert(true, invoked, t)
}

func assert(expected, actual interface{}, t *testing.T) {
	if expected != actual {
		t.Error(fmt.Sprintf("expected %v but got %v", expected, actual))
	}
}

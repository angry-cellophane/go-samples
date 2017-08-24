package commands

import (
	"strings"
	"fmt"
	"os"
)

type Command interface {
	Execute(args []string)
}

type IntermediateHandler func(params map[string]string, args []string) *Command
type TerminalHandler func(params map[string]string, args []string)

func processArgs(args []string) (map[string]string, []string) {
	params := make(map[string]string)
	newArgs := args
	// args[0] == commandName
	if len(args) > 1 {
		i := 1
		for ;i < len(args); {
			if !strings.HasPrefix(args[i], "-") {
				break
			}

			last := strings.LastIndex(args[i], "-")
			if last == len(args[i]) - 1 {
				fmt.Println("Empty paramaters (-- and -) are not allowed")
				os.Exit(1)
			}

			parameter := args[i][last + 1:]
			paramName := parameter
			paramValue := ""
			if eqLast := strings.LastIndex(parameter, "="); eqLast != -1 {
				paramName = paramName[:eqLast]
				if eqLast >= len(args[i]) {
					fmt.Printf("Parameter has no value followed =. %v. Please use --<<param_name>>=<<param_value>> or --<<param_name>>", args[i])
					os.Exit(1)
				}
				paramValue = parameter[eqLast+1:]
			}
			params[paramName] = paramValue
			i++
		}
		if i != 1 {
			newArgs = append(args[0:1], args[i:]...)
		}
	}
	return params, newArgs
}

type IntermediateCommand struct {
	Handler IntermediateHandler
}

func (this *IntermediateCommand) Execute(args []string) {
	params, newArgs := processArgs(args)
	next := this.Handler(params, newArgs)
	if len(newArgs) >= 2 {
		newArgs = newArgs[1:]
	}
	(*next).Execute(newArgs)
}

type TerminalCommand struct {
	Handler TerminalHandler
}

func (this *TerminalCommand) Execute(args []string) {
	params, newArgs := processArgs(args)
	this.Handler(params, newArgs)
}

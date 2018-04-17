package goash

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
	shellwords "github.com/mattn/go-shellwords"
)

var prompt string
var readLine *readline.Instance
var parser = shellwords.NewParser()

var completer = readline.NewPrefixCompleter(
	readline.PcItemDynamic(listFiles("/usr/bin")),
	readline.PcItemDynamic(listFiles(currentDir)),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func interactive(ctx context.Context) (err error) {
	parser.ParseEnv = true

	readLine, err = readline.NewEx(&readline.Config{
		Prompt:          getPrompt(),
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		return err
	}
	defer readLine.Close()

	for {
		line, err := readLine.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		// split by pipe
		cmds := make([]*exec.Cmd, 0)
		for _, cmdStr := range strings.Split(line, "|") {
			c, err := parser.Parse(cmdStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				continue
			}
			cmd, err := buildCommand(ctx, c)
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				continue
			}
			if cmd == nil { // this is builtin
				continue
			}
			cmds = append(cmds, cmd)
		}

		cmds = AssemblePipes(cmds, nil, os.Stdout)
		if cmds == nil {
			continue
		}

		if err := ExecutePipes(cmds); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}
	}
	return nil
}

func buildCommand(ctx context.Context, cmd []string) (*exec.Cmd, error) {
	if len(cmd) == 0 {
		return nil, nil
	}
	if isBuiltIn(cmd[0]) {
		builtin := Builtins[cmd[0]]
		return nil, builtin(ctx, cmd)
	}
	bin, err := exec.LookPath(cmd[0])
	if err != nil {
		return nil, err
	}

	return exec.CommandContext(ctx, bin, cmd[1:]...), nil
}

func isBuiltIn(cmd string) bool {
	_, ok := Builtins[cmd]
	return ok
}

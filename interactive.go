package goash

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/chzyer/readline"
	shellwords "github.com/mattn/go-shellwords"
)

var prompt string
var readLine *readline.Instance
var parser = shellwords.NewParser()

var completer = readline.PcItemDynamic(listFiles(currentDir))

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

		args, err := parser.Parse(line)
		if err != nil {
			PrintError(err)
			continue
		}
		if err := execute(ctx, args); err != nil {
			PrintError(err)
			continue
		}
	}
	return nil
}

func execute(ctx context.Context, cmd []string) error {
	if len(cmd) == 0 {
		return nil
	}
	if isBuiltIn(cmd[0]) {
		builtin := Builtins[cmd[0]]
		return builtin(ctx, cmd)
	} else {
		bin, err := exec.LookPath(cmd[0])
		if err != nil {
			return err
		}

		cmd := exec.CommandContext(ctx, bin, cmd[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return nil
}

func isBuiltIn(cmd string) bool {
	_, ok := Builtins[cmd]
	return ok
}

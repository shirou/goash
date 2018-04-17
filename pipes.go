package goash

import (
	"io"
	"os/exec"
)

// https://github.com/ebuchman/go-shell-pipes

func AssemblePipes(cmds []*exec.Cmd, stdin io.Reader, stdout io.Writer) []*exec.Cmd {
	if len(cmds) == 0 {
		return nil
	}
	cmds[0].Stdin = stdin
	cmds[0].Stderr = stdout
	// assemble pipes
	for i, c := range cmds {
		if i < len(cmds)-1 {
			cmds[i+1].Stdin, _ = c.StdoutPipe()
			cmds[i+1].Stderr = stdout
		} else {
			c.Stdout = stdout
			c.Stderr = stdout
		}
	}
	return cmds
}

func ExecutePipes(cmds []*exec.Cmd) error {
	if len(cmds) == 0 {
		return nil
	}

	// start processes in descending order
	for i := len(cmds) - 1; i > 0; i-- {
		if err := cmds[i].Start(); err != nil {
			return err
		}
	}
	// run the first process
	if err := cmds[0].Run(); err != nil {
		return err
	}
	// wait on processes in ascending order
	for i := 1; i < len(cmds); i++ {
		if err := cmds[i].Wait(); err != nil {
			return err
		}
	}

	return nil
}

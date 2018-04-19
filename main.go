package goash

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

func init() {
	if pwd, err := os.Getwd(); err == nil {
		currentDir = pwd
	}
}

func Main(stdout io.Writer, call []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(call) > 2 {
		call = call[0:1]
	}
	if len(call) == 2 {
		f, e := os.Open(call[1])
		if e != nil {
			return e
		}
		defer f.Close()
		in := bufio.NewReader(f)
		fmt.Println(ReadWholeLine(in))
	} else {
		return interactive(ctx)
	}

	return nil
}

func isComment(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, "#")
}

package goash

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type BuiltinHandler func(context.Context, []string) error

var (
	Builtins   map[string]BuiltinHandler
	currentDir = ""
)

func init() {
	Builtins = map[string]BuiltinHandler{
		"cd":       cd,
		"exit":     exit,
		"env":      env,
		"getenv":   getenv,
		"setenv":   setenv,
		"unsetenv": unsetenv,
		//		"fork":     fork,
		"refresh": refresh,
	}
}

func homedir() string {
	homedir := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if homedir == "" {
		homedir = os.Getenv("USERPROFILE")
	}
	if homedir == "" {
		homedir = os.Getenv("HOME")
	}
	return homedir
}

func refresh(ctx context.Context, call []string) error {
	readLine.Clean()
	return nil
}

func cd(ctx context.Context, call []string) error {
	var dst string
	if len(call) == 1 {
		dst = homedir()
	} else if call[1] == "~" && len(call) == 2 {
		dst = homedir()
	} else {
		dst = call[1]
	}
	if err := os.Chdir(dst); err != nil {
		return err
	}
	currentDir = dst
	if readLine != nil {
		readLine.SetPrompt(getPrompt())
	}

	return nil
}

func exit(ctx context.Context, call []string) (err error) {
	code := 0
	if len(call) >= 2 {
		code, err = strconv.Atoi(call[1])
		if err != nil {
			return err
		}
	}
	os.Exit(code)
	return nil
}

func env(ctx context.Context, call []string) error {
	for _, envvar := range os.Environ() {
		fmt.Println(envvar)
	}
	return nil
}

func getenv(ctx context.Context, call []string) error {
	if len(call) != 2 {
		return errors.New("`getenv <variable name>`")
	}
	fmt.Println(os.Getenv(call[1]))
	return nil
}

func setenv(ctx context.Context, call []string) error {
	if len(call) != 3 {
		return errors.New("`setenv <variable name> <value>`")
	}
	return os.Setenv(call[1], call[2])
}

func unsetenv(ctx context.Context, call []string) error {
	if len(call) != 2 {
		return errors.New("`unsetenv <variable name>`")
	}
	return os.Setenv(call[1], "")
}

func fork(ctx context.Context, call []string) error {
	if len(call) < 2 {
		return errors.New("`fork <command...>`")
	}
	go func(ctx context.Context, call []string) {
		cmd, err := buildCommand(ctx, call)
		if err != nil || cmd == nil {
			return
		}
		cmd.Start()
	}(ctx, call)
	return nil
}

package goash

import (
	"os"
	"os/user"
)

var shellVariables = make(map[string]string)

func getPrompt() string {
	if p := os.Getenv("PROMPT"); p != "" {
		return p
	}
	u, err := user.Current()
	if err != nil {
		// Current not implemented on linux/amd64
		return currentDir + " $ "
	}
	if u.Uid == "0" {
		return currentDir + "# "
	}
	return currentDir + " $ "
}

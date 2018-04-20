package goash

import (
	"io/ioutil"
	"strings"
)

func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0, 10)
		path := "/usr/bin"
		prefix := line
		if len(line) > 0 {
			path = currentDir
		}
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			if len(prefix) == 0 || strings.HasPrefix(f.Name(), prefix) {
				names = append(names, f.Name())
			}
		}
		return names
	}
}

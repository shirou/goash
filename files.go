package goash

import (
	"io/ioutil"
)

func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0, 10)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func readLines(fname string) []string {
	f, err := os.Open(fname)
	panicerr(err)
	b, err := ioutil.ReadAll(f)
	panicerr(err)
	links := strings.Split(string(b), "\n")
	return links
}

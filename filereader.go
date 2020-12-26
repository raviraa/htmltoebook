package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func readLines(fname string) []string {
	f, err := os.Open(fname)
	chkerr(err)
	b, err := ioutil.ReadAll(f)
	chkerr(err)
	links := strings.Split(string(b), "\n")
	return links
}

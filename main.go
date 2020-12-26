package main

import (
	"fmt"
	"os"
)

func main() {

	fmt.Println(os.Args)
	if len(os.Args) != 2 || os.Args[1] == "-h" {
		panic(usage)
	}

	// links := readLines(os.Args[1])
	// fetchLinks(links)

	writeMobi()

}

func chkerr(err error) {
	if err != nil {
		panic(err)
	}
}

var usage = `
Usage: htmltoepub  <file-with-links.(txt|html)>






`

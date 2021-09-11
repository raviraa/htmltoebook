package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	toml "github.com/pelletier/go-toml"
)

var separator = strings.Repeat("#", 60)

// marshals value to toml. writes toml, separator and tailmsg to file
func writeEditorFile(v interface{}, tailmsg string) error {
	tomlb, err := toml.Marshal(v)
	if err != nil {
		return err
	}
	f, err := os.Create(tmpfname)
	if err != nil {
		return err
	}
	f.Write(tomlb)
	fmt.Fprintf(f, "\n%s\n", separator)
	fmt.Fprintf(f, "%s\n", tailmsg)
	return f.Close()
}

// exec $EDITOR with tmpfile, split result, parse first segment as toml
func runEditor(v interface{}) (tailmsg string, err error) {
	err = runLinksEditorCommand()
	if err != nil {
		log.Fatalf("unable to open editor '%s' to edit links. please set correct EDITOR environment variable, %v", getEditor(), err)
	}
	f, err := os.Open(tmpfname)
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	spl := bytes.Split(b, []byte(separator))
	if len(spl) != 2 {
		return "", fmt.Errorf("could not find 2 segments in edited file. found %d segments", len(spl))
	}
	err = toml.Unmarshal(spl[0], v)
	return string(spl[1]), err
}

func runLinksEditorCommand() error {
	err := RunCmd(getEditor(), tmpfname)
	return err
}

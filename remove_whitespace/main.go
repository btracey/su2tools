// +build ignore

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func mainargs(args []string) error {
	// Check input arguments
	if len(os.Args) < 2 {
		return errors.New("No input arguments supplied. There must be one or two input arguments. The first argument is the file to read, and the second (optional) argument is the location to write the file")
	}
	if len(os.Args) > 3 {
		return errors.New("Too many arguments supplied. There must be one or two input arguments. The first argument is the file to read, and the second (optional) argument is the location to write the file")
	}

	// Open replacing file
	readFile := os.Args[1]
	b, err := ioutil.ReadFile(readFile)
	if err != nil {
		return errors.New("error opening file " + readFile + ": " + err.Error())
	}

	// Open write file
	var writeFile string
	if len(os.Args) == 2 {
		writeFile = readFile
	} else {
		writeFile = os.Args[2]
	}
	writer, err := os.Create(writeFile)
	if err != nil {
		return errors.New("error creating write file " + writeFile + " : " + err.Error())
	}
	return RemoveTrailingWhitespace(b, writer)
}

// Removes the trailing whitespace from the lines of a file. It takes in one or two arguments.
// The first argument is the file from white whitespace will be removed, and the
// second is where the resulting file will be written. If only one argument is specified,
// the file will be overwritten
func main() {
	err := mainargs(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

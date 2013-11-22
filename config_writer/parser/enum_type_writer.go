package main

import (
	"os"
)

type enumWriter struct {
	Filename string
	File     *os.File
}

func (e *enumWriter) init() (err error) {
	e.File, err = os.Create(e.Filename)
	if err != nil {
		return err
	}
	_, err = e.File.Write([]byte("package " + pkgName + "\n"))
	if err != nil {
		return err
	}
	e.File.Write(dynamicGeneration)
	e.File.Write([]byte("\n"))
	e.File.Write([]byte("// List of all the enumerable types and their options\n"))
	e.File.Write([]byte("\n"))
	return nil
}

func (e *enumWriter) finalize() {
	e.File.Close()
}

func (e *enumWriter) Name() string {
	return "enumWriter"
}

func (e *enumWriter) GetFilename() string {
	return e.Filename
}

func (e *enumWriter) add_option(option *pythonOption) error {
	if !is_enum_option(option.typeO) {
		return nil
	}
	// Print description
	e.File.Write([]byte("// "))
	e.File.Write([]byte(option.go_type))
	e.File.Write([]byte(" is the "))
	e.File.Write([]byte(option.description))
	e.File.Write([]byte("\n"))
	e.File.Write([]byte("type " + option.go_type + " int\n"))

	// Print constants with all the values
	return nil
}

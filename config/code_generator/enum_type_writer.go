package main

import (
//"errors"
//"os"
)

// NOTE: This is a good start, but it won't work nicely due to the way things are done.
// For example, all the flux routine codes take in the same set of strings. These would
// have to be hand-coded somehow because SU2 doesn't make that distinction

/*
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
	e.File.Write([]byte("const(\n"))
	if len(option.enumOptions) == 0 {
		return errors.New("no enum options")
	}
	// write the first type
	e.File.Write([]byte(option.enumGo[0] + " " + option.go_type + " = iota\n"))
	for i := 1; i < len(option.enumGo); i++ {
		e.File.Write([]byte(option.enumGo[i] + "\n"))
	}
	e.File.Write([]byte(")\n\n"))
	return nil
}
*/

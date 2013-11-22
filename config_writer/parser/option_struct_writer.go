package main

import (
	"os"
)

type optionFile struct {
	Filename string
	File     *os.File
}

func (o *optionFile) init() (err error) {
	o.File, err = os.Create(o.Filename)
	if err != nil {
		return err
	}
	_, err = o.File.Write([]byte("package " + pkgName + "\n"))
	if err != nil {
		return err
	}
	o.File.Write(dynamicGeneration)
	o.File.Write([]byte("\n"))
	o.File.Write([]byte("// Options is a struct containing all of the possible options in SU^2\n"))
	o.File.Write([]byte("type Options struct { \n"))
	return nil
}

func (o *optionFile) finalize() {
	o.File.Write([]byte("}"))
	o.File.Close()
}

func (o *optionFile) Name() string {
	return "optionFile"
}

func (o *optionFile) GetFilename() string {
	return o.Filename
}

func (o *optionFile) add_option(option *pythonOption) error {
	// write the field name
	o.File.Write([]byte(option.structName))
	// write a space
	o.File.Write([]byte(" "))
	// write type
	o.File.Write([]byte(option.go_type))

	// write the description
	o.File.Write([]byte(" // "))
	o.File.Write([]byte(option.description))
	o.File.Write([]byte("\n"))
	return nil
}

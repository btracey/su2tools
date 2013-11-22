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
	_, err = o.File.WriteString("package " + pkgName + "\n")
	if err != nil {
		return err
	}
	o.File.Write(dynamicGeneration)
	o.File.WriteString("\n")
	o.File.WriteString("// Options is a struct containing all of the possible options in SU^2\n")
	o.File.WriteString("type Options struct { \n")
	return nil
}

func (o *optionFile) finalize() {
	o.File.WriteString("}")
	o.File.Close()
}

func (o *optionFile) Name() string {
	return "optionFile"
}

func (o *optionFile) GetFilename() string {
	return o.Filename
}

func (o *optionFile) add_option(option *pythonOption) error {
	o.File.WriteString(" // " + option.description + "\n")
	o.File.WriteString(option.structName + " " + option.go_type + "\n")
	return nil
}

package main

import (
	"os"

	//"fmt"
	"github.com/btracey/su2tools/config/common"
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
	o.File.WriteString("import \"github.com/btracey/su2tools/config/common\"\n")
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
	//fmt.Println("go base type: ", option.goBaseType)
	//fmt.Println("kind string: " + common.GoTypeToKind[option.goBaseType])
	o.File.WriteString(string(option.optionsField) + " " + common.GoTypeToKind[option.goBaseType] + "\n")
	return nil
}

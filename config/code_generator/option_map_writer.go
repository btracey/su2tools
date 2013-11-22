package main

import (
	//"errors"
	"os"
	//"strconv"
)

type optionMapWriter struct {
	Filename string
	File     *os.File
}

func (d *optionMapWriter) init() (err error) {
	d.File, err = os.Create(d.Filename)
	if err != nil {
		return err
	}
	err = initializeFile(d.File)
	if err != nil {
		return err
	}
	d.File.WriteString("\n")
	d.File.WriteString("//optionMap is a map that stores the option printer for each object\n")
	d.File.WriteString("var optionMap map[string]*optionPrint = map[string]*optionPrint{\n")
	return nil
}

func (d *optionMapWriter) finalize() {
	d.File.WriteString("}\n")
	d.File.Close()
}

func (d *optionMapWriter) Name() string {
	return "optionMapWriter"
}

func (d *optionMapWriter) GetFilename() string {
	return d.Filename
}

func (d *optionMapWriter) add_option(option *pythonOption) error {
	d.File.WriteString("\"" + option.structName + "\"" + ": {\n")
	d.File.WriteString("Description: \"" + option.description + "\",\n")
	d.File.WriteString("Category: \"" + option.category + "\",\n")
	d.File.WriteString("SU2OptionName: \"" + option.name + "\",\n")
	d.File.WriteString("OptionTypeName: \"" + option.typeO + "\",\n")
	d.File.WriteString("Default: \"" + option.defaultO + "\",\n")
	d.File.WriteString("Type: \"" + option.go_type + "\",\n")
	d.File.WriteString("StructName: \"" + option.structName + "\",\n")
	d.File.WriteString("ValueString: \"" + option.values + "\",")
	d.File.WriteString("enumOptions: []string{\n")
	for _, val := range option.enumOptions {
		d.File.WriteString("\"" + val + "\"" + ",\n")
	}
	d.File.WriteString("},\n")
	d.File.WriteString("},\n")
	return nil
}

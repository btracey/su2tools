package main

import (
	"os"
)

type fieldMapWriter struct {
	Filename string
	File     *os.File
}

func (d *fieldMapWriter) init() (err error) {
	d.File, err = os.Create(d.Filename)
	if err != nil {
		return err
	}
	err = initializeFile(d.File)
	if err != nil {
		return err
	}
	d.File.WriteString("\n")
	d.File.WriteString("//su2ToGoFieldMap is a map that translates the SU2 field to a go field\n")
	d.File.WriteString("var su2ToGoFieldMap map[string]string = map[string]string{\n")
	return nil
}

func (d *fieldMapWriter) finalize() {
	d.File.WriteString("}\n")
	d.File.Close()
}

func (d *fieldMapWriter) Name() string {
	return "fieldMapWriter"
}

func (d *fieldMapWriter) GetFilename() string {
	return d.Filename
}

func (d *fieldMapWriter) add_option(option *pythonOption) error {
	d.File.WriteString("\"" + option.structName + "\"" + ": \"" + option.name + "\",\n")
	return nil
}

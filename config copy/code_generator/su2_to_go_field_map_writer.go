package main

import (
	"os"
)

type su2ToGoFieldMapWriter struct {
	Filename string
	File     *os.File
}

func (d *su2ToGoFieldMapWriter) init() (err error) {
	d.File, err = os.Create(d.Filename)
	if err != nil {
		return err
	}
	err = initializeFile(d.File)
	if err != nil {
		return err
	}
	d.File.WriteString("\n")
	d.File.WriteString("import \"github.com/btracey/su2tools/config/common\"\n")
	d.File.WriteString("//su2ToGoFieldMap is a map that translates the go field to an SU2 field\n")
	d.File.WriteString("var su2ToGoFieldMap map[common.ConfigfileOption]common.OptionsField = map[common.ConfigfileOption]common.OptionsField{\n")
	return nil
}

func (d *su2ToGoFieldMapWriter) finalize() {
	d.File.WriteString("}\n")
	d.File.Close()
}

func (d *su2ToGoFieldMapWriter) Name() string {
	return "su2ToGoFieldMapWriter"
}

func (d *su2ToGoFieldMapWriter) GetFilename() string {
	return d.Filename
}

func (d *su2ToGoFieldMapWriter) add_option(option *pythonOption) error {
	d.File.WriteString("\"" + string(option.configName) + "\"" + ": \"" + string(option.optionsField) + "\",\n")
	return nil
}

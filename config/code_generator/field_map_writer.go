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
	d.File.WriteString("import \"github.com/btracey/su2tools/config/common\"\n")
	d.File.WriteString("//goToSU2FieldMap is a map that translates the SU2 field to a go field\n")
	d.File.WriteString("var goToSU2FieldMap map[common.OptionsField]common.ConfigfileOption = map[common.OptionsField]common.ConfigfileOption{\n")
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
	d.File.WriteString("\"" + string(option.optionsField) + "\"" + ": \"" + string(option.configName) + "\",\n")
	return nil
}

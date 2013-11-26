package main

import (
	//"errors"
	"os"
	//"strconv"

	"github.com/btracey/su2tools/config/common"
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
	d.File.WriteString("import \"github.com/btracey/su2tools/config/common\"\n")
	d.File.WriteString("//optionMap is a map that stores the option printer for each object\n")
	d.File.WriteString("var optionMap map[common.OptionsField]*optionPrint = map[common.OptionsField]*optionPrint{\n")
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
	d.File.WriteString("\"" + string(option.optionsField) + "\"" + ": {\n")
	d.File.WriteString("Description: \"" + option.description + "\",\n")
	d.File.WriteString("ConfigCategory: \"" + string(option.configCategory) + "\",\n")
	d.File.WriteString("ConfigName: \"" + string(option.configName) + "\",\n")
	d.File.WriteString("ConfigType: \"" + string(option.configType) + "\",\n")
	d.File.WriteString("Default: \"" + option.defaultString + "\",\n")
	d.File.WriteString("GoBaseType: " + common.GoTypeToString[option.goBaseType] + ",\n")
	d.File.WriteString("OptionsField: \"" + string(option.optionsField) + "\",\n")
	//d.File.WriteString("ValueString: \"" + option.enumOptionsString + "\",\n")
	d.File.WriteString("enumOptions: []common.Enum{\n")
	for _, val := range option.enumOptions {
		d.File.WriteString("\"" + string(val) + "\"" + ",\n")
	}
	d.File.WriteString("},\n")
	d.File.WriteString("Value: " + common.ValueAsString(option.defaultValue, option.goBaseType))
	d.File.WriteString("},\n")
	return nil
}

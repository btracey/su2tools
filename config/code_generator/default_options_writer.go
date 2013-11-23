package main

import (
	//"errors"
	"os"
	//"strconv"

	"github.com/btracey/su2tools/config/common"
)

type defaultWriter struct {
	Filename string
	File     *os.File
}

func (d *defaultWriter) init() (err error) {
	d.File, err = os.Create(d.Filename)
	if err != nil {
		return err
	}
	_, err = d.File.WriteString("package " + pkgName + "\n")
	if err != nil {
		return err
	}
	d.File.Write(dynamicGeneration)
	d.File.WriteString("\n")
	d.File.WriteString("//NewOptions creates a new options structure with the default values\n")
	d.File.WriteString("func NewOptions() *Options{\n")
	d.File.WriteString("return &Options{\n")
	return nil
}

func (d *defaultWriter) finalize() {
	d.File.WriteString("}\n")
	d.File.WriteString("}\n")
	d.File.Close()
}

func (d *defaultWriter) Name() string {
	return "defaultWriter"
}

func (d *defaultWriter) GetFilename() string {
	return d.Filename
}

func (d *defaultWriter) add_option(option *pythonOption) error {
	d.File.WriteString(string(option.optionsField) + ": ")

	str := common.ValueAsString(option.defaultValue, option.goBaseType)
	d.File.WriteString(str)
	/*
		switch option.goBaseType {
		case common.StringType:
			d.File.WriteString("\"" + option.defaultString + "\"")
		case common.EnumType:
			d.File.WriteString("\"" + string(option.enumOptions[0]) + "\"")
		case common.BoolType:
			b := (option.defaultValue).(bool)
			if b {
				d.File.WriteString("false")
			} else {
				d.File.WriteString("true")
			}
		case common.Float64Type:
			f := (option.defaultValue).(float64)
			str := strconv.FormatFloat(f, 'g', -1, 64)
			d.File.WriteString(str)
		case common.Float64ArrayType:
			d.File.WriteString("[]float64{")

			slice := (option.defaultValue).([]float64)
			for i, f := range slice {
				str := strconv.FormatFloat(f, 'g', -1, 64)
				d.File.WriteString(str)
				if i != len(slice)-1 {
					d.File.WriteString(",")
				}
			}
			d.File.WriteString("}")
		case common.StringArrayType:
			d.File.WriteString("[]string{")
			strs := (option.defaultValue).([]string)
			for i, s := range strs {
				d.File.WriteString("\"" + s + "\"")
				if i != len(strs)-1 {
					d.File.WriteString(",")
				}
			}
			d.File.WriteString("}")
		case common.BadType:
			panic("bad type found")
		default:
			panic("unknown type ")
		}
	*/
	d.File.WriteString(",\n")
	return nil
}

package main

import (
	"errors"
	"os"
	"strconv"
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
	d.File.WriteString(option.structName + ": ")
	switch option.go_type {
	case "string":
		d.File.WriteString("\"" + option.defaultO + "\"")
	case "bool":
		switch option.defaultO {
		case "NO":
			d.File.WriteString("false")
		case "YES":
			d.File.WriteString("true")
		default:
			return errors.New("Bad boolean option. Default is " + option.defaultO)
		}
	case "float64":
		str := strconv.FormatFloat(option.defaultFloat, 'g', -1, 64)
		d.File.WriteString(str)
	case "[]float64":
		d.File.WriteString("[]float64{")
		for i, f := range option.defaultFloatArray {
			str := strconv.FormatFloat(f, 'g', -1, 64)
			d.File.WriteString(str)
			if i != len(option.defaultFloatArray)-1 {
				d.File.WriteString(",")
			}
		}
		d.File.WriteString("}")
	case "[]string":
		d.File.WriteString("[]string{")
		for i, s := range option.defaultStringArray {
			d.File.WriteString("\"" + s + "\"")
			if i != len(option.defaultStringArray)-1 {
				d.File.WriteString(",")
			}
		}
		d.File.WriteString("}")
	default:
		panic("unknown type " + option.go_type)
	}
	d.File.WriteString(",\n")
	return nil
}

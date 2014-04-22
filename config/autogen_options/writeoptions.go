package main

import (
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func writeoptions(categories []*ConfigCategory, options []*ConfigOption) {
	optionConstFilename := filepath.Join(configPath, "option_consts.go")

	optionConst, err := os.Create(optionConstFilename)
	if err != nil {
		log.Fatalf("error creating %s: %v", optionConstFilename, err)
	}

	optionConst.WriteString(packageHeader)
	optionConst.WriteString(autogenmessage)

	optionConst.WriteString(
		`type option struct{
	Name string
	Config string
	Category int
	Description string
	Type string 
	ExtraType string
	Default string
	OptionConst Option
	}

	// Option is a specific option of SU2. See the Options struct comments for
// descriptions of the options 
type Option string

//func (o Option) String() string{
//	return optionMap[o].Name
//}

const(
	All Option = "AllOptions"		// All is a shortcut for printing all of the config options
	`)

	optionMap := `var optionMap = map[Option]option {
`

	otherMap := `var stringToOption = map[string]Option{
`

	for _, option := range options {
		optionConst.WriteString(option.Value + " = \"" + option.Value + "\"\n")

		optionMap += option.Value + ": {\n" +
			"Name: \"" + option.Value + "\",\n" +
			"Config: \"" + option.ConfigString + "\",\n" +
			"Category: " + strconv.Itoa(option.Category) + ",\n" +
			"Description: \"" + option.Description + "\",\n" +
			"Type: \"" + option.Type + "\",\n" +
			"ExtraType: \"" + option.ExtraType + "\",\n" +
			//"Default: \"" + option.Default + "\",\n" +
			"OptionConst:" + option.Value + ",\n" +
			"},\n"

		otherMap += "\"" + option.Value + "\": " + option.Value + ",\n"
	}
	optionConst.WriteString(")\n\n")
	optionMap += "}\n\n"
	otherMap += "}\n\n"
	optionConst.WriteString(optionMap)
	optionConst.WriteString(otherMap)

	optionConst.Close()

	b, err := ioutil.ReadFile(optionConstFilename)
	if err != nil {
		log.Fatalf("error creating %s: %v", optionConstFilename, err)
	}
	b, err = format.Source(b)
	if err != nil {
		log.Fatalf("error formatting %s: %v", optionConstFilename, err)
	}
	err = ioutil.WriteFile(optionConstFilename, b, 0700)
	if err != nil {
		log.Fatalf("error writing %s: %v", optionConstFilename, err)
	}
}

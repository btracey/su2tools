package config

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"reflect"
	"strings"

	"github.com/btracey/su2tools/config/common"

	"fmt"
)

/*
// IsConfigOption returns true if configString is a valid option string
// in the config file
func (o *Options) IsConfigOption(configString common.ConfigOptionType) bool {
	goField := su2ToGoFieldMap[configString]
	_, ok := optionMap[goField]
	return ok
}
*/

// IsEnum returns true if the options field is an enumable
func (o *Options) IsEnum(field common.OptionsField) bool {
	option := optionMap[field]
	return common.IsEnumOption(option.ConfigType)
}

// SetEnum  sets the enumerable field, checking that it is a valid option.
// Returns an error if it is not a valid option or if
func (o *Options) SetEnum(field common.OptionsField, val common.Enum) error {
	// First, check that it's a real field
	option, ok := optionMap[field]
	if !ok {
		return errors.New("setenum: " + string(field) + " is not a valid field")
	}
	// check that the value is any of the possible options
	for _, str := range option.enumOptions {
		if str == val {
			reflect.ValueOf(o).Elem().FieldByName(string(field)).Set(reflect.ValueOf(val))
			return nil
		}
	}
	return errors.New("setenum: " + string(val) + " is not a valid option ")
}

// SetField sets the field with the value. It calls SetEnum if it is an enumerable option
func (o *Options) SetField(field common.OptionsField, val interface{}) error {
	option, ok := optionMap[field]
	if !ok {
		return errors.New("setfields: " + string(field) + " is not a field")
	}
	if common.IsEnumOption(option.ConfigType) {
		err := o.SetEnum(field, val.(common.Enum))
		if err != nil {
			errors.New("setfelids: error setting field " + string(field) + ": " + err.Error())
		}
		return nil
	}
	reflect.ValueOf(o).Elem().FieldByName(string(field)).Set(reflect.ValueOf(val))
	return nil
}

type FieldMap map[common.OptionsField]interface{}

// SetFields sets the fields of the options structure
func (o *Options) SetFields(fieldMap FieldMap) error {
	for field, value := range fieldMap {
		o.SetField(field, value)
	}
	return nil
}

// WriteConfig writes a config file to the writer with the options given in list
func (o *Options) WriteConfig(writer io.Writer, list OptionList) error {
	buf := &bytes.Buffer{}
	buf.Write(configHeader)
	printAll := list["All"]
	currentHeading := -1
	for _, option := range optionOrder {
		catNumber := categoryOrder[option.ConfigCategory]
		if catNumber > currentHeading {
			buf.WriteString("\n\n")
			buf.WriteString("%" + categoryBookend + string(option.ConfigCategory) + categoryBookend + "% \n")
			currentHeading = catNumber
		}
		printOption, ok := list[option.OptionsField]
		if printAll {
			printOption = true
			ok = true
		}
		if !ok {
			continue
		}

		if printOption {
			r := reflect.ValueOf(o).Elem()
			i := r.FieldByName(string(option.OptionsField)).Interface()
			//option.Value = i

			commentOut := common.ShouldCommentOut(i, option.defaultValue, option.GoBaseType, option.ConfigType)
			b, _ := option.MarshalSU2Config(commentOut)
			b2 := common.ValueAsConfigString(i, option.GoBaseType)
			buf.Write(b)
			buf.WriteString(b2)
			buf.WriteString("\n")
		}
	}
	writer.Write(buf.Bytes())
	return nil
}

func Read(reader io.Reader) (*Options, OptionList, error) {
	//setter := make(map[string]interface{})
	// turn it into a scanner
	optionList := make(OptionList)
	scanner := bufio.NewScanner(reader)

	options := NewOptions()
	for scanner.Scan() {
		if shouldcontinue(scanner) {
			continue
		}
		goFieldName, value, err := getoption(scanner)
		if err != nil {
			return nil, nil, err
		}
		_, ok := optionList[goFieldName]
		if ok {
			return nil, nil, fmt.Errorf("field %s set multiple times", goFieldName)
		}
		optionList[goFieldName] = true
		options.SetField(goFieldName, value)
	}
	err := scanner.Err()
	if err != nil {
		return nil, nil, errors.New("readconfig: " + err.Error())
	}
	return options, optionList, nil
}

func shouldcontinue(scanner *bufio.Scanner) bool {
	if len(scanner.Bytes()) == 0 {
		return true
	}
	if scanner.Bytes()[0] == '%' {
		return true
	}
	return false
}

func getoption(scanner *bufio.Scanner) (common.OptionsField, interface{}, error) {
	line := string(scanner.Bytes())
	// split the line at the equals sign
	parts := strings.Split(line, "=")
	if len(parts) > 2 {
		return "", nil, errors.New("readconfig: option line has two equals signs")
	}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	goFieldName, ok := su2ToGoFieldMap[common.ConfigfileOption(parts[0])]
	if !ok {
		return "", nil, errors.New("readconfig: option " + parts[0] + " does not exist")
	}
	option := optionMap[goFieldName]
	value, err := common.InterfaceFromString(parts[1], option.GoBaseType)
	if err != nil {
		return "", nil, errors.New("readconfig: error parsing " + parts[0] + ": " + err.Error())
	}
	return goFieldName, value, nil
}

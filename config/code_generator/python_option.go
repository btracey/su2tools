package main

import (
	"errors"
	//"strconv"
	"strings"

	"github.com/btracey/su2tools/config/common"

	//"fmt"
)

// typeS and defaultS because type and default are keywords
type pythonOption struct {
	configName        common.ConfigfileOption
	configType        common.ConfigOptionType
	configCategory    common.ConfigCategory
	enumOptionsString string
	defaultString     string // string representing the default value
	description       string
	optionsField      common.OptionsField
	goBaseType        common.GoBaseType // string representing the base type
	enumOptions       []common.Enum     // list of options for the enumerable

	defaultValue interface{}
	/*
		defaultFloat       float64
		defaultStringArray []string
		defaultFloatArray  []float64
		defaultBool        bool
	*/
}

// ParseEnumOptions parses the enum type options, and puts the default option
// as the first entry in the list
func (o *pythonOption) ParseEnumOptions() error {
	if !common.IsEnumOption(o.configType) {
		return nil
	}
	values := strings.TrimSpace(o.enumOptionsString)
	// check beginning and end are '[' and ']' signifying array list
	if values[0] != '[' {
		return errors.New("parse enum: leading descriptor not [")
	}
	if values[len(values)-1] != ']' {
		return errors.New("parse enum: leading descriptor not ]")
	}
	values = values[1 : len(values)-1]
	// Parse the values at the commas

	list := strings.Split(values, ",")
	for i, s := range list {
		s = strings.TrimSpace(s)
		if s[0] != byte("'"[0]) {
			return errors.New("first character isn't '")
		}
		if s[len(s)-1] != byte("'"[0]) {
			return errors.New("last character isn't '")
		}
		list[i] = s[1 : len(s)-1]

	}
	// reorganize the list so that the default option is first
	orderedList := make([]string, len(list))
	defaultInd := -1
	for i, s := range list {
		if s == o.defaultString {
			if defaultInd != -1 {
				return errors.New("More than one case of enum option")
			}
			defaultInd = i
		}
	}
	for i, s := range list {
		if i == defaultInd {
			orderedList[0] = s
			continue
		}
		if i < defaultInd {
			orderedList[i+1] = s
			continue
		}
		orderedList[i] = s
	}

	o.enumOptions = make([]common.Enum, len(orderedList))
	for i := range o.enumOptions {
		o.enumOptions[i] = common.Enum(orderedList[i])
	}

	//o.enumOptions = orderedList
	/*
		o.enumGo = make([]string, len(o.enumOptions))
		for i, s := range o.enumOptions {
			o.enumGo[i] = to_camel_case(s)
		}
	*/
	return nil
}

// OptionTypeToGoType
func (o *pythonOption) SetOptionTypeAndDefault() error {
	gotype, err := common.ConfigTypeToGoType(o.configType, o.defaultString)
	if err != nil {
		return err
	}
	o.goBaseType = gotype
	o.defaultValue, err = common.OptionStringToInterface(gotype, o.defaultString)
	return err
}

func (o *pythonOption) ConfigNameToStructName() error {
	s := to_camel_case(string(o.configName))
	o.optionsField = common.OptionsField(s)
	return nil //to campel case might get error checking evenually
}

func (o *pythonOption) Process() error {
	err := o.ConfigNameToStructName()
	if err != nil {
		return err
	}
	err = o.SetOptionTypeAndDefault()
	if err != nil {
		return err
	}
	err = o.ParseEnumOptions()
	if err != nil {
		return err
	}
	//fmt.Printf("option is: %#v\n", o)
	return nil
}

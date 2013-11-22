package main

import (
	"errors"
	"strconv"
	"strings"

	//"fmt"
)

// typeS and defaultS because type and default are keywords
type pythonOption struct {
	name, typeO, category, values, defaultO, description string
	structName                                           string
	go_type                                              string
	enumOptions                                          []string
	//enumGo                                               []string
	defaultFloat       float64
	defaultStringArray []string
	defaultFloatArray  []float64
}

// ParseEnumOptions parses the enum type options, and puts the default option
// as the first entry in the list
func (o *pythonOption) ParseEnumOptions() error {
	if !is_enum_option(o.typeO) {
		return nil
	}
	values := strings.TrimSpace(o.values)
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
		if s == o.defaultO {
			if defaultInd != -1 {
				return errors.New("More than one case of default option")
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

	o.enumOptions = orderedList
	/*
		o.enumGo = make([]string, len(o.enumOptions))
		for i, s := range o.enumOptions {
			o.enumGo[i] = to_camel_case(s)
		}
	*/
	return nil
}

func (o *pythonOption) OptionTypeToGoType() error {
	// write the type

	if is_enum_option(o.typeO) {
		if o.typeO == "EnumListOption" {
			//o.go_type = "[]" + o.structName + "Enum"
			o.go_type = "[]string"
			return nil
		}
		//o.go_type = o.structName + "Enum"
		o.go_type = "string"
		return nil
	}

	switch o.typeO {
	case "ArrayOption":
		strs, err := o.SplitArrayOption(o.defaultO)
		if err != nil {
			return err
		}
		// check if it's a float array
		fs := o.IsFloatArray(strs)
		if fs != nil {
			o.defaultFloatArray = fs
			o.go_type = "[]float64"
		} else {
			o.defaultStringArray = strs
			o.go_type = "[]string"
		}
	case "ListOption":
		o.go_type = "string"
	case "DVParamOption":
		o.go_type = "string"
	case "MarkerOption":
		o.go_type = "string"
	case "MarkerDirichlet":
		o.go_type = "string"
	case "MarkerPeriodic":
		o.go_type = "string"
	case "MarkerInlet":
		o.go_type = "string"
	case "MarkerOutlet":
		o.go_type = "string"
	case "MarkerDisplacement":
		o.go_type = "string"
	case "MarkerLoad":
		o.go_type = "string"
	case "MarkerFlowLoad":
		o.go_type = "string"
	case "ScalarOption":
		// See if it looks like a float 64
		f, err := strconv.ParseFloat(o.defaultO, 64)
		if err != nil {
			o.go_type = "string"
		} else {
			o.go_type = "float64"
			o.defaultFloat = f
		}
	case "SpecialOption":
		o.go_type = "bool"
	default:
		return errors.New("option type " + o.typeO + " not implemented")
	}
	return nil
}

func (o *pythonOption) SplitArrayOption(str string) ([]string, error) {
	// Check that the beginning and ending are ( and ]
	if str[0] != '(' {
		return nil, errors.New("no ( at start of string")
	}
	if str[len(str)-1] != ')' {
		return nil, errors.New("no ) at end of string")
	}
	str = str[1 : len(str)-1]
	strs := strings.Split(str, ",")
	for i := range strs {
		strs[i] = strings.TrimSpace(strs[i])
	}
	return strs, nil
}

func (o *pythonOption) IsFloatArray(strs []string) []float64 {
	fs := make([]float64, len(strs))
	for i, str := range strs {
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil
		}
		fs[i] = f
	}
	return fs
}

func (o *pythonOption) ConfigNameToStructName() error {
	o.structName = to_camel_case(o.name)
	return nil //to campel case might get error
}

func (o *pythonOption) Process() error {
	err := o.ConfigNameToStructName()
	if err != nil {
		return err
	}
	err = o.OptionTypeToGoType()
	if err != nil {
		return err
	}
	err = o.ParseEnumOptions()
	if err != nil {
		return err
	}
	return nil
}

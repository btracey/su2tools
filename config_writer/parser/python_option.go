package main

import (
	"errors"
	"strings"

	//"fmt"
)

// typeS and defaultS because type and default are keywords
type pythonOption struct {
	name, typeO, category, values, defaultO, description string
	structName                                           string
	go_type                                              string
	enumOptions                                          []string
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
	// TODO: Add in parsing
	return nil
}

func (o *pythonOption) OptionTypeToGoType() error {
	// write the type

	if is_enum_option(o.typeO) {
		if o.typeO == "EnumListOption" {
			o.go_type = "[]" + o.structName + "Enum"
		}
		o.go_type = o.structName + "Enum"
		return nil
	}

	switch o.typeO {
	case "ArrayOption":
		o.go_type = "[]float64"
	case "DVParamOption":
		o.go_type = "string"
	case "ListOption":
		o.go_type = "[]float64"
	case "MarkerOption":
		o.go_type = "Marker"
	case "MarkerDirichlet":
		o.go_type = "MarkerDirchlet"
	case "MarkerPeriodic":
		o.go_type = "MarkerPeriodic"
	case "MarkerInlet":
		o.go_type = "MarkerInlet"
	case "MarkerOutlet":
		o.go_type = "MarkerOutlet"
	case "MarkerDisplacement":
		o.go_type = "MarkerDisplacement"
	case "MarkerLoad":
		o.go_type = "MarkerLoad"
	case "MarkerFlowLoad":
		o.go_type = "MarkerFlowLoad"
	case "ScalarOption":
		o.go_type = "float64"
	case "SpecialOption":
		o.go_type = "bool"
	default:
		return errors.New("option type " + o.typeO + " not implemented")
	}
	return nil
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

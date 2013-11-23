package config

import (
	"bytes"
	//"fmt"
	"sort"

	"github.com/btracey/su2tools/config/common"
)

func init() {
	makeOptionOrder()
}

var configHeader []byte = []byte(`
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%                                                               %
% Stanford University unstructured (SU2) configuration file     %
%                                                               %
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
`)

var categoryBookend string = string(` ----------- `)

type optionOrderSorter []*optionPrint

func (o optionOrderSorter) Len() int {
	return len(o)
}

func (o optionOrderSorter) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o optionOrderSorter) Less(i, j int) bool {
	iCategory := o[i].ConfigCategory
	jCategory := o[j].ConfigCategory
	iOrder, ok := categoryOrder[iCategory]
	if !ok {
		panic("heading not in map")
	}
	jOrder, ok := categoryOrder[jCategory]
	if !ok {
		panic("heading not in map")
	}
	if iOrder < jOrder {
		return true
	}
	if iOrder > jOrder {
		return false
	}
	// Same category, sort by alphabetical order
	return o[i].ConfigName < o[j].ConfigName
}

var optionOrder optionOrderSorter

func makeOptionOrder() {
	optionOrder = make(optionOrderSorter, len(optionMap))
	i := 0
	for _, val := range optionMap {
		optionOrder[i] = val
		i++
	}
	sort.Sort(optionOrder)
}

// OptionPrint is a type for printing an option in SU2 format
type optionPrint struct {
	Description    string
	ConfigCategory common.ConfigCategory
	ConfigName     common.ConfigfileOption
	ConfigType     common.ConfigOptionType
	Default        string            // default value as a string
	GoBaseType     common.GoBaseType // go type string
	OptionsField   common.OptionsField
	enumOptions    []common.Enum
	//ValueString    string //enum string list
	Value interface{} // Default value as a value
}

// OptionList is a list of options to print while writing the config file
type OptionList map[common.OptionsField]bool

var PrintAll OptionList = OptionList{"All": true}

//var categoryOrder []string =

func (o *optionPrint) MarshalSU2Config() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteString("\n")
	buf.WriteString("% " + o.Description + "\n")
	buf.WriteString("% Type: " + string(o.ConfigType) + " Default: " + o.Default + "\n")
	if len(o.enumOptions) > 0 {
		buf.WriteString("% Options: (")
		for i, str := range o.enumOptions {
			buf.WriteString(string(str))
			if i != len(o.enumOptions)-1 {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(" )\n")
	}
	buf.WriteString(string(o.ConfigName) + "= ")
	return buf.Bytes(), nil
}

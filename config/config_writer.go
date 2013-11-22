package config

import (
	"bytes"
	//"fmt"
	"sort"
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

type optionOrderSorter []*OptionPrint

func (o optionOrderSorter) Len() int {
	return len(o)
}

func (o optionOrderSorter) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o optionOrderSorter) Less(i, j int) bool {
	iCategory := o[i].Category
	jCategory := o[j].Category
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
	return o[i].SU2OptionName < o[j].SU2OptionName
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
type OptionPrint struct {
	Description    string
	Category       string
	SU2OptionName  string
	OptionTypeName string
	Default        string
	Type           string // go type string
	StructName     string
	enumOptions    []string
}

// OptionList is a list of options to print while writing the config file
type OptionList map[string]bool

var PrintAll OptionList = OptionList{"All": true}

//var categoryOrder []string =

func (o *OptionPrint) MarshalSU2Config() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteString("\n")
	buf.WriteString("% " + o.Description + "\n")
	buf.WriteString("% Type: " + o.OptionTypeName + " Default: " + o.Default + "\n")
	if len(o.enumOptions) > 0 {
		buf.WriteString("% Options: (")
		for i, str := range o.enumOptions {
			buf.WriteString(str)
			if i != len(o.enumOptions)-1 {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(" )\n")
	}
	buf.WriteString(o.SU2OptionName + "= ")
	return buf.Bytes(), nil
}

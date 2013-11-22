package su2config

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"sort"
	"strconv"

	"fmt"
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

// WriteConfig writes a config file to the writer with the options given in list
func (o *Options) WriteConfig(writer io.Writer, list OptionList) error {

	buf := &bytes.Buffer{}
	buf.Write(configHeader)
	printAll := list["All"]
	currentHeading := -1
	for _, option := range optionOrder {
		catNumber := categoryOrder[option.Category]
		fmt.Println("cat # = ", catNumber)
		if catNumber > currentHeading {
			buf.WriteString("\n\n")
			buf.WriteString("%" + categoryBookend + option.Category + categoryBookend + "% \n")
			currentHeading = catNumber
		}
		printOption := list[option.SU2OptionName]
		if printAll {
			printOption = true
		}
		if printOption {
			b, _ := option.MarshalSU2Config()
			b2, err := o.marshalSU2StructValue(option.StructName)
			if err != nil {
				return errors.New("WriteConfig " + option.StructName + ": " + err.Error())
			}
			buf.Write(b)
			buf.Write(b2)
			buf.WriteString("\n")
		}
	}
	writer.Write(buf.Bytes())
	return nil
}

// This may need to change to be a switch on the SU2Type
func (o *Options) marshalSU2StructValue(structName string) ([]byte, error) {
	// Get the value of the field
	reflectvalue := reflect.ValueOf(*o)
	fieldvalue := reflectvalue.FieldByName(structName)
	switch t := fieldvalue.Interface().(type) {
	case float64:
		str := strconv.FormatFloat(t, 'g', -1, 64)
		return []byte(str), nil
	case bool:
		if t == false {
			return []byte("NO"), nil
		}
		return []byte("YES"), nil
	case string:
		return []byte(t), nil
	case []float64:
		buf := &bytes.Buffer{}
		buf.WriteString("( ")
		for i, val := range t {
			str := strconv.FormatFloat(val, 'g', -1, 64)
			buf.WriteString(str)
			if i != len(t)-1 {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(" )")
		return buf.Bytes(), nil
	case []string:
		buf := &bytes.Buffer{}
		buf.WriteString("(")
		for i, val := range t {
			buf.WriteString(val)
			if i != len(t)-1 {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(" )")
		return buf.Bytes(), nil
	default:
		panic("type not implemented for " + structName)
	}
}

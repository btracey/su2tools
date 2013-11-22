package config

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"

	"fmt"
)

func is_enum_option(t string) bool {
	switch t {
	default:
		return false
	case "ConvectOption", "EnumListOption", "MathProblem", "EnumOption":
		return true

	}
}

// IsField returns true if field is a field of SU2
func (o *Options) IsConfigOption(configString string) bool {
	goField := su2ToGoFieldMap[configString]
	_, ok := optionMap[goField]
	return ok
}

// IsEnum returns true if the options field is an enumable
func (o *Options) IsEnum(field string) bool {
	option := optionMap[field]
	return is_enum_option(option.OptionTypeName)
}

// SetEnum  sets the enumerable field, checking that it is a valid option.
// Returns an error if it is not a valid option or if
func (o *Options) SetEnum(field string, val string) error {
	return nil
}

// SetFields sets the fields of the options structure
func (o *Options) SetFields(fields map[string]interface{}) error {
	for field, value := range fields {
		option, ok := optionMap[field]
		if !ok {
			return errors.New("setfields: " + field + " is not a field")
		}
		if is_enum_option(option.OptionTypeName) {
			err := o.SetEnum(field, value)
			if err != nil {
				errors.New("setfelids: error setting field " + field + ": " + err.Error())
			}
			continue
		}
		mut := reflect.ValueOf(o).Elem()
		switch t := value.(type) {
		case float64:
			mut.FieldByName(field).SetFloat(t)
		case string:
			mut.FieldByName(field).SetString(t)
		case bool:
			mut.FieldByName(field).SetBool(t)
		case []float64:
			iface := mut.FieldByName(field).Interface()
			slice := iface.([]float64)
			slice = t
			fmt.Println("confirm this works")
		case []string:
			iface := mut.FieldByName(field).Interface()
			slice := iface.([]string)
			slice = t
			fmt.Println("confirm this works")
		default:
			panic("field type not implemented")
		}
	}
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

func ReadConfig(reader io.Reader) (*Options, error) {
	setter := make(map[string]interface{})
	// turn it into a scanner
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		if scanner.Bytes()[0] == '%' {
			continue
		}
		line = string(scanner.Bytes())
		// split the line at the equals sign
		parts := strings.Split(line, "=")
		if len(parts) > 2 {
			return errors.New("readconfig: option line has two equals signs")
		}
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		goFieldName, ok = su2ToGoFieldMap[parts[0]]
		if !ok {
			return errors.New("readconfig: option " + parts[0] + " not in struct")
		}
		option := optionMap[goFieldName]
		value, err := valueFromString(option, parts[1])

		//LEFT OFF HERE AND VALUE_FROM_STRING
	}

	options := NewOptions()
	_ = scanner
	_ = setter
	return nil
}

// TODO: Avoid replication in the generator

func valueFromString(option OptionPrint, str string) (interface{}, error) {

	//LEFT OFF HERE

	switch option.Type {
	case "float64":
		return strconv.ParseFloat(str, 64)
	default:
		return nil, "Unknown case"
	}
}

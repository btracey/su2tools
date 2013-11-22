package config

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"

	//"fmt"
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
	// First, check that it's a real field
	option, ok := optionMap[field]
	if !ok {
		return errors.New("setenum: " + field + " is not a valid field")
	}
	// check that the value is any of the possible options
	for _, str := range option.enumOptions {
		if str == val {
			reflect.ValueOf(o).Elem().FieldByName(field).SetString(val)
			return nil
		}
	}
	return errors.New("setenum: " + " val is not a valid option " + option.ValueString)
}

func (o *Options) SetField(field string, val interface{}) error {
	option, ok := optionMap[field]
	if !ok {
		return errors.New("setfields: " + field + " is not a field")
	}
	if is_enum_option(option.OptionTypeName) && option.OptionTypeName != "EnumListOption" {
		err := o.SetEnum(field, val.(string))
		if err != nil {
			errors.New("setfelids: error setting field " + field + ": " + err.Error())
		}
		return nil
	}
	reflect.ValueOf(o).Elem().FieldByName(field).Set(reflect.ValueOf(val))
	return nil
}

// SetFields sets the fields of the options structure
func (o *Options) SetFields(fieldMap map[string]interface{}) error {
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
		catNumber := categoryOrder[option.Category]
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
	//setter := make(map[string]interface{})
	// turn it into a scanner
	scanner := bufio.NewScanner(reader)

	options := NewOptions()
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		if scanner.Bytes()[0] == '%' {
			continue
		}
		line := string(scanner.Bytes())
		// split the line at the equals sign
		parts := strings.Split(line, "=")
		if len(parts) > 2 {
			return nil, errors.New("readconfig: option line has two equals signs")
		}
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		goFieldName, ok := su2ToGoFieldMap[parts[0]]
		if !ok {
			return nil, errors.New("readconfig: option " + parts[0] + " not in struct")
		}
		option := optionMap[goFieldName]
		value, err := valueFromString(option, parts[1])
		if err != nil {
			return nil, errors.New("readconfig: error parsing " + parts[0] + ": " + err.Error())
		}
		options.SetField(goFieldName, value)
	}
	err := scanner.Err()
	if err != nil {
		return nil, errors.New("readconfig: " + err.Error())
	}
	//_ = setter
	return options, nil
}

func valueFromString(option *OptionPrint, str string) (interface{}, error) {
	switch option.Type {
	case "float64":
		return strconv.ParseFloat(str, 64)
	case "bool":
		switch str {
		case "NO":
			return false, nil
		case "YES":
			return true, nil
		default:
			return false, errors.New("boolean string not YES or NO")
		}
	case "string":
		return str, nil
	case "[]float64":
		// parse the string
		strs, err := splitStringList(str)
		if err != nil {
			return nil, err
		}
		fs := make([]float64, len(strs))
		for i, str := range strs {
			f, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return nil, err
			}
			fs[i] = f
		}
		return fs, nil
	case "[]string":
		return splitStringList(str)
	default:
		return nil, errors.New("Unknown case " + option.Type)
	}
}

// CODE DUPLICATION HERE WITH CODE GENERATOR
func splitStringList(str string) ([]string, error) {
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

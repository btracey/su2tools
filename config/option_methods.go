package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	//"github.com/btracey/su2tools/config copy/common"
	"github.com/btracey/su2tools/config/common"
	"github.com/btracey/su2tools/config/su2types"
)

var delimiters = " ()[]{}:,\t\n\v\f\r"

// ForceAll is a convenience variable for forcing all of the config options to
// be printed
var ForceAll = map[Option]bool{All: true}

var configHeader = []byte(`
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%                                                               %
% Stanford University unstructured (SU2) configuration file     %
%                                                               %
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
`)

// Map from string to bool. If is nil than print non-default

// var WriteDifferent = []string{"Write_Different_Config_Options"} //

// Need to have a slice of fields ordered by category so when printing can iterate
// over the list in order and use reflect FieldByName

// WriteConfigTo writes the config file as a slice of bytes and puts it into the
// writer. WriteConfigTo will print all fields that are either different from
// the default vaule or are in the second argument. As a special case, if
// forcePrint contains All, all options will be printed.
func (o *Options) WriteTo(writer io.Writer, forcePrint map[Option]bool) (int, error) {
	// Loop over the config options

	var printAll bool

	if forcePrint != nil {
		_, printAll = forcePrint[All]
	}

	writer.Write(configHeader)

	nWritten := 0

	for i, options := range optionList {
		if len(options) != 0 {
			cat := categoryList[i]
			// Print the category name and description
			n, err := writer.Write([]byte("\n\n%%%%% " + cat.Name + "\n"))
			nWritten += n
			if err != nil {
				return nWritten, err
			}
			n, err = writer.Write([]byte("% " + cat.Description + "\n\n"))
			nWritten += n
			if err != nil {
				return nWritten, err
			}
		}
		optionsValue := reflect.ValueOf(o).Elem()
		defaultValue := reflect.ValueOf(defaultOptions).Elem()
		for _, opt := range options {

			optStruct := optionMap[opt]
			fieldValue := optionsValue.FieldByName(optStruct.Name)
			optStr := su2types.ConfigString(fieldValue.Interface())

			mustPrint := printAll
			if !mustPrint {
				// See if it's in the list of things we have to print
				if forcePrint != nil {
					_, mustPrint = forcePrint[optStruct.OptionConst]
				}
			}

			if !mustPrint {
				//fmt.Println("Checking if strings are the same")
				// See if the value is different from default
				defaultFieldValue := defaultValue.FieldByName(optStruct.Name)
				defStr := su2types.ConfigString(defaultFieldValue.Interface())

				mustPrint = defStr != optStr

			}

			if !mustPrint {
				// Value is same as default and nothing is forcing it
				continue
			}

			n, err := writer.Write([]byte("%  " + optStruct.Description + "\n"))
			nWritten += n
			if err != nil {
				return nWritten, err
			}
			n, err = writer.Write([]byte(optStruct.Config + "= " + optStr + "\n\n"))
			nWritten += n
			if err != nil {
				return nWritten, err
			}
		}
	}
	return nWritten, nil
}

func (o *Options) Copy() *Options {
	// Write to a buffer
	b := &bytes.Buffer{}
	_, err := o.WriteTo(b, ForceAll)
	if err != nil {
		str := "copy: error writing: " + err.Error()
		panic(str)
	}

	// now read
	options, _, err := Read(b)
	if err != nil {
		str := "copy: error reading " + err.Error()
		panic(str)
	}
	return options
}

/*
// NewOptions returns a new Options struct with all of the default values
func NewOptions() *Options {
	var o Options
	o = *defaultOptions
	return &o
}
*/

func Read(reader io.Reader) (*Options, map[Option]bool, error) {
	// First, get a new options struct populated with default values
	o := NewOptions()

	optionList := make(map[Option]bool)

	// Now, read through and parse all the values
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {

		if shouldcontinue(scanner) {
			continue
		}
		field, optionValues, err := getoption(scanner)
		if err != nil {
			return nil, nil, err
		}

		opt, ok := stringToOption[field]
		if !ok {
			return nil, nil, fmt.Errorf("Unknown field: " + field)
		}
		//optStr := optionMap[opt]

		_, ok = optionList[opt]
		if ok {
			return nil, nil, fmt.Errorf("field %s set multiple times", field)
		}

		// Get the existing value in the field
		v := reflect.ValueOf(o).Elem().FieldByName(field).Addr().Interface()

		// Using the type and the string, get the actual value
		err = su2types.FromConfigString(v, optionValues)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: error reading: %v string is: %v", field, err, optionValues)
		}

		optionList[opt] = true
		// Set the feild
		//reflect.ValueOf(o).Elem().FieldByName(string(field)).Set(reflect.ValueOf(value))
	}
	err := scanner.Err()
	if err != nil {
		return nil, nil, errors.New("readconfig: " + err.Error())
	}
	return o, optionList, nil
}

func ReadFromFile(filename string) (*Options, map[Option]bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	return Read(f)
}

func shouldcontinue(scanner *bufio.Scanner) bool {
	if len(scanner.Bytes()) == 0 {
		return true
	}
	if strings.TrimSpace(scanner.Text())[0] == '%' {
		return true
	}
	return false
}

func getoption(scanner *bufio.Scanner) (fieldString string, optionValues []string, err error) {
	line := string(scanner.Bytes())
	// split the line at the equals sign
	parts := strings.Split(line, "=")
	if len(parts) > 2 {
		return "", nil, errors.New("readconfig: option line has two equals signs")
	}
	if len(parts) == 1 {
		return "", nil, errors.New("readconfig: line \"" + parts[0] + "\" is not commented and has no equals sign")
	}

	//fmt.Println(line)
	// Split the second part by all of the delimiters of SU2
	stringRemains := true
	str := strings.TrimSpace(parts[1])
	strs := make([]string, 0)
	for i := strings.IndexAny(str, delimiters); i != -1; i = strings.IndexAny(str, delimiters) {
		newstr := str[:i]
		newstr = strings.TrimSpace(newstr)
		if len(newstr) > 0 {
			strs = append(strs, newstr)
		}
		if len(str)-1 == i {
			stringRemains = false
			break
		}
		str = str[i+1:]
	}
	if stringRemains {
		strs = append(strs, str)
	}
	fieldString = strings.TrimSpace(parts[0])
	fieldString = common.FixOptionId(fieldString)
	return fieldString, strs, nil
}

/*
// SetField sets the field with the value. It calls SetEnum if it is an enumerable option
func (o *Options) setField(field common.OptionsField, val interface{}) error {
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
*/

// Diff returns the options whose values are different
func Diff(options, options2 *Options) []string {
	if reflect.DeepEqual(options, options2) {
		return nil
	}

	var strs []string
	for _, cat := range optionList {
		for _, opt := range cat {
			v1 := reflect.ValueOf(options).Elem()
			v2 := reflect.ValueOf(options2).Elem()
			item1 := v1.FieldByName(string(opt)).Interface()
			item2 := v2.FieldByName(string(opt)).Interface()
			if !reflect.DeepEqual(item1, item2) {
				strs = append(strs, fmt.Sprintf("%s: item1 is %v, item2 is %v", opt, item1, item2))
			}
		}
	}
	return strs
}

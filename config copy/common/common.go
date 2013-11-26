package common

import (
	"errors"
	"strconv"
	"strings"

	//"fmt"
)

// ValueAsConfigString sets the value to the string that the config file needs
func ValueAsConfigString(value interface{}, gotype GoBaseType) string {
	switch gotype {
	case StringType:
		return value.(string)
	case EnumType:
		return string(value.(Enum))
	case BoolType:
		b := (value).(bool)
		if b {
			return "YES"
		} else {
			return "NO"
		}
	case Float64Type:
		f := (value).(float64)
		str := strconv.FormatFloat(f, 'g', -1, 64)
		return str
	case Float64ArrayType:
		str := "( "

		slice := (value).([]float64)
		for i, f := range slice {
			newStr := strconv.FormatFloat(f, 'g', -1, 64)
			str += newStr
			if i != len(slice)-1 {
				str += ", "
			}
		}
		str += " )"
		return str
	case StringArrayType:
		str := "( "
		strs := (value).([]string)
		for i, s := range strs {
			str += "\"" + s + "\""
			if i != len(strs)-1 {
				str += ","
			}
		}
		str += " )"
		return str
	case BadType:
		panic("bad type found")
	default:
		panic("unknown type ")
	}
}

// ValueAsString gets the value as a string such that it could be written (i.e. "[]float64{1,2}")
func ValueAsString(value interface{}, gotype GoBaseType) string {
	switch gotype {
	case StringType:
		return "\"" + value.(string) + "\""
	case EnumType:
		return "\"" + string(value.(Enum)) + "\""
	case BoolType:
		b := (value).(bool)
		if b {
			return "false"
		} else {
			return "true"
		}
	case Float64Type:
		f := (value).(float64)
		str := strconv.FormatFloat(f, 'g', -1, 64)
		return str
	case Float64ArrayType:
		str := "[]float64{"

		slice := (value).([]float64)
		for i, f := range slice {
			newStr := strconv.FormatFloat(f, 'g', -1, 64)
			str += newStr
			if i != len(slice)-1 {
				str += ","
			}
		}
		str += "}"
		return str
	case StringArrayType:

		strs := (value).([]string)
		if len(strs) == 0 {
			return "nil"
		}
		str := "[]string{"
		for i, s := range strs {
			str += "\"" + s + "\""
			if i != len(strs)-1 {
				str += ","
			}
		}
		str += "}"
		return str

	case BadType:
		panic("bad type found")
	default:
		panic("unknown type ")
	}
}

// InterfaceFromString returns the value as an interface{} from a string
func InterfaceFromString(str string, gotype GoBaseType) (interface{}, error) {
	switch gotype {
	case Float64Type:
		return strconv.ParseFloat(str, 64)
	case BoolType:
		switch str {
		case "NO":
			return false, nil
		case "YES":
			return true, nil
		default:
			return false, errors.New("boolean string not YES or NO")
		}
	case StringType:
		return str, nil
	case EnumType:
		return Enum(str), nil
	case Float64ArrayType:
		// parse the string
		strs, err := SplitArrayOption(str)
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
	case StringArrayType:
		return SplitArrayOption(str)
	default:
		return nil, errors.New("Unknown case " + string(gotype))
	}
}

// IsEnumOption returns true if the config type represents
// an enumerable option, i.e. a field that can take one of
// a list of strings
func IsEnumOption(t ConfigOptionType) bool {
	switch t {
	default:
		return false
	case "ConvectOption", "MathProblem", "EnumOption":
		return true
	}
}

// SplitArrayOption splits an option array string into its components
func SplitArrayOption(str string) ([]string, error) {
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
	if len(strs) == 1 && strs[0] == "" {
		return nil, nil
	}
	return strs, nil
}

func IsFloat(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return false
	}
	return true
}

func IsFloatArray(strs []string) bool {
	if len(strs) == 0 {
		return false
	}
	for _, str := range strs {
		if !IsFloat(str) {
			return false
		}
	}
	return true
}

func StringsToFloatArray(strs []string) []float64 {
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

const (
	BadType GoBaseType = iota
	BoolType
	StringType
	Float64Type
	EnumType
	Float64ArrayType
	StringArrayType
)

var GoTypeToKind map[GoBaseType]string = map[GoBaseType]string{

	BoolType:         "bool",
	StringType:       "string",
	Float64Type:      "float64",
	EnumType:         "common.Enum",
	Float64ArrayType: "[]float64",
	StringArrayType:  "[]string",
}

var GoTypeToString map[GoBaseType]string = map[GoBaseType]string{

	BoolType:         "common.BoolType",
	StringType:       "common.StringType",
	Float64Type:      "common.Float64Type",
	EnumType:         "common.EnumType",
	Float64ArrayType: "common.Float64ArrayType",
	StringArrayType:  "common.StringArrayType",
}

func OptionStringToInterface(gotype GoBaseType, entry string) (interface{}, error) {
	switch gotype {
	default:
		panic("type not implemented " + string(gotype))
	case StringType:
		return entry, nil
	case EnumType:
		return Enum(entry), nil
	case BoolType:
		switch entry {
		case "NO":
			return false, nil
		case "YES":
			return true, nil
		default:
			return nil, errors.New("bad bool option")
		}
	case Float64Type:
		f, err := strconv.ParseFloat(entry, 64)
		return f, err
	case Float64ArrayType:
		// parse the floats
		strs, err := SplitArrayOption(entry)
		if err != nil {
			return nil, err
		}
		fs := StringsToFloatArray(strs)
		if fs == nil {
			return nil, errors.New("bad float parse")
		}
		return fs, nil
	case StringArrayType:
		strs, _ := SplitArrayOption(entry)
		return strs, nil
	}
}

// ConfigTypeToGoType takes in a type from the config file
// and returns the go base type as represented by a string
// panics if the type is not in the list
// TODO: It would be nice if this could be a map, but it can't be because
// of overloading of Scalar and List options
func ConfigTypeToGoType(t ConfigOptionType, defaultValue string) (GoBaseType, error) {
	// write the type

	if IsEnumOption(t) {
		return EnumType, nil
	}

	switch t {
	case "ArrayOption":
		strs, err := SplitArrayOption(defaultValue)
		if err != nil {
			return BadType, errors.New(err.Error())
		}
		// check if it's a float array
		if IsFloatArray(strs) {
			return Float64ArrayType, nil
		}
		return StringArrayType, nil
	case "EnumListOption":
		return StringArrayType, nil
	case "ListOption":
		return StringType, nil
	case "DVParamOption":
		return StringType, nil
	case "MarkerOption":
		return StringType, nil
	case "MarkerDirichlet":
		return StringType, nil
	case "MarkerPeriodic":
		return StringType, nil
	case "MarkerInlet":
		return StringType, nil
	case "MarkerOutlet":
		return StringType, nil
	case "MarkerDisplacement":
		return StringType, nil
	case "MarkerLoad":
		return StringType, nil
	case "MarkerFlowLoad":
		return StringType, nil
	case "ScalarOption":
		// See if it looks like a float 64
		_, err := strconv.ParseFloat(defaultValue, 64)
		if err != nil {
			return StringType, nil
		}
		return Float64Type, nil
	case "SpecialOption":
		return BoolType, nil
	default:
		return BadType, errors.New("option type " + string(t) + " not implemented")
	}
}

// ConfigfileOption is a string in the config file
type ConfigfileOption string

// OptionsField is a field of the options struct
type OptionsField string

// ConfigOptionType is a string representing the type of config option (EnumOption, ArrayOption, etc.)
type ConfigOptionType string

// ConfigCategory is the configuration category to which the option belongs
type ConfigCategory string

// GoBaseType is a string representing the underlying base type of the value
type GoBaseType int

// An EnumOption is a valid setting for an EnumOption
type Enum string

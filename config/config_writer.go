package su2config

import ()

type Enum string //This could be better

type OptionKind interface {
	Base() OptionBase
}

type OptionBase struct {
	Name        string // name of the option in cfg file
	Description string // description of the option
	Heading     string // under which heading does this go
	FieldName   string // name of the options field
}

// EnumOption represents an option that has a fixed number of values
type EnumOption struct {
	Value          Enum
	PossibleValues string
	standard       string
	OptionBase
}

func (e EnumOption) Base() OptionBase {
	return e.OptionBase
}

// ScalarOption represents an option that can be a float64
type ScalarOption struct {
	Value    float64
	standard float64
	OptionBase
}

func (s ScalarOption) Base() OptionBase {
	return s.OptionBase
}

// BoolOption represents an option that can be a float64
type BoolOption struct {
	Value    bool
	standard bool
	OptionBase
}

func (s ScalarOption) Base() OptionBase {
	return s.OptionBase
}

var allSU2Options map[string]Option

func init() {
	addSU2Options()
}

// need funciton to get default options

// list of options
func addSU2Options() {}

// use list of strings for specific printing of options

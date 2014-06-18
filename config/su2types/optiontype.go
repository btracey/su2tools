// Provides a set of marshalers and unmarshalers for the option types

package su2types

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/btracey/su2tools/config/enum"
)

func init() {
	i := enum.Compressible
	_ = ConfigMarshaler(i)
}

var delimiters = " ()[]{}:,\t\n\v\f\r"

type ConfigMarshaler interface {
	ConfigString() string
}

type ConfigUnmmarshaler interface {
	FromConfigString([]string) error
}

// ConfigString converts the option type into a string for going into SU2
func ConfigString(v interface{}) string {
	switch t := v.(type) {
	default:
		// If the type is a ConfigMarshaler (enum types, many of the custom ones,
		// etc.), print what it is.
		cm, ok := t.(ConfigMarshaler)
		if ok {
			return cm.ConfigString()
		}

		// If it is a slice type and the length is zero, print none, or call
		// this recursively to get the actual value
		val := reflect.ValueOf(t)
		kind := val.Kind()
		if kind == reflect.Slice || kind == reflect.Array {
			if val.Len() == 0 {
				return "NONE"
			}
			sliceStr := "("
			for i := 0; i < val.Len(); i++ {
				if i != 0 {
					sliceStr += ", "
				}
				sliceStr += ConfigString(val.Index(i).Interface())
			}
			sliceStr += ")"
			return sliceStr
		}

		// Otherwise, return the sprintf version. This works for all the base types

		return fmt.Sprintf("%v", t)
		/*
			fmt.Println("Not a configMarshaler")
			fmt.Println("Printing v: ")
			fmt.Println(v)
			fmt.Println(reflect.TypeOf(v))
			panic("unknown type")
		*/
	/*
		case float64:
			return strconv.FormatFloat(t, 'e', -1, 64)
		case string:
			return t
		case int32:
			return strconv.Itoa(int(t))
		case uint64:
			return strconv.FormatUint(t, 10)
		case uint16:
			return strconv.FormatUint(uint64(t), 10)
		case int:
			return strconv.Itoa(t)
	*/
	case bool:
		b := v.(bool)
		if b {
			return "YES"
		}
		return "NO"
	case []string:
		if len(t) == 0 {
			return "NONE"
		}
		str := t[0]
		for i := 1; i < len(t); i++ {
			str += ", " + t[i]
		}
		return str
	}

}

//func FromConfigString(typename string, extratype string, values []string) (interface{}, error) {
func FromConfigString(v interface{}, values []string) error {
	var err error
	switch pt := v.(type) {
	default:
		// See if it is a config unmarshaler
		configUnmarshaler, ok := pt.(ConfigUnmmarshaler)
		if ok {
			err := configUnmarshaler.FromConfigString(values)
			return err
		}

		// The field may be a pointer, in which case we have a pointer to a pointer
		// Dereference the outer pointer, and see if that is a configUnmarshaler
		valueOfPT := reflect.ValueOf(pt)
		derefPT := reflect.Indirect(valueOfPT)
		if derefPT.Kind() == reflect.Ptr {
			configUnmarshaler, ok = derefPT.Interface().(ConfigUnmmarshaler)
			if ok {
				err := configUnmarshaler.FromConfigString(values)
				return err
			}
			// Shouldn't need to assign, because only the type changed
			return nil
		}

		// If we have a pointer to a slice type, then range over the slice
		// and unmarshal each element
		if derefPT.Kind() == reflect.Array {
			l := derefPT.Len()
			if len(values) != l {
				str := "bad array length: " + strconv.Itoa(l) + " found, " + strconv.Itoa(len(values)) + " expected."
				return errors.New(str)
			}
			for i := 0; i < l; i++ {
				// Get the address of the ith element
				arrayElementPtr := derefPT.Index(i).Addr().Interface()

				// Call this function with the string in order to dereference it
				err := FromConfigString(arrayElementPtr, values[i:i+1])
				if err != nil {
					return err
				}
			}
			return nil
		}

		// If we have a pointer to a slice type, then range over the slice
		// and unmarshal each element
		if derefPT.Kind() == reflect.Slice {
			elementType := reflect.TypeOf(derefPT.Interface())
			// If it's a slice, then if the length is one and the thing is none,
			// assign the pointer to zero
			if len(values) == 1 && values[0] == "NONE" {
				// Make a new slice of zero length of that type
				newSlice := reflect.MakeSlice(elementType, 0, 0)

				derefPT.Set(newSlice)
				return nil
			}

			l := len(values)
			// Otherwise, the length of the slice is however many values there are
			newSlice := reflect.MakeSlice(elementType, l, l)

			for i := 0; i < l; i++ {
				// Get the address of the ith element
				sliceElementPtr := newSlice.Index(i).Addr().Interface()

				// Call this function with the string in order to dereference it
				err := FromConfigString(sliceElementPtr, values[i:i+1])
				if err != nil {
					return err
				}
			}
			derefPT.Set(newSlice)
			return nil
		}

		fmt.Println("Not an unmarshaler")
		fmt.Println(valueOfPT)
		str := "unknown type"
		return errors.New(str)

	case *float64:
		if err = oneValue("float64", values); err != nil {
			return err
		}
		float, err := strconv.ParseFloat(values[0], 64)
		if err != nil {
			return err
		}
		*pt = float
		return nil
	case *string:
		if err = oneValue("string", values); err != nil {
			return err
		}
		*pt = values[0]
		return nil
	case *int32:
		if err = oneValue("int32", values); err != nil {
			return err
		}
		val, err := strconv.ParseInt(values[0], 10, 32)
		if err != nil {
			return err
		}
		*pt = int32(val)
		return nil
	case *int:
		if err = oneValue("int", values); err != nil {
			return err
		}
		val, err := strconv.Atoi(values[0])
		if err != nil {
			return err
		}
		*pt = val
		return nil
	case *uint64:
		if err = oneValue("uint64", values); err != nil {
			return err
		}
		val, err := strconv.ParseUint(values[0], 10, 64)
		if err != nil {
			return err
		}
		*pt = val
		return nil
	case *uint16:
		if err = oneValue("int16", values); err != nil {
			return err
		}
		val, err := strconv.ParseUint(values[0], 10, 16)
		if err != nil {
			return err
		}
		*pt = uint16(val)
		return nil
	case *bool:
		if err = oneValue("bool", values); err != nil {
			return err
		}
		if values[0] == "YES" {
			*pt = true
			return nil
		}
		if values[0] == "NO" {
			*pt = false
			return nil
		}
		return errors.New("bad boolean value: " + values[0])
		/*
			case *[]string:
				if values[0] == "NONE" {
					*pt = nil
					return nil
				}
				*pt = values
				return nil
		*/
		/*
			case **Convect:
				c := &Convect{}
				err := c.FromConfigString(values)
				if err != nil {
					return err
				}
				*pt = c
				return nil
			case **DVParam:
				c := &DVParam{}
				err := c.FromConfigString(values)
				if err != nil {
					return err
				}
				*pt = c
				return nil
			case **StringDoubleList:
				c := &StringDoubleList{}
				err := c.FromConfigString(values)
				if err != nil {
					return err
				}
				*pt = c
				return nil
			case **Periodic:
				c := &Periodic{}
				err := c.FromConfigString(values)
				if err != nil {
					return err
				}
				*pt = c
				return nil
			case **Inlet:
				c := &Inlet{}
				err := c.FromConfigString(values)
				if err != nil {
					return err
				}
				*pt = c
				return nil
			case **InletFixed:
				c := &InletFixed{}
				err := c.FromConfigString(values)
				if err != nil {
					return err
				}
				*pt = c
				return nil
		*/
		/*
			case *[2]float64:
				if len(values) != 2 {
					return errors.New("Wrong number of values")
				}
				var err error
				for i := range *pt {
					(*pt)[i], err = strconv.ParseFloat(values[i], 64)
					if err != nil {
						return err
					}
				}
				return nil
			case *[3]float64:
				if len(values) != 3 {
					return errors.New("Wrong number of values")
				}
				var err error
				for i := range *pt {
					(*pt)[i], err = strconv.ParseFloat(values[i], 64)
					if err != nil {
						return err
					}
				}
				return nil
			case *[6]float64:
				if len(values) != 6 {
					return errors.New("Wrong number of values")
				}
				var err error
				for i := range *pt {
					(*pt)[i], err = strconv.ParseFloat(values[i], 64)
					if err != nil {
						return err
					}
				}
				return nil
		*/
		/*
			case *[]float64:
				fmt.Println("In slice of floats")
				if len(values) == 1 && values[0] == "NONE" {
					*pt = nil
					return nil
				}
				*pt = make([]float64, len(values))
				for i := range *pt {
					(*pt)[i], err = strconv.ParseFloat(values[i], 64)
					if err != nil {
						return err
					}
				}
				return nil
		*/
	}
}

// oneValue checks that one value is present
func oneValue(typename string, values []string) error {
	if len(values) == 1 {
		return nil
	}
	nValues := strconv.Itoa(len(values))

	return fmt.Errorf("type: ", typename, " ", nValues, " found, 1 expected")
}

/*
// This should be like the SU2 version
func tokenizeString(str string) []string {
	if str == "" {
		return []string{}
	}
	strs := make([]string, 0)
	for i := strings.IndexAny(str, delimiters); i != -1; {
		newstr := str[:i]
		newstr = strings.TrimSpace(newstr)
		strs = append(strs, newstr)
		str = str[i+1:]
	}
	return strs
}
*/

type Convect struct {
	String string
}

func (c *Convect) ConfigString() string {
	if c == nil || c.String == "" {
		return "NONE"
	}
	return c.String
}

func (c *Convect) FromConfigString(values []string) error {
	if len(values) == 1 && values[0] == "NONE" {
		c.String = ""
		return nil
	}
	for i, s := range values {
		c.String += s
		if i != len(values)-1 {
			c.String += " "
		}
	}
	return nil
}

type DVParam struct {
	String string
}

func (c *DVParam) ConfigString() string {
	if c == nil || c.String == "" {
		return "NONE"
	}
	return c.String
}

func (c *DVParam) FromConfigString(values []string) error {
	if len(values) == 1 && values[0] == "NONE" {
		c.String = ""
		return nil
	}
	for i, s := range values {
		c.String += s
		if i != len(values)-1 {
			c.String += " "
		}
	}
	return nil
}

type StringDoubleList struct {
	Strings []string
	Doubles []float64
}

func (c *StringDoubleList) ConfigString() string {
	if c == nil || len(c.Strings) == 0 {
		return "NONE"
	}
	if len(c.Strings) != len(c.Doubles) {
		panic("lengths must match")
	}
	var str string
	for i, s := range c.Strings {
		str += s + ", "
		str += strconv.FormatFloat(c.Doubles[i], 'g', 20, 64)
		if i != len(c.Strings)-1 {
			str += ", "
		}
	}
	return str
}

func (c *StringDoubleList) FromConfigString(values []string) error {
	if len(values) == 1 && values[0] == "NONE" {
		c.Strings = nil
		c.Doubles = nil
		return nil
	}
	if (len(values) % 2) != 0 {
		return errors.New("must have even number of values")
	}
	c.Strings = make([]string, len(values)/2)
	c.Doubles = make([]float64, len(values)/2)
	for i := range c.Strings {
		c.Strings[i] = values[2*i]
		f64, err := strconv.ParseFloat(values[2*i+1], 64)
		if err != nil {
			return err
		}
		c.Doubles[i] = f64
	}
	return nil
}

type Inlet struct {
	Strings []string
}

func (c *Inlet) ConfigString() string {
	if c == nil || len(c.Strings) == 0 {
		return "NONE"
	}
	printstring := []byte("(")
	for i, str := range c.Strings {
		printstring = append(printstring, []byte(str)...)
		if i != len(c.Strings)-1 {
			printstring = append(printstring, []byte(", ")...)
		}
	}
	printstring = append(printstring, []byte(")")...)
	return string(printstring)
}

func (c *Inlet) FromConfigString(values []string) error {
	if len(values) == 1 && values[0] == "NONE" {
		c.Strings = nil
		return nil
	}
	c.Strings = make([]string, len(values))
	for i, s := range values {
		c.Strings[i] = s
	}
	return nil
}

type InletFixed struct {
	String string
}

func (c *InletFixed) ConfigString() string {
	if c == nil || c.String == "" {
		return "NONE"
	}
	return c.String
}

func (c *InletFixed) FromConfigString(values []string) error {

	if len(values) == 1 && values[0] == "NONE" {
		c.String = ""
		return nil
	}
	for i, s := range values {
		c.String += s
		if i != len(values)-1 {
			c.String += " "
		}
	}
	return nil
}

type ActuatorDisk struct {
	String string
}

func (c *ActuatorDisk) ConfigString() string {
	if c == nil || c.String == "" {
		return "NONE"
	}
	return c.String
}

func (c *ActuatorDisk) FromConfigString(values []string) error {
	if len(values) == 1 && values[0] == "NONE" {
		c.String = ""
		return nil
	}
	for i, s := range values {
		c.String += s
		if i != len(values)-1 {
			c.String += " "
		}
	}
	return nil
}

type Periodic struct {
	String string
}

func (c *Periodic) ConfigString() string {
	if c == nil || c.String == "" {
		return "NONE"
	}
	return c.String
}

func (c *Periodic) FromConfigString(values []string) error {
	if len(values) == 1 && values[0] == "NONE" {
		c.String = ""
		return nil
	}
	for i, s := range values {
		c.String += s
		if i != len(values)-1 {
			c.String += " "
		}
	}
	return nil
}

type MathProblem struct {
	String string
}

func (c *MathProblem) ConfigString() string {
	if c == nil || c.String == "" {
		return "NONE"
	}
	return c.String
}

func (c *MathProblem) FromConfigString(values []string) error {
	if len(values) == 1 && values[0] == "NONE" {
		c.String = ""
		return nil
	}
	for i, s := range values {
		c.String += s
		if i != len(values)-1 {
			c.String += " "
		}
	}
	return nil
}

type Python struct {
	String string
}

func (c *Python) ConfigString() string {
	if c == nil || c.String == "" {
		return "NONE"
	}
	return c.String
}

func (c *Python) FromConfigString(values []string) error {
	if len(values) == 1 && values[0] == "NONE" {
		c.String = ""
		return nil
	}
	for i, s := range values {
		c.String += s
		if i != len(values)-1 {
			c.String += " "
		}
	}
	return nil
}

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/btracey/su2tools/config/common"
	enumautogen "github.com/btracey/su2tools/config/enum/autogen"
	enumautogencommon "github.com/btracey/su2tools/config/enum/autogen/autogencommon"
)

var su2home string
var gopath string
var savepath string // where to save the finished JSON files
var configPath string
var configImportPath string

func init() {
	// Get gopath and su2home
	gopath = os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatal("GOPATH must be set")
	}
	su2home = os.Getenv("SU2_HOME")
	if su2home == "" {
		log.Fatal("SU2_HOME must be set")
	}
	configImportPath = "github.com/btracey/su2tools/config"
	savepath = filepath.Join(gopath, "src", "github.com", "btracey", "su2tools", "config", "autogen_options")
	configPath = filepath.Join(gopath, "src", "github.com", "btracey", "su2tools", "config")

}

type ConfigOption struct {
	Value        string // Gofield name of the string
	ConfigString string // string to print to config file
	Category     int    // What category is it
	Description  string
	Type         string // The kind of option it is (as a string)
	ExtraType    string // Extra info needed for the type (enum type, array size, etc.)
	Default      string // The default value for go expressed as a string (will be literally copied into the defaultStruct)
}

type ConfigCategory struct {
	Id          int
	Name        string
	Description string
}

/*
type OptionType struct {
	//Id   int
	Name string
}
*/

func main() {
	fmt.Println("Be sure to run the enum generation script before this one")
	categories, options := parseconfig()

	// Write all of these to a file
	s := struct {
		Categories []*ConfigCategory
		Options    []*ConfigOption
	}{
		Categories: categories,
		Options:    options,
	}

	optionFilename := filepath.Join(savepath, "option_file.json")
	f, err := os.Create(optionFilename)
	if err != nil {
		log.Fatal(err.Error())
	}

	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		log.Fatal(err.Error())
	}
	f.Write(b)

	writeConfigAndDefault(categories, options)

	writeCategoriesAndOptionOrder(categories, options)

	writeoptions(categories, options)

}

func parseconfig() ([]*ConfigCategory, []*ConfigOption) {
	configFilename := filepath.Join(su2home, "Common", "src", "config_structure.cpp")
	configFile, err := os.Open(configFilename)
	if err != nil {
		log.Fatal("error opening option: " + err.Error())
	}
	defer configFile.Close()

	scanner := bufio.NewScanner(configFile)

	lines := getlines(scanner)

	var categoryId int = -1
	var categories []*ConfigCategory

	//var optionTypes []*OptionType
	//optionTypeMap := make(map[string]*OptionType)

	optionMap := make([]*ConfigOption, 0)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.Contains(line, "CONFIG_CATEGORY:") {

			i, categories, categoryId = addConfigCategory(i, categories, lines)
			continue
		}
		if strings.HasPrefix(line, "add") {
			line = strings.TrimPrefix(line, "add")
			strs := strings.Split(line, "Option(")
			if len(strs) != 2 {
				log.Fatalf("bad assumption in addoption: %s", strs)
			}

			optionTypeName := strs[0]
			/*
				_, ok := optionTypeMap[optionTypeName]
				if !ok {
					//optionTypeId = len(optionTypes)

					newOpt := &OptionType{
						//	Id:   optionTypeId,
						Name: optionTypeName,
					}
					//	optionTypes = append(optionTypes, newOpt)
					optionTypeMap[optionTypeName] = newOpt
				}
			*/

			// The description is all of the commented lines before it
			// (don't modify i directly because we don't want to rewind the next parse)
			prevIdx := 1
			description := lines[i-prevIdx]
			if !strings.HasSuffix(description, "*/") {
				fmt.Println("desc: ", description)
				fmt.Println("lines: ", lines[i])
				log.Fatal("Config type has no comment")
			}
			for !strings.HasPrefix(description, "/*") {
				prevIdx++
				description = lines[i-prevIdx] + description
			}

			description = strings.TrimPrefix(description, "/*")
			description = strings.TrimSuffix(description, "*/")
			description = strings.TrimSpace(description)
			description = strings.TrimPrefix(description, "DESCRIPTION:")
			description = strings.TrimSpace(description)

			funcCall := strs[1]
			for strings.HasSuffix(funcCall, ",") {
				// The function call continues to the next line
				i++
				funcCall += lines[i]
			}

			if !strings.HasSuffix(funcCall, ");") {
				log.Fatalf("weird function format, does not end with );")
			}

			funcCall = strings.TrimSuffix(funcCall, ");")
			// Break the function call into arguments
			args := strings.Split(funcCall, ",")
			for i, str := range args {
				args[i] = strings.TrimSpace(str)
			}

			// The first argument is always the type name
			configString := args[0]
			configString = strings.TrimPrefix(configString, "\"")
			configString = strings.TrimSuffix(configString, "\"")

			goString := common.FixOptionId(configString)

			defaultString, extraInfo, err := defaultValueFromArgs(args, optionTypeName, lines, i-prevIdx)
			if err != nil {
				log.Fatalf("Error getting default string from args: %s\n args: %v", err, args)
			}

			/*
				_, ok := optionMap[configString]
				if ok {
					log.Fatal("Config string appears twice: ", configString)
				}
			*/

			optionMap = append(optionMap, &ConfigOption{
				Value:        goString,
				ConfigString: configString,
				Category:     categoryId,
				Description:  description,
				Type:         optionTypeName,
				ExtraType:    extraInfo,
				Default:      defaultString,
			})

		}
	}
	return categories, optionMap
}

func getlines(scanner *bufio.Scanner) []string {
	var beginFound bool
	for scanner.Scan() {
		line := scanner.Text()
		// Keep scanning until we find the BEGIN_CONFIG_OPTIONS line
		if strings.Contains(line, "BEGIN_CONFIG_OPTIONS") {
			beginFound = true
			break
		}
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err().Error())
	}
	if !beginFound {
		log.Fatal("BEGIN_CONFIG_OPTIONS string not found")
	}

	lines := make([]string, 800)
	var endFound bool
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "END_CONFIG_OPTIONS") {
			endFound = true
			break
		}
		line = strings.TrimSpace(line)
		if len(line) != 0 {
			lines = append(lines, line)
		}
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err().Error())
	}
	if !endFound {
		log.Fatal("END_CONFIG_OPTIONS string not found")
	}
	return lines
}

func addConfigCategory(i int, categories []*ConfigCategory, lines []string) (int, []*ConfigCategory, int) {
	line := lines[i]
	strs := strings.Split(line, "CONFIG_CATEGORY:")
	if len(strs) != 2 {
		log.Fatalf("bad config category line: %s", line)
	}
	// Don't care about before
	str := strs[1]
	configId := len(categories)

	if !strings.Contains(str, "*/") {
		log.Fatalf("No comment ender in config category")
	}

	strs = strings.Split(str, "*/")

	name := strings.TrimSpace(strs[0])

	// Comment is the next line
	i++
	line = lines[i]
	if !strings.HasPrefix(line, "/*") {
		log.Fatalf("comment description wrong %s", line)
	}
	if !strings.HasSuffix(line, "*/") {
		log.Fatalf("comment description wrong %s", line)
	}
	line = strings.TrimPrefix(line, "/*")
	line = strings.TrimSuffix(line, "*/")

	categories = append(categories, &ConfigCategory{
		Name:        name,
		Id:          configId,
		Description: line,
	})
	return i, categories, configId
}

// this returns a go string of the typed value (what will be written as a type)
// prevIdx is the line where the description started
func defaultValueFromArgs(args []string, optionType string, lines []string, prevIdx int) (string, string, error) {
	switch optionType {
	default:
		return "", "", errors.New("UnknownCategory: " + optionType)
	case "Enum":
		// TODO: Find some way to check this
		enumTypeVal := enumautogencommon.FixEnumId(args[3])

		typename, err := enumExtra(args[2])
		if err != nil {
			return "", "", err
		}
		return enumTypeVal, typename, nil
	case "Bool":
		return boolstring(args[2]), "", nil
	case "MathProblem":
		/*
			str := "[]bool{"
			str += boolstring(args[2]) + ", "
			str += boolstring(args[4]) + ", "
			str += boolstring(args[6]) + ", "
			str += boolstring(args[8])
			str += "}"
			return str, "", nil
		*/
		return "enum.DirectProblem", "", nil
	case "String":
		str := args[2]
		str = strings.TrimPrefix(str, "string(")
		str = strings.TrimSuffix(str, ")")
		return str, "", nil
	case "Double", "UnsignedLong", "UnsignedShort", "Long":
		// TODO: Add check it's a real double
		return args[2], "", nil
	case "StringList":
		return "[]string{}", "", nil
	case "DoubleList":
		return "[]float64{}", "", nil
	case "UShortList":
		return "[]uint16{}", "", nil
	case "StringDoubleList", "Periodic", "Inlet", "InletFixed", "DVParam":
		return "&su2types." + optionType + "{}", "", nil
	//	return "&su2types.Periodic{}", "", nil
	//case "Inlet", "InletFixed","DoubleList", "UShortList", "DVParam":
	//	return "NONE", "", nil
	case "Convect":
		return "enum.NoConvective", "", nil
	//case "DVParam":
	//	return "&su2types.DVParam{}", "", nil
	case "EnumList":
		typename, err := enumExtra(args[3])
		if err != nil {
			return "", "", err
		}
		return "NONE", typename, nil
	case "DoubleArray":
		arraySizeStr := args[1]

		arraySize, err := strconv.Atoi(arraySizeStr)
		if err != nil {
			return "", "", errors.New("bad parsing of array size")
		}

		// Go to the line before the description and find all of the lines which
		// initialize the default vector
		defVecString := ""
		for {
			prevIdx--
			line := lines[prevIdx]
			if !strings.HasPrefix(line, "default_vec") {
				break
			}
			defVecString = line + defVecString
		}
		defVecString = strings.TrimSpace(defVecString)

		// Split the vector into each vector component
		vecComps := strings.Split(defVecString, ";")
		if vecComps[len(vecComps)-1] == "" {
			vecComps = vecComps[:len(vecComps)-1]
		}
		if len(vecComps) != arraySize {
			return "", "", errors.New("Wrong number of array elements parsed")
		}

		arrayString := "[" + arraySizeStr + "]float64{"

		for i, comp := range vecComps {
			comp = strings.TrimSpace(comp)
			if !strings.HasPrefix(comp, "default_vec_") {
				return "", "", errors.New("Bad vec string")
			}
			comp = strings.TrimPrefix(comp, "default_vec_")

			// Get the number of default_vec
			split := strings.Split(comp, "d")

			integer, err := strconv.Atoi(split[0])
			if err != nil {
				return "", "", errors.New("bad parse")
			}

			if integer != arraySize {
				fmt.Println("integer = ", integer)
				fmt.Println("arraySize = ", arraySize)
				return "", "", errors.New("Default vec integer length doesn't match actual length")
			}

			// Get the index which is after the d
			comp = split[1]

			split = strings.Split(comp, "]")
			idxStr := strings.TrimPrefix(split[0], "[")
			idx, err := strconv.Atoi(idxStr)
			if err != nil {
				return "", "", errors.New("Bad parsing of index")
			}
			if idx != i {
				return "", "", errors.New("Parsing script assumes that indices are in order")
			}

			// Lastly, get the actual value
			comp = split[1]
			split = strings.Split(comp, "=")

			comp = split[1]
			comp = strings.TrimSpace(comp)

			arrayString += comp
			if i != arraySize-1 {
				arrayString += ", "
			}
		}
		arrayString += "}"

		extraInfo := args[1]

		return arrayString, extraInfo, nil
	}
}

func enumExtra(rawmap string) (string, error) {
	mapName := enumautogencommon.FixMapname(rawmap)
	typename, ok := enumautogen.MapnameToTypename[mapName]
	if !ok {
		return "", errors.New("Unknown enum map " + mapName)
	}
	return typename, nil
}

func boolstring(str string) string {
	if str == "true" {
		return "true"
	}
	if str == "false" {
		return "false"
	}
	panic("bad bool")
}

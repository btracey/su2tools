package common

import "strings"

// toCamelCase takes a string, makes it lower case, capitalizes the first letter,
// removes all underscores and capitalizes the letter after the underscore
func ToCamelCase(str string) string {

	str = strings.ToLower(str)
	strs := strings.Split(str, "_")

	newstring := ""
	for _, s := range strs {
		newstring += strings.ToUpper(string(s[0]))
		if len(s) > 1 {
			newstring += s[1:]
		}
	}
	return newstring
}

func FixOptionId(option string) string {
	option = strings.TrimSpace(option)
	option = strings.TrimPrefix(option, "\"")
	option = strings.TrimSuffix(option, "\"")
	option = strings.ToLower(option)
	option = ToCamelCase(option)
	return option
}

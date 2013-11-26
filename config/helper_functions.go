package config

import (
	"fmt"
	"reflect"
)

func Diff(options, options2 *Options) []string {
	if reflect.DeepEqual(options, options2) {
		return nil
	}

	var strs []string
	for key := range goToSU2FieldMap {
		v1 := reflect.ValueOf(options).Elem()
		v2 := reflect.ValueOf(options2).Elem()
		item1 := v1.FieldByName(string(key)).Interface()
		item2 := v2.FieldByName(string(key)).Interface()
		if !reflect.DeepEqual(item1, item2) {
			strs = append(strs, fmt.Sprintf(string(key)+": item1 is %v, item2 is %v", item1, item2))
		}
	}
	return strs
}

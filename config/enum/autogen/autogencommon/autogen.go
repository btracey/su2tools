package autogencommon

import (
	"strings"

	"github.com/btracey/su2tools/config/common"
)

type EnumOption struct {
	Identifier string // Variable name
	//Value      int    // Value given in enum
	ConfigString string // What is the lookup in the map
}

type EnumType struct {
	Typename string
	Option   map[string]*EnumOption
	Mapname  string
}

func FixEnumId(enumId string) string {
	enumId = strings.TrimSpace(enumId)
	enumId = strings.ToLower(enumId)
	enumId = common.ToCamelCase(enumId)
	return enumId
}

func FixEnumType(enumId string) string {
	enumId = strings.TrimSpace(enumId)
	enumId = strings.ToLower(enumId)
	enumId = strings.TrimPrefix(enumId, "enum")
	enumId = strings.TrimPrefix(enumId, "_")
	enumId = common.ToCamelCase(enumId)
	return enumId
}

func FixMapname(mapname string) string {
	mapname = strings.TrimSpace(mapname)
	mapname = common.ToCamelCase(mapname)
	return mapname
}

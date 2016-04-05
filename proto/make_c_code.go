package proto

import (
// "fmt"
// "github.com/toophy/easybuff/help"
// "os"
// "os/exec"
// "sort"
// "strconv"
// "strings"
)

func TypeToClang(s string) string {
	type_name := s
	switch s {
	case "int8":
		type_name = "char"
	case "int16":
		type_name = "short"
	case "int24":
		type_name = "int"
	case "int32":
		type_name = "int"
	case "int40":
		type_name = "long long"
	case "int48":
		type_name = "long long"
	case "int56":
		type_name = "long long"
	case "int64":
		type_name = "long long"
	case "uint8":
		type_name = "unsigned char"
	case "uint16":
		type_name = "unsigned short"
	case "uint24":
		type_name = "unsigned int"
	case "uint32":
		type_name = "unsigned int"
	case "uint40":
		type_name = "unsigned long long"
	case "uint48":
		type_name = "unsigned long long"
	case "uint56":
		type_name = "unsigned long long"
	case "uint64":
		type_name = "unsigned long long"
	case "float32":
		type_name = "float"
	case "float64":
		type_name = "double"
	case "string":
		type_name = "std::string"
	}
	return type_name
}

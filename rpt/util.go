package rpt

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
)

func Typeof(v interface{}) string {

	switch v.(type) {
	case int8:
		return "int8"
	case uint8:
		return "int"
	case int16:
		return "int16"
	case uint16:
		return "uint16"
	case int32:
		return "int32"
	case uint32:
		return "uint32"
	case int64:
		return "int64"
	case uint64:
		return "uint64"
	case uint:
		return "uint"
	case uintptr:
		return "uintptr"
	case int:
		return "int"
	case float64:
		return "float64"
	case float32:
		return "float32"
	case complex64:
		return "complex64"
	case complex128:
		return "complex128"
	case string:
		return "string"
	case bool:
		return "bool"
	case map[string]string:
		return "map[string]string"
	case map[int]string:
		return "map[int]string"
	case map[int]int:
		return "map[int]int"
	case map[string]int:
		return "map[string]int"
	case map[string]bool:
		return "map[string]bool"
	case []byte:
		return "[]byte"
	case []string:
		return "[]string"
	case []int:
		return "[]int"
	case error:
		return "error"
	default:
		return "Unknown"
	}
}

func StringTo(s string, to string) (interface{}, error) {
	return nil, nil
}

func ToJSON(i interface{}) []byte {
	out, _ := json.MarshalIndent(i, "", "  ")
	return out
}

func NewGUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	guid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return guid
}

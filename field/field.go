package field

import (
	"fmt"
	"strconv"
	"strings"
)

type Field struct {
	Name        string
	Type        string //"string", "float", "int", "bool"
	Description string
}

func ParseStringToType(v string) string {
	VType := "string"
	if v == "true" || v == "false" {
		VType = "bool"
	}
	if _, err := strconv.Atoi(v); err == nil {
		VType = "int"
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		VType = "float"
	}
	return VType
}
func (f *Field) Validate(v string) bool {
	v = strings.ToLower(v)
	pstt := ParseStringToType(v)
	return ParseStringToType(v) == f.Type || (f.Type == "string" && pstt == "bool")
}

func (f *Field) TypeValidationError() string {
	return fmt.Sprintf("FieldTypeValidationError: The field %s was given the wrong type. Expected: %s.", f.Name, f.Type)
}

func (f *Field) MissingFieldError() string {
	return fmt.Sprintf("FieldMissingError: The field %s was not provided.", f.Name)
}

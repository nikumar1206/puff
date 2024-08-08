package puff

import "fmt"

func FieldTypeError(value string, expectedType string) error {
	return fmt.Errorf("type error: the value %s cant be used as the expected type %s.", value, expectedType)
}

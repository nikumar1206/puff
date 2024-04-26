package field

import "reflect"

type Field struct {
	Name        string
	Type        reflect.Type
	Description string
}

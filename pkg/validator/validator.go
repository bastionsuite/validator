package validator

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrNotImplemented = errors.New("validator: function not implemented")

func ParseRequestForm() {
}

func ParseRequestJSONBody() {
}

func ParseStruct(target interface{}, source interface{}) error {
	t := reflect.TypeOf(target)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fmt.Println(t)

	s := reflect.TypeOf(source)

	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}

	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		fmt.Println(i, field.Name, field.Type, field.Tag)
	}

	numFields = s.NumField()
	for i := 0; i < numFields; i++ {
		field := s.Field(i)
		fmt.Println(i, field.Name, field.Type)
	}

	fmt.Println(target, source, t)
	return ErrNotImplemented
}

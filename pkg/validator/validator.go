// Copyright 2020 BastionSuite. All rights resreved.
// Use of this source code is governed by a MIT license that can be
// found in the LICENSE file

// Package validator implements input validation, allowing to verify if
// input in a JSON body, HTTP Form values or any struct to be checked
// against certain rules.
package validator

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrNotImplemented = errors.New("validator: function not implemented")
var ErrNotAStruct = errors.New("validator: type must be a struct")

func ParseRequestForm() {
}

func ParseRequestJSONBody() {
}

// ParseStruct parses the target struct tags and adopt values it can find from
// source.
func ParseStruct(target interface{}, source interface{}) error {
	t := reflect.TypeOf(target)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	//fmt.Println(t)

	s := reflect.TypeOf(source)

	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	if s.Kind() != reflect.Struct {
		return ErrNotAStruct
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

	//fmt.Println(target, source, t)
	return nil
	// return ErrNotImplemented
}

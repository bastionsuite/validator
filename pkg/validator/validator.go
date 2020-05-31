// Copyright 2020 BastionSuite. All rights resreved.
// Use of this source code is governed by a MIT license that can be
// found in the LICENSE file

// Package validator implements input validation, allowing to verify if
// input in a JSON body, HTTP Form values or any struct to be checked
// against certain rules.
package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

var ErrNotImplemented = errors.New("validator: function not implemented")
var ErrNotAStruct = errors.New("validator: type must be a struct")

func ParseRequestForm(target interface{}, source url.Values) error {
	t := reflect.ValueOf(target).Elem()

	fields, err := listFields(target)
	if err != nil {
		return err
	}

	for _, field := range fields {
		var tag ValidationOptions

		valField := t.FieldByName(field.Name)

		sourceField := field.Name
		required := false

		if string(field.Tag) != "" {
			tag, err = parseTag(string(field.Tag))
			if err != nil {
				return err
			}
			sourceField = tag.Source
			required = tag.Required
		}

		value := source.Get(sourceField)
		if value == "" && required {
			return fmt.Errorf("ParseRequestForm: missing or empty required field '%s'", sourceField)
		}

		switch field.Type.Kind() {
		case reflect.String:
			if !execValidate(value, tag.Validators) {
				return fmt.Errorf("Invalid value for %s: %v", field.Name, value)
			}
			valField.SetString(value)
		case reflect.Int:
			i, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("Could not convert value to Int: %+v", err)
			}
			if !execValidate(i, tag.Validators) {
				return fmt.Errorf("Invalid value for %s: %v", field.Name, value)
			}
			valField.SetInt(int64(i))
		case reflect.Float64:
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("Could not convert value to float: %+v", err)
			}
			if !execValidate(f, tag.Validators) {
				return fmt.Errorf("Invalid value for %s: %v", field.Name, value)
			}
			valField.SetFloat(f)
		case reflect.Struct:
			// This case happens when the object is not a simple primitive.
			switch field.Type.Name() {
			case "time.Time":
				parsedTime, err := time.Parse("2006-01-02", value)
				if err != nil {
					return fmt.Errorf("Could not conver value to time.Time: %+v", err)
				}
				valField.Set(reflect.ValueOf(parsedTime))
			default:
				fmt.Errorf("ParseJSONBody: Unsupport type in target struct: %+v", field.Type.Name())
			}
		default:
			return fmt.Errorf("ParseJSONBody: Unsupported type in target struct: %+v", field.Type.Kind())
		}
	}

	return nil
}

func ParseJSONBody(target interface{}, source []byte) error {
	var src map[string]interface{}
	err := json.Unmarshal(source, &src)
	if err != nil {
		return err
	}

	t := reflect.ValueOf(target).Elem()

	fields, err := listFields(target)
	if err != nil {
		return err
	}

	for _, field := range fields {
		valField := t.FieldByName(field.Name)
		sourceField := field.Name
		required := false
		var tag ValidationOptions
		if string(field.Tag) != "" {
			tag, err = parseTag(string(field.Tag))
			if err != nil {
				return err
			}
			sourceField = tag.Source
			required = tag.Required
		}

		value, ok := src[sourceField]
		if required && !ok {
			return fmt.Errorf("ParseJSONBody: missing required field '%s'", sourceField)
		}

		switch field.Type.Kind() {
		case reflect.String:
			strVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("Couldn't cast a thing to string. Odd.")
			}
			if len(tag.Validators) > 0 {
				if !execValidate(strVal, tag.Validators) {
					return fmt.Errorf("Invalid value for %s: %v", field.Name, value)
				}
			}
			valField.SetString(strVal)
		case reflect.Int:
			i, ok := value.(float64) // We cast to float, because JSON.
			if !ok {
				return fmt.Errorf("ParseJSONBody: Error during cast to int. Value: %+v", value)
			}
			if len(tag.Validators) > 0 {
				if !execValidate(int(i), tag.Validators) {
					return fmt.Errorf("Invalid value for %s: %v", field.Name, value)
				}
			}
			iVal := int64(i) // And then we convert it back to int, because JSON.
			valField.SetInt(iVal)
		case reflect.Float64:
			i, ok := value.(float64)
			if !ok {
				return fmt.Errorf("ParseJSONBody: Error during cast to float. Value: %+v", value)
			}
			if len(tag.Validators) > 0 {
				if !execValidate(i, tag.Validators) {
					return fmt.Errorf("Invalid value for %s: %v", field.Name, value)
				}
			}
			valField.SetFloat(i)
		default:
			return fmt.Errorf("ParseJSONBody: Unsupported type in target struct: %+v", field.Type.Kind())
		}
	}

	return nil
}

func listFields(target interface{}) ([]reflect.StructField, error) {
	var fields []reflect.StructField

	t := reflect.TypeOf(target)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return fields, ErrNotAStruct
	}

	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i))
	}
	return fields, nil
}

func execValidate(in interface{}, validators []ValidateFunc) bool {
	for _, validate := range validators {
		res := validate(in)
		if !res {
			return false
		}
	}
	return true
}

// Copyright 2020 BastionSuite. All rights resreved.
// Use of this source code is governed by a MIT license that can be
// found in the LICENSE file

// Most of this is an interpretation of the golang json tags
// implementation. Because that's exactly what we need, but their
// implementation is private

package validator

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var ErrInvalidQualifier = errors.New("validator/tags: invalid qualifier")
var ErrDuplicateQualifier = errors.New("validator/tags: duplicate qualifier")
var ErrInvalidValidation = errors.New("validator/tags: unknown validation")

var reQualifiers = regexp.MustCompile(`([a-z]+)(:"([^"]+)")?`)
var reValidators = regexp.MustCompile(`([a-z]+)(=([^,]+))?`)

type ValidateFunc func(string) bool

type ValidationOptions struct {
	Source     string
	Required   bool
	Validators []ValidateFunc
}

func parseTag(tag string) (ValidationOptions, error) {
	r := ValidationOptions{}

	qualifiers := reQualifiers.FindAllStringSubmatch(tag, -1)
	hadSource := false
	hadValidations := false
	for _, q := range qualifiers {
		switch {
		case q[1] == "source":
			if hadSource {
				return r, ErrDuplicateQualifier
			}
			r.Source = q[3]
			hadSource = true
		case q[1] == "validate":
			if hadValidations {
				return r, ErrDuplicateQualifier
			}

			hadValidations = true

			validators, err := parseValidations(q[3])
			if err != nil {
				return r, err
			}
			r.Validators = validators
		case q[1] == "required":
			r.Required = true
		default:
			return r, ErrInvalidQualifier
		}
	}

	return r, nil
}

func parseValidations(validations string) ([]ValidateFunc, error) {
	var v []ValidateFunc
	valis := reValidators.FindAllStringSubmatch(validations, -1)

	for _, validation := range valis {
		switch {
		case validation[1] == "length":
			range_ := strings.SplitN(validation[3], "-", 2)
			min, err := strconv.Atoi(range_[0])
			if err != nil {
				return v, err
			}
			max, err := strconv.Atoi(range_[1])
			if err != nil {
				return v, err
			}
			v = append(v, func(in string) bool {
				length := len(in)
				return length >= min && length <= max
			})
		case validation[1] == "range":
			range_ := strings.SplitN(validation[3], "-", 2)
			min, err := strconv.Atoi(range_[0])
			if err != nil {
				return v, err
			}
			max, err := strconv.Atoi(range_[1])
			if err != nil {
				return v, err
			}

			v = append(v, func(in string) bool {
				n, err := strconv.Atoi(in)
				if err != nil {
					return false
				}

				return n >= min && n <= max
			})
		case validation[1] == "string":
			v = append(v, func(in string) bool {
				return true
			})
		case validation[1] == "int":
			v = append(v, func(in string) bool {
				_, err := strconv.Atoi(in)
				if err != nil {
					return false
				}
				return true
			})
		case validation[1] == "float":
			v = append(v, func(in string) bool {
				_, err := strconv.ParseFloat(in, 64)
				return err == nil
			})
		default:
			return v, ErrInvalidValidation
		}
	}

	return v, nil
}
